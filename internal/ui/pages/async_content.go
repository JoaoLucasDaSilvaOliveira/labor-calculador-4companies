package pages

import (
	"labor-calculador-4companies/internal/ui/components"

	"fyne.io/fyne/v2"
)

func newAsyncPrototypeContent(loadingMessage string, load func() fyne.CanvasObject) fyne.CanvasObject {
	content := components.NewAnimatedContent()
	content.SetContent(components.NewLoadingState(loadingMessage), nil)

	go func() {
		loadedContent := load()
		update := func() { content.SetContent(loadedContent, nil) }
		if fyne.CurrentApp() == nil {
			update()
			return
		}
		fyne.Do(update)
	}()

	return content.View()
}
