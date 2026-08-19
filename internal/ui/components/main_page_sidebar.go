package components

import (
	"strings"
	"time"

	query "labor-calculador-4companies/internal/application/query/company"
	"labor-calculador-4companies/internal/ui/assets"
	uiUtils "labor-calculador-4companies/internal/ui/utils"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

// MainPageSidebarDeps contains the data access and navigation callback shared
// by the opened and search sidebar states.
type MainPageSidebarDeps struct {
	Companies         CompanyFinder
	CompanyListReload *CompanyListReloadHandle
	OnCompanySelected func(companyID int)
	OnAddCompany      func()
}

func NewOpenedMainPageSidebar(deps MainPageSidebarDeps, onCloseSidebar func(), onSearch func()) *fyne.Container {
	toolbar := newSidebarToolbar(
		onCloseSidebar,
		assets.MagnifingGlass,
		onSearch,
		deps.OnAddCompany,
	)
	reloadableCompaniesList := NewReloadableCompaniesListComponent(
		deps.Companies,
		query.GetCompanyWithFilter{},
		companiesEmptyMessage,
		deps.OnCompanySelected,
	)
	if deps.CompanyListReload != nil {
		deps.CompanyListReload.Set(reloadableCompaniesList.Reload)
	}
	companiesList := reloadableCompaniesList.View()

	return container.NewBorder(
		toolbar,
		nil, nil, nil,
		container.NewBorder(nil, nil, uiUtils.HPadding(8), uiUtils.HPadding(8), companiesList),
	)
}

func NewClosedMainPageSidebar(onOpenSidebar func()) *fyne.Container {
	openIcon := newToolbarIcon(assets.CloseTabIcon, onOpenSidebar)
	iconBox := container.NewHBox(
		uiUtils.HPadding(8),
		openIcon,
		uiUtils.HPadding(8),
	)
	topSection := container.NewVBox(
		iconBox,
		widget.NewSeparator(),
	)

	return container.NewBorder(topSection, nil, nil, nil, nil)
}

func NewSearchMainPageSidebar(deps MainPageSidebarDeps, onCloseSidebar func(), onCloseSearch func()) *fyne.Container {
	toolbar := newSidebarToolbar(
		onCloseSidebar,
		assets.CloseMagnifingGlass,
		onCloseSearch,
		deps.OnAddCompany,
	)
	searchInput := widget.NewEntry()
	searchInput.SetPlaceHolder("Buscar empresa...")

	searchResults := NewAnimatedContent()
	searchResults.SetContent(widget.NewLabel("Digite um nome para buscar."), nil)

	searchInput.OnChanged = func(companyName string) {
		trimmedName := strings.TrimSpace(companyName)
		if trimmedName == "" {
			searchResults.SetContent(widget.NewLabel("Digite um nome para buscar."), nil)
			return
		}

		result := NewReloadableCompaniesListComponent(
			deps.Companies,
			query.GetCompanyWithFilter{Name: trimmedName},
			companiesSearchEmptyMessage,
			deps.OnCompanySelected,
		)
		if deps.CompanyListReload != nil {
			deps.CompanyListReload.Set(result.Reload)
		}
		searchResults.SetContent(result.View(), nil)
	}

	searchContent := container.NewBorder(searchInput, nil, nil, nil, searchResults.View())
	return container.NewBorder(
		toolbar,
		nil, nil, nil,
		searchContent,
	)
}

func newToolbarIcon(resource fyne.Resource, tapped func()) *widget.Button {
	return widget.NewButtonWithIcon("", resource, func() {
		if tapped == nil {
			return
		}

		// Let the button render its tap animation before replacing the sidebar view.
		time.AfterFunc(canvas.DurationShort, func() {
			fyne.Do(tapped)
		})
	})
}

func newSidebarToolbar(onCloseSidebar func(), actionIcon fyne.Resource, onAction func(), onAddCompany func()) *fyne.Container {
	closeIcon := newToolbarIcon(assets.OpenTabIcon, onCloseSidebar)
	actionIconButton := newToolbarIcon(actionIcon, onAction)
	iconBox := container.NewHBox(
		uiUtils.HPadding(8),
		closeIcon,
		layout.NewSpacer(),
		actionIconButton,
		uiUtils.HPadding(8),
	)

	companiesLabel := widget.NewLabel("Empresas")
	addCompanyIcon := newToolbarIcon(assets.AddCompany, onAddCompany)
	return container.NewVBox(
		iconBox,
		widget.NewSeparator(),
		container.NewHBox(
			uiUtils.HPadding(8),
			companiesLabel,
			layout.NewSpacer(),
			addCompanyIcon,
			uiUtils.HPadding(8),
		),
	)
}
