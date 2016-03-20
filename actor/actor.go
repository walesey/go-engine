package actor

// Actor - An actor is a wrapper for a game Entity that defines the behaviour
// of that entity based on controllers and game events,
// as well as other game integration points such as physics.
type Actor interface {
	Update(dt float64)
}

type ActorStore struct {
	actors []Actor
}

func NewActorStore() *ActorStore {
	return &ActorStore{
		actors: make([]Actor, 0, 0),
	}
}

func (store *ActorStore) UpdateAll(dt float64) {
	for _, actor := range store.actors {
		actor.Update(dt)
	}
}

func (store *ActorStore) Add(actors ...Actor) {
	store.actors = append(store.actors, actors...)
}

func (store *ActorStore) Remove(actors ...Actor) {
	for _, remove := range actors {
		for index, actor := range store.actors {
			if actor == remove {
				if index+1 == len(store.actors) {
					store.actors = store.actors[:index]
				} else {
					store.actors = append(store.actors[:index], store.actors[index+1:]...)
				}
				break
			}
		}
	}
}
