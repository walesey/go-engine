package controller

type Controller interface {
	SetKeyAction(function func(key Key, action Action))
	BindKeyAction(function func(), key Key, action Action)
	SetMouseAction(function func(button MouseButton, action Action))
	BindMouseAction(function func(), button MouseButton, action Action)
	BindAxisAction(function func(xpos, ypos float32))
	BindScrollAction(function func(xoffset, yoffset float32))
}
