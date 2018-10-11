// Package usermapper returns a mapper that maps user ids to roles.
package usermapper

import (
	"github.com/stackrox/rox/central/role"
	"github.com/stackrox/rox/pkg/auth/permissions"
	"github.com/stackrox/rox/pkg/auth/tokenbased"
)

// Currently, we don't really have a notion of identities for human users.
// So we return a mapper that gives any human user all access to the system.
type allAccessMapper struct {
	roleStore permissions.RoleStore
}

func (a *allAccessMapper) Role(id string) permissions.Role {
	return a.roleStore.RoleByName(role.Admin)
}

// New returns a new instance of the mapper.
func New(roleStore permissions.RoleStore) tokenbased.RoleMapper {
	return &allAccessMapper{roleStore: roleStore}
}
