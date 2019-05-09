package authn

import (
	"context"

	"github.com/stackrox/rox/pkg/grpc/requestinfo"
)

// IdentityExtractor extracts the identity of a user making a request from a request info.
type IdentityExtractor interface {
	IdentityForRequest(ctx context.Context, ri requestinfo.RequestInfo) (Identity, error)
}

type extractorList []IdentityExtractor

func (l extractorList) IdentityForRequest(ctx context.Context, ri requestinfo.RequestInfo) (Identity, error) {
	for _, extractor := range l {
		if id, err := extractor.IdentityForRequest(ctx, ri); id != nil || err != nil {
			return id, err
		}
	}
	return nil, nil
}

// CombineExtractors combines the given identity extractors.
func CombineExtractors(extractors ...IdentityExtractor) IdentityExtractor {
	if len(extractors) == 1 {
		return extractors[0]
	}
	return extractorList(extractors)
}
