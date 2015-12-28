package dynamics

import vmath "github.com/walesey/go-engine/vectormath"

type Event struct {
	Name string
	Data interface{}
}

func CollisionEvent(globalContact vmath.Vector3) Event {
	return Event{
		Name: "collision",
		Data: map[string]interface{}{
			"globalContact": globalContact,
		},
	}
}
