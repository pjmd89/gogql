package http

import (
	"net/http"

	"github.com/gorilla/mux"
)
type Gql interface{
	GetServerName() string
	GQLRender(w http.ResponseWriter,r *http.Request, session *Session) string
}
type pathConfig struct {
	Mode 			string 					`json:"mode,omitempty"`
	Endpoint 		string 					`json:"endpoint,omitempty"`
	Path 			string 					`json:"path,omitempty"`
	AllowOrigin 	string 					`json:"allowOrigin,omitempty"`
	FileDefault		string					`json:"default,omitempty"`
	RewriteTo		string					`json:"RewriteTo,omitempty"`
	Rewrite 		bool					`json:"rewrite,omitempty"`
	Url 			string 					
	httpsPort		string
	redirect 		bool
	enableHttps		bool	
	gqlRender		map[string]Gql
	serverName		string
	session 		*Session
}
type server struct {
	ServerName 		string	 				`json:"serverName,omitempty"`
	Path 	 		[]*pathConfig 			`json:"path,omitempty"`
	Cert	   		string 					`json:"cert,omitempty"`
	Key	   			string 					`json:"key,omitempty"`
	Redirect		bool					`json:"redirectToHttps,omitempty"`
	EnableHttps		bool					`json:"enableHttps,omitempty"`
	subrouter		*mux.Router
}
type Http struct{
	HttpPort		string					`json:"httpPort,omitempty"`
	HttpsPort		string					`json:"httpsPort,omitempty"`
	Server 			[]server 				`json:"server,omitempty"`
	Path			map[string]string
	HTTPService		*http.Server
	HTTPSService	*http.Server
	router 			*mux.Router
	gql				map[string]Gql
}