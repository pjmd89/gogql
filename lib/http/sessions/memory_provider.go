package sessions

import (
	"fmt"
	"sync"
)

type MemoryProvider struct {
	sessions map[string]interface{}
	lock     sync.Mutex
}

func newMemoryProvider() SessionProvider {
	return &MemoryProvider{
		sessions: map[string]interface{}{},
	}
}

func (o *MemoryProvider) Init(sessionID string, sessionData interface{}) (err error) {
	o.lock.Lock()
	o.sessions[sessionID] = sessionData
	o.lock.Unlock()
	return
}

func (o *MemoryProvider) Get(sessionID string, dataReceiver any) (r interface{}, err error) {
	if _, ok := o.sessions[sessionID]; ok {
		r = o.sessions[sessionID]
	} else {
		err = fmt.Errorf("session is no seted")
	}
	return
}

func (o *MemoryProvider) Count() (r int, err error) {
	return len(o.sessions), nil
}
func (o *MemoryProvider) Set(sessionID string, sessionData interface{}) (err error) {
	o.lock.Lock()
	defer o.lock.Unlock()
	if _, ok := o.sessions[sessionID]; ok {
		o.sessions[sessionID] = sessionData
	} else {
		err = fmt.Errorf("session is no seted")
	}
	return
}

func (o *MemoryProvider) Destroy(sessionID string) (err error) {
	o.lock.Lock()
	defer o.lock.Unlock()
	if _, ok := o.sessions[sessionID]; !ok {
		err = fmt.Errorf("session is no seted")
		return
	}

	delete(o.sessions, sessionID)
	return
}
