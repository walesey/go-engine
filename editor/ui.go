package editor

import (
	"os"
	"strings"

	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/walesey/go-engine/engine"
	"github.com/walesey/go-engine/ui"
	vmath "github.com/walesey/go-engine/vectormath"
)

func (e *Editor) setupUI() {
	e.openOverviewMenu()
	e.openProgressBar()
	e.setProgressBar(15)

	e.customController.BindAction(func() {
		if e.mainMenuOpen {
			e.mainMenuOpen = false
			e.closeMainMenu()
		} else {
			e.mainMenuOpen = true
			e.openMainMenu()
		}
	}, glfw.KeyEscape, glfw.Press)

	e.uiAssets.AddCallback("open", func(element ui.Element, args ...interface{}) {

	})
	e.uiAssets.AddCallback("save", func(element ui.Element, args ...interface{}) {

	})
	e.uiAssets.AddCallback("exit", func(element ui.Element, args ...interface{}) {
		os.Exit(0)
	})
}

func (e *Editor) closeMainMenu() {
	e.gameEngine.RemoveOrtho(e.mainMenu, false)
}

func (e *Editor) openMainMenu() {
	if e.mainMenu == nil {
		window, container, _ := e.defaultWindow()
		window.SetTranslation(vmath.Vector3{200, 200, 1})
		window.SetScale(vmath.Vector3{400, 0, 1})

		e.controllerManager.AddController(ui.NewUiController(window))
		ui.LoadHTML(container, strings.NewReader(mainMenuHtml), strings.NewReader(globalCss), e.uiAssets)

		e.gameEngine.AddUpdatable(engine.UpdatableFunc(func(dt float64) {
			window.Render()
		}))
		e.mainMenu = window
	}
	e.gameEngine.AddOrtho(e.mainMenu)
}

func (e *Editor) openOverviewMenu() {
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

func (e *Editor) defaultWindow() (window *ui.Window, container *ui.Container, tab *ui.Container) {
	window = ui.NewWindow()

	tab = ui.NewContainer()
	tab.SetBackgroundColor(70, 70, 70, 255)
	tab.SetHeight(40)

	mainContainer := ui.NewContainer()
	window.SetElement(mainContainer)
	container = ui.NewContainer()
	container.SetBackgroundColor(200, 200, 200, 255)
	mainContainer.AddChildren(tab, container)
	ui.ClickAndDragWindow(window, tab.Hitbox, e.customController)
	return
}
