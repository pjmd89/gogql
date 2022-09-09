package http

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func (o *SessionManager) Init(sessionName string, sessionLifetime int, w http.ResponseWriter, r *http.Request, sessionData interface{}) (id string, err error) {
	o.lock.Lock()
	cookie, err := r.Cookie(sessionName)
	secureCookie := false
	if err == nil {
		id = cookie.Value
	}

	if r.TLS != nil {
		secureCookie = true
	}
	if o.sessions == nil {
		o.sessions = make(map[string]interface{})
	}
	if cookie != nil && o.sessions[cookie.Value] != nil {
		fmt.Printf("session %v:", o.sessions[cookie.Value])
	}

	if err != nil {
		id = o.sessionID()
		cookie = &http.Cookie{
			Name:  sessionName,
			Value: id,
		}
		cookie.Path = "/"
		cookie.HttpOnly = true
		cookie.MaxAge = sessionLifetime
		cookie.Secure = secureCookie
		cookie.SameSite = http.SameSiteNoneMode
		o.sessions[cookie.Value] = sessionData
		http.SetCookie(w, cookie)
	}
	o.lock.Unlock()

	return
}
func (o *SessionManager) Get(sessionID string) (r interface{}) {
	if o.sessions[sessionID] != nil {
		r = o.sessions[sessionID]
	}
	return
}
func (o *SessionManager) Set(sessionID string, sessionData interface{}) {
	o.sessions[sessionID] = sessionData
}
func (o *SessionManager) sessionID() string {
	id := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, id); err != nil {
		return ""
	}
	return url.QueryEscape(base64.URLEncoding.EncodeToString(id))
}
