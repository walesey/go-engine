package engine

type Updatable interface {
	Update(dt float64)
}

type updatableImpl struct {
	updateFunc func(dt float64)
}

func (updatable *updatableImpl) Update(dt float64) {
	updatable.updateFunc(dt)
}

func UpdatableFunc(update func(dt float64)) Updatable {
	return &updatableImpl{update}
}

func UpdatableFanIn(updatables ...Updatable) Updatable {
	return UpdatableFunc(func(dt float64) {
		for _, u := range updatables {
			u.Update(dt)
		}
	})
}

type UpdatableStore struct {
	updatables []Updatable
}

func NewUpdatableStore() *UpdatableStore {
	return &UpdatableStore{
		updatables: make([]Updatable, 0),
	}
}

func (store *UpdatableStore) UpdateAll(dt float64) {
	for _, updatable := range store.updatables {
		if updatable != nil {
			updatable.Update(dt)
		}
	}
}

func (store *UpdatableStore) Add(updatable Updatable) {
	store.updatables = append(store.updatables, updatable)
}

func (store *UpdatableStore) Remove(updatable Updatable) {
	for i, u := range store.updatables {
		if updatable == u {
			store.updatables[i] = store.updatables[len(store.updatables)-1]
			store.updatables[len(store.updatables)-1] = nil
			store.updatables = store.updatables[:len(store.updatables)-1]
			break
		}
	}
}
