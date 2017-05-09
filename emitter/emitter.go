package emitter

import "sync"

type Event interface{}

type EventChan <-chan Event

type EventEmitter interface {
	On(topic string, handlers ...func(Event)) EventChan
	Off(topic string, channels ...EventChan)
	Listeners(topic string) []EventChan
	Emit(topic string, event Event)
	Close()
	FlushAll()
	Flush(topic string)
}

type listener struct {
	handlers []func(Event)
	ch       chan Event
}

type listeners []listener

type Emitter struct {
	topics      map[string]listeners
	channelSize int
	mux         *sync.Mutex
}

func New(channelSize int) *Emitter {
	return &Emitter{
		topics:      make(map[string]listeners),
		channelSize: channelSize,
		mux:         &sync.Mutex{},
	}
}

// On - creates a new topic listener with optional handlers and returns the
func (e *Emitter) On(topic string, handlers ...func(Event)) EventChan {
	l := listener{
		handlers: handlers,
		ch:       make(chan Event, e.channelSize),
	}
	if topicListeners, ok := e.getListeners(topic); ok {
		e.setListeners(topic, append(topicListeners, l))
	} else {
		e.setListeners(topic, listeners{l})
	}
	return l.ch
}

func (e *Emitter) Off(topic string, channels ...EventChan) {
	if topicListeners, ok := e.getListeners(topic); ok {
		if len(channels) > 0 {
			for _, ch := range channels {
				for i := len(topicListeners) - 1; i >= 0; i-- {
					if topicListeners[i].ch == ch {
						close(topicListeners[i].ch)
						topicListeners = append(topicListeners[:i], topicListeners[i+1:]...)
					}
				}
			}
			e.setListeners(topic, topicListeners)
		} else {
			for _, l := range topicListeners {
				close(l.ch)
			}
			e.deleteListeners(topic)
		}
	}
}

func (e *Emitter) Listeners(topic string) []EventChan {
	if topicListeners, ok := e.getListeners(topic); ok {
		listeners := make([]EventChan, len(topicListeners))
		for i, l := range topicListeners {
			listeners[i] = l.ch
		}
		return listeners
	}
	return []EventChan{}
}

// Emit - writes an event to the given topic
func (e *Emitter) Emit(topic string, event Event) {
	e.mux.Lock()
	defer e.mux.Unlock()
	if topicListeners, ok := e.topics[topic]; ok {
		for _, l := range topicListeners {
			l.ch <- event
		}
	}
}

// Close - closes all channels and removes all topics
func (e *Emitter) Close() {
	e.mux.Lock()
	defer e.mux.Unlock()
	for _, topicListeners := range e.topics {
		for _, l := range topicListeners {
			close(l.ch)
		}
	}
	e.topics = make(map[string]listeners)
}

// FlushAll - flushes all events from all channels and calls handlers for them
func (e *Emitter) FlushAll() {
	e.mux.Lock()
	defer e.mux.Unlock()
	for topic := range e.topics {
		e.mux.Unlock()
		e.Flush(topic)
		e.mux.Lock()
	}
}

// Flush - flushes all events from the topic channels and calls handlers for them
func (e *Emitter) Flush(topic string) {
	topicListeners, ok := e.getListeners(topic)
	if ok {
		for _, l := range topicListeners {
			l.flush()
		}
	}
}

func (l listener) flush() {
Loop:
	for {
		select {
		case event := <-l.ch:
			if event == nil {
				break Loop
			}
			for _, h := range l.handlers {
				h(event)
			}
		default:
			break Loop
		}
	}
}

func (e *Emitter) getListeners(topic string) (l listeners, ok bool) {
	e.mux.Lock()
	defer e.mux.Unlock()
	l, ok = e.topics[topic]
	return
}

func (e *Emitter) setListeners(topic string, l listeners) {
	e.mux.Lock()
	defer e.mux.Unlock()
	e.topics[topic] = l
}

func (e *Emitter) deleteListeners(topic string) {
	e.mux.Lock()
	defer e.mux.Unlock()
	delete(e.topics, topic)
}
