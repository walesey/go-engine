package ui

import (
	"fmt"

	"image/color"

	"github.com/go-gl/mathgl/mgl32"
)

const nbBars = 20

func NewProgressBar(label string) *Window {
	window := NewWindow()
	window.SetTranslation(mgl32.Vec3{500, 200, 0})
	window.SetScale(mgl32.Vec3{325, 0, 1})

	container := NewContainer()
	window.SetElement(container)
	container.SetBackgroundColor(200, 200, 200, 255)
	container.SetWidth(325)

	labelElem := NewTextElement(label, color.RGBA{10, 10, 20, 255}, 12, nil)
	labelContainer := NewContainer()
	labelContainer.SetWidth(100)
	labelContainer.UsePercentWidth(true)
	labelContainer.SetMargin(NewMargin(5))
	labelContainer.AddChildren(labelElem)

	progressBar := NewContainer()
	progressBar.SetBackgroundColor(50, 50, 50, 255)
	progressBar.SetMargin(NewMargin(5))
	progressBar.SetPadding(Margin{0, 0, 0, 5})
	progressBar.SetHeight(40)
	progressBar.SetWidth(315)

	container.AddChildren(labelContainer, progressBar)

	for i := 1; i <= nbBars; i++ {
		box := NewContainer()
		progressBar.AddChildren(box)
		box.SetId(fmt.Sprintf("progress%v", i))
		box.SetBackgroundColor(50, 220, 80, 255)
		box.SetMargin(Margin{10, 0, 10, 5})
		box.SetHeight(20)
		box.SetWidth(10)
	}

	window.Render()

	return window
}

func SetProgressBar(pb *Window, progress int) {
	for i := 1; i <= nbBars; i++ {
		container, ok := pb.ElementById(fmt.Sprintf("progress%v", i)).(*Container)
		if ok {
			if i > progress {
				container.SetBackgroundColor(0, 0, 0, 0)
			} else {
				container.SetBackgroundColor(0, 255, 0, 255)
			}
		}
	}
}
