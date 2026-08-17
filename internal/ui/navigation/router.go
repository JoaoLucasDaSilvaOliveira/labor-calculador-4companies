package navigation

import (
	"errors"
	"fmt"

	"labor-calculador-4companies/internal/ui/components"

	"fyne.io/fyne/v2"
)

var (
	ErrEmptyRouteID      = errors.New("o identificador da rota não pode ser vazio")
	ErrNilPageFactory    = errors.New("a factory da página não pode ser nula")
	ErrRouteAlreadyAdded = errors.New("a rota já foi registrada")
	ErrRouteNotFound     = errors.New("a rota não foi registrada")
	ErrNilPage           = errors.New("a factory retornou uma página nula")
)

// PageFactory creates a page for the supplied route parameters.
type PageFactory func(params any) (fyne.CanvasObject, error)

// Navigator is the small navigation contract injected into pages.
// Pages can request navigation without depending on the concrete Router.
type Navigator interface {
	Push(route RouteID, params any) error
	Replace(route RouteID, params any) error
	Reset(route RouteID, params any) error
	Back() bool
	CanGoBack() bool
	Current() RouteID
}

var _ Navigator = (*Router)(nil)

// historyEntry keeps the semantic route together with its rendered view.
// Keeping the view allows Back to restore local widget state, such as form input.
type historyEntry struct {
	route  RouteID
	params any
	view   fyne.CanvasObject
}

// Router coordinates application-level page navigation.
// The AnimatedContent outlet is responsible only for the visual transition.
type Router struct {
	outlet        *components.AnimatedContent
	fallbackRoute RouteID
	factories     map[RouteID]PageFactory
	history       []historyEntry
}

// NewRouter creates an empty router with the route used as navigation fallback.
func NewRouter(fallbackRoute RouteID) *Router {
	return &Router{
		outlet:        components.NewAnimatedContent(),
		fallbackRoute: fallbackRoute,
		factories:     make(map[RouteID]PageFactory),
	}
}

// Register associates a stable route identifier with a page factory.
func (r *Router) Register(route RouteID, factory PageFactory) error {
	if route == "" {
		return ErrEmptyRouteID
	}
	if factory == nil {
		return fmt.Errorf("%w: %s", ErrNilPageFactory, route)
	}
	if _, exists := r.factories[route]; exists {
		return fmt.Errorf("%w: %s", ErrRouteAlreadyAdded, route)
	}

	r.factories[route] = factory
	return nil
}

// View returns the stable container that must be mounted in the main window.
func (r *Router) View() fyne.CanvasObject {
	return r.outlet.View()
}

// Push opens a page and adds it to the navigation history.
func (r *Router) Push(route RouteID, params any) error {
	entry, err := r.buildEntry(route, params)
	if err != nil {
		return r.handleNavigationError(err)
	}

	r.history = append(r.history, entry)
	r.show(entry)
	return nil
}

// Replace opens a page in place of the current history entry.
func (r *Router) Replace(route RouteID, params any) error {
	entry, err := r.buildEntry(route, params)
	if err != nil {
		return r.handleNavigationError(err)
	}

	r.replaceCurrentEntry(entry)
	r.show(entry)
	return nil
}

// Back restores the previous history entry without rebuilding its widgets.
func (r *Router) Back() bool {
	if !r.CanGoBack() {
		return false
	}

	r.history = r.history[:len(r.history)-1]
	r.show(r.history[len(r.history)-1])
	return true
}

// CanGoBack reports whether Back has a previous page to restore.
func (r *Router) CanGoBack() bool {
	return len(r.history) > 1
}

// Current returns the active semantic route, or an empty ID before navigation starts.
func (r *Router) Current() RouteID {
	if len(r.history) == 0 {
		return ""
	}

	return r.history[len(r.history)-1].route
}

func (r *Router) Reset(route RouteID, params any) error {
	entry, err := r.buildEntry(route, params)
	if err == nil {
		r.history = []historyEntry{entry}
		r.show(entry)
		return nil
	}

	// Reset starts a new navigation session. If its destination cannot be
	// built, the previous session must not be restored behind the fallback.
	fallbackEntry, fallbackErr := r.buildEntry(r.fallbackRoute, nil)
	if fallbackErr == nil {
		r.history = []historyEntry{fallbackEntry}
		r.show(fallbackEntry)
		return err
	}

	return errors.Join(err, fmt.Errorf("não foi possível abrir a rota fallback: %w", fallbackErr))
}

func (r *Router) buildEntry(route RouteID, params any) (historyEntry, error) {
	factory, exists := r.factories[route]
	if !exists {
		return historyEntry{}, fmt.Errorf("%w: %s", ErrRouteNotFound, route)
	}

	view, err := factory(params)
	if err != nil {
		return historyEntry{}, fmt.Errorf("não foi possível construir a rota %s: %w", route, err)
	}
	if view == nil {
		return historyEntry{}, fmt.Errorf("%w: %s", ErrNilPage, route)
	}

	return historyEntry{
		route:  route,
		params: params,
		view:   view,
	}, nil
}

func (r *Router) handleNavigationError(navigationErr error) error {
	fallbackEntry, fallbackErr := r.buildEntry(r.fallbackRoute, nil)
	if fallbackErr != nil {
		return errors.Join(navigationErr, fmt.Errorf("não foi possível abrir a rota fallback: %w", fallbackErr))
	}

	r.replaceCurrentEntry(fallbackEntry)
	r.show(fallbackEntry)
	return navigationErr
}

func (r *Router) replaceCurrentEntry(entry historyEntry) {
	if len(r.history) == 0 {
		r.history = append(r.history, entry)
		return
	}

	r.history[len(r.history)-1] = entry
}

func (r *Router) show(entry historyEntry) {
	r.outlet.SetContent(entry.view, nil)
}
