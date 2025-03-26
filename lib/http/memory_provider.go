package http

import (
	"fmt"
	"sync"
)

type memoryProvider struct {
	sessions map[string]interface{}
	lock     sync.Mutex
}

func newMemoryProvider() sessionProvider {
	return &memoryProvider{
		sessions: map[string]interface{}{},
		lock:     sync.Mutex{},
	}
}

func (o *memoryProvider) init(sessionID string, sessionData interface{}) (err error) {
	o.lock.Lock()
	o.sessions[sessionID] = sessionData
	o.lock.Unlock()
	return
}

func (o *memoryProvider) Get(sessionID string) (r interface{}, err error) {
	if _, ok := o.sessions[sessionID]; ok {
		r = o.sessions[sessionID]
	} else {
		err = fmt.Errorf("session is no seted")
	}
	return
}

func (o *memoryProvider) Count() (r int, err error) {
	return len(o.sessions), nil
}
func (o *memoryProvider) Set(sessionID string, sessionData interface{}) (err error) {
	o.lock.Lock()
	defer o.lock.Unlock()
	if _, ok := o.sessions[sessionID]; ok {
		o.sessions[sessionID] = sessionData
	} else {
		err = fmt.Errorf("session is no seted")
	}
	return
}

func (o *memoryProvider) Destroy(sessionID string) (err error) {
	o.lock.Lock()
	defer o.lock.Unlock()
	if _, ok := o.sessions[sessionID]; !ok {
		err = fmt.Errorf("session is no seted")
		return
	}

	delete(o.sessions, sessionID)
	return
}

func (o *memoryProvider) garbageCollector(sessMaxLifeTime int64) {

}
func (o *memoryProvider) updateSessionAccess(sessionID string) {

}
