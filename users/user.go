package users

import openapi_types "github.com/oapi-codegen/runtime/types"

type User struct {
	ID                  string     `json:"-"`
	FirstName           string     `json:"firstName,omitempty"`
	LastName            string     `json:"lastName,omitempty"`
	Username            string     `json:"username,omitempty"`
	Roles               []UserRole `json:"roles,omitempty"`
	AllAppsVisible      bool       `json:"allAppsVisible,omitempty"`
	ProvisioningAllowed bool       `json:"provisioningAllowed,omitempty"`
	HasAcceptedInvite   bool       `json:"userHasAcceptedInvitation,omitempty"`

	VisibleAppIDs []string `json:"-"`
}

func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func derefEmail(e *openapi_types.Email) string {
	if e == nil {
		return ""
	}
	return string(*e)
}

func derefBool(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}

func derefRoles(r *[]UserRole) []UserRole {
	if r == nil {
		return nil
	}
	return *r
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

func rolesOrNil(roles []UserRole) *[]UserRole {
	if len(roles) == 0 {
		return nil
	}
	return &roles
}
