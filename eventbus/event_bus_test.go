package eventbus

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestEventBus(t *testing.T) {
	eb, err := Setup(1)
	if err != nil {
		panic(err.Error())
	}
	defer eb.Release()
	suite.Run(t, &eventBusSuite{})
}

type eventBusSuite struct {
	suite.Suite
}

func (s *eventBusSuite) SetupTest() {
	eb := GetEventBus().(*eventBus)
	eb.decorators = nil
	eb.handlers = make(map[string][]EventHandler)
}

// TODO: multiple events for multiple handlers

func (s *eventBusSuite) TestSubscribe() {
	var wg sync.WaitGroup

	es := GetEventSubscriber()

	ha := newEventHandler(&wg)
	es.Subscribe(eventA{}, ha)

	hb := newEventHandler(&wg)
	es.Subscribe(eventB{}, hb)

	wg.Add(2)

	ep := GetEventPublisher()

	ea := eventA{Message: "Hello"}
	ep.Publish(ea)

	eb := eventB{Message: "World"}
	ep.Publish(eb)

	wg.Wait()

	s.Len(ha.Events, 1)
	s.Contains(ha.Events, ea)

	s.Len(hb.Events, 1)
	s.Contains(hb.Events, eb)
}

func (s *eventBusSuite) TestSubscribePtr() {
	var wg sync.WaitGroup

	es := GetEventSubscriber()

	ha := newEventHandler(&wg)
	es.Subscribe(&eventA{}, ha)

	hb := newEventHandler(&wg)
	es.Subscribe(&eventB{}, hb)

	wg.Add(2)

	ep := GetEventPublisher()

	ea := eventA{Message: "Hello"}
	ep.Publish(&ea)

	eb := eventB{Message: "World"}
	ep.Publish(&eb)

	wg.Wait()

	// handler should only get value of event, not pointer of event

	s.Len(ha.Events, 1)
	s.Contains(ha.Events, ea)

	s.Len(hb.Events, 1)
	s.Contains(hb.Events, eb)
}

func (s *eventBusSuite) TestSubscribeAll() {
	var wg sync.WaitGroup

	es := GetEventSubscriber()

	h := newEventHandler(&wg)
	es.SubscribeAll(h)

	wg.Add(2)

	ep := GetEventPublisher()

	ea := eventA{Message: "Hello"}
	ep.Publish(ea)

	eb := eventB{Message: "World"}
	ep.Publish(eb)

	wg.Wait()

	s.Len(h.Events, 2)
	s.Contains(h.Events, ea)
	s.Contains(h.Events, eb)
}

type eventDecorator struct {
	handler EventHandler
}

func (d *eventDecorator) Handle(e Event) error {
	ev := e.(event)
	ev.Message = "Decorated " + ev.Message
	return d.handler.Handle(ev)
}

func (s *eventBusSuite) TestUse() {
	var wg sync.WaitGroup

	es := GetEventSubscriber()

	es.Use(func(h EventHandler) EventHandler {
		return &eventDecorator{handler: h}
	})

	h := newEventHandler(&wg)
	es.Subscribe(event{}, h)

	wg.Add(1)

	ep := GetEventPublisher()

	e := event{Message: "Hello"}
	ep.Publish(e)

	wg.Wait()

	s.Len(h.Events, 1)
	got, _ := h.Events[0].(event)
	s.Contains(got.Message, "Decorated")
}

// event for test

type event struct {
	Message string
}

type eventA struct {
	Message string
}

type eventB struct {
	Message string
}

// event handler for test

func newEventHandler(wg *sync.WaitGroup) *eventHandler {
	return &eventHandler{wg: wg}
}

type eventHandler struct {
	wg     *sync.WaitGroup
	Events []Event
}

func (h *eventHandler) Handle(event Event) error {
	h.Events = append(h.Events, event)
	h.wg.Done()
	return nil
}
