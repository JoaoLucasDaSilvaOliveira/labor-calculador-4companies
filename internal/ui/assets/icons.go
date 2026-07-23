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

	//go:embed "icons/close-m-glass.png"
	CloseMGlassIconpng []byte

	//go:embed "icons/add-company-icon.png"
	AddCompanypng []byte

	OpenTabIcon         = fyne.NewStaticResource("open-tab-icon.png", openTabIconpng)
	CloseTabIcon        = fyne.NewStaticResource("close-tab-icon.png", closeTabIconpng)
	MagnifingGlass      = fyne.NewStaticResource("m-glass.png", mGlassIconpng)
	CloseMagnifingGlass = fyne.NewStaticResource("close-m-glass.png", CloseMGlassIconpng)
	AddCompany          = fyne.NewStaticResource("add-company-icon.png", AddCompanypng)
)
