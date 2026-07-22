package assets

import (
	_ "embed"

	"fyne.io/fyne/v2"
)

var (
	//go:embed "icons/open-tab-icon.png"
	openTabIconpng []byte

	//go:embed "icons/close-tab-icon.png"
	closeTabIconpng []byte

	//go:embed "icons/m-glass.png"
	mGlassIconpng []byte

	OpenTabIcon  = fyne.NewStaticResource("open-tab-icon.png", openTabIconpng)
	CloseTabIcon = fyne.NewStaticResource("close-tab-icon.png", closeTabIconpng)
	MagnifingGlass = fyne.NewStaticResource("m-glass.png", mGlassIconpng)
)
