package components

import (
	"time"

	"labor-calculador-4companies/internal/ui/application"
	"labor-calculador-4companies/internal/ui/assets"
	uiUtils "labor-calculador-4companies/internal/ui/utils"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func NewOpenedHomePageSideBar(app *application.MainApplication, onCloseSideBarCallback, onSearchCallback func()) *fyne.Container {
	var sideBarBox *fyne.Container
	topSection := topSection(onCloseSideBarCallback, assets.MagnifingGlass, onSearchCallback)
	//companies list component
	companiesList := NewCompaniesListComponent(app.GetCompaniesUC)
	sideBarBox = container.NewBorder(
		topSection,
		nil, nil, nil,
		companiesList,
	)

	return sideBarBox
}

func NewClosedHomePageSideBar(app *application.MainApplication, onOpenSideBarCallback func()) *fyne.Container {
	var sideBarBox *fyne.Container
	//top section components
	// icons
	openIcon := newToolbarIcon(assets.CloseTabIcon, onOpenSideBarCallback)
	// box
	iconBox := container.NewHBox(
		uiUtils.HPadding(8), // padding esquerdo
		openIcon,
		uiUtils.HPadding(8), // padding direito
	)
	topSection := container.NewVBox(
		iconBox,
		widget.NewSeparator(),
	)
	sideBarBox = container.NewBorder(
		topSection,
		nil, nil, nil, nil,
	)

	return sideBarBox
}

func NewSearchHomePage(app *application.MainApplication, onCloseSideBarCallback, onCloseSearchCallback func()) *fyne.Container {
	topSection := topSection(onCloseSideBarCallback, assets.CloseMagnifingGlass, onCloseSearchCallback)
	// Entry for searching companies by name.
	searchInput := widget.NewEntry()
	searchInput.SetPlaceHolder("Buscar empresa...")
	searchResults := container.NewStack(widget.NewLabel("Digite um nome para buscar."))

	searchInput.OnChanged = func(companyName string) {
		companiesList := NewCompaniesSearchByNameListComponent(app.GetCompaniesUC, companyName)
		if companiesList == nil {
			searchResults.Objects = []fyne.CanvasObject{widget.NewLabel("Digite um nome para buscar.")}
			searchResults.Refresh()
			return
		}

		searchResults.Objects = []fyne.CanvasObject{companiesList}
		searchResults.Refresh()
	}

	searchContent := container.NewBorder(searchInput, nil, nil, nil, searchResults)
	return container.NewBorder(
		topSection,
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

func topSection(onCloseSideBarCallback func(), actionIcon fyne.Resource, onActionCallback func()) *fyne.Container {
	//top section components
	// icons
	closeIcon := newToolbarIcon(assets.OpenTabIcon, onCloseSideBarCallback)
	actionIconButton := newToolbarIcon(actionIcon, onActionCallback)
	// box
	iconBox := container.NewHBox(
		uiUtils.HPadding(8), // padding esquerdo
		closeIcon,
		layout.NewSpacer(),
		actionIconButton,
		uiUtils.HPadding(8), // padding direito
	)
	// label
	companiesLabel := widget.NewLabel("Empresas")
	// mother box
	return container.NewVBox(
		iconBox,
		widget.NewSeparator(),
		companiesLabel,
	)
}
