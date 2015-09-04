package controller

import ()

type Controller struct {
	ActionMap map[string]func()
}

func (c Controller) BindAction(keyCode string, action func()) {
	c.ActionMap[keyCode] = action
}

func (c Controller) TriggerAction(keyCode string) {
	c.ActionMap[keyCode]()
}