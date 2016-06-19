package ui

import (
	"encoding/hex"
	"fmt"
	"image/color"
	"io"
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	"github.com/aymerick/douceur/css"
	"github.com/aymerick/douceur/parser"
	vmath "github.com/walesey/go-engine/vectormath"
	"golang.org/x/net/html"
)

// LoadHTML - load the html/css code into the container
func LoadHTML(container *Container, htmlInput, cssInput io.Reader, assets HtmlAssets) ([]Activatable, error) {
	document, err := html.Parse(htmlInput)
	if err != nil {
		log.Printf("Error parsing html: %v", err)
		return []Activatable{}, err
	}

	css, err := ioutil.ReadAll(cssInput)
	if err != nil {
		log.Printf("Error reading css: %v", err)
		return []Activatable{}, err
	}

	styles, err := parser.Parse(string(css))
	if err != nil {
		log.Printf("Error parsing css: %v", err)
		return []Activatable{}, err
	}

	activatables := renderNode(container, document.FirstChild, styles, assets)

	return activatables, nil
}

func renderNode(container *Container, node *html.Node, styles *css.Stylesheet, assets HtmlAssets) []Activatable {
	activatables := []Activatable{}
	nextNode := node
	for nextNode != nil {
		if nextNode.Type == 1 {
			// create a text node
			text := nextNode.Data
			text = strings.TrimSpace(text)
			if len(text) > 0 {
				createText(text, nextNode.Parent, container, styles, assets)
			}
		} else {
			// Create a container
			newContainer := NewContainer()
			newContainer.id = getAttribute(nextNode, "id")
			container.AddChildren(newContainer)

			//Parse other html tag types
			var textField *TextElement
			var imageElement *ImageElement
			tagType := nextNode.DataAtom.String()
			switch {
			case tagType == "input":
				inputType := getAttribute(nextNode, "type")
				if inputType == "text" || inputType == "password" {
					textField = createText("", nextNode, newContainer, styles, assets)
					textField.SetHidden(inputType == "password")
					activatables = append(activatables, textField)
					newContainer.Hitbox.AddOnClick(func(button int, release bool, position vmath.Vector2) {
						if !release {
							textField.Activate()
						}
					})
				}
			case tagType == "img":
				imgSrc := getAttribute(nextNode, "src")
				img, ok := assets.imageMap[imgSrc]
				if ok {
					imageElement = NewImageElement(img)
					newContainer.AddChildren(imageElement)
				}
			}

			//Parse Styles
			normalStyles := getStyles(styles, nextNode, "")
			applyStyles(newContainer, normalStyles, assets)
			hoverStyles := getStyles(styles, nextNode, ":hover")
			activeStyles := getStyles(styles, nextNode, ":active")
			hover := false
			active := false
			updateImage := func() {
				if imageElement != nil {
					imageElement.UsePercentWidth(newContainer.percentWidth)
					if newContainer.percentWidth {
						imageElement.SetWidth(100)
					} else {
						imageElement.SetWidth(newContainer.width)
					}
					imageElement.UsePercentHeight(newContainer.percentHeight)
					if newContainer.percentHeight {
						imageElement.SetHeight(100)
					} else {
						imageElement.SetHeight(newContainer.height)
					}
				}
			}
			updateImage()
			updateState := func() {
				applyDefaultStyles(newContainer)
				applyStyles(newContainer, normalStyles, assets)
				if hover {
					applyStyles(newContainer, hoverStyles, assets)
				}
				if active {
					applyStyles(newContainer, activeStyles, assets)
				}
				updateImage()
			}
			if len(hoverStyles) > 0 {
				newContainer.Hitbox.AddOnHover(func() {
					hover = true
					updateState()
				})
				newContainer.Hitbox.AddOnUnHover(func() {
					hover = false
					updateState()
				})
			}
			if len(activeStyles) > 0 {
				newContainer.Hitbox.AddOnClick(func(button int, release bool, position vmath.Vector2) {
					active = !release
					updateState()
				})
			}

			//Parse html Props
			for _, attr := range nextNode.Attr {
				switch {
				case attr.Key == "onclick":
					callback, ok := assets.callbackMap[attr.Val]
					if ok {
						newContainer.Hitbox.AddOnClick(func(button int, release bool, position vmath.Vector2) {
							callback(newContainer, button, release, position)
						})
					}
				case attr.Key == "onhover":
					callback, ok := assets.callbackMap[attr.Val]
					if ok {
						newContainer.Hitbox.AddOnHover(func() {
							callback(newContainer)
						})
					}
				case attr.Key == "onfocus":
					callback, ok := assets.callbackMap[attr.Val]
					if ok && textField != nil {
						textField.AddOnFocus(func() {
							callback(newContainer)
						})
					}
				case attr.Key == "onblur":
					callback, ok := assets.callbackMap[attr.Val]
					if ok && textField != nil {
						textField.AddOnBlur(func() {
							callback(newContainer)
						})
					}
				case attr.Key == "onkeypress":
					callback, ok := assets.callbackMap[attr.Val]
					if ok && textField != nil {
						textField.AddOnKeyPress(func(key string, release bool) {
							callback(newContainer, key, release)
						})
					}
				case attr.Key == "placeholder":
					if textField != nil {
						textField.SetPlaceholder(attr.Val)
					}
				}
			}

			//Render children
			activatables = append(activatables, renderNode(newContainer, nextNode.FirstChild, styles, assets)...)
		}
		if nextNode == nextNode.NextSibling {
			break
		}
		nextNode = nextNode.NextSibling
	}
	return activatables
}

