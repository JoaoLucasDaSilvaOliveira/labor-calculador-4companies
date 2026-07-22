package components

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

// ClickableText wraps a text primitive to add click support.
type ClickableText struct {
	widget.BaseWidget
	TextObj      *canvas.Text
	DefaultColor color.Color
	HoverColor   color.Color
	PressedColor color.Color
	OnTapped     func()

	hovered        bool
	pressed        bool
	colorAnimation *fyne.Animation
}

// NewClickableText creates a clickable text with animated default and hover colors.
func NewClickableText(text string, defaultColor, hoverColor color.Color, tapped func()) *ClickableText {
	c := &ClickableText{
		TextObj:      canvas.NewText(text, defaultColor),
		DefaultColor: defaultColor,
		HoverColor:   hoverColor,
		PressedColor: lighterColor(hoverColor),
		OnTapped:     tapped,
	}

	c.ExtendBaseWidget(c)

	return c
}

// CreateRenderer draws the text primitive.
func (c *ClickableText) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(c.TextObj)
}

// Cursor returns a pointer cursor while the text is hovered.
func (c *ClickableText) Cursor() desktop.Cursor {
	if c.hovered {
		return desktop.PointerCursor
	}

	return desktop.DefaultCursor
}

// MouseIn starts the transition to the hover color.
func (c *ClickableText) MouseIn(*desktop.MouseEvent) {
	c.hovered = true
	c.animateColor(c.HoverColor)
}

// MouseMoved is required by desktop.Hoverable.
func (c *ClickableText) MouseMoved(*desktop.MouseEvent) {}

// MouseDown starts the pressed-state color transition.
func (c *ClickableText) MouseDown(event *desktop.MouseEvent) {
	if event.Button != desktop.MouseButtonPrimary {
		return
	}

	c.pressed = true
	c.animateColor(c.PressedColor)
}

// MouseUp restores the hover or default color after a click.
func (c *ClickableText) MouseUp(*desktop.MouseEvent) {
	if !c.pressed {
		return
	}

	c.pressed = false
	c.restoreInteractionColor()
}

// MouseOut starts the transition back to the default color.
func (c *ClickableText) MouseOut() {
	c.hovered = false
	c.pressed = false
	c.animateColor(c.DefaultColor)
}

func (c *ClickableText) restoreInteractionColor() {
	if c.hovered {
		c.animateColor(c.HoverColor)
		return
	}

	c.animateColor(c.DefaultColor)
}

func (c *ClickableText) animateColor(target color.Color) {
	if c.colorAnimation != nil {
		c.colorAnimation.Stop()
	}

	c.colorAnimation = canvas.NewColorRGBAAnimation(c.TextObj.Color, target, canvas.DurationShort, func(current color.Color) {
		c.TextObj.Color = current
		c.TextObj.Refresh()
	})
	c.colorAnimation.Curve = fyne.AnimationEaseInOut
	c.colorAnimation.Start()
}

func lighterColor(source color.Color) color.Color {
	current := color.NRGBAModel.Convert(source).(color.NRGBA)
	redMultiplier, greenMultiplier, blueMultiplier := 1.2, 1.2, 1.2

	switch {
	case current.R == current.G && current.G == current.B:
		redMultiplier, greenMultiplier, blueMultiplier = 1.5, 1.5, 1.5
	case current.R >= current.G && current.R >= current.B:
		redMultiplier = 1.5
	case current.G >= current.B:
		greenMultiplier = 1.5
	default:
		blueMultiplier = 1.5
	}

	return color.NRGBA{
		R: multiplyColorChannel(current.R, redMultiplier),
		G: multiplyColorChannel(current.G, greenMultiplier),
		B: multiplyColorChannel(current.B, blueMultiplier),
		A: current.A,
	}
}

func multiplyColorChannel(value uint8, multiplier float64) uint8 {
	brightened := float64(value) * multiplier
	if brightened > 255 {
		return 255
	}

	return uint8(brightened)
}

// Tapped intercepts the primary mouse button click.
func (c *ClickableText) Tapped(event *fyne.PointEvent) {
	if c.OnTapped != nil {
		c.OnTapped()
	}
}

// TappedSecondary intercepts the secondary mouse button click.
func (c *ClickableText) TappedSecondary(event *fyne.PointEvent) {
	// Intentionally empty.
}
