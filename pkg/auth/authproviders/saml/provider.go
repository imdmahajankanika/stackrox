package saml

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/russellhaering/gosaml2"
	"github.com/stackrox/rox/pkg/auth/tokens"
	"github.com/stackrox/rox/pkg/httputil"
	"github.com/stackrox/rox/pkg/logging"
)

var (
	log = logging.LoggerForModule()
)

type provider struct {
	acsURLPath string
	sp         saml2.SAMLServiceProvider
	id         string
}

func (p *provider) loginURL(clientState string) (string, error) {
	doc, err := p.sp.BuildAuthRequestDocument()
	if err != nil {
		return "", fmt.Errorf("could not construct auth request: %v", err)
	}
	authURL, err := p.sp.BuildAuthURLRedirect(makeState(p.id, clientState), doc)
	if err != nil {
		return "", fmt.Errorf("could not construct auth URL: %v", err)
	}
	return authURL, nil
}

func newProvider(ctx context.Context, acsURLPath string, id, uiEndpoint string, config map[string]string) (*provider, map[string]string, error) {
	p := &provider{
		acsURLPath: acsURLPath,
		id:         id,
	}

	acsURL := &url.URL{
		Scheme: "https",
		Host:   uiEndpoint,
		Path:   acsURLPath,
	}
	p.sp.AssertionConsumerServiceURL = acsURL.String()

	spIssuer := config["sp_issuer"]
	if spIssuer == "" {
		return nil, nil, errors.New("no ServiceProvider issuer specified")
	}
	p.sp.ServiceProviderIssuer = spIssuer

	effectiveConfig := map[string]string{
		"sp_issuer": spIssuer,
	}

	if config["idp_metadata_url"] != "" {
		if config["idp_issuer"] != "" || config["idp_cert_pem"] != "" || config["idp_sso_url"] != "" {
			return nil, nil, errors.New("if IdP metadata URL is set, IdP issuer, SSO URL and certificate data must be left blank")
		}
		if err := configureIDPFromMetadataURL(ctx, &p.sp, config["idp_metadata_url"]); err != nil {
			return nil, nil, fmt.Errorf("could not configure auth provider from IdP metadata URL: %v", err)
		}
		effectiveConfig["idp_metadata_url"] = config["idp_metadata_url"]
	} else {
		if config["idp_issuer"] == "" || config["idp_sso_url"] == "" || config["idp_cert_pem"] == "" {
			return nil, nil, errors.New("if IdP metadata URL is not set, IdP issuer, SSO URL, and certificate data must be specified")
		}
		if err := configureIDPFromSettings(&p.sp, config["idp_issuer"], config["idp_sso_url"], config["idp_cert_pem"]); err != nil {
			return nil, nil, fmt.Errorf("could not configure auth provider from settings: %v", err)
		}
		effectiveConfig["idp_issuer"] = config["idp_issuer"]
		effectiveConfig["idp_sso_url"] = config["idp_sso_url"]
		effectiveConfig["idp_cert_pem"] = config["idp_cert_pem"]
	}

	return p, effectiveConfig, nil
}

func (p *provider) consumeSAMLResponse(samlResponse string) (*tokens.ExternalUserClaim, []tokens.Option, error) {
	ai, err := p.sp.RetrieveAssertionInfo(samlResponse)
	if err != nil {
		return nil, nil, err
	}

	var opts []tokens.Option
	if ai.SessionNotOnOrAfter != nil {
		opts = append(opts, tokens.WithExpiry(*ai.SessionNotOnOrAfter))
	}

	claim := saml2AssertionInfoToExternalClaim(ai)
	return claim, opts, nil
}

func (p *provider) ProcessHTTPRequest(w http.ResponseWriter, r *http.Request) (*tokens.ExternalUserClaim, []tokens.Option, string, error) {
	if r.URL.Path != p.acsURLPath {
		return nil, nil, "", httputil.NewError(http.StatusNotFound, "Not Found")
	}
	if r.Method != http.MethodPost {
		return nil, nil, "", httputil.NewError(http.StatusMethodNotAllowed, "Method Not Allowed")
	}

	samlResponse := r.FormValue("SAMLResponse")
	if samlResponse == "" {
		return nil, nil, "", httputil.NewError(http.StatusBadRequest, "no SAML response transmitted")
	}

	claims, opts, err := p.consumeSAMLResponse(samlResponse)
	if err != nil {
		return nil, nil, "", err
	}

	relayState := r.FormValue("RelayState")
	_, clientState := splitState(relayState)

	return claims, opts, clientState, err
}

func (p *provider) ExchangeToken(ctx context.Context, externalToken, state string) (*tokens.ExternalUserClaim, []tokens.Option, string, error) {
	return nil, nil, "", errors.New("not implemented")
}

func (p *provider) RefreshURL() string {
	return ""
}

func (p *provider) LoginURL(clientState string) string {
	url, err := p.loginURL(clientState)
	if err != nil {
		log.Errorf("could not obtain the login URL: %v", err)
	}
	return url
}

// Helpers
//////////

func saml2AssertionInfoToExternalClaim(assertionInfo *saml2.AssertionInfo) *tokens.ExternalUserClaim {
	claim := &tokens.ExternalUserClaim{
		UserID: assertionInfo.NameID,
	}
	claim.Attributes = make(map[string][]string)
	claim.Attributes["userid"] = []string{claim.UserID}

	// We store claims as both friendly name and name for easy of use.
	for _, value := range assertionInfo.Values {
		for _, innerValue := range value.Values {
			if value.Name != "" {
				claim.Attributes[value.Name] = append(claim.Attributes[value.Name], innerValue.Value)
			}
			if value.FriendlyName != "" {
				claim.Attributes[value.FriendlyName] = append(claim.Attributes[value.FriendlyName], innerValue.Value)
			}
		}
	}
	return claim
}