func applyDefaultStyles(container *Container) {
	container.SetBackgroundColor(0, 0, 0, 0)
	container.SetHeight(0)
	container.SetWidth(0)
	container.SetMargin(NewMargin(0))
	container.SetPadding(NewMargin(0))
}

func applyStyles(container *Container, styles map[string]string, assets HtmlAssets) {
	for prop, value := range styles {
		switch {
		case prop == "padding":
			paddings, units := parseDimensions(value)
			if len(paddings) == 1 {
				container.SetPadding(NewMargin(paddings[0]))
				container.SetPaddingPercent(NewMarginPercentages(len(units) == 1 && units[0] == "%"))
			} else if len(paddings) == 4 && len(units) == 4 {
				container.SetPadding(Margin{paddings[0], paddings[1], paddings[2], paddings[3]})
				container.SetPaddingPercent(MarginPercentages{units[0] == "%", units[1] == "%", units[2] == "%", units[3] == "%"})
			}
		case prop == "margin":
			margins, units := parseDimensions(value)
			if len(margins) == 1 {
				container.SetMargin(NewMargin(margins[0]))
				container.SetMarginPercent(NewMarginPercentages(len(units) == 1 && units[0] == "%"))
			} else if len(margins) == 4 && len(units) == 4 {
				container.SetMargin(Margin{margins[0], margins[1], margins[2], margins[3]})
				container.SetMarginPercent(MarginPercentages{units[0] == "%", units[1] == "%", units[2] == "%", units[3] == "%"})
			}
		case prop == "background-color":
			color := parseColor(value)
			container.SetBackgroundColor(color[0], color[1], color[2], color[3])
		case prop == "background-image":
			img, ok := assets.imageMap[value]
			if ok {
				container.SetBackgroundImage(img)
			}
		case prop == "width":
			width, units := parseDimensions(value)
			if len(width) == 1 {
				container.SetWidth(width[0])
				container.UsePercentWidth(len(units) == 1 && units[0] == "%")
			}
		case prop == "height":
			height, units := parseDimensions(value)
			if len(height) == 1 {
				container.SetHeight(height[0])
				container.UsePercentHeight(len(units) == 1 && units[0] == "%")
			}
		}
	}
}

func createText(text string, node *html.Node, container *Container, styles *css.Stylesheet, assets HtmlAssets) *TextElement {
	textElement := NewTextElement(text, color.Black, 16, assets.fontMap["default"])
	textElement.SetAlign(LEFT_ALIGN)
	textStyles := getStyles(styles, node, "")
	applyTextStyles(textElement, textStyles, assets)
	hoverTextStyles := getStyles(styles, node, ":hover")
	activeTextStyles := getStyles(styles, node, ":active")
	hover := false
	active := false
	updateState := func() {
		applyDefaultTextStyles(textElement, assets)
		applyTextStyles(textElement, textStyles, assets)
		if hover {
			applyTextStyles(textElement, hoverTextStyles, assets)
		}
		if active {
			applyTextStyles(textElement, activeTextStyles, assets)
		}
		textElement.ReRender()
	}
	if len(hoverTextStyles) > 0 {
		container.Hitbox.AddOnHover(func() {
			hover = true
			updateState()
		})
		container.Hitbox.AddOnUnHover(func() {
			hover = false
			updateState()
		})
	}
	if len(activeTextStyles) > 0 {
		container.Hitbox.AddOnClick(func(button int, release bool, position vmath.Vector2) {
			active = !release
			updateState()
		})
	}
	container.AddChildren(textElement)
	return textElement
}

func applyDefaultTextStyles(textField *TextElement, assets HtmlAssets) {
	textField.SetTextColor(color.Black)
	textField.SetTextSize(16)
	textField.SetFont(assets.fontMap["default"])
	textField.SetAlign(CENTER_ALIGN)
}

func applyTextStyles(textField *TextElement, textStyles map[string]string, assets HtmlAssets) {
	for prop, value := range textStyles {
		switch {
		case prop == "color":
			c := parseColor(value)
			textField.SetTextColor(color.RGBA{c[0], c[1], c[2], c[3]})
		case prop == "font-size":
			size, _ := parseDimensions(value)
			if len(size) == 1 {
				textField.SetTextSize(size[0])
			}
		case prop == "font-family":
			fontStyle, ok := assets.fontMap[value]
			if ok {
				textField.SetFont(fontStyle)
			}
		case prop == "text-align":
			if value == "center" {
				// TODO:
				textField.SetAlign(CENTER_ALIGN)
			}
		}
	}
}

