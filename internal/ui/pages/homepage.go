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

func NewHomePage(app *application.MainApplication) {
	// sidebarHost remains mounted in the split while its single child changes between sidebar states.
	sidebarHost := container.NewStack()
	// showSidebar is the source of truth for the sidebar mode.
	showSidebar := true

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
	renderSidebar = func() {
		if showSidebar {
			openedSidebar := components.NewOpenedHomePageSideBar(app, func() {
				showSidebar = false
				renderSidebar()
			})

			if len(sidebarHost.Objects) == 0 {
				// Initial render does not need an animation.
				setSidebar(openedSidebar)
				homepageBox.SetOffset(homePageOpenedSplitOffset)
				return
			}

			// Expand to the opened view minimum width before mounting it to avoid a layout jump.
			animateSplitOffset(minimumOffset(openedSidebar), func() {
				setSidebar(openedSidebar)
				// Continue from the minimum width to the normal opened sidebar width.
				animateSplitOffset(homePageOpenedSplitOffset, nil)
			})
			return
		}

		closedSidebar := components.NewClosedHomePageSideBar(app, func() {
			showSidebar = true
			renderSidebar()
		})
		// Shrink the current view to its minimum before replacing it with the closed sidebar.
		animateSplitOffset(minimumOffset(sidebarHost), func() {
			setSidebar(closedSidebar)
			// Finish the transition using the minimum width required by the closed view.
			animateSplitOffset(minimumOffset(closedSidebar), nil)
		})
	}

	renderSidebar()

	// Set the homepage content.
	app.MasterWindow.SetContent(homepageBox)
}

func openTabIconFunction() {

}

func closeTabIconFunction() {

}

func mGlassIconFunction() {

}
