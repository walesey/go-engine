package engine

type Updatable interface {
	Update(dt float64)
}

// ActorStore - is a collection of actors.
type UpdatableStore struct {
	updatables []Updatable
}

func NewUpdatableStore() *UpdatableStore {
	return &UpdatableStore{
		updatables: make([]Updatable, 0, 64),
	}
}

func (store *UpdatableStore) UpdateAll(dt float64) {
	for _, updatable := range store.updatables {
		updatable.Update(dt)
	}
}

func (store *UpdatableStore) Add(updatables ...Updatable) {
	store.updatables = append(store.updatables, updatables...)
}

func (store *UpdatableStore) Remove(updatables ...Updatable) {
	for _, remove := range updatables {
		for index, updatable := range store.updatables {
			if updatable == remove {
				if index+1 == len(store.updatables) {
					store.updatables = store.updatables[:index]
				} else {
					store.updatables = append(store.updatables[:index], store.updatables[index+1:]...)
				}
				break
			}
		}
	}
}
