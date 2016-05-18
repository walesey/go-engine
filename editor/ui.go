package editor

import (
	"bytes"
	"os"
	"strings"

	"github.com/walesey/go-engine/assets"
	"github.com/walesey/go-engine/controller"
	"github.com/walesey/go-engine/glfwController"
	"github.com/walesey/go-engine/ui"
	vmath "github.com/walesey/go-engine/vectormath"
)

func (e *Editor) setupUI() {
	//images
	loadImageAsset("file", FileIconData, e.uiAssets)
	loadImageAsset("copy", CopyIconData, e.uiAssets)
	loadImageAsset("reference", LinkIconData, e.uiAssets)
	loadImageAsset("folderOpen", FolderOpenData, e.uiAssets)
	loadImageAsset("folderClosed", FolderClosedData, e.uiAssets)
	loadImageAsset("planetOpen", PlanetOpenData, e.uiAssets)
	loadImageAsset("planetClosed", PlanetClosedData, e.uiAssets)
	loadImageAsset("trash", TrashIconData, e.uiAssets)
	loadImageAsset("geometry", GeometryIconData, e.uiAssets)
	loadImageAsset("scale", ScaleIconData, e.uiAssets)
	loadImageAsset("translate", TranslateIconData, e.uiAssets)
	loadImageAsset("rotate", RotateIconData, e.uiAssets)
	loadImageAsset("reset", ResetIconData, e.uiAssets)

	e.initOverviewMenu()

	e.customController.BindAction(func() {
		if e.mainMenuOpen {
			e.mainMenuOpen = false
			e.closeMainMenu()
		} else {
			e.mainMenuOpen = true
			e.openMainMenu()
		}
	}, controller.KeyEscape, controller.Press)
}

func loadImageAsset(key, data string, uiAssets ui.HtmlAssets) {
	img, _ := assets.DecodeImage(bytes.NewBuffer(assets.Base64ToBytes(data)))
	uiAssets.AddImage(key, img)
}

func (e *Editor) closeMainMenu() {
	e.gameEngine.RemoveOrtho(e.mainMenu, false)
}

func (e *Editor) openMainMenu() {
	if e.mainMenu == nil {

		e.uiAssets.AddCallback("open", func(element ui.Element, args ...interface{}) {
			if len(args) >= 2 && !args[1].(bool) { // not on release
				e.openFileBrowser("Open", func(filePath string) {
					e.loadMap(filePath)
					e.closeFileBrowser()
				}, ".map")
			}
		})
		e.uiAssets.AddCallback("save", func(element ui.Element, args ...interface{}) {
			if len(args) >= 2 && !args[1].(bool) { // not on release
				e.openFileBrowser("Save", func(filePath string) {
					e.saveMap(filePath)
					e.closeFileBrowser()
				}, ".map")
			}
		})
		e.uiAssets.AddCallback("exit", func(element ui.Element, args ...interface{}) {
			os.Exit(0)
		})

		window, container, _ := e.defaultWindow()
		window.SetTranslation(vmath.Vector3{200, 200, 1})
		window.SetScale(vmath.Vector3{400, 0, 1})

		e.controllerManager.AddController(ui.NewUiController(window).(glfwController.Controller))
		ui.LoadHTML(container, window, strings.NewReader(mainMenuHtml), strings.NewReader(globalCss), e.uiAssets)

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

	e.customController.BindMouseAction(func() {
		ui.DeactivateAllTextElements(container)
	}, controller.MouseButton1, controller.Press)

	return
}
