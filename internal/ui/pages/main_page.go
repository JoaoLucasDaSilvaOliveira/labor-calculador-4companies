package pages

import (
	"labor-calculador-4companies/internal/ui/components"
	"labor-calculador-4companies/internal/ui/navigation"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
)

const mainPageOpenedSplitOffset = 0.3

type mainPageSidebarState uint8

const (
	mainPageSidebarCompanies mainPageSidebarState = iota
	mainPageSidebarSearch
	mainPageSidebarClosed
)

// MainPageDeps lists the external capabilities required by the mainpage.
type MainPageDeps struct {
	Companies     components.CompanyFinder
	Navigator     navigation.Navigator
	WorkspaceView fyne.CanvasObject
}

// mainPage owns only the presentation state local to the mainpage.
// Application-level navigation remains in navigation.Router.
type mainPage struct {
	deps MainPageDeps

	sidebarContent *components.AnimatedContent
	sidebarHost    fyne.CanvasObject
	split          *container.Split
	splitAnimation *fyne.Animation

	sidebarState           mainPageSidebarState
	lastOpenedSidebarState mainPageSidebarState
	renderedSidebarState   mainPageSidebarState
	hasRenderedSidebar     bool
}

// NewMainPage builds and returns the mainpage without changing the main window.
func NewMainPage(deps MainPageDeps) fyne.CanvasObject {
	page := &mainPage{
		deps:                   deps,
		sidebarState:           mainPageSidebarCompanies,
		lastOpenedSidebarState: mainPageSidebarCompanies,
	}

	return page.build()
}

func (p *mainPage) build() fyne.CanvasObject {
	p.sidebarContent = components.NewAnimatedContent()
	p.sidebarHost = p.sidebarContent.View()

	p.split = container.NewHSplit(p.sidebarHost, p.deps.WorkspaceView)
	p.split.SetOffset(mainPageOpenedSplitOffset)

	p.renderSidebar()
	return p.split
}

func (p *mainPage) sidebarDeps() components.MainPageSidebarDeps {
	return components.MainPageSidebarDeps{
		Companies:         p.deps.Companies,
		OnCompanySelected: p.openCompanyDetails,
		OnAddCompany:      p.openCompanyRegistration,
	}
}

func (p *mainPage) closeSidebar() {
	if p.sidebarState == mainPageSidebarClosed {
		return
	}

	p.lastOpenedSidebarState = p.sidebarState
	p.sidebarState = mainPageSidebarClosed
	p.renderSidebar()
}

func (p *mainPage) openSidebar() {
	if p.sidebarState != mainPageSidebarClosed {
		return
	}

	p.sidebarState = p.lastOpenedSidebarState
	p.renderSidebar()
}

func (p *mainPage) showSearch() {
	if p.sidebarState == mainPageSidebarClosed ||
		p.sidebarState == mainPageSidebarSearch {
		return
	}

	p.sidebarState = mainPageSidebarSearch
	p.lastOpenedSidebarState = mainPageSidebarSearch
	p.renderSidebar()
}

func (p *mainPage) showCompanies() {
	if p.sidebarState != mainPageSidebarSearch {
		return
	}

	p.sidebarState = mainPageSidebarCompanies
	p.lastOpenedSidebarState = mainPageSidebarCompanies
	p.renderSidebar()
}

func (p *mainPage) renderSidebar() {
	targetState := p.sidebarState
	nextSidebar := p.buildSidebar(targetState)

	if !p.hasRenderedSidebar {
		p.setInitialSidebar(nextSidebar, targetState)
		return
	}

	if p.renderedSidebarState == mainPageSidebarClosed &&
		targetState != mainPageSidebarClosed {
		p.expandSidebar(nextSidebar, targetState)
		return
	}

	if p.renderedSidebarState != mainPageSidebarClosed &&
		targetState == mainPageSidebarClosed {
		p.collapseSidebar(nextSidebar, targetState)
		return
	}

	p.swapSidebar(nextSidebar, targetState)
}

