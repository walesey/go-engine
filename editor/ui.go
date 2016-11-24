package editor

import (
	"bytes"
	"os"
	"strings"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/walesey/go-engine/assets"
	"github.com/walesey/go-engine/controller"
	"github.com/walesey/go-engine/glfwController"
	"github.com/walesey/go-engine/ui"
	"github.com/walesey/go-engine/util"
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
	loadImageAsset("cog", CogIconData, e.uiAssets)

	// callbacks used to highlight text active text fields
	e.uiAssets.AddCallback("inputfocus", func(element ui.Element, args ...interface{}) {
		container, ok := element.(*ui.Container)
		if ok {
			container.SetBackgroundColor(255, 255, 50, 255)
		}
	})
	e.uiAssets.AddCallback("inputblur", func(element ui.Element, args ...interface{}) {
		container, ok := element.(*ui.Container)
		if ok {
			container.SetBackgroundColor(255, 255, 255, 255)
		}
	})

	e.initOverviewMenu()
	e.gameEngine.InitFpsDial()

	e.customController.BindKeyAction(func() {
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
	img, _ := assets.DecodeImage(bytes.NewBuffer(util.Base64ToBytes(data)))
	uiAssets.AddImage(key, img)
}

func (e *Editor) closeMainMenu() {
	e.gameEngine.RemoveSpatial(e.mainMenu, false)
	e.controllerManager.RemoveController(e.mainMenuController)
}

func (e *Editor) openMainMenu() {
	if e.mainMenu == nil {

		e.uiAssets.AddCallback("open", func(element ui.Element, args ...interface{}) {
			if len(args) >= 2 && !args[1].(bool) { // not on release
				e.openFileBrowser("Open", func(filePath string) {
					e.loadMap(filePath)
					e.closeFileBrowser()
				}, ".json")
			}
		})
		e.uiAssets.AddCallback("save", func(element ui.Element, args ...interface{}) {
			if len(args) >= 2 && !args[1].(bool) { // not on release
				e.openFileBrowser("Save", func(filePath string) {
					e.saveMap(filePath)
					e.closeFileBrowser()
				}, ".json")
			}
		})
		e.uiAssets.AddCallback("exit", func(element ui.Element, args ...interface{}) {
			os.Exit(0)
		})

		window, container, _ := e.defaultWindow()
		window.SetTranslation(mgl32.Vec3{200, 200, 1})
		window.SetScale(mgl32.Vec3{400, 0, 1})

		ui.LoadHTML(container, strings.NewReader(mainMenuHtml), strings.NewReader(globalCss), e.uiAssets)
		window.Render()

		e.mainMenuController = ui.NewUiController(window).(glfwController.Controller)
		e.mainMenu = window
	}
	e.gameEngine.AddOrtho(e.mainMenu)
	e.controllerManager.AddController(e.mainMenuController)
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
