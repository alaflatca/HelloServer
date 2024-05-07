package measure

import (
	"fmt"
)

func init() {
	evtBus = &eventBus{
		events: make(map[string]eventHandler),
	}
}

type eventHandler func(interface{})

type eventBus struct {
	events map[string]eventHandler
}

var evtBus *eventBus

func Subscribe(key string, e eventHandler) {
	evtBus.events[key] = e
}

func Publish(key string, data interface{}) error {
	e, ok := evtBus.events[key]
	if !ok {
		return fmt.Errorf("invalid event (key: '%s')", key)
	}

	go e(data)

	return nil
}
