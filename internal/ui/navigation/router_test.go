package navigation

import (
	"errors"
	"fmt"
	"slices"
	"testing"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
)

const routeTestDetails RouteID = "test.details"

func TestRouterPushAndBackRestoreHistory(t *testing.T) {
	fyneApp := test.NewApp()
	t.Cleanup(fyneApp.Quit)

	router := newTestRouter(t)
	if err := router.Replace(RouteMain, nil); err != nil {
		t.Fatalf("Replace(RouteMain) returned an error: %v", err)
	}

	params := CompanyDetailsParams{CompanyID: 42}
	if err := router.Push(routeTestDetails, params); err != nil {
		t.Fatalf("Push(routeTestDetails) returned an error: %v", err)
	}

	if current := router.Current(); current != routeTestDetails {
		t.Fatalf("Current() = %q, want %q", current, routeTestDetails)
	}
	if !router.CanGoBack() {
		t.Fatal("CanGoBack() = false, want true")
	}
	if stored := router.history[len(router.history)-1].params; stored != params {
		t.Fatalf("stored params = %#v, want %#v", stored, params)
	}

	if wentBack := router.Back(); !wentBack {
		t.Fatal("Back() = false, want true")
	}
	if current := router.Current(); current != RouteMain {
		t.Fatalf("Current() after Back = %q, want %q", current, RouteMain)
	}
	if router.CanGoBack() {
		t.Fatal("CanGoBack() after Back = true, want false")
	}
}

func TestRouterReplaceDoesNotAddHistory(t *testing.T) {
	fyneApp := test.NewApp()
	t.Cleanup(fyneApp.Quit)

	router := newTestRouter(t)
	if err := router.Replace(RouteMain, nil); err != nil {
		t.Fatalf("Replace(RouteMain) returned an error: %v", err)
	}
	if err := router.Replace(
		routeTestDetails,
		CompanyDetailsParams{CompanyID: 7},
	); err != nil {
		t.Fatalf("Replace(routeTestDetails) returned an error: %v", err)
	}

	if len(router.history) != 1 {
		t.Fatalf("history length = %d, want 1", len(router.history))
	}
	if router.Back() {
		t.Fatal("Back() = true after Replace, want false")
	}
}

func TestRouterUsesFallbackForUnknownRoute(t *testing.T) {
	fyneApp := test.NewApp()
	t.Cleanup(fyneApp.Quit)

	router := newTestRouter(t)
	if err := router.Replace(RouteMain, nil); err != nil {
		t.Fatalf("Replace(RouteMain) returned an error: %v", err)
	}

	err := router.Push(RouteID("unknown"), nil)
	if !errors.Is(err, ErrRouteNotFound) {
		t.Fatalf("Push(unknown) error = %v, want ErrRouteNotFound", err)
	}
	if current := router.Current(); current != RouteMain {
		t.Fatalf("Current() = %q, want fallback %q", current, RouteMain)
	}
	if len(router.history) != 1 {
		t.Fatalf("history length = %d, want 1", len(router.history))
	}
}

func TestRouterLatestNavigationWinsDuringTransition(t *testing.T) {
	fyneApp := test.NewApp()
	t.Cleanup(fyneApp.Quit)

	router := newTestRouter(t)
	if err := router.Replace(RouteMain, nil); err != nil {
		t.Fatalf("Replace(RouteMain) returned an error: %v", err)
	}
	if err := router.Push(
		routeTestDetails,
		CompanyDetailsParams{CompanyID: 10},
	); err != nil {
		t.Fatalf("Push(routeTestDetails) returned an error: %v", err)
	}

	router.Back()
	time.Sleep(canvas.DurationShort * 3)

	if current := router.Current(); current != RouteMain {
		t.Fatalf("Current() = %q, want latest route %q", current, RouteMain)
	}

	outlet := router.View().(*fyne.Container)
	renderedLabel, ok := outlet.Objects[0].(*widget.Label)
	if !ok {
		t.Fatalf("rendered object type = %T, want *widget.Label", outlet.Objects[0])
	}
	if renderedLabel.Text != "Main" {
		t.Fatalf("rendered label = %q, want Main", renderedLabel.Text)
	}
}

