package measure

import "sync"

var measureSync map[string]*Measure = make(map[string]*Measure)
var measureMtx sync.Mutex

func Get(key string) *Measure {
	measureMtx.Lock()
	defer measureMtx.Unlock()

	val, ok := measureSync[key]
	if ok {
		return val
	}
	return &Measure{}
}

func Set(key string, val *Measure) {
	measureMtx.Lock()
	defer measureMtx.Unlock()

	measureSync[key] = val
}
