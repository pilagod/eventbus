package eventbus

import (
	"reflect"

	"github.com/panjf2000/ants/v2"
)

// EventBus event bus
type EventBus interface {
	EventPublisher
	EventSubscriber
	Release()
}

// EventPublisher event publisher
type EventPublisher interface {
	Publish(events ...Event) error
}

// EventSubscriber event subscriber
type EventSubscriber interface {
	Subscribe(event Event, handlers ...EventHandler)
	SubscribeAll(handlers ...EventHandler)
	Use(decorator EventHandlerDecorator)
}

// EventHandler event handler
type EventHandler interface {
	Handle(e Event) error
}

// EventHandlerDecorator decorator for event handler
type EventHandlerDecorator func(h EventHandler) EventHandler

type eventBus struct {
	pool       *ants.Pool
	handlers   map[string][]EventHandler
	decorators []EventHandlerDecorator
}

// EventBus

func (eb *eventBus) Release() {
	eb.pool.Release()
}

// EventPublisher

func (eb *eventBus) Publish(events ...Event) error {
	for _, event := range events {
		e := event
		if reflect.ValueOf(event).Kind() == reflect.Ptr {
			e = reflect.ValueOf(event).Elem().Interface()
		}
		for _, handler := range eb.getHandlers(e) {
			h := handler
			if err := eb.pool.Submit(func() { h.Handle(e) }); err != nil {
				return err
			}
		}
	}
	return nil
}

// EventSubscriber

func (eb *eventBus) Subscribe(event Event, handlers ...EventHandler) {
	var hs []EventHandler
	for _, h := range handlers {
		hs = append(hs, eb.applyDecorators(h))
	}
	eb.subscribe(getEventName(event), hs...)
}

func (eb *eventBus) SubscribeAll(handlers ...EventHandler) {
	eb.subscribe("*", handlers...)
}

func (eb *eventBus) Use(decorator EventHandlerDecorator) {
	eb.decorators = append(eb.decorators, nil)
	copy(eb.decorators[1:], eb.decorators)
	eb.decorators[0] = decorator
}

// util

func (eb *eventBus) getHandlers(event Event) []EventHandler {
	return append(
		eb.handlers[getEventName(event)],
		eb.handlers["*"]...,
	)
}

func (eb *eventBus) applyDecorators(handler EventHandler) EventHandler {
	result := handler
	for _, d := range eb.decorators {
		result = d(result)
	}
	return result
}

func (eb *eventBus) subscribe(eventName string, handlers ...EventHandler) {
	eb.handlers[eventName] = append(eb.handlers[eventName], handlers...)
}

func getEventName(event Event) string {
	t := reflect.TypeOf(event)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t.Name()
}