func TestRouterRejectsDuplicateRegistration(t *testing.T) {
	router := NewRouter(RouteMain)
	factory := func(_ any) (fyne.CanvasObject, error) {
		return widget.NewLabel("Main"), nil
	}

	if err := router.Register(RouteMain, factory); err != nil {
		t.Fatalf("first Register returned an error: %v", err)
	}
	if err := router.Register(RouteMain, factory); !errors.Is(err, ErrRouteAlreadyAdded) {
		t.Fatalf("second Register error = %v, want ErrRouteAlreadyAdded", err)
	}
}

func TestRouterResetClearsHistory(t *testing.T) {
	fyneApp := test.NewApp()
	t.Cleanup(fyneApp.Quit)

	router := newTestRouter(t)
	router.history = []historyEntry{
		{route: RouteMain, params: nil, view: nil},
		{route: routeTestDetails, params: nil, view: nil},
	}

	fmt.Printf("Antes do reset: %v", router.history)

	if err := router.Reset(RouteMain, nil); err != nil {
		t.Fatalf("Reset(RouteMain) returned an error: %v", err)
	}
	fmt.Printf("Depois do reset: %v", router.history)

	if got := len(router.history); got != 1 {
		t.Fatalf("history length after Reset = %d, want 1", got)
	}
	if got := router.Current(); got != RouteMain {
		t.Fatalf("Current() after Reset = %q, want %q", got, RouteMain)
	}
}

func TestRouterResetWithNoRegisteredRoute(t *testing.T) {
	fyneApp := test.NewApp()
	t.Cleanup(fyneApp.Quit)

	router := newTestRouter(t)
	router.history = []historyEntry{
		{route: RouteMain, params: nil, view: nil},
		{route: routeTestDetails, params: nil, view: nil},
	}

	fmt.Printf("Antes do reset: %v\n\n", router.history)

	err := router.Reset(RouteCalculationCreate, nil)
	if err == nil { //non registered
		t.Fatalf("Reset with a non registered route didn't returned an error")
	}
	fmt.Println(err)
	//in this case the history needs to be restored an putted the fallback in the history
	// EXPECTED: [main, test.details, main] -> len = 3
	expected := []historyEntry{
		{route: RouteMain, params: nil, view: nil},
		{route: routeTestDetails, params: nil, view: nil},
		{route: RouteMain, params: nil, view: nil},
	}

	if got := len(router.history); got != 3 {
		t.Fatalf("history length after Reset = %d, want 3", got)
	}
	if got := router.Current(); got != RouteMain {
		t.Fatalf("Current() after Reset = %q, want %q", got, RouteMain)
	}
	if equal := slices.EqualFunc(
		router.history,
		expected,
		func(r, ex historyEntry) bool {
			return r.route == ex.route
		},
	); equal == false {
		t.Fatalf("Hitory received is diferent from the expected after Reset.\nGot: %q, want %q", router.history, expected)
	}
	fmt.Printf("Depois do reset: %v\n\n", router.history)
}

func newTestRouter(t *testing.T) *Router {
	t.Helper()

	router := NewRouter(RouteMain)
	if err := router.Register(
		RouteMain,
		func(_ any) (fyne.CanvasObject, error) {
			return widget.NewLabel("Main"), nil
		},
	); err != nil {
		t.Fatalf("Register(RouteMain) returned an error: %v", err)
	}
	if err := router.Register(
		routeTestDetails,
		func(_ any) (fyne.CanvasObject, error) {
			return widget.NewLabel("Details"), nil
		},
	); err != nil {
		t.Fatalf("Register(routeTestDetails) returned an error: %v", err)
	}

	return router
}
