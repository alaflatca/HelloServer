package measure

import "sync"

var measureCache map[string]*Measure = make(map[string]*Measure)
var measureMtx sync.Mutex

func Get(key string) *Measure {
	measureMtx.Lock()
	defer measureMtx.Unlock()

	val, ok := measureCache[key]
	if ok {
		return val
	}
	return nil
}

func Set(key string, val *Measure) {
	measureMtx.Lock()
	defer measureMtx.Unlock()

	measureCache[key] = val
}
