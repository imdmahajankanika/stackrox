package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/pkg/clientconn"
	"golang.org/x/net/context"
)

const tokenEnv = "STACKROX_TOKEN"

var (
	passFail = flag.Bool("pass-fail", true, "exit 1 on any critical policy failures")
	central  = flag.String("central", "localhost:8080", "endpoint where central is available")
	digest   = flag.String("digest", "", "the sha256 digest for the image")
	registry = flag.String("registry", "", "registry where the image is uploaded")
	remote   = flag.String("remote", "", "the remote name of the image")
	tag      = flag.String("tag", "", "the tag for the image")
)

func main() {
	// Parse the input flags.
	flag.Parse()

	// Read token from ENV.
	token, exists := os.LookupEnv(tokenEnv)
	if !exists {
		fmt.Println(fmt.Errorf("the STACKROX_TOKEN environment variable must be set to a token generated by stackrox for a Remote host"))
		return
	}

	// Get the violated policies for the input data.
	violatedPolicies, err := getViolatedPolicies(token)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Run requested operations.
	if passFail != nil && *passFail {
		err = runPassFail(violatedPolicies)
	} else {
		err = runInformational(violatedPolicies)
	}
	if err != nil {
		os.Exit(1)
	}
}

// runPassFail runs the pipeline in pass/fail mode.
func runPassFail(violatedPolicies []*v1.Policy) error {
	fmt.Println("----------------BEGIN STACKROX CI---------------")
	defer fmt.Println("----------------END STACKROX CI---------------")

	// Print violated policy names.
	var failed bool
	if len(violatedPolicies) > 0 {
		fmt.Println("Policies failed: ")
		for _, policy := range violatedPolicies {
			if policy.GetSeverity() >= v1.Severity_CRITICAL_SEVERITY {
				fmt.Println(policy.Name, " is a critical vulnerability")
				failed = true
			} else {
				fmt.Println(policy.Name)
			}
		}
	} else {
		fmt.Println("no policy violations found")
	}

	// If any are failing conditions (CRITICAL severity), print a message and exit.
	if failed {
		return fmt.Errorf("critical vulnerability encountered, failing")
	}
	return nil
}

// runInformational just pipes failed policies out on STDOUT in JSON format.
func runInformational(violatedPolicies []*v1.Policy) error {
	// Just pipe out the violated policies as JSON.
	for _, policy := range violatedPolicies {
		jsonified, err := json.Marshal(policy)
		if err == nil {
			os.Stdout.Write([]byte(jsonified))
		} else {
			return fmt.Errorf("recieved unexpected output, failing")
		}
	}
	return nil
}

// Fetch the alerts for the inputs and convert them to a list of Policies that are violated.
func getViolatedPolicies(token string) ([]*v1.Policy, error) {
	alerts, err := getAlerts(token)
	if err != nil {
		return nil, err
	}

	var policies []*v1.Policy
	for _, alert := range alerts {
		policies = append(policies, alert.GetPolicy())
	}
	return policies, nil
}

// Get the alerts for the command line inputs.
func getAlerts(token string) ([]*v1.Alert, error) {
	// Attempt to construct the request first since it is the cheapest op.
	image, err := buildRequest()
	if err != nil {
		return nil, err
	}

	// Create the connection to the central detection service.
	conn, err := clientconn.UnauthenticatedGRPCConnection(*central)
	if err != nil {
		return nil, err
	}
	service := v1.NewDetectionServiceClient(conn)

	// Build context with token header.
	md := metautils.NiceMD{}
	md = md.Add("authorization", token)
	ctx := md.ToOutgoing(context.Background())

	// Call detection and return the returned alerts.
	response, err := service.DetectBuildTime(ctx, image)
	if err != nil {
		return nil, err
	}
	return response.GetAlerts(), nil
}

// Use inputs to generate an image name for request.
func buildRequest() (*v1.Image, error) {
	im := v1.ImageName{
		Remote:   *remote,
		Registry: *registry,
	}

	switch {
	case *registry == "":
		return nil, fmt.Errorf("image registry must be set, or we don't know where to get metadata from")
	case *remote == "":
		return nil, fmt.Errorf("image remote must be set, or we don't know which image in the registry to process")
	case (*digest == "" && *tag == "") || (*digest != "" && *tag != ""):
		return nil, fmt.Errorf("one of image digest or tag must be set, or we don't know which version of the image to process")
	case *digest != "":
		im.Sha = *digest
	case *tag != "":
		im.Tag = *tag
	}

	return &v1.Image{
		Name: &im,
	}, nil
}
