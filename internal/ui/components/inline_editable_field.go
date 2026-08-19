package components

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// InlineEditableField displays a value as a label and swaps to an Entry when
// its owning form activates the edit session.
type InlineEditableField struct {
	widget.BaseWidget

	value           *widget.Label
	entry           *widget.Entry
	validation      *widget.Label
	readView        *fyne.Container
	editView        *fyne.Container
	content         *fyne.Container
	original        string
	editing         bool
	validationError error
	onActivate      func(*InlineEditableField)
	onChanged       func(*InlineEditableField)
}

func NewInlineEditableField(text string, width float32) *InlineEditableField {
	field := &InlineEditableField{
		value:    widget.NewLabel(text),
		entry:    widget.NewEntry(),
		original: text,
	}
	field.value.Truncation = fyne.TextTruncateEllipsis
	field.entry.SetText(text)
	field.entry.TextStyle = fyne.TextStyle{Bold: true}
	field.entry.AlwaysShowValidationError = true
	field.entry.OnChanged = func(string) {
		if field.onChanged != nil {
			field.onChanged(field)
		}
	}

	readBackground := newRoundedInformationRectangle(
		color.NRGBA{R: 240, G: 240, B: 240, A: 255},
		width,
		field.entry.MinSize().Height,
	)
	field.readView = container.NewStack(readBackground, field.value)

	editBackground := newRoundedInformationRectangle(
		color.NRGBA{R: 224, G: 237, B: 255, A: 255},
		width,
		field.entry.MinSize().Height,
	)
	editBackground.StrokeColor = theme.Color(theme.ColorNamePrimary)
	editBackground.StrokeWidth = 2
	field.editView = container.NewStack(editBackground, field.entry)

	field.validation = widget.NewLabelWithStyle(
		"",
		fyne.TextAlignLeading,
		fyne.TextStyle{},
	)
	field.validation.Hide()
	field.content = container.NewVBox(field.readView, field.validation)
	field.entry.Hide()
	field.ExtendBaseWidget(field)
	return field
}

func (field *InlineEditableField) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(field.content)
}

// SetOnActivate registers the callback used by the containing form to enter
// edit mode for every field in the screen.
func (field *InlineEditableField) SetOnActivate(callback func(*InlineEditableField)) {
	field.onActivate = callback
}

// SetOnChanged registers a callback for changes to the active input value.
// The callback is also notified when a value is restored or committed so the
// owning form can keep its dirty state synchronized.
func (field *InlineEditableField) SetOnChanged(callback func(*InlineEditableField)) {
	field.onChanged = callback
}

// BeginEdit shows the input control without changing its original value.
func (field *InlineEditableField) BeginEdit() {
	if field.editing {
		return
	}

	field.editing = true
	field.value.Hide()
	field.entry.Show()
	field.setPresentation(field.editView)
	field.Refresh()
}

// EndEdit commits the current entry value as the new display and snapshot.
func (field *InlineEditableField) EndEdit() {
	field.original = field.entry.Text
	field.value.SetText(field.entry.Text)
	field.entry.Hide()
	field.value.Show()
	field.editing = false
	field.setPresentation(field.readView)
	field.SetValidationError(nil)
	field.Refresh()
	field.notifyChanged()
}

// CancelEdit restores the value captured when the form entered edit mode.
func (field *InlineEditableField) CancelEdit() {
	field.entry.SetText(field.original)
	field.value.SetText(field.original)
	field.entry.Hide()
	field.value.Show()
	field.editing = false
	field.setPresentation(field.readView)
	field.SetValidationError(nil)
	field.Refresh()
	field.notifyChanged()
}

// Text returns the value currently being edited or displayed.
func (field *InlineEditableField) Text() string {
	if field.editing {
		return field.entry.Text
	}

	return field.value.Text
}

// SetText updates both the editable value and the displayed value while
// preserving edit mode.
func (field *InlineEditableField) SetText(text string) {
	field.original = text
	field.entry.SetText(text)
	field.value.SetText(text)
	field.SetValidationError(nil)
	field.Refresh()
	field.notifyChanged()
}

// SetEditingText changes only the active input value. It is useful to keep
// the original snapshot intact until the user saves or cancels the form.
func (field *InlineEditableField) SetEditingText(text string) {
	if !field.editing {
		field.SetText(text)
		return
	}

	field.entry.SetText(text)
	field.Refresh()
}

// SetValidationError renders the error next to the field and uses Fyne's
// standard entry validation state as an additional visual cue.
func (field *InlineEditableField) SetValidationError(err error) {
	field.validationError = err
	field.entry.SetValidationError(err)
	if err == nil {
		field.validation.SetText("")
		field.validation.Hide()
		field.Refresh()
		return
	}

	field.validation.SetText(err.Error())
	field.validation.Show()
	field.Refresh()
}

// ValidationError returns the currently displayed field validation error.
func (field *InlineEditableField) ValidationError() error {
	return field.validationError
}

// HasChanges reports whether the active entry differs from the value captured
// when the field entered edit mode.
func (field *InlineEditableField) HasChanges() bool {
	return field.editing && field.entry.Text != field.original
}

func (field *InlineEditableField) IsEditing() bool {
	return field.editing
}

// Tapped activates the containing form when present, otherwise it behaves as
// a standalone inline field.
func (field *InlineEditableField) Tapped(*fyne.PointEvent) {
	if field.editing {
		return
	}

	if field.onActivate != nil {
		field.onActivate(field)
		return
	}

	field.BeginEdit()
	field.focus()
}

func (field *InlineEditableField) focus() {
	if app := fyne.CurrentApp(); app != nil {
		if canvas := app.Driver().CanvasForObject(field.entry); canvas != nil {
			canvas.Focus(field.entry)
		}
	}
}

func (field *InlineEditableField) notifyChanged() {
	if field.onChanged != nil {
		field.onChanged(field)
	}
}

func (field *InlineEditableField) setPresentation(presentation fyne.CanvasObject) {
	field.content.Objects[0] = presentation
	field.content.Refresh()
}

func newRoundedInformationRectangle(fillColor color.Color, width, height float32) *canvas.Rectangle {
	rectangle := canvas.NewRectangle(fillColor)
	rectangle.SetMinSize(fyne.NewSize(width, height))
	rectangle.CornerRadius = 5
	rectangle.StrokeColor = color.Black
	rectangle.StrokeWidth = 0.5
	return rectangle
}

func newTransparentInformationRectangle(width, height float32) *canvas.Rectangle {
	rectangle := canvas.NewRectangle(color.Transparent)
	rectangle.SetMinSize(fyne.NewSize(width, height))
	return rectangle
}
