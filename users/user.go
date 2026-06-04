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



