package components

import (
	"labor-calculador-4companies/internal/ui/application"
	"labor-calculador-4companies/internal/ui/assets"
	uiUtils "labor-calculador-4companies/internal/ui/utils"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func NewOpenedHomePageSideBar(app *application.MainApplication, onCloseSideBarCallback func()) *fyne.Container {
	var sideBarBox *fyne.Container
	//top section components
	// icons
	closeIcon := newToolbarIcon(assets.OpenTabIcon, onCloseSideBarCallback)
	mGlassIcon := newToolbarIcon(assets.MagnifingGlass, func() {})          //we'll gonna  set the function later
	// box
	iconBox := container.NewHBox(
		uiUtils.HPadding(8), // padding esquerdo
		closeIcon,
		layout.NewSpacer(),
		mGlassIcon,
		uiUtils.HPadding(8), // padding direito
	)
	// label
	companiesLabel := widget.NewLabel("Empresas")
	// mother box
	topSection := container.NewVBox(
		iconBox,
		widget.NewSeparator(),
		companiesLabel,
	)

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

func newToolbarIcon(resource fyne.Resource, tapped func()) *widget.Button {
	return widget.NewButtonWithIcon("", resource, tapped)
}
