package users

import "github.com/oliver-binns/appstore-go/openapi"

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

// visibleAppsLinkages returns the visible-apps relationship payload used in
// write requests. It returns nil when all apps are visible and no explicit
// list is required.
func (u *User) visibleAppsLinkages() *openapi.UserVisibleAppsLinkages {
	if len(u.VisibleAppIDs) == 0 && u.AllAppsVisible {
		return nil
	}

	data := make([]openapi.AppRelationship, len(u.VisibleAppIDs))
	for i, id := range u.VisibleAppIDs {
		data[i] = openapi.AppRelationship{Id: id, Type: "apps"}
	}
	return &openapi.UserVisibleAppsLinkages{Data: data}
}

// visibleAppIDs extracts the slice of app IDs from a UserRelationships
// response object, returning an empty slice when no data is present.
func visibleAppIDs(r *openapi.UserRelationships) []string {
	if r == nil || r.VisibleApps == nil || r.VisibleApps.Data == nil {
		return []string{}
	}
	ids := make([]string, len(*r.VisibleApps.Data))
	for i, app := range *r.VisibleApps.Data {
		ids[i] = app.Id
	}
	return ids
}

// invitationVisibleAppIDs extracts app IDs from UserInvitationRelationships.
func invitationVisibleAppIDs(r *openapi.UserInvitationRelationships) []string {
	if r == nil || r.VisibleApps == nil || r.VisibleApps.Data == nil {
		return []string{}
	}
	ids := make([]string, len(*r.VisibleApps.Data))
	for i, app := range *r.VisibleApps.Data {
		ids[i] = app.Id
	}
	return ids
}

func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
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
