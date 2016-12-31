package engine

func Timeout(e Engine, time float64, fn func()) Updatable {
	timer := time
	var updater Updatable
	updater = UpdatableFunc(func(dt float64) {
		if timer -= dt; timer <= 0 {
			fn()
			e.RemoveUpdatable(updater)
		}
	})
	e.AddUpdatable(updater)
	return updater
}

func Interval(e Engine, time float64, fn func()) Updatable {
	timer := 0.0
	var updater Updatable
	updater = UpdatableFunc(func(dt float64) {
		for timer += dt; timer >= time; timer -= time {
			fn()
		}
	})
	e.AddUpdatable(updater)
	return updater
}
