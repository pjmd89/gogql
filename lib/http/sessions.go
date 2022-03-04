package http

import (
	"net/http"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

var Session *Cookie;
var store = sessions.NewCookieStore(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32),
);
func SessionStart(w http.ResponseWriter,r *http.Request, sessionName *[]byte, cookieName string ) *Cookie{
	Session = &Cookie{};
	Session.r = r;
	Session.w = w;
	Session.cookieName = cookieName;
	Session.sessionName = *sessionName;
	if sessionName != nil{
		session, err := store.Get(r, cookieName)
		if err == nil {
			Session.Start = true;
			Session.session = session
		}
	}
	return Session
}
func(o *Cookie) New( values map[interface{}]interface{} ){
	o.Start = true;
	o.session = sessions.NewSession(store, o.cookieName);
	o.session.Values = values;
	o.session.Save(o.r, o.w);
}
func(o *Cookie)Set(values map[interface{}]interface{}){
	for k,v :=range values{
		o.session.Values[k] = v;
	}
	o.session.Save(o.r, o.w);
}
func(o *Cookie)Get() map[interface{}]interface{}{
	r := make(map[interface{}]interface{});
	if o.session != nil{
		r = o.session.Values
	}
	return r;
}