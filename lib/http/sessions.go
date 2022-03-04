package http

import (
	"fmt"
	"net/http"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)
type Session struct{
	sessionName []byte
	session *sessions.Session
	cookieName string
	w http.ResponseWriter
	r *http.Request
	Start bool
}
func SessionStart(w http.ResponseWriter,r *http.Request, sessionName *[]byte, cookieName string )(re *Session){
	re = &Session{};
	re.r = r;
	re.w = w;
	re.cookieName = cookieName;
	re.sessionName = *sessionName;
	fmt.Println("KEY ",re.sessionName);
	if sessionName != nil{
		store = sessions.NewCookieStore(*sessionName)
		session, err := store.Get(r, cookieName)
		if err == nil {
			re.Start = true;
			re.session = session
		}
	}
	return re
}
func(o *Session) New( values map[interface{}]interface{} ){
	o.sessionName = securecookie.GenerateRandomKey(32);
	o.Start = true;
	fmt.Println("KEY ",o.sessionName);
	store = sessions.NewCookieStore(o.sessionName)
	o.session = sessions.NewSession(store, o.cookieName);
	o.session.Values = values;
}
func(o *Session)Set(values map[interface{}]interface{}){
	for k,v :=range values{
		o.session.Values[k] = v;
	}
	o.session.Save(o.r, o.w);
}
func(o *Session)Get() map[interface{}]interface{}{
	return o.session.Values;
}