package users

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

func (u *User) relationships() *userRelationships {
	if len(u.VisibleAppIDs) == 0 && u.AllAppsVisible {
		return nil
	}

	apps := make([]relationshipObject, len(u.VisibleAppIDs))
	for index, id := range u.VisibleAppIDs {
		apps[index] = appReference(id)
	}

	return &userRelationships{
		VisibleApps: visibleApps{
			AppReferences: apps,
		},
	}
}

type userRelationships struct {
	VisibleApps visibleApps `json:"visibleApps"`
}

func (r *userRelationships) ids() []string {
	if r == nil || len(r.VisibleApps.AppReferences) == 0 {
		return []string{}
	}
	ids := make([]string, len(r.VisibleApps.AppReferences))
	for index, app := range r.VisibleApps.AppReferences {
		ids[index] = app.ID
	}
	return ids

}

type visibleApps struct {
	AppReferences []relationshipObject `json:"data"`
}

type relationshipObject struct {
	ID   string `json:"id,omitempty"`
	Type string `json:"type,omitempty"` // must be "apps"
}

func appReference(id string) relationshipObject {
	data := relationshipObject{}
	data.ID = id
	data.Type = "apps"
	return data
}
