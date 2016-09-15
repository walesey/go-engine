package editor

import (
	"fmt"
	"strings"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/walesey/go-engine/glfwController"
	"github.com/walesey/go-engine/ui"
)

func (e *Editor) closeProgressBar() {
	e.gameEngine.RemoveOrtho(e.progressBar, false)
}

func (e *Editor) openProgressBar() {
	if e.progressBar == nil {
		window := ui.NewWindow()
		window.SetTranslation(mgl32.Vec3{500, 0, 0})
		window.SetScale(mgl32.Vec3{330, 0, 1})

		container := ui.NewContainer()
		container.SetBackgroundColor(200, 200, 200, 255)
		window.SetElement(container)

		e.controllerManager.AddController(ui.NewUiController(window).(glfwController.Controller))
		ui.LoadHTML(container, strings.NewReader(progressBarHtml), strings.NewReader(globalCss), e.uiAssets)
		window.Render()

		e.progressBar = window
	}
	e.gameEngine.AddOrtho(e.progressBar)
}

func (e *Editor) setProgressBar(progress int) {
	for i := 1; i <= 20; i++ {
		container, ok := e.progressBar.ElementById(fmt.Sprintf("progress%v", i)).(*ui.Container)
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
	container, ok := e.progressBar.ElementById("progressBarMessage").(*ui.Container)
	if ok {
		container.RemoveAllChildren()
		html := fmt.Sprintf("<p>%v</p>", message)
		css := "p { font-size: 8px; }"
		ui.LoadHTML(container, strings.NewReader(html), strings.NewReader(css), e.uiAssets)
		e.progressBar.Render()
	}
}
