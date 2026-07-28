package navigation

// RouteID identifies a page registered in the application router.
type RouteID string

const (
	RouteMain              RouteID = "main"
	RouteCompanyDetails    RouteID = "company.details"
	RouteCompanyCreate     RouteID = "company.create"
	RouteCalculationCreate RouteID = "calculation.create"
)

// CompanyDetailsParams contains the data required to open a company page.
// Only the identifier travels between pages so company data can be loaded from
// the application's source of truth when the real details screen is added.
type CompanyDetailsParams struct {
	CompanyID int
}
