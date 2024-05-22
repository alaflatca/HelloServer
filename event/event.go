package event

import (
	"fmt"
	"sync"
)

func init() {
	evtBus = &eventBus{
		events: make(map[string]eventHandler),
	}
	evtMutex = &sync.Mutex{}
}

type eventHandler func(interface{})

type eventBus struct {
	events map[string]eventHandler
}

var evtBus *eventBus
var evtMutex *sync.Mutex

func Subscribe(key string, e eventHandler) {
	evtMutex.Lock()
	evtBus.events[key] = e
	evtMutex.Unlock()
}

func Publish(key string, data interface{}) error {
	evtMutex.Lock()
	defer evtMutex.Unlock()
	e, ok := evtBus.events[key]
	if !ok {
		return fmt.Errorf("invalid event (key: '%s')", key)
	}

	go e(data)

	return nil
}
