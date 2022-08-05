package store

import (
	"context"

	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
)

// Store provides access and update functions for secrets.
//
//go:generate mockgen-wrapper
type Store interface {
	Count(ctx context.Context) (int, error)
	Get(ctx context.Context, id string) (*storage.Secret, bool, error)
	GetByQuery(ctx context.Context, q *v1.Query) ([]*storage.Secret, error)
	GetMany(ctx context.Context, ids []string) ([]*storage.Secret, []int, error)
	Walk(context.Context, func(secret *storage.Secret) error) error

	Upsert(ctx context.Context, secret *storage.Secret) error
	Delete(ctx context.Context, id string) error
}
