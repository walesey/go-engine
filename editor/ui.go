package editor

import (
	"log"
	"strings"

	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/walesey/go-engine/engine"
	"github.com/walesey/go-engine/ui"
	vmath "github.com/walesey/go-engine/vectormath"
)

func (e *Editor) setupUI() {
	e.uiAssets = ui.NewHtmlAssets()
	audiowideFont, err := ui.LoadFont("TestAssets/Audiowide-Regular.ttf")
	if err != nil {
		log.Printf("Error loading ui font: %v", err)
	}
	e.uiAssets.AddFont("default", audiowideFont)

	e.customController.BindAction(func() {
		if e.mainMenuOpen {
			e.mainMenuOpen = false
			e.closeMainMenu()
		} else {
			e.mainMenuOpen = true
			e.openMainMenu()
		}
	}, glfw.KeyEscape, glfw.Press)
}

func (e *Editor) closeMainMenu() {
	e.gameEngine.RemoveOrtho(e.mainMenu, false)
}

func (e *Editor) openMainMenu() {
	if e.mainMenu == nil {
		window, container, _ := e.defaultWindow()
		window.SetTranslation(vmath.Vector3{0, 0, 1})
		window.SetScale(vmath.Vector3{300, 0, 1})

		e.controllerManager.AddController(ui.NewUiController(window))
		ui.LoadPage(container, strings.NewReader(mainMenuHtml), strings.NewReader(globalCss), e.uiAssets)

		e.gameEngine.AddUpdatable(engine.UpdatableFunc(func(dt float64) {
			window.Render()
		}))

		e.mainMenu = window
	}
	e.gameEngine.AddOrtho(e.mainMenu)
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
