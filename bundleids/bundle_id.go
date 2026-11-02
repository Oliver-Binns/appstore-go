package bundleids

import (
	"github.com/oliver-binns/appstore-go/internal/ptr"
	"github.com/oliver-binns/appstore-go/openapi"
)

// BundleID is the domain struct for a bundle ID.
type BundleID struct {
	ID         string
	Name       string
	Identifier string
	Platform   Platform
	SeedID     string
}

// Platform is an alias for the generated openapi.BundleIdPlatform type.
type Platform = openapi.BundleIdPlatform

func fromResponse(b openapi.BundleId) *BundleID {
	if b.Attributes == nil {
		return &BundleID{ID: b.Id}
	}
	return &BundleID{
		ID:         b.Id,
		Name:       ptr.Deref(b.Attributes.Name),
		Identifier: ptr.Deref(b.Attributes.Identifier),
		Platform:   ptr.Deref(b.Attributes.Platform),
		SeedID:     ptr.Deref(b.Attributes.SeedId),
	}
}
