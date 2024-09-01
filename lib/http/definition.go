package http

import (
	"embed"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

type Gql interface {
	GQLRender(w http.ResponseWriter, r *http.Request, sessionID string) (isErr bool)
	GQLRenderSubscription(mt int, message []byte, socketId string, sessionID string)
}
type Rest interface {
	RestRender(w http.ResponseWriter, r *http.Request, sessionID string) (isErr bool)
}
type URL struct {
	Scheme     string
	Port       string
	Host       string
	RequestURI string
	TLS        bool
	Method     string
	URL        string
	Referer    string
	Origin     struct {
		Scheme string
		Host   string
		Port   string
		URL    string
	}
}
type Path struct {
	Mode        string   `json:"mode,omitempty"`
	Path        string   `json:"path,omitempty"`
	Endpoint    string   `json:"endpoint,omitempty"`
	Rewrite     bool     `json:"rewrite,omitempty"`
	RewriteTo   string   `json:"rewriteTo,omitempty"`
	Redirect    Redirect `json:"redirect,omitempty"`
	FileDefault string   `json:"fileDefault,omitempty"`
	len         int
	pathURL     string
	host        string
	origin      string
}
type Redirect struct {
	From string `json:"from,omitempty"`
	To   string `json:"to,omitempty"`
}
type server struct {
	Host            string   `json:"host,omitempty"`
	Reject          []string `json:"reject,omitempty"`
	Cert            string   `json:"cert,omitempty"`
	Key             string   `json:"key,omitempty"`
	LetsEncrypt     bool     `json:"letsEncrypt,omitempty"`
	RedirectToHttps bool     `json:"redirectToHttps,omitempty"`
	EnableHttps     bool     `json:"enableHttps,omitempty"`
	Path            []Path   `json:"path,omitempty"`
}
type Http struct {
	HttpPort     string   `json:"httpPort,omitempty"`
	HttpsPort    string   `json:"httpsPort,omitempty"`
	CookieName   string   `json:"cookieName,omitempty"`
	Server       []server `json:"server,omitempty"`
	embed        *embed.FS
	httpService  mux.Router
	httpsService mux.Router
	router       *mux.Router
	gql          Gql
	rest         Rest
	CheckOrigin  func(url URL) (bool, interface{})
	OnBegin      func(url URL, httpPath *Path, originData interface{}, uid string) bool
	OnFinish     func(isErr bool, uid string)
	OnSession    func() (r interface{})
	originData   any
	http404      func()
	http405      func()
}
type SessionManager struct {
	lock        sync.Mutex
	maxLifetime int
	sessions    map[string]interface{}
}
