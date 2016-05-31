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

type UpdatableStore struct {
	updatables map[string]Updatable
}

func NewUpdatableStore() *UpdatableStore {
	return &UpdatableStore{
		updatables: make(map[string]Updatable),
	}
}

func (store *UpdatableStore) UpdateAll(dt float64) {
	for _, updatable := range store.updatables {
		updatable.Update(dt)
	}
}

func (store *UpdatableStore) Add(key string, updatable Updatable) {
	store.updatables[key] = updatable
}

func (store *UpdatableStore) Remove(key string) {
	delete(store.updatables, key)
}

func (store *UpdatableStore) RemoveUpdatable(updatable Updatable) {
	for key, u := range store.updatables {
		if updatable == u {
			delete(store.updatables, key)
		}
	}
}
