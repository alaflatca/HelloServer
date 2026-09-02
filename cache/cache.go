package cache

import (
	"helloServer/measure"
	"sync"
)

var measureCache map[string]measure.Measure = make(map[string]measure.Measure)
var measureMutex sync.RWMutex

var anyCache map[string]interface{} = make(map[string]interface{})
var anyMutex sync.Mutex

func Get(key string) *measure.Measure {
	measureMutex.RLock()
	defer measureMutex.RUnlock()

	val, ok := measureCache[key]
	if ok {
		snapshot := val
		return &snapshot
	}
	return nil
}

func Set(key string, val *measure.Measure) {
	if val == nil {
		return
	}

	measureMutex.Lock()
	defer measureMutex.Unlock()

	measureCache[key] = *val
}
