package splunk

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/golang/protobuf/jsonpb"
	"github.com/stackrox/rox/central/notifiers"
	"github.com/stackrox/rox/generated/internalapi/wrapper"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/logging"
	"github.com/stackrox/rox/pkg/protoutils"
	"github.com/stackrox/rox/pkg/urlfmt"
)

const (
	source                    = "stackrox"
	sourceType                = "_json"
	splunkHECDefaultDataLimit = 10000
)

var (
	log = logging.LoggerForModule()

	// blacklist of annotations to be scrubbed
	scrubAnnotations = map[string]bool{
		"kubectl.kubernetes.io/last-applied-configuration": true,
	}
)

type splunk struct {
	token    string
	endpoint string
	insecure bool
	truncate int
	*storage.Notifier
}

func (s *splunk) AlertNotify(alert *storage.Alert) error {
	return s.postAlert(alert)
}

func (s *splunk) NetworkPolicyYAMLNotify(yaml string, clusterName string) error {
	return nil
}

func (s *splunk) ProtoNotifier() *storage.Notifier {
	return s.Notifier
}

func (s *splunk) Test() error {
	alert := &storage.Alert{
		Policy:     &storage.Policy{Name: "Test Policy"},
		Deployment: &storage.Deployment{Name: "Test Deployment"},
		Violations: []*storage.Alert_Violation{
			{Message: "This is a sample Splunk alert message created to test integration with StackRox."},
		},
	}
	return s.postAlert(alert)
}

func (s *splunk) AckAlert(alert *storage.Alert) error {
	return nil
}

func (s *splunk) ResolveAlert(alert *storage.Alert) error {
	return nil
}

func (s *splunk) postAlert(alert *storage.Alert) error {
	clonedAlert := protoutils.CloneStorageAlert(alert)

	// Splunk's HEC by default has a limitation of data size == 10KB
	// Removing some of the fields here to make it smaller
	// More details on HEC limitation: https://developers.perfectomobile.com/display/TT/Splunk+-+Configure+HTTP+Event+Collector
	// Check section on "Increasing the Event Data Truncate Limit"
	clonedAlert.GetDeployment().Risk = nil
	for i := range clonedAlert.GetDeployment().GetContainers() {
		clonedAlert.GetDeployment().Containers[i].GetImage().Metadata = nil
		clonedAlert.GetDeployment().Containers[i].GetImage().Scan = nil
	}

	processViolations := clonedAlert.GetProcessViolation().GetProcesses()
	if len(processViolations) > 5 {
		clonedAlert.ProcessViolation.Processes = clonedAlert.ProcessViolation.Processes[0:5]
	}

	// Scrub black listed annotations
	for needScrubbing := range clonedAlert.GetDeployment().GetAnnotations() {
		if _, ok := scrubAnnotations[needScrubbing]; ok {
			delete(clonedAlert.Deployment.Annotations, needScrubbing)
		}
	}

	splunkEvent := &wrapper.SplunkEvent{
		Event:      clonedAlert,
		Source:     source,
		Sourcetype: sourceType,
	}

	var jsonPayload bytes.Buffer
	err := new(jsonpb.Marshaler).Marshal(&jsonPayload, splunkEvent)
	if err != nil {
		return err
	}

	if s.truncate == 0 {
		s.truncate = splunkHECDefaultDataLimit
	}

	if jsonPayload.Len() > s.truncate {
		return fmt.Errorf("Splunk HEC truncate data limit (%d bytes) exceeded: %d", s.truncate, jsonPayload.Len())
	}

	req, err := s.createSplunkHTTPRequest(&jsonPayload)
	if err != nil {
		return err
	}

	resp, err := s.sendHTTPPayload(req)
	if err != nil {
		return err
	}
	if resp != nil {
		defer func() {
			_ = resp.Body.Close()
		}()
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP Status Code: %d", resp.StatusCode)
	}
	return nil
}

func (s *splunk) createSplunkHTTPRequest(jsonPayload *bytes.Buffer) (*http.Request, error) {
	req, err := http.NewRequest("POST", s.endpoint, jsonPayload)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Splunk %s", s.token))
	return req, err
}

func (s *splunk) sendHTTPPayload(req *http.Request) (*http.Response, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: s.insecure},
	}

	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)
	return resp, err
}

func init() {
	notifiers.Add("splunk", func(notifier *storage.Notifier) (notifiers.Notifier, error) {
		s, err := newSplunk(notifier)
		return s, err
	})
}

func newSplunk(notifier *storage.Notifier) (*splunk, error) {
	splunkConfig, ok := notifier.GetConfig().(*storage.Notifier_Splunk)
	if !ok {
		return nil, fmt.Errorf("Splunk configuration required")
	}
	conf := splunkConfig.Splunk
	if err := validate(conf); err != nil {
		return nil, err
	}
	endpoint, err := urlfmt.FormatURL(conf.GetHttpEndpoint(), urlfmt.HTTPS, urlfmt.NoTrailingSlash)
	if err != nil {
		return nil, err
	}

	truncate := 0
	if conf.GetTruncate() == 0 {
		truncate = splunkHECDefaultDataLimit
	}
	return &splunk{
		conf.HttpToken,
		endpoint,
		conf.GetInsecure(),
		int(truncate),
		notifier}, nil
}

func validate(conf *storage.Splunk) error {
	if len(conf.HttpToken) == 0 {
		return fmt.Errorf("Splunk HTTP Event Collector(HEC) token must be specified")
	}
	if len(conf.HttpEndpoint) == 0 {
		return fmt.Errorf("Splunk HTTP endpoint must be specified")
	}
	return nil
}
