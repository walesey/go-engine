package editor

import (
	"strings"

	"github.com/walesey/go-engine/engine"
	"github.com/walesey/go-engine/ui"
)

func (e *Editor) initOverviewMenu() {
	e.uiAssets.AddCallback("import", func(element ui.Element, args ...interface{}) {
		if len(args) >= 2 && !args[1].(bool) { // not on release
			e.openFileBrowser("import", func(filePath string) {
				// e.loadMap(filePath) // TODO
				e.closeFileBrowser()
			})
		}
	})

	window, container, _ := e.defaultWindow()
	window.SetTranslation(vmath.Vector3{10, 10, 1})
	window.SetScale(vmath.Vector3{500, 0, 1})

	e.controllerManager.AddController(ui.NewUiController(window))
	ui.LoadHTML(container, strings.NewReader(overviewMenuHtml), strings.NewReader(globalCss), e.uiAssets)

	e.gameEngine.AddUpdatable(engine.UpdatableFunc(func(dt float64) {
		window.Render()
	}))
	e.gameEngine.AddOrtho(window)
}
