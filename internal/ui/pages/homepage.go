package pages

import (
	"fmt"
	"image/color"
	"labor-calculador-4companies/internal/ui/application"
	"labor-calculador-4companies/internal/ui/components"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/JoaoLucasDaSilvaOliveira/learning-fyne/classes/utils"
)

const (
	homePageOpenedSplitOffset = 0.3
)

type homePageSidebarState uint8

const (
	homePageSidebarCompanies homePageSidebarState = iota
	homePageSidebarSearch
	homePageSidebarClosed
)

func NewHomePage(app *application.MainApplication) {
	// sidebarHost remains mounted in the split while its single child changes between sidebar states.
	sidebarHost := container.NewStack()
	// sidebarState is the requested view, while lastOpenedSidebarState is restored after reopening.
	sidebarState := homePageSidebarCompanies
	lastOpenedSidebarState := homePageSidebarCompanies

	// Quick access component - center component.
	fastAccLabel := widget.NewLabel("ACESSO RÁPIDO")
	fastAccLabel.TextStyle.Bold = true
	addCompanyBtn := components.NewClickableText("+ CADASTRAR EMPRESA", color.Black, utils.NewColor(0, 105, 204, 255), func() { fmt.Println("tapped1") })
	newCalculusBtn := components.NewClickableText("+ NOVO CÁLCULO", color.Black, utils.NewColor(0, 105, 204, 255), func() { fmt.Println("tapped2") })
	centerBox := container.NewCenter(container.NewVBox(
		fastAccLabel,
		addCompanyBtn,
		newCalculusBtn,
	))

	homepageBox := container.NewHSplit(sidebarHost, centerBox)
	homepageBox.SetOffset(homePageOpenedSplitOffset)

	// Keep the active animation so a new click can interrupt an in-progress transition.
	var splitAnimation *fyne.Animation
	var animateSplitOffset func(targetOffset float64, onFinished func())
	animateSplitOffset = func(targetOffset float64, onFinished func()) {
		if splitAnimation != nil {
			splitAnimation.Stop()
		}

		initialOffset := homepageBox.Offset
		if initialOffset == targetOffset {
			if onFinished != nil {
				onFinished()
			}
			return
		}

		var animation *fyne.Animation
		animation = fyne.NewAnimation(canvas.DurationStandard, func(progress float32) {
			offset := initialOffset + (targetOffset-initialOffset)*float64(progress)
			homepageBox.SetOffset(offset)
			// Only the current animation may start the next transition phase.
			if progress == 1 && splitAnimation == animation && onFinished != nil {
				onFinished()
			}
		})
		animation.Curve = fyne.AnimationEaseInOut
		splitAnimation = animation
		animation.Start()
	}
	// Convert a view minimum width into the equivalent split offset.
	minimumOffset := func(view fyne.CanvasObject) float64 {
		dividerWidth := theme.CurrentForWidget(homepageBox).Size(theme.SizeNameSplitThickness)
		availableWidth := homepageBox.Size().Width - dividerWidth
		if availableWidth <= 0 {
			return 0
		}

		return float64(view.MinSize().Width / availableWidth)
	}
	// Replaces the variable sidebar content without recreating the HSplit.
	setSidebar := func(view fyne.CanvasObject) {
		sidebarHost.Objects = []fyne.CanvasObject{view}
		sidebarHost.Refresh()
	}

	var renderSidebar func()
	var renderedSidebarState homePageSidebarState
	renderSidebar = func() {
		closeSidebar := func() {
			if sidebarState == homePageSidebarClosed {
				return
			}

			lastOpenedSidebarState = sidebarState
			sidebarState = homePageSidebarClosed
			renderSidebar()
		}
		openSidebar := func() {
			if sidebarState != homePageSidebarClosed {
				return
			}

			sidebarState = lastOpenedSidebarState
			renderSidebar()
		}
		showSearch := func() {
			if sidebarState == homePageSidebarClosed || sidebarState == homePageSidebarSearch {
				return
			}

			sidebarState = homePageSidebarSearch
			lastOpenedSidebarState = homePageSidebarSearch
			renderSidebar()
		}
		showCompanies := func() {
			if sidebarState != homePageSidebarSearch {
				return
			}

			sidebarState = homePageSidebarCompanies
			lastOpenedSidebarState = homePageSidebarCompanies
			renderSidebar()
		}

		var nextSidebar fyne.CanvasObject
		switch sidebarState {
		case homePageSidebarCompanies:
			nextSidebar = components.NewOpenedHomePageSideBar(app, closeSidebar, showSearch)
		case homePageSidebarSearch:
			nextSidebar = components.NewSearchHomePage(app, closeSidebar, showCompanies)
		case homePageSidebarClosed:
			nextSidebar = components.NewClosedHomePageSideBar(app, openSidebar)
		}

		if len(sidebarHost.Objects) == 0 {
			// Initial render does not need an animation.
			setSidebar(nextSidebar)
			renderedSidebarState = sidebarState
			homepageBox.SetOffset(homePageOpenedSplitOffset)
			return
		}

		if renderedSidebarState == homePageSidebarClosed && sidebarState != homePageSidebarClosed {
			// Expand to the incoming view minimum width before mounting it to avoid a layout jump.
			animateSplitOffset(minimumOffset(nextSidebar), func() {
				setSidebar(nextSidebar)
				renderedSidebarState = sidebarState
				animateSplitOffset(homePageOpenedSplitOffset, nil)
			})
			return
		}

		if renderedSidebarState != homePageSidebarClosed && sidebarState == homePageSidebarClosed {
			// Shrink the current view before replacing it with the closed sidebar.
			animateSplitOffset(minimumOffset(sidebarHost), func() {
				setSidebar(nextSidebar)
				renderedSidebarState = sidebarState
				animateSplitOffset(minimumOffset(nextSidebar), nil)
			})
			return
		}

		setSidebar(nextSidebar)
		renderedSidebarState = sidebarState
	}

	renderSidebar()

	// Set the homepage content.
	app.MasterWindow.SetContent(homepageBox)
}
