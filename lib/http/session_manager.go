package http

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"sync"
	"time"

	"github.com/pjmd89/goutils/systemutils"
)

func newSessionManager(providerName ProviderKind, providerConf any, sessMaxLifeTime int64) (sm *sessionManager) {
	var providerObj sessionProvider
	switch providerName {
	case DATABASE_PROVIDER:
	case MEMORY_PROVIDER:
		providerObj = newMemoryProvider()
	case FILE_PROVIDER:
		providerObj = newFileProvider(providerConf.(string))
	default:
		panic("unknown provider")
	}

	sm = &sessionManager{
		lock:            sync.Mutex{},
		maxLifetime:     sessMaxLifeTime,
		sessionProvider: providerObj,
	}

	return
}

func (o *sessionManager) startSession(sessionName string, w http.ResponseWriter, r *http.Request, sessionData interface{}) (id string, err error) {
	o.lock.Lock()
	cookie, err := r.Cookie(sessionName)
	secureCookie := false
	if r.TLS != nil {
		secureCookie = true
	}

	if err == nil {
		id = cookie.Value
	}

	_, getErr := o.Get(id)
	if err != nil || getErr != nil {
		id = o.sessionID()
		cookie = &http.Cookie{
			Name:  sessionName,
			Value: id,
		}
		cookie.Path = "/"
		cookie.HttpOnly = true
		cookie.Secure = secureCookie
		cookie.SameSite = http.SameSiteNoneMode
		o.init(id, sessionData)
		http.SetCookie(w, cookie)
	} else {
		o.updateSessionAccess(id)
	}
	cookie.Expires = time.Now().Add(time.Second * time.Duration(o.maxLifetime))
	o.lock.Unlock()
	return
}

func (o *sessionManager) checkSession(sessionID string) {
	providerName := reflect.TypeOf(o.sessionProvider).Elem().Name()
	if providerName == "MemoryProvider" {
		sessVal, sessErr := o.Get(sessionID)
		if sessErr == nil && sessVal == nil {
			o.Destroy(sessionID)
		}
	}
	routineID := systemutils.GetRoutineID()
	o.lock.Lock()
	delete(o.routineSessions, routineID)
	o.lock.Unlock()
}

func (o *sessionManager) GetSessionByRoutine() (r any, err error) {
	routineID := systemutils.GetRoutineID()
	if sessID, isIn := o.routineSessions[routineID]; isIn {
		r, err = o.sessionProvider.Get(sessID)
	}

	return
}

func (o *sessionManager) setRoutineSession(sid string) {
	routineID := systemutils.GetRoutineID()
	o.lock.Lock()
	if o.routineSessions == nil {
		o.routineSessions = make(map[uint64]string)
	}

	o.routineSessions[routineID] = sid
	o.lock.Unlock()
}

func (o *sessionManager) sessionID() string {
	id := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, id); err != nil {
		return ""
	}
	return url.QueryEscape(base64.URLEncoding.EncodeToString(id))
}
