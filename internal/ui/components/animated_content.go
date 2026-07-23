package components

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
)

// AnimatedContent replaces its content behind a short fade transition.
type AnimatedContent struct {
	container  *fyne.Container
	overlay    *canvas.Rectangle
	animation  *fyne.Animation
	hasContent bool
}

// NewAnimatedContent creates an empty transition host.
func NewAnimatedContent() *AnimatedContent {
	overlay := canvas.NewRectangle(color.Transparent)
	return &AnimatedContent{
		container: container.NewStack(overlay),
		overlay:   overlay,
	}
}

// View returns the container that must be mounted in the parent layout.
func (a *AnimatedContent) View() fyne.CanvasObject {
	return a.container
}

// SetContent fades out the current view, swaps it, and fades in the next view.
func (a *AnimatedContent) SetContent(content fyne.CanvasObject, onSwapped func()) {
	if !a.hasContent {
		a.replace(content)
		if onSwapped != nil {
			onSwapped()
		}
		return
	}

	a.fadeTo(255, func() {
		a.replace(content)
		if onSwapped != nil {
			onSwapped()
		}
		a.fadeTo(0, nil)
	})
}

func (a *AnimatedContent) replace(content fyne.CanvasObject) {
	a.container.Objects = []fyne.CanvasObject{content, a.overlay}
	a.container.Refresh()
	a.hasContent = true
}

func (a *AnimatedContent) fadeTo(alpha uint8, onFinished func()) {
	if a.animation != nil {
		a.animation.Stop()
	}

	startAlpha := color.NRGBAModel.Convert(a.overlay.FillColor).(color.NRGBA).A
	if startAlpha == alpha {
		if onFinished != nil {
			onFinished()
		}
		return
	}

	var animation *fyne.Animation
	animation = fyne.NewAnimation(canvas.DurationShort, func(progress float32) {
		currentAlpha := uint8(float32(startAlpha) + (float32(alpha)-float32(startAlpha))*progress)
		a.setOverlayAlpha(currentAlpha)
		if progress == 1 && a.animation == animation && onFinished != nil {
			onFinished()
		}
	})
	animation.Curve = fyne.AnimationEaseInOut
	a.animation = animation
	animation.Start()
}

func (a *AnimatedContent) setOverlayAlpha(alpha uint8) {
	background := color.NRGBAModel.Convert(theme.Color(theme.ColorNameBackground)).(color.NRGBA)
	background.A = alpha
	a.overlay.FillColor = background
	a.overlay.Refresh()
}
