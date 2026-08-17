package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func NewLoadingState(message string) fyne.CanvasObject {
	return centeredState(message, nil)
}

func NewEmptyState(title, description string) fyne.CanvasObject {
	return centeredState(title+"\n"+description, nil)
}

func NewRecoverableErrorState(message string, onBack func()) fyne.CanvasObject {
	var action fyne.CanvasObject
	if onBack != nil {
		action = widget.NewButton("Voltar", onBack)
	}
	return centeredState(message, action)
}

func centeredState(message string, action fyne.CanvasObject) fyne.CanvasObject {
	label := widget.NewLabel(message)
	label.Wrapping = fyne.TextWrapWord

	objects := []fyne.CanvasObject{label}
	if action != nil {
		objects = append(objects, action)
	}

	return container.NewCenter(container.NewVBox(objects...))
}