func (p *mainPage) buildSidebar(state mainPageSidebarState) fyne.CanvasObject {
	switch state {
	case mainPageSidebarCompanies:
		return components.NewOpenedMainPageSidebar(
			p.sidebarDeps(),
			p.closeSidebar,
			p.showSearch,
		)
	case mainPageSidebarSearch:
		return components.NewSearchMainPageSidebar(
			p.sidebarDeps(),
			p.closeSidebar,
			p.showCompanies,
		)
	case mainPageSidebarClosed:
		return components.NewClosedMainPageSidebar(p.openSidebar)
	default:
		return components.NewOpenedMainPageSidebar(
			p.sidebarDeps(),
			p.closeSidebar,
			p.showSearch,
		)
	}
}

func (p *mainPage) setInitialSidebar(sidebar fyne.CanvasObject, state mainPageSidebarState) {
	p.sidebarContent.SetContent(sidebar, nil)
	p.renderedSidebarState = state
	p.hasRenderedSidebar = true
	p.split.SetOffset(mainPageOpenedSplitOffset)
}

func (p *mainPage) expandSidebar(nextSidebar fyne.CanvasObject, targetState mainPageSidebarState) {
	p.animateSplitOffset(p.minimumOffset(nextSidebar), func() {
		p.sidebarContent.SetContent(nextSidebar, func() {
			p.renderedSidebarState = targetState
			p.animateSplitOffset(mainPageOpenedSplitOffset, nil)
		})
	})
}

func (p *mainPage) collapseSidebar(nextSidebar fyne.CanvasObject, targetState mainPageSidebarState) {
	p.animateSplitOffset(p.minimumOffset(p.sidebarHost), func() {
		p.sidebarContent.SetContent(nextSidebar, func() {
			p.renderedSidebarState = targetState
			p.animateSplitOffset(p.minimumOffset(nextSidebar), nil)
		})
	})
}

func (p *mainPage) swapSidebar(nextSidebar fyne.CanvasObject, targetState mainPageSidebarState) {
	p.sidebarContent.SetContent(nextSidebar, func() {
		p.renderedSidebarState = targetState
	})
}

func (p *mainPage) animateSplitOffset(targetOffset float64, onFinished func()) {
	if p.splitAnimation != nil {
		p.splitAnimation.Stop()
	}

	initialOffset := p.split.Offset
	if initialOffset == targetOffset {
		if onFinished != nil {
			onFinished()
		}
		return
	}

	var animation *fyne.Animation
	animation = fyne.NewAnimation(canvas.DurationStandard, func(progress float32) {
		offset := initialOffset + (targetOffset-initialOffset)*float64(progress)
		p.split.SetOffset(offset)

		if progress == 1 &&
			p.splitAnimation == animation &&
			onFinished != nil {
			onFinished()
		}
	})
	animation.Curve = fyne.AnimationEaseInOut
	p.splitAnimation = animation
	animation.Start()
}

func (p *mainPage) minimumOffset(view fyne.CanvasObject) float64 {
	dividerWidth := theme.CurrentForWidget(p.split).Size(theme.SizeNameSplitThickness)
	availableWidth := p.split.Size().Width - dividerWidth
	if availableWidth <= 0 {
		return 0
	}

	return float64(view.MinSize().Width / availableWidth)
}

func (p *mainPage) openCompanyDetails(companyID int) {
	if err := p.deps.Navigator.Reset(
		navigation.RouteCompanyDetails,
		navigation.CompanyDetailsParams{CompanyID: companyID},
	); err != nil {
		fyne.LogError("Não foi possível abrir a empresa selecionada.", err)
	}
}

func (p *mainPage) openCompanyRegistration() {
	p.navigate(navigation.RouteCompanyCreate, nil)
}

func (p *mainPage) openCalculationCreation() {
	p.navigate(navigation.RouteCalculationCreate, nil)
}

func (p *mainPage) navigate(route navigation.RouteID, params any) {
	if err := p.deps.Navigator.Push(route, params); err != nil {
		fyne.LogError("Não foi possível navegar para a página solicitada.", err)
	}
} // VAI SER USADO SÓ NO FUTURO, WORKSPACE VAI PRECISAR DE UM MÉTODO COMO ESSE
