package http

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/pjmd89/goutils/jsonutils"
	"github.com/pjmd89/goutils/systemutils"
	"golang.org/x/exp/slices"
)

var upgrader = websocket.Upgrader{
	EnableCompression: true,
}
var sessionIndex = map[uint64]string{}
var WsIds map[string]chan bool = make(map[string]chan bool)
var WsChannels map[string]*websocket.Conn = make(map[string]*websocket.Conn)
var Session = SessionManager{}
var logs systemutils.Logs

func Init(sytemlogs systemutils.Logs, configFile string) *Http {
	logs = sytemlogs
	o := &Http{HttpPort: "8080", HttpsPort: "8443"}

	jsonutils.GetJson(configFile, &o)

	if strings.Trim(o.HttpPort, " ") == "" {
		o.HttpPort = "8080"
	}
	if strings.Trim(o.HttpsPort, " ") == "" {
		o.HttpsPort = "8443"
	}
	if strings.Trim(o.CookieName, " ") == "" {
		o.CookieName = "GOGQL_SESSION"
	}
	return o
}
func (o *Http) SetRest(rest Rest) *Http {
	/*
		mapRest := make(map[string]Rest)
		for _, v := range rest {
			mapRest[v.GetServerName()] = v
		}
	*/
	o.rest = rest
	return o
}
func (o *Http) SetGql(gql Gql) *Http {
	/*
		mapGql := make(map[string]Gql)
		for _, v := range gql {
			mapGql[v.GetServerName()] = v
		}
	*/
	o.gql = gql
	return o
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
	sessionID, err := Session.Init(o.CookieName, 0, w, r, sessionData)
	o.setSessionIndex(sessionID)
	defer o.sessionDestroy()
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
	if httpPathMode != nil {
		if o.OnBegin != nil {
			o.OnBegin(urlInfo, httpPathMode, o.originData)
		}
		switch httpPathMode.Mode {
		case "file":
			o.fileServeHTTP(w, r, httpPathMode, sessionID)
		case "gql":
			o.gqlServeHTTP(w, r, httpPathMode, sessionID)
		case "websocket":
			o.websocketServeHTTP(w, r, httpPathMode, sessionID)
		case "rest":
			o.restServeHTTP(w, r, httpPathMode, sessionID)
		}
	} else {
		logs.System.Error().Println("mode not found")
		w.WriteHeader(http.StatusNotFound)
	}
	if o.OnFinish != nil {
		o.OnFinish()
	}
	return
}
func (o *Http) fileServeHTTP(w http.ResponseWriter, r *http.Request, httpPath *Path, sessionID string) {
	path := httpPath.Path
	var filePath = path + httpPath.pathURL
	if strings.Trim(httpPath.Redirect.From, " ") != "" && httpPath.Redirect.From == httpPath.pathURL {
		http.Redirect(w, r, httpPath.Redirect.To, http.StatusSeeOther)
		return
	}
	if strings.Trim(httpPath.Redirect.From, " ") == "" && strings.Trim(httpPath.Redirect.To, " ") != "" {
		http.Redirect(w, r, httpPath.Redirect.To, http.StatusSeeOther)
		return
	}
	if strings.Trim(path, " ") == "" {
		path = "."
	}
	file, fErr := os.Open(filePath)
	fileStat, _ := file.Stat()

	if fileStat == nil || fileStat.IsDir() {
		filePath = httpPath.Path + httpPath.pathURL + "/" + httpPath.FileDefault
		file, fErr = os.Open(filePath)
		fileStat, _ = file.Stat()
		if fErr != nil && httpPath.Rewrite {
			file, fErr = os.Open(httpPath.Path + httpPath.RewriteTo)
			fileStat, _ = file.Stat()
		}
	}

	if fErr != nil {
		logs.System.Error().Printf("el archivo %s no se encuentra.", filePath)
		fmt.Fprint(w, "file not found, archivo no se encuentra")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	http.ServeContent(w, r, httpPath.Path+"/"+r.RequestURI, time.Time{}, file)
}
func (o *Http) gqlServeHTTP(w http.ResponseWriter, r *http.Request, httpPath *Path, sessionID string) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	//por favor, revisa que o.serverName exista, si no existe entonces devuelvele un dedito
	/*
		if o.gql[httpPath.host] != nil {
			rx := o.gql[httpPath.host].GQLRender(w, r, sessionID)
			fmt.Fprint(w, rx)
		} else {
			logs.System.Fatal().Printf("%s domain do not exists.", httpPath.host)
		}
	*/

	rx := o.gql.GQLRender(w, r, sessionID)
	fmt.Fprint(w, rx)
}
func (o *Http) restServeHTTP(w http.ResponseWriter, r *http.Request, httpPath *Path, sessionID string) {
	//w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	//por favor, revisa que o.serverName exista, si no existe entonces devuelvele un dedito
	/*
		if o.rest[httpPath.host] != nil {
			o.rest[httpPath.host].RestRender(w, r, sessionID)
		} else {
			logs.System.Fatal().Printf("%s domain do not exists.", httpPath.host)
		}
	*/
	o.rest.RestRender(w, r, sessionID)
}
func (o *Http) websocketServeHTTP(w http.ResponseWriter, r *http.Request, httpPath *Path, sessionID string) {
	headers := http.Header{}
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
			break
		}
		go o.WebSocketMessage(mt, message, id, httpPath, sessionID)
	}
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
func (o *Http) setSessionIndex(sessionID string) {
	routineID := systemutils.GetRoutineID()
	sessionIndex[routineID] = sessionID
}
func (o *Http) sessionDestroy() {
	routineID := systemutils.GetRoutineID()
	delete(sessionIndex, routineID)
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
