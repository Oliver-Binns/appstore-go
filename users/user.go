package users

type User struct {
	FirstName           string     `json:"firstName"`
	LastName            string     `json:"lastName"`
	Username            string     `json:"username"`
	Roles               []UserRole `json:"roles"`
	AllAppsVisible      bool       `json:"allAppsVisible,omitempty"`
	ProvisioningAllowed bool       `json:"provisioningAllowed,omitempty"`
}
