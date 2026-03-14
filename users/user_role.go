package users

import "github.com/oliver-binns/appstore-go/openapi"

// UserRole is an alias for the generated openapi.UserRole type.
type UserRole = openapi.UserRole

// Role constants re-exported with idiomatic Go names.
const (
	Admin                       = openapi.ADMIN
	Finance                     = openapi.FINANCE
	AccountHolder               = openapi.ACCOUNTHOLDER
	Sales                       = openapi.SALES
	Marketing                   = openapi.MARKETING
	AppManager                  = openapi.APPMANAGER
	Developer                   = openapi.DEVELOPER
	AccessToReports             = openapi.ACCESSTOREPORTS
	CustomerSupport             = openapi.CUSTOMERSUPPORT
	CreateApps                  = openapi.CREATEAPPS
	CloudManagedDeveloperID     = openapi.CLOUDMANAGEDDEVELOPERID
	CloudManagedAppDistribution = openapi.CLOUDMANAGEDAPPDISTRIBUTION
	GenerateIndividualKeys      = openapi.GENERATEINDIVIDUALKEYS
)
