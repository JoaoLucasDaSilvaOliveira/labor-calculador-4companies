package components

import (
	"strings"
	"sync"

	query "labor-calculador-4companies/internal/application/query/company"
	"labor-calculador-4companies/internal/domain/entity"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

const (
	companiesLoadErrorMessage   = "Não foi possível carregar as empresas."
	companiesEmptyMessage       = "Nenhuma empresa cadastrada."
	companiesSearchEmptyMessage = "Nenhuma empresa encontrada para essa busca."
)

// CompanyFinder is the query contract required by the company list.
// GetCompanyUsecase implements this interface without an adapter.
type CompanyFinder interface {
	Execute(query.GetCompanyWithFilter) ([]*entity.Company, error)
}

// CompanyListReloadHandle points to the currently mounted company list. The
// mainpage keeps the handle stable while the sidebar replaces its visual
// state, allowing a saved company to request a reload without knowing the
// sidebar layout.
type CompanyListReloadHandle struct {
	mu     sync.RWMutex
	reload func()
}

func (handle *CompanyListReloadHandle) Set(reload func()) {
	if handle == nil {
		return
	}

	handle.mu.Lock()
	handle.reload = reload
	handle.mu.Unlock()
}

func (handle *CompanyListReloadHandle) Reload() {
	if handle == nil {
		return
	}

	handle.mu.RLock()
	reload := handle.reload
	handle.mu.RUnlock()
	if reload != nil {
		reload()
	}
}

// ReloadableCompaniesList owns one query-backed list and can refresh it
// without rebuilding the sidebar container.
type ReloadableCompaniesList struct {
	content  *AnimatedContent
	finder   CompanyFinder
	filter   query.GetCompanyWithFilter
	empty    string
	onSelect func(companyID int)

	mu        sync.Mutex
	requestID uint64
}

func NewReloadableCompaniesListComponent(
	finder CompanyFinder,
	filter query.GetCompanyWithFilter,
	emptyMessage string,
	onCompanySelected func(companyID int),
) *ReloadableCompaniesList {
	list := &ReloadableCompaniesList{
		content:  NewAnimatedContent(),
		finder:   finder,
		filter:   filter,
		empty:    emptyMessage,
		onSelect: onCompanySelected,
	}
	list.content.SetContent(NewLoadingState("Carregando empresas..."), nil)
	list.Reload()
	return list
}

func (list *ReloadableCompaniesList) View() fyne.CanvasObject {
	return list.content.View()
}

// Reload queries in a goroutine and applies only the newest response on the
// Fyne UI thread. This keeps the sidebar responsive during persistence.
func (list *ReloadableCompaniesList) Reload() {
	if list == nil {
		return
	}

	list.mu.Lock()
	list.requestID++
	requestID := list.requestID
	list.mu.Unlock()

	go func() {
		companies, err := list.finder.Execute(list.filter)
		apply := func() {
			list.mu.Lock()
			isCurrent := requestID == list.requestID
			list.mu.Unlock()
			if !isCurrent {
				return
			}
			list.setCompanies(companies, err)
		}

		if fyne.CurrentApp() == nil {
			apply()
			return
		}
		fyne.Do(apply)
	}()
}

func (list *ReloadableCompaniesList) setCompanies(companies []*entity.Company, err error) {
	if err != nil {
		list.content.SetContent(newStatusLabel(companiesLoadErrorMessage), nil)
		return
	}
	if len(companies) == 0 {
		list.content.SetContent(newStatusLabel(list.empty), nil)
		return
	}

	companiesList := widget.NewList(
		func() int { return len(companies) },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(id widget.ListItemID, object fyne.CanvasObject) {
			object.(*widget.Label).SetText(companies[id].Name())
		},
	)
	companiesList.OnSelected = func(id widget.ListItemID) {
		companiesList.Unselect(id)
		if list.onSelect != nil {
			list.onSelect(companies[id].GetId())
		}
	}
	list.content.SetContent(companiesList, nil)
}

// NewCompaniesListComponent creates the complete company list.
func NewCompaniesListComponent(finder CompanyFinder, onCompanySelected func(companyID int)) fyne.CanvasObject {
	return NewReloadableCompaniesListComponent(
		finder,
		query.GetCompanyWithFilter{},
		companiesEmptyMessage,
		onCompanySelected,
	).View()
}

// NewCompaniesSearchByNameListComponent creates a company list filtered by name.
func NewCompaniesSearchByNameListComponent(finder CompanyFinder, companyName string, onCompanySelected func(companyID int)) fyne.CanvasObject {
	trimmedName := strings.TrimSpace(companyName)
	if trimmedName == "" {
		return newStatusLabel("Digite um nome para buscar.")
	}

	return NewReloadableCompaniesListComponent(
		finder,
		query.GetCompanyWithFilter{Name: trimmedName},
		companiesSearchEmptyMessage,
		onCompanySelected,
	).View()
}

func newStatusLabel(message string) *widget.Label {
	label := widget.NewLabel(message)
	label.Wrapping = fyne.TextWrapWord
	return label
}
