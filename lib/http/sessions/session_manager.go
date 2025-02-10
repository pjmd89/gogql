package sessions

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

const (
	DATABASE_PROVIDER = "database"
	FILE_PROVIDER     = "file"
	MEMORY_PROVIDER   = "memory"
)

type SessionProvider interface {
	Init(sessionID string, sessionData interface{}) (err error)
	Set(sessionID string, sessionData interface{}) (err error)
	Get(sessionID string, dataReceiver any) (r interface{}, err error)
	Destroy(sessionID string) (err error)
	Count() (r int, err error)
}

type SessionManager struct {
	lock        sync.Mutex
	maxLifetime int
	SessionProvider
	routineSessions map[uint64]string
}

func NewSessionManager(providerName string, providerConf any, sessMaxLifeTime int) (sm *SessionManager) {
	var providerObj SessionProvider
	switch providerName {
	case DATABASE_PROVIDER:
	case MEMORY_PROVIDER:
		providerObj = newMemoryProvider()
	case FILE_PROVIDER:
		providerObj = newFileProvider(providerConf.(string))
	default:
		panic("unknown provider")
	}

	sm = &SessionManager{
		lock:            sync.Mutex{},
		maxLifetime:     sessMaxLifeTime,
		SessionProvider: providerObj,
	}

	return
}

func (o *SessionManager) StartSession(sessionName string, w http.ResponseWriter, r *http.Request, sessionData interface{}) (id string, err error) {
	o.lock.Lock()
	cookie, err := r.Cookie(sessionName)
	secureCookie := false
	if r.TLS != nil {
		secureCookie = true
	}

	if err == nil {
		id = cookie.Value
	}

	_, getErr := o.Get(id, nil)
	if err != nil || cookie == nil || getErr != nil {
		id = o.sessionID()
		cookie = &http.Cookie{
			Name:  sessionName,
			Value: id,
		}
		cookie.Path = "/"
		cookie.HttpOnly = true
		cookie.Expires = time.Now().AddDate(2, 0, 0)
		cookie.Secure = secureCookie
		cookie.SameSite = http.SameSiteNoneMode
		o.Init(id, sessionData)
		http.SetCookie(w, cookie)
	}
	o.lock.Unlock()
	return
}

func (o *SessionManager) CheckSession(sessionID string) {
	providerName := reflect.TypeOf(o.SessionProvider).Elem().Name()
	if providerName == "MemoryProvider" {
		sessVal, sessErr := o.Get(sessionID, nil)
		if sessErr == nil && sessVal == nil {
			o.Destroy(sessionID)
		}
	}
	routineID := systemutils.GetRoutineID()
	delete(o.routineSessions, routineID)
}

func (o *SessionManager) GetSessionByRoutine(dataReceiver any) (r interface{}, err error) {
	routineID := systemutils.GetRoutineID()
	if sessID, isIn := o.routineSessions[routineID]; isIn {
		r, err = o.SessionProvider.Get(sessID, dataReceiver)
	}
	return
}

func (o *SessionManager) SetRoutineSession(sid string) {
	routineID := systemutils.GetRoutineID()
	if o.routineSessions == nil {
		o.routineSessions = make(map[uint64]string)
	}
	o.routineSessions[routineID] = sid
}

func (o *SessionManager) sessionID() string {
	id := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, id); err != nil {
		return ""
	}
	return url.QueryEscape(base64.URLEncoding.EncodeToString(id))
}
