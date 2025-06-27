package users

type UserRole string

const (
	Admin                       UserRole = "ADMIN"
	Finance                     UserRole = "FINANCE"
	AccountHolder               UserRole = "ACCOUNT_HOLDER"
	Sales                       UserRole = "SALES"
	Marketing                   UserRole = "MARKETING"
	AppManager                  UserRole = "APP_MANAGER"
	Developer                   UserRole = "DEVELOPER"
	AccessToReports             UserRole = "ACCESS_TO_REPORTS"
	CustomerSupport             UserRole = "CUSTOMER_SUPPORT"
	CreateApps                  UserRole = "CREATE_APPS"
	CloudManagedDeveloperID     UserRole = "CLOUD_MANAGED_DEVELOPER_ID"
	CloudManagedAppDistribution UserRole = "CLOUD_MANAGED_APP_DISTRIBUTION"
	GenerateIndividualKeys      UserRole = "GENERATE_INDIVIDUAL_KEYS"
)
