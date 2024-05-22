package cache

import (
	"helloServer/measure"
	"sync"
)

var measureCache map[string]*measure.Measure = make(map[string]*measure.Measure)
var measureMutex sync.Mutex

var anyCache map[string]interface{} = make(map[string]interface{})
var anyMutex sync.Mutex

func Get(key string) *measure.Measure {
	measureMutex.Lock()
	defer measureMutex.Unlock()

	val, ok := measureCache[key]
	if ok {
		return val
	}
	return nil
}

func Set(key string, val *measure.Measure) {
	measureMutex.Lock()
	defer measureMutex.Unlock()

	measureCache[key] = val
}
