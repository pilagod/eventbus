package eventbus

import (
	"github.com/panjf2000/ants/v2"
)

var eb *eventBus = nil

// Setup setup event bus
func Setup(workerPoolSize int) (EventBus, error) {
	pool, err := ants.NewPool(workerPoolSize)
	if err != nil {
		return nil, err
	}
	eb = &eventBus{
		pool:     pool,
		handlers: make(map[string][]EventHandler),
	}
	return eb, nil
}

// GetEventBus returns event bus
func GetEventBus() EventBus {
	if eb == nil {
		panic("event bus is not setup yet")
	}
	return eb
}

// GetEventPublisher returns event publisher
func GetEventPublisher() EventPublisher {
	return GetEventBus()
}

// GetEventSubscriber returns event subscriber
func GetEventSubscriber() EventSubscriber {
	return GetEventBus()
}
