package http

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)
type Gql interface{
	GetServerName() string
	GQLRender(w http.ResponseWriter,r *http.Request) string
	GQLRenderSubscription(mt int,message []byte, socketId string)
}

type Cookie struct{
	sessionName 	[]byte
	session 		*sessions.Session
	cookieName 		string
	w 				http.ResponseWriter
	r 				*http.Request
	Start 			bool
}
type URL struct{
	Scheme			string
	Port			string
	Host			string
	RequestURI		string
	TLS				bool
	Method			string
	URL				string
	Referer			string
	Origin	 		struct{
		Scheme		string
		Host		string
		Port		string
		URL			string
	}
}
type Path struct{
	Mode			string					`json:"mode,omitempty"`
	Path			string					`json:"path,omitempty"`
	Endpoint		string					`json:"endpoint,omitempty"`
	len				int
	pathURL			string
	host			string
	origin			string
}
type server struct{
	Host			string					`json:"host,omitempty"`
	Reject			[]string				`json:"reject,omitempty"`
	Cert			string					`json:"cert,omitempty"`
	Key				string					`json:"key,omitempty"`
	LetsEncrypt 	bool					`json:"letsEncrypt,omitempty"`
    RedirectToHttps	bool					`json:"redirectToHttps,omitempty"`
    EnableHttps		bool					`json:"enableHttps,omitempty"`
	RewriteTo		string					`json:"rewriteTo,omitempty"`
	Rewrite			bool					`json:"rewrite,omitempty"`
	FileDefault		string					`json:"fileDefault,omitempty"`
	Path			[]Path					`json:"path,omitempty"`
}
type Http struct{
	HttpPort		string					`json:"httpPort,omitempty"`
	HttpsPort		string					`json:"httpsPort,omitempty"`
	CookieName		string					`json:"cookieName,omitempty"`
	Server 			[]server 				`json:"server,omitempty"`
	httpService		mux.Router
	httpsService	mux.Router
	router 			*mux.Router
	gql				map[string]Gql
	CheckOrigin		func( url URL ) bool
	OnBegin			func( url URL, httpPath *Path ) bool
	OnFinish		func()
	http404			func()
	http405			func()
}