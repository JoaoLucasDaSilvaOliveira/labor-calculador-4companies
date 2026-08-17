package components

import "fyne.io/fyne/v2"

// fixedHeightLayout keeps a list's viewport at the height calculated by the
// prototype while still allowing it to use all available horizontal space.
// This is deliberately a small visual layout primitive, not a screen-level
// abstraction.
type fixedHeightLayout struct {
	height float32
}

func (l fixedHeightLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	for _, object := range objects {
		object.Move(fyne.NewPos(0, 0))
		object.Resize(fyne.NewSize(size.Width, l.height))
	}
}

func (l fixedHeightLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	minSize := fyne.NewSize(0, l.height)
	for _, object := range objects {
		minSize.Width = max(minSize.Width, object.MinSize().Width)
	}
	return minSize
}

func newFixedHeightObject(object fyne.CanvasObject, height float32) fyne.CanvasObject {
	return fyne.NewContainerWithLayout(fixedHeightLayout{height: height}, object)
}
