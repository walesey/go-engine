package controller

var defaultConstructor func() Controller

func CreateController() Controller {
	if defaultConstructor != nil {
		return defaultConstructor()
	}
	return nil
}

func SetDefaultConstructor(constructor func() Controller) {
	defaultConstructor = constructor
}
