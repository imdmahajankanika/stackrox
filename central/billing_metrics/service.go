package billingmetrics

import (
	"context"

	v1 "github.com/stackrox/rox/generated/api/v1"
	pkgGRPC "github.com/stackrox/rox/pkg/grpc"
)

// Service provides the interface to the svc that handles API keys.
type Service interface {
	pkgGRPC.APIService
	AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error)

	v1.MaximumValueServiceServer
}
