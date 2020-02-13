package authproviders

import (
	"context"
	"net/http"
)

// BackendFactory is responsible for creating Backends.
type BackendFactory interface {
	// CreateAuthProviderBackend creates a new backend instance for the given auth provider, using the specified
	// configuration.
	CreateBackend(ctx context.Context, id string, uiEndpoints []string, config map[string]string) (Backend, error)

	// ProcessHTTPRequest is the dispatcher for HTTP/1.1 requests to `<sso-prefix>/<provider-type>/...`. The envisioned
	// workflow consists of extracting the specific auth provider ID from the request, usually via a `state` parameter,
	// and returning this provider ID from the function (with the Registry taking care of forwarding the request to that
	// provider's HTTP handler). If there are any provider-independent HTTP endpoints (such as the SP metadata for
	// SAML), this can be handled in this function as well - an empty provider ID along with a nil error needs to be
	// returned in that case.
	ProcessHTTPRequest(w http.ResponseWriter, r *http.Request) (providerID string, err error)

	// ResolveProvider takes care of looking up the provider ID from an (opaque) state string.
	ResolveProvider(state string) (providerID string, err error)
}

// BackendFactoryCreator is a function for creating a BackendFactory, given a base URL (excluding trailing slashes) for
// those factory's HTTP handlers.
type BackendFactoryCreator func(urlPathPrefix string) BackendFactory
