package users

type userInvitation struct {
	ID                  string     `json:"id,omitempty"`
	FirstName           string     `json:"firstName"`
	LastName            string     `json:"lastName"`
	Email               string     `json:"email"`
	Roles               []UserRole `json:"roles"`
	AllAppsVisible      bool       `json:"allAppsVisible,omitempty"`
	ProvisioningAllowed bool       `json:"provisioningAllowed,omitempty"`
}
