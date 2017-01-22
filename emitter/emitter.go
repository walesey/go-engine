package emitter

import "sync"

type Event interface{}

type EventEmitter interface {
	On(topic string, handler func(Event)) chan Event
	Off(topic string)
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
func (e *Emitter) On(topic string, handlers ...func(Event)) <-chan Event {
	e.mux.Lock()
	l := listener{
		handlers: handlers,
		ch:       make(chan Event, e.channelSize),
	}
	if topicListeners, ok := e.topics[topic]; ok {
		e.topics[topic] = append(topicListeners, l)
	} else {
		e.topics[topic] = listeners{l}
	}
	e.mux.Unlock()
	return l.ch
}

func (e *Emitter) Off(topic string) {
	e.mux.Lock()
	if topicListeners, ok := e.topics[topic]; ok {
		for _, l := range topicListeners {
			close(l.ch)
		}
	}
	delete(e.topics, topic)
	e.mux.Unlock()
}

// Emit - writes an event to the given topic
func (e *Emitter) Emit(topic string, event Event) {
	e.mux.Lock()
	if topicListeners, ok := e.topics[topic]; ok {
		for _, l := range topicListeners {
			l.ch <- event
		}
	}
	e.mux.Unlock()
}

func (e *Emitter) Close() {
	e.mux.Lock()
	for _, topicListeners := range e.topics {
		for _, l := range topicListeners {
			close(l.ch)
		}
	}
	e.topics = make(map[string]listeners)
	e.mux.Unlock()
}

// FlushAll - flushes all events from all channels and calls handlers for them
func (e *Emitter) FlushAll() {
	for topic, _ := range e.topics {
		e.Flush(topic)
	}
}

// Flush - flushes all events from the topic channels and calls handlers for them
func (e *Emitter) Flush(topic string) {
	if ll, ok := e.topics[topic]; ok {
		for _, l := range ll {
			l.flush()
		}
	}
}

func (l listener) flush() {
Loop:
	for {
		select {
		case event := <-l.ch:
			for _, h := range l.handlers {
				h(event)
			}
		default:
			break Loop
		}
	}
}
