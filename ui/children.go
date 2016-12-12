package ui

type Children []Element

func (c Children) GetChild(index int) Element {
	if index >= len(c) {
		return nil
	}
	return c[index]
}

func (c Children) GetChildById(id string) Element {
	for _, child := range c {
		if child.GetId() == id {
			return child
		}
		if elem := child.GetChildren().GetChildById(id); elem != nil {
			return elem
		}
	}
	return nil
}

func (c Children) TextElementById(id string) *TextElement {
	elem := c.GetChildById(id)
	if elem != nil && len(elem.GetChildren()) > 0 {
		if textElement, ok := elem.GetChildren().GetChild(0).(*TextElement); ok {
			return textElement
		}
		if textField, ok := elem.GetChildren().GetChild(0).(*TextField); ok {
			return textField.text
		}
	}
	return nil
}
