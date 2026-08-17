package pages

import (
	"testing"
	"time"

	companyCommand "labor-calculador-4companies/internal/application/command/company"
	employeeCommand "labor-calculador-4companies/internal/application/command/employee"
	employeeQuery "labor-calculador-4companies/internal/application/query/employee"
	receiptQuery "labor-calculador-4companies/internal/application/query/receipt"
	"labor-calculador-4companies/internal/domain/entity"
	"labor-calculador-4companies/internal/ui/navigation"

	"fyne.io/fyne/v2/test"
)

type pageNavigator struct{}

func (pageNavigator) Push(navigation.RouteID, any) error    { return nil }
func (pageNavigator) Replace(navigation.RouteID, any) error { return nil }
func (pageNavigator) Reset(navigation.RouteID, any) error   { return nil }
func (pageNavigator) Back() bool                            { return true }
func (pageNavigator) CanGoBack() bool                       { return true }
func (pageNavigator) Current() navigation.RouteID           { return navigation.RouteEmployeeDetails }

type recordingNavigator struct {
	operations []string
	routes     []navigation.RouteID
}

func (n *recordingNavigator) Push(route navigation.RouteID, _ any) error {
	n.operations = append(n.operations, "push")
	n.routes = append(n.routes, route)
	return nil
}

func (n *recordingNavigator) Replace(route navigation.RouteID, _ any) error {
	n.operations = append(n.operations, "replace")
	n.routes = append(n.routes, route)
	return nil
}

func (n *recordingNavigator) Reset(route navigation.RouteID, _ any) error {
	n.operations = append(n.operations, "reset")
	n.routes = append(n.routes, route)
	return nil
}

func (n *recordingNavigator) Back() bool                  { return true }
func (n *recordingNavigator) CanGoBack() bool             { return true }
func (n *recordingNavigator) Current() navigation.RouteID { return navigation.RouteCompanyDetails }

func TestResetToQuickAccessResetsAndReplacesWorkspaceRoute(t *testing.T) {
	navigator := &recordingNavigator{}

	resetToQuickAccess(navigator)

	if got, want := len(navigator.operations), 2; got != want {
		t.Fatalf("operation count = %d, want %d", got, want)
	}
	if navigator.operations[0] != "reset" || navigator.operations[1] != "replace" {
		t.Fatalf("operations = %#v, want reset then replace", navigator.operations)
	}
	for _, route := range navigator.routes {
		if route != navigation.RouteQuickAccess {
			t.Fatalf("route = %q, want %q", route, navigation.RouteQuickAccess)
		}
	}
}

type companyPageCompanyFinder struct {
	called  chan companyCommand.GetCompanyById
	company *entity.Company
}

func (f companyPageCompanyFinder) Execute(cmd companyCommand.GetCompanyById) (*entity.Company, error) {
	f.called <- cmd
	return f.company, nil
}

type companyPageEmployeeFinder struct {
	called chan employeeQuery.GetEmployeeWithFilter
}

func (f companyPageEmployeeFinder) Execute(query employeeQuery.GetEmployeeWithFilter) ([]*entity.Employee, error) {
	f.called <- query
	return nil, nil
}

type employeePageEmployeeFinder struct {
	called   chan employeeCommand.GetEmployeeById
	employee *entity.Employee
}

func (f employeePageEmployeeFinder) Execute(cmd employeeCommand.GetEmployeeById) (*entity.Employee, error) {
	f.called <- cmd
	return f.employee, nil
}

type employeePageReceiptFinder struct {
	called chan receiptQuery.GetReceiptWithFilter
}

func (f employeePageReceiptFinder) Execute(query receiptQuery.GetReceiptWithFilter) ([]*entity.Receipt, error) {
	f.called <- query
	return nil, nil
}

func TestCompanyPageQueriesEmployeesByCompanyID(t *testing.T) {
	app := test.NewApp()
	t.Cleanup(app.Quit)

	company, err := entity.LoadCompany(10, "Acme", "11222333000181")
	if err != nil {
		t.Fatalf("LoadCompany returned an error: %v", err)
	}

	companyCalls := make(chan companyCommand.GetCompanyById, 1)
	employeeCalls := make(chan employeeQuery.GetEmployeeWithFilter, 1)
	NewCompanyPage(CompanyPageDeps{
		CompanyID: 10,
		Company: companyPageCompanyFinder{
			called:  companyCalls,
			company: company,
		},
		Employees: companyPageEmployeeFinder{called: employeeCalls},
		Navigator: pageNavigator{},
	})

	select {
	case cmd := <-companyCalls:
		if cmd.IDCompany != 10 {
			t.Fatalf("company ID = %d, want 10", cmd.IDCompany)
		}
	case <-time.After(time.Second):
		t.Fatal("company query was not executed")
	}

	select {
	case query := <-employeeCalls:
		if query.CompanyID != 10 {
			t.Fatalf("employee company ID = %d, want 10", query.CompanyID)
		}
	case <-time.After(time.Second):
		t.Fatal("employee query was not executed")
	}
}

func TestEmployeePageQueriesReceiptsByEmployeeID(t *testing.T) {
	app := test.NewApp()
	t.Cleanup(app.Quit)

	employee, err := entity.LoadEmployee(25, 10, "Ana", "Silva", "52998224725")
	if err != nil {
		t.Fatalf("LoadEmployee returned an error: %v", err)
	}

	employeeCalls := make(chan employeeCommand.GetEmployeeById, 1)
	receiptCalls := make(chan receiptQuery.GetReceiptWithFilter, 1)
	NewEmployeePage(EmployeePageDeps{
		EmployeeID: 25,
		Employee: employeePageEmployeeFinder{
			called:   employeeCalls,
			employee: employee,
		},
		Receipts:  employeePageReceiptFinder{called: receiptCalls},
		Navigator: pageNavigator{},
	})

	select {
	case cmd := <-employeeCalls:
		if cmd.IDEmployee != 25 {
			t.Fatalf("employee ID = %d, want 25", cmd.IDEmployee)
		}
	case <-time.After(time.Second):
		t.Fatal("employee query was not executed")
	}

	select {
	case query := <-receiptCalls:
		if query.IDEmployee != 25 {
			t.Fatalf("receipt employee ID = %d, want 25", query.IDEmployee)
		}
	case <-time.After(time.Second):
		t.Fatal("receipt query was not executed")
	}
}
