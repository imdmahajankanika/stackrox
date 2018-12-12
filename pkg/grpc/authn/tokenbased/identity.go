package tokenbased

import (
	"time"

	"github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/auth/authproviders"
)

type roleBasedIdentity struct {
	uid          string
	friendlyName string
	role         *v1.Role
	expiry       time.Time
}

func (i *roleBasedIdentity) UID() string {
	return i.uid
}

func (i *roleBasedIdentity) FriendlyName() string {
	return i.friendlyName
}

func (i *roleBasedIdentity) Role() *v1.Role {
	return i.role
}

func (i *roleBasedIdentity) Service() *storage.ServiceIdentity {
	return nil
}

func (i *roleBasedIdentity) Expiry() time.Time {
	return i.expiry
}

func (i *roleBasedIdentity) ExternalAuthProvider() authproviders.AuthProvider {
	return nil
}
