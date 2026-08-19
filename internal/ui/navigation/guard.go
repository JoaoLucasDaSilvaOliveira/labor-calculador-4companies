package navigation

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// NavigationGuard can defer a navigation operation until the user resolves a
// pending edit session.
type NavigationGuard interface {
	ShouldBlock() bool
	Confirm(onContinue func())
}

// UnsavedChangesGuard is the workspace guard used by the desktop UI. It owns
// only the confirmation flow; the page remains responsible for restoring its
// own fields through EditSession.
type UnsavedChangesGuard struct {
	Session *EditSession
	Parent  fyne.Window

	dialogOpen bool
}

func (g *UnsavedChangesGuard) ShouldBlock() bool {
	return g != nil && g.Session != nil && g.Session.Active()
}

func (g *UnsavedChangesGuard) Confirm(onContinue func()) {
	if g == nil || !g.ShouldBlock() || g.dialogOpen {
		return
	}

	if g.Parent == nil {
		g.Session.Discard()
		onContinue()
		return
	}

	g.dialogOpen = true
	dialog.NewCustomConfirm(
		"Alterações não salvas",
		"Continuar",
		"Voltar à edição",
		widget.NewLabel("Há alterações não salvas. Deseja continuar e descartá-las?"),
		func(confirmed bool) {
			g.dialogOpen = false
			if !confirmed {
				return
			}

			g.Session.Discard()
			onContinue()
		},
		g.Parent,
	).Show()
}
