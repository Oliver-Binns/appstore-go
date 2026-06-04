package users

import (
	"github.com/oliver-binns/appstore-go/openapi"
)

type User struct {
	ID                  string           `json:"-"`
	FirstName           string           `json:"firstName,omitempty"`
	LastName            string           `json:"lastName,omitempty"`
	Username            string           `json:"username,omitempty"`
	Roles               []openapi.UserRole `json:"roles,omitempty"`
	AllAppsVisible      bool             `json:"allAppsVisible,omitempty"`
	ProvisioningAllowed bool             `json:"provisioningAllowed,omitempty"`
	HasAcceptedInvite   bool             `json:"userHasAcceptedInvitation,omitempty"`

	VisibleAppIDs []string `json:"-"`
}

// boolPtrOrNil returns a pointer to b when b is true and nil when b is false.
// This mirrors the previous behaviour of `bool` fields with `omitempty` JSON tags,
// where zero (false) values were omitted from requests. A nil pointer on a
// `*bool, omitempty` field is also omitted, so false values are not sent.
// If explicit false values need to be included in future, this helper should
// return &b unconditionally.
func boolPtrOrNil(b bool) *bool {
	if b {
		return &b
	}
	return nil
}

func rolesOrNil(roles []openapi.UserRole) *[]openapi.UserRole {
	if len(roles) == 0 {
		return nil
	}
	return &roles
}
