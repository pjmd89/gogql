package http

import (
	"crypto/tls"
	"embed"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/pjmd89/gogql/lib/http/sessions"
	"github.com/pjmd89/goutils/jsonutils"
	"github.com/pjmd89/goutils/systemutils"
	"golang.org/x/exp/slices"
)

var upgrader = websocket.Upgrader{
	EnableCompression: true,
}

var (
	WsIds          map[string]chan bool       = make(map[string]chan bool)
	WsChannels     map[string]*websocket.Conn = make(map[string]*websocket.Conn)
	SessionManager                            = sessions.NewSessionManager(sessions.FILE_PROVIDER, "/var/tmp/gogql", 17520)
	logs           systemutils.Logs
)

func Init(sytemlogs systemutils.Logs, configFile string) *Http {
	logs = sytemlogs
	o := &Http{HttpPort: "80", HttpsPort: "443"}

	jsonutils.GetJson(configFile, &o)

	if strings.Trim(o.HttpPort, " ") == "" {
		o.HttpPort = "80"
	}
	if strings.Trim(o.HttpsPort, " ") == "" {
		o.HttpsPort = "443"
	}
	if strings.Trim(o.CookieName, " ") == "" {
		o.CookieName = "GOGQL_SESSION"
	}

	return o
}

func (o *Http) SetRest(rest Rest) *Http {
	o.rest = rest
	return o
}
func (o *Http) SetGql(gql Gql) *Http {
	o.gql = gql
	return o
}
func (o *Http) FrontEmbed(filesystem embed.FS) {
	o.embed = &filesystem
}
func (o *Http) Start() {
	channel := make(chan bool)
	stop := false
	isTls := false
	tlsConfig := &tls.Config{}
	tlsConfig.Certificates = []tls.Certificate{}
	/*
		xhttp:=&Http{Url:v.ServerName}
		o.router[counter].NotFoundHandler = xhttp;
		o.router[counter].MethodNotAllowedHandler = xhttp;
	*/
	o.router = mux.NewRouter()
	o.router.Use(handlers.CompressHandler)
	subrouter := o.router.MatcherFunc(o.MatcherFunc).Subrouter()
	subrouter.Methods("POST", "GET", "OPTIONS").Handler(o)

	hTTPService := &http.Server{Addr: ":" + o.HttpPort, Handler: o.router}
	hTTPSService := &http.Server{Addr: ":" + o.HttpsPort, Handler: o.router, TLSConfig: tlsConfig}

	for _, server := range o.Server {
		if server.EnableHttps {
			isTls = true
		}
		if server.EnableHttps {
			tmp, certErr := tls.LoadX509KeyPair(server.Cert, server.Key)
			if certErr != nil {
				logs.System.Error().Println(certErr.Error())
				logs.System.Fatal().Fatal("Error on certificate. " + server.Cert)
				stop = true
			}
			tlsConfig.Certificates = append(tlsConfig.Certificates, tmp)
		}
	}
	logs.System.Info().Println("http server start")
	go o.listenHttp(channel, hTTPService)
	if isTls {
		go o.listenHttps(channel, hTTPSService)
	}
	tlsConfig.BuildNameToCertificate()
	for !stop {
		x := <-channel
		if x == true {
			stop = true
			hTTPService.Shutdown(nil)
		}
	}
}
func (o *Http) MatcherFunc(hr *http.Request, hm *mux.RouteMatch) (r bool) {
	r = true
	var url URL
	url.Split(hr)
	if o.CheckOrigin != nil {
		r, o.originData = o.CheckOrigin(url)
	}
	logs.Access.Info().Printf("Access to %s is %v", url.URL, r)
	return r
}
func (o *Http) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var urlInfo URL
	urlInfo.Split(r)
	urlPath := "/"
	var httpPathMode *Path
	var sessionData interface{}
	if o.OnSession != nil {
		sessionData = o.OnSession()
	}
	sessionID, err := SessionManager.StartSession(o.CookieName, w, r, sessionData)
	SessionManager.SetRoutineSession(sessionID)
	if err != nil {
		logs.System.Error().Println(err.Error())
	}
	if urlInfo.Origin.URL != "" {
		w.Header().Set("Access-Control-Allow-Origin", urlInfo.Origin.URL)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
	}
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	//w.Header().Set("Access-Control-Max-Age", "86400");
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	parseUrl, _ := url.Parse(urlInfo.URL + urlInfo.RequestURI)

	if parseUrl.Path != "/" {
		urlPath = strings.TrimSuffix(parseUrl.Path, "/")
	}

	for i, svr := range o.Server {
		if strings.Trim(svr.Host, " ") == "*" {
			tmp := make([]server, 0)
			tmp = append(tmp, o.Server[:i]...)
			tmp = append(tmp, o.Server[i+1:]...)
			o.Server = append(tmp, svr)
		}
	}
	for _, server := range o.Server {
		var httpModes []Path
		host := strings.Replace(server.Host, "*", `[^.]+`, -1)
		if strings.Trim(server.Host, " ") == "*" {
			host = `[^/]+`
		}
		match, _ := regexp.MatchString(host, urlInfo.Host)
		serverBreak := false
		if match && !slices.Contains(server.Reject, urlInfo.Host) {
			for _, serverPath := range server.Path {
				matchPath, _ := regexp.MatchString(`^`+serverPath.Endpoint, urlPath)
				if matchPath {
					serverBreak = true
					if server.RedirectToHttps && urlInfo.Scheme == "http" {
						reservedPort := ""
						if o.HttpsPort != "443" {
							reservedPort = ":" + o.HttpsPort
						}
						http.Redirect(w, r, "https://"+urlInfo.Host+reservedPort+urlInfo.RequestURI, http.StatusSeeOther)
						return
					}
					serverPath.len = len(serverPath.Endpoint)
					httpModes = append(httpModes, serverPath)
				}
			}
		}
		if serverBreak {
			mode := 0
			if len(httpModes) > 0 {
				for i, httpMode := range httpModes {
					if httpMode.len > httpModes[mode].len {
						mode = i
					}
				}
			}
			if len(httpModes) >= mode {
				httpPathMode = &Path{}
				*httpPathMode = httpModes[mode]
				httpPathMode.pathURL = urlPath
				httpPathMode.host = server.Host
				httpPathMode.origin = urlInfo.Origin.URL
			}
			break
		}
	}
	uID := uuid.New().String()
	isErr := false
	if httpPathMode != nil {
		if o.OnBegin != nil {
			o.OnBegin(urlInfo, httpPathMode, o.originData, uID)
		}
		switch httpPathMode.Mode {
		case "file", "embed":
			isErr = o.fileServeHTTP(w, r, httpPathMode, sessionID)
		case "gql":
			isErr = o.gqlServeHTTP(w, r, httpPathMode, sessionID)
		case "websocket":
			isErr = o.websocketServeHTTP(w, r, httpPathMode, sessionID)
		case "rest":
			isErr = o.restServeHTTP(w, r, httpPathMode, sessionID)
		}
	} else {
		logs.System.Error().Println("mode not found")
		w.WriteHeader(http.StatusNotFound)
	}
	if o.OnFinish != nil {
		o.OnFinish(isErr, uID)
	}

	SessionManager.CheckSession(sessionID)

}
func (o *Http) fileServeHTTP(w http.ResponseWriter, r *http.Request, httpPath *Path, sessionID string) (isErr bool) {
	var fileStat fs.FileInfo
	isErr = false
	if o.OnSession != nil {
		//sessionData = o.OnSession()
	}
	//o.setSessionIndex(sessionID)

	if strings.Trim(httpPath.Redirect.From, " ") != "" && httpPath.Redirect.From == httpPath.pathURL {
		http.Redirect(w, r, httpPath.Redirect.To, http.StatusSeeOther)
		return
	}
	if strings.Trim(httpPath.Redirect.From, " ") == "" && strings.Trim(httpPath.Redirect.To, " ") != "" {
		http.Redirect(w, r, httpPath.Redirect.To, http.StatusSeeOther)
		return
	}
	filePath := filepath.Join(httpPath.Path, httpPath.pathURL)
	file, fErr := o.fileOpen(filePath, httpPath.Mode)
	if fErr == nil {
		fileStat, _ = file.Stat()
	}

	if fileStat == nil || fileStat.IsDir() {
		filePath = filepath.Join(httpPath.Path, httpPath.pathURL, "/", httpPath.FileDefault)
		file, fErr = o.fileOpen(filePath, httpPath.Mode)
		if fErr != nil && httpPath.Rewrite {
			file, fErr = o.fileOpen(httpPath.Path+httpPath.RewriteTo, httpPath.Mode)
		}
	}

	if fErr != nil {
		logs.System.Error().Printf("el archivo %s no se encuentra.", filePath)
		fmt.Fprint(w, "file not found, archivo no se encuentra")
		isErr = true
		w.WriteHeader(http.StatusNotFound)
		return
	}
	http.ServeContent(w, r, httpPath.Path+"/"+r.RequestURI, time.Time{}, file.(io.ReadSeeker))
	return
}
func (o *Http) fileOpen(name string, mode string) (fs.File, error) {
	if o.embed != nil && mode == "embed" {
		return o.embed.Open(name)
	}
	return os.Open(name)
}
func (o *Http) gqlServeHTTP(w http.ResponseWriter, r *http.Request, httpPath *Path, sessionID string) (isErr bool) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	isErr = o.gql.GQLRender(w, r, sessionID)
	return
}
func (o *Http) restServeHTTP(w http.ResponseWriter, r *http.Request, httpPath *Path, sessionID string) (isErr bool) {
	isErr = o.rest.RestRender(w, r, sessionID)
	return
}
func (o *Http) websocketServeHTTP(w http.ResponseWriter, r *http.Request, httpPath *Path, sessionID string) (isErr bool) {
	headers := http.Header{}
	isErr = false
	headers.Set("Sec-WebSocket-Protocol", "graphql-transport-ws")
	headers.Set("Sec-WebSocket-Version", "13")
	headers.Set("Content-Type", "application/json; charset=UTF-8")
	headers.Set("Access-Control-Allow-Credentials", "true")
	headers.Set("Access-Control-Allow-Origin", httpPath.origin)

	id := uuid.New().String()
	WsIds[id] = make(chan bool, 1)
	var upgraderError error
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	WsChannels[id], upgraderError = upgrader.Upgrade(w, r, headers)
	if upgraderError != nil {
		select {
		case WsIds[id] <- false:
		}
		close(WsIds[id])
		delete(WsIds, id)
		delete(WsChannels, id)
		fmt.Println(len(WsChannels))
	}
	defer WsChannels[id].Close()
	while := true
	for while {
		mt, message, err := WsChannels[id].ReadMessage()
		if err != nil {
			select {
			case WsIds[id] <- false:
			}
			close(WsIds[id])
			delete(WsIds, id)
			delete(WsChannels, id)
			while = false
			isErr = true
			break
		}
		go o.WebSocketMessage(mt, message, id, httpPath, sessionID)
	}
	return
}
func (o *Http) listenHttp(channel chan bool, handler *http.Server) {
	channel <- false
	err := handler.ListenAndServe()
	if err != nil {
		logs.System.Error().Println("http server start error: " + err.Error())
		channel <- true
	}
}
func (o *Http) listenHttps(channel chan bool, handler *http.Server) {
	channel <- false
	listener, tlsErr := tls.Listen("tcp", handler.Addr, handler.TLSConfig)
	if tlsErr != nil {
		logs.System.Error().Println("http server start error: " + tlsErr.Error())
		channel <- true
	}

	err := handler.Serve(listener)
	if err != nil {
		logs.System.Error().Println("http server start error: " + err.Error())
		channel <- true
	}
}
func (o *Http) WebSocketMessage(mt int, message []byte, id string, httpPath *Path, sessionID string) {

	//o.gql[httpPath.host].GQLRenderSubscription(mt, message, id, sessionID)
	o.gql.GQLRenderSubscription(mt, message, id, sessionID)
}

func WriteWebsocketMessage(mt int, id string, message []byte) {
	if WsChannels[id] != nil {
		WsChannels[id].WriteMessage(mt, message)
	}
}
func contains(s []interface{}, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