func parseDimensions(dimensionsStr string) (values []float64, units []string) {
	dimensions := strings.Fields(dimensionsStr)
	values = make([]float64, len(dimensions))
	units = make([]string, len(dimensions))
	for i, dimension := range dimensions {
		var err error
		if strings.HasSuffix(dimension, "px") {
			values[i], err = strconv.ParseFloat(strings.Replace(dimension, "px", "", 1), 64)
			units[i] = "px"
		} else if strings.HasSuffix(dimension, "%") {
			values[i], err = strconv.ParseFloat(strings.Replace(dimension, "%", "", 1), 64)
			units[i] = "%"
		} else {
			values[i], err = strconv.ParseFloat(dimension, 64)
			units[i] = "px"
		}
		if err != nil {
			log.Printf("Error parsing dimensions: %v;\n", err)
		}
	}
	return
}

func parseColor(colorStr string) [4]uint8 {
	data := []byte(colorStr)
	r, g, b, a := "0", "0", "0", "ff"
	if len(data) == 4 {
		r, g, b = string(data[1:2]), string(data[2:3]), string(data[3:])
		r, g, b = fmt.Sprintf("%v0", r), fmt.Sprintf("%v0", g), fmt.Sprintf("%v0", b)
	} else if len(data) == 5 {
		r, g, b, a = string(data[1:2]), string(data[2:3]), string(data[3:4]), string(data[4:])
		r, g, b, a = fmt.Sprintf("%v0", r), fmt.Sprintf("%v0", g), fmt.Sprintf("%v0", b), fmt.Sprintf("%v0", a)
	} else if len(data) == 7 {
		r, g, b = string(data[1:3]), string(data[3:5]), string(data[5:])
	} else if len(data) == 9 {
		r, g, b, a = string(data[1:3]), string(data[3:5]), string(data[5:7]), string(data[7:])
	}
	rd, _ := hex.DecodeString(r)
	gd, _ := hex.DecodeString(g)
	bd, _ := hex.DecodeString(b)
	ad, _ := hex.DecodeString(a)
	result := [4]uint8{0, 0, 0, 0}
	if len(rd) > 0 {
		result[0] = uint8(rd[0])
	}
	if len(gd) > 0 {
		result[1] = uint8(gd[0])
	}
	if len(bd) > 0 {
		result[2] = uint8(bd[0])
	}
	if len(ad) > 0 {
		result[3] = uint8(ad[0])
	}
	return result
}

func getStyles(styles *css.Stylesheet, node *html.Node, modifier string) map[string]string {
	hierarchy := []*html.Node{node}
	parent := node.Parent
	for parent != nil {
		hierarchy = append(hierarchy, parent)
		parent = parent.Parent
	}

	// css styles
	rules := make(map[string]string)
	for _, rule := range styles.Rules {
		for _, sel := range rule.Selectors {
			if selectorMatch(sel, modifier, hierarchy) {
				for _, declaration := range rule.Declarations {
					rules[declaration.Property] = declaration.Value
				}
			}
		}
	}

	return rules
}

func selectorMatch(sel, modifier string, hierarchy []*html.Node) bool {
	selectors := strings.Fields(sel)
	// reverse the order
	for i, j := 0, len(selectors)-1; i < j; i, j = i+1, j-1 {
		selectors[i], selectors[j] = selectors[j], selectors[i]
	}

	if len(selectors) == 0 || len(hierarchy) == 0 {
		return false
	}

	firstSelector := selectors[0]
	firstNode := hierarchy[0]
	if strings.HasPrefix(firstSelector, "#") {
		if fmt.Sprintf("#%v%v", getAttribute(firstNode, "id"), modifier) != firstSelector {
			return false
		}
	} else if strings.HasPrefix(firstSelector, ".") {
		if !matchClasses(firstNode, firstSelector, modifier) {
			return false
		}
	} else if fmt.Sprintf("%v%v", firstNode.DataAtom.String(), modifier) != firstSelector {
		return false
	}

	i := 1
	for j := 1; j < len(hierarchy) && i < len(selectors); j++ {
		selector := selectors[i]
		nextNode := hierarchy[j]
		if strings.HasPrefix(selector, "#") {
			if fmt.Sprintf("#%v", getAttribute(nextNode, "id")) == selector {
				i++
				continue
			}
		} else if strings.HasPrefix(selector, ".") {
			if matchClasses(nextNode, selector, "") {
				i++
				continue
			}
		} else if nextNode.DataAtom.String() == selector {
			i++
			continue
		}
	}
	return i == len(selectors)
}

func matchClasses(node *html.Node, selector, modifier string) bool {
	nodeSelectors := strings.Fields(getAttribute(node, "class"))
	for _, nodeSelector := range nodeSelectors {
		if fmt.Sprintf(".%v%v", nodeSelector, modifier) == selector {
			return true
		}
	}
	return false
}

func getAttribute(node *html.Node, key string) string {
	for _, attr := range node.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}
