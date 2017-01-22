package emitter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type eventChecker struct {
	event Event
}

func (e *eventChecker) handler(event Event) {
	e.event = event
}

func TestSyncronousEvents(t *testing.T) {
	e := New(5)
	helloEvent := "helloEvent"
	otherEvent := "otherEvent"
	checker := &eventChecker{}
	otherChecker := &eventChecker{}
	e.On("hello", checker.handler)
	e.On("other", otherChecker.handler)

	e.Emit("hello", helloEvent)
	e.FlushAll()
	assert.EqualValues(t, checker.event, helloEvent)
	assert.NotEqual(t, otherChecker.event, helloEvent)

	e.Emit("other", otherEvent)
	e.FlushAll()
	assert.NotEqual(t, checker.event, otherEvent)
	assert.EqualValues(t, otherChecker.event, otherEvent)
}

func TestASyncronousEvents(t *testing.T) {
	e := New(5)
	helloEvent := "helloEvent"
	otherEvent := "otherEvent"
	helloChan := e.On("hello")
	otherChan := e.On("other")

	go e.Emit("hello", helloEvent)
	result := <-helloChan
	assert.EqualValues(t, result, helloEvent)

	go e.Emit("other", otherEvent)
	result = <-otherChan
	assert.EqualValues(t, result, otherEvent)
}
