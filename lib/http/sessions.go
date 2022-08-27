package http

import (
	"net/http"
	"sync"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)
var Session *Cookie;
var mutex = &sync.Mutex{}
var store = sessions.NewCookieStore(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32),
);
func SessionStart(w http.ResponseWriter,r *http.Request, sessionName *[]byte, cookieName string ) *Cookie{
	mutex.Lock()
	Session = &Cookie{};
	Session.r = r;
	Session.w = w;
	Session.cookieName = cookieName;
	Session.sessionName = *sessionName;
	if sessionName != nil{
		secureCookie := false;
		if r.TLS != nil {
			secureCookie = true;
		}
		store.Options = &sessions.Options{ Path: "/", HttpOnly: true, Secure: secureCookie, MaxAge: 0,SameSite: http.SameSiteNoneMode};
		session, err := store.Get(r, cookieName)
		Session.Start = true;
		Session.session = session
		
		if session.IsNew || err == nil{
			Session.session = sessions.NewSession(store, Session.cookieName);
			Session.session.Options = store.Options
			Session.session.Values = make(map[interface{}]interface{});
			Session.session.Save(Session.r, Session.w);
		}
	}
	mutex.Unlock()
	return Session
}
func(o *Cookie)Set(values map[interface{}]interface{}){
	mutex.Lock()
	for k,v :=range values{
		o.session.Values[k] = v;
	}
	o.session.Save(o.r, o.w);
	mutex.Unlock()
}
func(o *Cookie)Get() map[interface{}]interface{}{
	r := make(map[interface{}]interface{});
	if o != nil && o.session != nil{
		r = o.session.Values
	}
	return r;
}