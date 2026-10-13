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

func visibleAppsRelationship(ids []string) *openapi.VisibleAppsRelationship {
	refs := make([]openapi.AppReference, len(ids))
	for i, id := range ids {
		refs[i] = openapi.AppReference{Id: id, Type: openapi.Apps}
	}
	return &openapi.VisibleAppsRelationship{Data: &refs}
}

func invitationCreateRelationships(ids []string, allAppsVisible bool) *openapi.UserInvitationCreateRelationships {
	if len(ids) == 0 && allAppsVisible {
		return nil
	}
	return &openapi.UserInvitationCreateRelationships{VisibleApps: visibleAppsRelationship(ids)}
}

func userUpdateRelationships(ids []string, allAppsVisible bool) *openapi.UserUpdateRelationships {
	if len(ids) == 0 && allAppsVisible {
		return nil
	}
	return &openapi.UserUpdateRelationships{VisibleApps: visibleAppsRelationship(ids)}
}

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
