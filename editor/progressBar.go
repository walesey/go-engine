package editor

import (
	"fmt"
	"strings"

	"github.com/walesey/go-engine/engine"
	"github.com/walesey/go-engine/ui"
	vmath "github.com/walesey/go-engine/vectormath"
)

func (e *Editor) closeProgressBar() {
	e.gameEngine.RemoveOrtho(e.progressBar, false)
}

func (e *Editor) openProgressBar() {
	if e.progressBar == nil {
		window := ui.NewWindow()
		window.SetTranslation(vmath.Vector3{500, 0, 0})
		window.SetScale(vmath.Vector3{330, 0, 1})

		container := ui.NewContainer()
		container.SetBackgroundColor(200, 200, 200, 255)
		window.SetElement(container)

		e.controllerManager.AddController(ui.NewUiController(window))
		ui.LoadHTML(container, strings.NewReader(progressBarHtml), strings.NewReader(globalCss), e.uiAssets)

		e.gameEngine.AddUpdatable(engine.UpdatableFunc(func(dt float64) {
			window.Render()
		}))
		e.progressBar = window
	}
	e.gameEngine.AddOrtho(e.progressBar)
}

func (e *Editor) setProgressBar(progress int) {
	for i := 1; i <= 20; i++ {
		elem := e.progressBar.ElementById(fmt.Sprintf("progress%v", i))
		if elem == nil {
			continue
		}
		container, ok := elem.(*ui.Container)
		if ok {
			if i > progress {
				container.SetBackgroundColor(0, 0, 0, 0)
			} else {
				container.SetBackgroundColor(0, 255, 0, 255)
			}
		}
	}
}

func (e *Editor) setProgressTime(message string) {
	elem := e.progressBar.ElementById("progressBarMessage")
	container, ok := elem.(*ui.Container)
	if ok {
		container.RemoveAllChildren()
		html := fmt.Sprintf("<p>%v</p>", message)
		css := "p { font-size: 8px; }"
		ui.LoadHTML(container, strings.NewReader(html), strings.NewReader(css), e.uiAssets)
	}
}
