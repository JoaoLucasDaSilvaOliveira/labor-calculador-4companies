package components

import (
	"testing"
	"time"

	query "labor-calculador-4companies/internal/application/query/company"
	"labor-calculador-4companies/internal/domain/entity"

	"fyne.io/fyne/v2/test"
)

type reloadCompanyFinder struct {
	calls chan query.GetCompanyWithFilter
}

func (finder reloadCompanyFinder) Execute(filter query.GetCompanyWithFilter) ([]*entity.Company, error) {
	finder.calls <- filter
	company, err := entity.LoadCompany(10, "Empresa", "11222333000181")
	return []*entity.Company{company}, err
}

func waitForCompanyListQuery(t *testing.T, calls <-chan query.GetCompanyWithFilter) query.GetCompanyWithFilter {
	t.Helper()
	select {
	case filter := <-calls:
		return filter
	case <-time.After(time.Second):
		t.Fatal("company list query was not executed")
		return query.GetCompanyWithFilter{}
	}
}

func TestReloadableCompaniesListReloadsInBackground(t *testing.T) {
	app := test.NewApp()
	t.Cleanup(app.Quit)

	calls := make(chan query.GetCompanyWithFilter, 2)
	list := NewReloadableCompaniesListComponent(
		reloadCompanyFinder{calls: calls},
		query.GetCompanyWithFilter{},
		companiesEmptyMessage,
		nil,
	)
	if got := waitForCompanyListQuery(t, calls); got != (query.GetCompanyWithFilter{}) {
		t.Fatalf("initial filter = %#v, want empty filter", got)
	}

	list.Reload()
	if got := waitForCompanyListQuery(t, calls); got != (query.GetCompanyWithFilter{}) {
		t.Fatalf("reload filter = %#v, want empty filter", got)
	}
}
