package service

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	rolesDataStore "github.com/stackrox/rox/central/rbac/k8srole/datastore"
	roleBindingsDataStore "github.com/stackrox/rox/central/rbac/k8srolebinding/datastore"
	"github.com/stackrox/rox/central/role/resources"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/auth/permissions"
	"github.com/stackrox/rox/pkg/grpc/authz"
	"github.com/stackrox/rox/pkg/grpc/authz/perrpc"
	"github.com/stackrox/rox/pkg/grpc/authz/user"
	"github.com/stackrox/rox/pkg/search"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	authorizer = perrpc.FromMap(map[authz.Authorizer][]string{
		user.With(permissions.View(resources.K8sRole)): {
			"/v1.RbacService/GetRole",
			"/v1.RbacService/ListRoles",
		},
		user.With(permissions.View(resources.K8sRoleBinding)): {
			"/v1.RbacService/GetRoleBinding",
			"/v1.RbacService/ListRoleBindings",
		},
	})
)

// serviceImpl provides APIs for k8s rbac objects.
type serviceImpl struct {
	roles    rolesDataStore.DataStore
	bindings roleBindingsDataStore.DataStore
}

// RegisterServiceServer registers this service with the given gRPC Server.
func (s *serviceImpl) RegisterServiceServer(grpcServer *grpc.Server) {
	v1.RegisterRbacServiceServer(grpcServer, s)
}

// RegisterServiceHandler registers this service with the given gRPC Gateway endpoint.
func (s *serviceImpl) RegisterServiceHandler(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return v1.RegisterRbacServiceHandler(ctx, mux, conn)
}

// AuthFuncOverride specifies the auth criteria for this API.
func (s *serviceImpl) AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error) {
	return ctx, authorizer.Authorized(ctx, fullMethodName)
}

// GetRole returns the k8s role for the id.
func (s *serviceImpl) GetRole(ctx context.Context, request *v1.ResourceByID) (*v1.GetRoleResponse, error) {
	role, exists, err := s.roles.GetRole(request.GetId())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if !exists {
		return nil, status.Errorf(codes.NotFound, "k8s role with id '%q' does not exist", request.GetId())
	}

	return &v1.GetRoleResponse{Role: role}, nil
}

// ListRoles returns all roles that match the query.
func (s *serviceImpl) ListRoles(ctx context.Context, rawQuery *v1.RawQuery) (*v1.ListRolesResponse, error) {
	var roles []*storage.K8SRole
	var err error
	if rawQuery.GetQuery() == "" {
		roles, err = s.roles.ListRoles()
	} else {
		var q *v1.Query
		q, err = search.ParseRawQueryOrEmpty(rawQuery.GetQuery())
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		roles, err = s.roles.SearchRawRoles(q)
	}
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to retrieve k8s roles: %v", err)
	}

	return &v1.ListRolesResponse{Roles: roles}, nil
}

// GetRole returns the k8s role binding for the id.
func (s *serviceImpl) GetRoleBinding(ctx context.Context, request *v1.ResourceByID) (*v1.GetRoleBindingResponse, error) {
	binding, exists, err := s.bindings.GetRoleBinding(request.GetId())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if !exists {
		return nil, status.Errorf(codes.NotFound, "k8s role binding with id '%q' does not exist", request.GetId())
	}

	return &v1.GetRoleBindingResponse{Binding: binding}, nil
}

// ListRoleBindings returns all role bindings that match the query.
func (s *serviceImpl) ListRoleBindings(ctx context.Context, rawQuery *v1.RawQuery) (*v1.ListRoleBindingsResponse, error) {
	var bindings []*storage.K8SRoleBinding
	var err error
	if rawQuery.GetQuery() == "" {
		bindings, err = s.bindings.ListRoleBindings()
	} else {
		var q *v1.Query
		q, err = search.ParseRawQueryOrEmpty(rawQuery.GetQuery())
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		bindings, err = s.bindings.SearchRawRoleBindings(q)
	}
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to retrieve k8s role bindings: %v", err)
	}

	return &v1.ListRoleBindingsResponse{Bindings: bindings}, nil
}
