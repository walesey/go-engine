package editor

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/walesey/go-engine/assets"
	"github.com/walesey/go-engine/ui"
	vmath "github.com/walesey/go-engine/vectormath"
)

const maxFileDisplayed = 37

type FileBrowser struct {
	window          *ui.Window
	assets          ui.HtmlAssets
	callback        func(filePath string)
	root            string
	scrollOffset    int
	openFolders     map[string]bool
	selectedFile    string
	extensionFilter []string
}

func (e *Editor) closeFileBrowser() {
	if e.fileBrowserOpen {
		e.gameEngine.RemoveOrtho(e.fileBrowser.window, false)
		e.fileBrowserOpen = false
	}
}

func (e *Editor) openFileBrowser(heading string, callback func(filePath string), filters ...string) {
	if e.fileBrowserOpen {
		return
	}
	if e.fileBrowser == nil {
		fileImg, err := assets.DecodeImage(bytes.NewBuffer(assets.Base64ToBytes(FileIconData)))
		if err == nil {
			e.uiAssets.AddImage("file", fileImg)
		}

		folderOpenImg, err := assets.DecodeImage(bytes.NewBuffer(assets.Base64ToBytes(FolderOpenData)))
		if err == nil {
			e.uiAssets.AddImage("folderOpen", folderOpenImg)
		}

		folderClosedImg, err := assets.DecodeImage(bytes.NewBuffer(assets.Base64ToBytes(FolderClosedData)))
		if err == nil {
			e.uiAssets.AddImage("folderClosed", folderClosedImg)
		}

		e.uiAssets.AddCallback("fileBrowserOpen", func(element ui.Element, args ...interface{}) {
			if len(args) >= 2 && !args[1].(bool) { // not on release
				e.fileBrowser.callback(e.fileBrowser.selectedFile)
			}
		})

		e.uiAssets.AddCallback("fileBrowserCancel", func(element ui.Element, args ...interface{}) {
			e.closeFileBrowser()
		})

		e.uiAssets.AddCallback("fileBrowserScrollup", func(element ui.Element, args ...interface{}) {
			if len(args) >= 2 && !args[1].(bool) { // not on release
				if e.fileBrowser.scrollOffset > 0 {
					e.fileBrowser.scrollOffset = e.fileBrowser.scrollOffset - 1
					e.fileBrowser.UpdateFileSystem()
				}
			}
		})

		e.uiAssets.AddCallback("fileBrowserScrollDown", func(element ui.Element, args ...interface{}) {
			if len(args) >= 2 && !args[1].(bool) { // not on release
				e.fileBrowser.scrollOffset = e.fileBrowser.scrollOffset + 1
				e.fileBrowser.UpdateFileSystem()
			}
		})

		window, container, _ := e.defaultWindow()
		window.SetTranslation(vmath.Vector3{100, 100, 1})
		window.SetScale(vmath.Vector3{800, 0, 1})

		e.controllerManager.AddController(ui.NewUiController(window))
		ui.LoadHTML(container, window, strings.NewReader(fileBrowserHtml), strings.NewReader(globalCss), e.uiAssets)

		e.fileBrowser = &FileBrowser{
			window:       window,
			assets:       e.uiAssets,
			callback:     callback,
			root:         ".",
			scrollOffset: 0,
			openFolders:  make(map[string]bool),
		}
	}
	e.fileBrowser.callback = callback
	e.fileBrowser.SetHeading(heading)
	e.fileBrowser.extensionFilter = filters
	e.gameEngine.AddOrtho(e.fileBrowser.window)
	e.fileBrowserOpen = true
	e.fileBrowser.UpdateFileSystem()
}

func (fb *FileBrowser) checkExtensionFilter(filename string) bool {
	for _, filter := range fb.extensionFilter {
		if strings.HasSuffix(filename, filter) {
			return true
		}
	}
	return false
}

func (fb *FileBrowser) UpdateFileSystem() {
	fb.ClearFiles()
	fileCounter := 0
	inClosedDir := false
	closedDepth := 0
	filepath.Walk(fb.root, func(path string, info os.FileInfo, err error) error {
		usePath := strings.Replace(path, "\\", "/", -1)
		depth := strings.Count(usePath, "/")
		if inClosedDir {
			if depth > closedDepth {
				return nil
			}
			inClosedDir = false
		}
		if !info.IsDir() && len(fb.extensionFilter) > 0 && !fb.checkExtensionFilter(info.Name()) {
			return nil
		}
		open, ok := fb.openFolders[path]
		isOpen := ok && open
		if !isOpen {
			inClosedDir = true
			closedDepth = depth
		}
		if err == nil && fileCounter >= fb.scrollOffset && fileCounter < fb.scrollOffset+maxFileDisplayed {
			if info.IsDir() {
				if isOpen {
					fb.RenderFile(info.Name(), path, "folderOpen", depth)
				} else {
					fb.RenderFile(info.Name(), path, "folderClosed", depth)
				}
			} else {
				fb.RenderFile(info.Name(), path, "file", depth)
			}
		}
		fileCounter++
		return nil
	})
}

func (fb *FileBrowser) ClearFiles() {
	elem := fb.window.ElementById("fileView")
	container, ok := elem.(*ui.Container)
	if ok {
		container.RemoveAllChildren()
	}
}

func (fb *FileBrowser) RenderFile(name, path, img string, depth int) {
	elem := fb.window.ElementById("fileView")
	container, ok := elem.(*ui.Container)
	if ok {
		onclickName := fmt.Sprintf("onClick_%v", path)
		fb.assets.AddCallback(onclickName, func(element ui.Element, args ...interface{}) {
			if len(args) >= 2 && !args[1].(bool) { // not on release
				fb.selectedFile = path
				open, openOk := fb.openFolders[path]
				fb.openFolders[path] = !openOk || !open
				fb.UpdateFileSystem()
			}
		})

		html := fmt.Sprintf("<div onclick=%v><img src=%v></img><p>%v</p></div>", onclickName, img, name)
		css := fmt.Sprintf(`
			p { font-size: 12px; width: 80%%; padding: 0 0 0 5px; }
			img { width: 16px; height: 16px; }
			div { padding: 0 0 0 %vpx; }
			div:hover { background-color: #00f2 }`, depth*8)
		if fb.selectedFile == path {
			css = fmt.Sprintf("%v div { background-color: #ff5 }", css)
		}

		ui.LoadHTML(container, fb.window, strings.NewReader(html), strings.NewReader(css), fb.assets)
	}
}

func (fb *FileBrowser) SetHeading(heading string) {
	elem := fb.window.ElementById("heading")
	container, ok := elem.(*ui.Container)
	if ok {
		container.RemoveAllChildren()
		html := fmt.Sprintf("<h1>%v</h1>", heading)
		css := "h1 { font-size: 16px }"
		ui.LoadHTML(container, fb.window, strings.NewReader(html), strings.NewReader(css), fb.assets)
	}
}
