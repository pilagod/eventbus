# go-eventbus [![Build Status](https://travis-ci.com/pilagod/go-eventbus.svg?branch=master)](https://travis-ci.com/pilagod/go-eventbus) [![Coverage Status](https://coveralls.io/repos/github/pilagod/go-eventbus/badge.svg?branch=master)](https://coveralls.io/github/pilagod/go-eventbus?branch=master)

Event bus for Go, leveraging [ants](https://github.com/panjf2000/ants) for worker pool management.

## Installation

```shell
$ go get -u github.com/pilagod/go-eventbus
```

## Usage

You should import `eventbus` module under `go-eventbus`:

```go
import "github.com/pilagod/go-eventbus/eventbus"
```

### Initialization

Event bus must be setup beforehead:

```go
var workerPoolSize = 4

func main() {
    // ...

    eb, err := eventbus.Setup(workerPoolSize)
    if err != nil {
        panic(err.Error())
    }
    // don't forget to release workers
    defer eb.Release()

    // ...
}
```

### Event Subscriber

```go
// Event

type Event struct {
    Message string
}

// Event handler

type EventHandler struct {}

func (h *EventHandler) Handle(event eventbus.Event) error {
    e, ok := event.(Event)
    if !ok {
        return nil
    }
    fmt.Println(e.Message)
}

es := eventbus.GetEventSubscriber() // GetEventSubscriber will panic if event bus hasn't setup

// Subscribe handler to specific event
es.Subscribe(Event{}, &EventHandler{})

// Subscribe handler to all events
es.SubscribeAll(&EventHandler{})
```

### Event Publisher

```go
// Event

type Event struct {
    Message string
}

ep := eventbus.GetEventPublisher() // GetEventPublisher will panic if event bus hasn't setup

// Publish specific event
e := Event{
    Message: "Hello World",
}
ep.Publish(e)

// Publish multiple events
es := []eventbus.Event{
   Event{Message: "A"}, 
   Event{Message: "B"}, 
   Event{Message: "C"}, 
}
ep.Publish(es...)

```

## License

© Cyan Ho (pilagod), 2021-NOW

Released under the [MIT License](https://github.com/pilagod/go-eventbus/blob/master/LICENSE)
