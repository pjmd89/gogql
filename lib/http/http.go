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
	"github.com/pjmd89/gogql/lib"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"golang.org/x/exp/slices"
)
type CheckOrigin func(host *url.URL, referer *url.URL)( r bool)
type OnBegin func(host *url.URL, referer *url.URL, mode string)( r bool)
type OnFinish func()

var upgrader = websocket.Upgrader{
	EnableCompression: true,
}
var WsIds map[string]chan bool = make(map[string]chan bool);
var WsChannels map[string] *websocket.Conn = make(map[string] *websocket.Conn);

func Init(gql... Gql) *Http {
	mapGQL := make(map[string]Gql);
	for _,v := range gql{
		mapGQL[v.GetServerName()] = v;
	}
	o := &Http{HttpPort: "8080", HttpsPort: "8443", gql: mapGQL};
	lib.GetJson("http/http.json", &o);
	
	if  strings.Trim(o.HttpPort," ") == "" {
		o.HttpPort = "8080";
	}
	if  strings.Trim(o.HttpsPort," ") == "" {
		o.HttpsPort = "8443";
	}
	if  strings.Trim(o.CookieName," ") == "" {
		o.CookieName = "GOGQL_SESSION";
	}
	return o;
}
func (o *Http) Start(){
	var err error;

	channel := make(chan bool);
	stop := false;
	isTls := false;
	tlsConfig := &tls.Config{};
	tlsConfig.Certificates = []tls.Certificate{};
	/*
	xhttp:=&Http{Url:v.ServerName}
	o.router[counter].NotFoundHandler = xhttp;
	o.router[counter].MethodNotAllowedHandler = xhttp;
	*/
	o.router = mux.NewRouter();
	o.router.Use(handlers.CompressHandler);
	subrouter := o.router.MatcherFunc(o.MatcherFunc).Subrouter();
	subrouter.Methods("POST","GET", "OPTIONS").Handler(o);

	hTTPService := &http.Server{Addr: ":" + o.HttpPort, Handler: o.router};
	hTTPSService := &http.Server{Addr: ":" + o.HttpsPort, Handler: o.router, TLSConfig: tlsConfig};

	for _,server := range o.Server{
		if server.EnableHttps {
			isTls = true;
		}
		if server.EnableHttps{
			tmp, certErr := tls.LoadX509KeyPair("etc/http/certs/"+server.Cert,"etc/http/certs/"+server.Key);
			if certErr !=nil{
				fmt.Println(certErr);
				err = fmt.Errorf("Error on certificate. " + server.Cert);
				stop = true;
			}
			tlsConfig.Certificates = append(tlsConfig.Certificates, tmp);
		}
	}
	fmt.Println("http server start");
	go o.listenHttp(channel, hTTPService);
	if isTls {
		go o.listenHttps(channel,hTTPSService);
	}
	tlsConfig.BuildNameToCertificate();
	for !stop{
		x := <-channel;
		if(x == true){
			stop = true;
			hTTPService.Shutdown(nil);
		}
	}
	
	if err != nil {
		fmt.Println(err);
	}
	
}
func(o *Http) MatcherFunc(hr *http.Request, hm *mux.RouteMatch)(r bool){
	r = true;
	var url URL;
	url.Split(hr);
	if o.CheckOrigin != nil{
		r = o.CheckOrigin(url);
	}
	return r;
}
func (o *Http) ServeHTTP(w http.ResponseWriter,r *http.Request){
	var urlInfo URL;
	urlInfo.Split(r);
	urlPath := "/";
	var httpPathMode *Path;
	

	if urlInfo.Origin.URL != "" {
		w.Header().Set("Access-Control-Allow-Origin", urlInfo.Origin.URL);
		w.Header().Set("Access-Control-Allow-Credentials", "true");
	}
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS");
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization");
	//w.Header().Set("Access-Control-Max-Age", "86400");
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusNoContent);
		return;
	}
	parseUrl,_ := url.Parse(urlInfo.URL+urlInfo.RequestURI)
	
	if parseUrl.Path != "/"{
		urlPath = strings.TrimSuffix(parseUrl.Path,"/");
	}
	
	for i,svr := range o.Server{
		if strings.Trim(svr.Host," ") == "*" {
			tmp := make([]server,0)
			tmp = append(tmp, o.Server[:i]...)
			tmp = append(tmp, o.Server[i+1:]...)
			o.Server = append(tmp, svr)
		}
	}
	for _,server := range o.Server{
		var httpModes []Path;
		host := strings.Replace(server.Host,"*",`[^.]+`,-1);
		if strings.Trim(server.Host," ") == "*" {
			host = `[^/]+`;
		}
		match,_ := regexp.MatchString(host,urlInfo.Host);
		serverBreak := false;
		if match && !slices.Contains(server.Reject, urlInfo.Host){
			for _,serverPath := range server.Path{
				matchPath,_ := regexp.MatchString(`^`+serverPath.Endpoint,urlPath);
				if matchPath{
					serverBreak = true;
					serverPath.len = len(serverPath.Endpoint);
					httpModes = append(httpModes,serverPath);
				}
			}
		}
		if serverBreak {
			mode := 0;
			if len(httpModes) > 0{
				for i,httpMode := range httpModes[1:]{
					if httpMode.len < httpModes[mode].len{
						mode = i;
					}
				}
			}
			if len(httpModes) >= mode  {
				httpPathMode = &Path{};
				*httpPathMode = httpModes[mode];
				httpPathMode.pathURL = urlPath;
				httpPathMode.host = server.Host;
				httpPathMode.origin = urlInfo.Origin.URL;
			}
			break;
		}
	}
	if httpPathMode != nil{
		cookie,_ := r.Cookie(o.CookieName);
		var cookieValue []byte;
		if cookie != nil {
			cookieValue = []byte(cookie.Value);
		}
		SessionStart(w,r,&cookieValue,o.CookieName)
		if o.OnBegin != nil {
			o.OnBegin( urlInfo, httpPathMode );
		}
		switch httpPathMode.Mode{
		case "file":
			o.fileServeHTTP(w,r,httpPathMode);
			break;
		case "gql":
			o.gqlServeHTTP(w,r,httpPathMode);
			break;
		case "websocket":
			o.websocketServeHTTP(w,r,httpPathMode);
			break;
		}
		if o.OnFinish != nil {
			o.OnFinish();
		}
	} else {
		fmt.Println("mode not found");
		w.WriteHeader(http.StatusNotFound);
	}
	if o.OnFinish != nil {
		o.OnFinish();
	}
	return;
}
func(o *Http) fileServeHTTP(w http.ResponseWriter, r *http.Request, httpPath *Path){
	file, fErr := os.Open(httpPath.Path+"/"+httpPath.pathURL);
	fileStat,_ := file.Stat();
	
	if fErr != nil || fileStat.IsDir(){
		file, fErr = os.Open(httpPath.Path+httpPath.pathURL+"/index.html");
		fileStat,_ = file.Stat();
	}
	
	if false && fErr != nil{
		file, fErr = os.Open(httpPath.Path+"/index.html");
		fileStat,_ = file.Stat();
	}
	
	if fErr != nil || fileStat.IsDir(){
		w.WriteHeader(http.StatusNotFound);
		fmt.Fprint(w,"file not found, archivo no se encuentra");

		return;
	}
	
	http.ServeContent(w,r,httpPath.Path+"/"+r.RequestURI,time.Time{},file);
}
func(o *Http) gqlServeHTTP(w http.ResponseWriter, r *http.Request, httpPath *Path){
	w.Header().Set("Content-Type", "application/json; charset=UTF-8");
	//por favor, revisa que o.serverName exista, si no existe entonces devuelvele un dedito
	rx := o.gql[httpPath.host].GQLRender(w,r);
	fmt.Fprint(w,rx);
}
func(o *Http) websocketServeHTTP(w http.ResponseWriter, r *http.Request, httpPath *Path){
		headers := http.Header{};
		headers.Set("Sec-WebSocket-Protocol", "graphql-transport-ws");
		headers.Set("Sec-WebSocket-Version", "13");
		headers.Set("Content-Type", "application/json; charset=UTF-8");
		headers.Set("Access-Control-Allow-Credentials", "true");
		headers.Set("Access-Control-Allow-Origin", httpPath.origin);

		id := uuid.New().String();
		WsIds[id] = make(chan bool,1);
		var upgraderError error;
		WsChannels[id], upgraderError = upgrader.Upgrade(w, r, headers)
		if upgraderError != nil{

			//fmt.Println(upgraderError);
			select{
				case WsIds[id] <- false:
			}
			close(WsIds[id]);
			delete(WsIds,id);
			delete(WsChannels,id);
			//fmt.Println(upgraderError);
		}
		defer WsChannels[id].Close();
		while:=true;
		for while{
			mt, message, err := WsChannels[id].ReadMessage();
			if err != nil {
				select{
				case WsIds[id] <- false:
				}
				close(WsIds[id]);
				delete(WsIds,id);
				delete(WsChannels,id);
				while = false;
				break;
			}
			go o.WebSocketMessage(mt, message, id, httpPath);
		}
}
func (o *Http) listenHttp(channel chan bool, handler *http.Server){
	channel <- false;
	err := handler.ListenAndServe();
	if(err != nil){
		fmt.Println("http server start error: " + err.Error());
		channel <- true;
	}
	
}
func(o *Http) listenHttps(channel chan bool, handler *http.Server){
	channel <- false;
	listener , tlsErr := tls.Listen("tcp",handler.Addr,handler.TLSConfig);
	if tlsErr != nil{
		fmt.Println("https server start error: " + tlsErr.Error());
		channel <- true;
	}
	
	err := handler.Serve(listener);
	if err != nil{
		fmt.Println("https server error: " + err.Error());
		channel <- true;
	}
}
func (o *Http) WebSocketMessage(mt int, message []byte, id string, httpPath *Path ){
	
	o.gql[httpPath.host].GQLRenderSubscription(mt,message,id);
}
func WriteWebsocketMessage(mt int , id string,message []byte){
	if WsChannels[id] != nil{
		WsChannels[id].WriteMessage(mt,message);
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
func in_array(val string, array []string) (ok bool) {
    for i := range array {
        if ok = array[i] == val; ok {
            return
        }
    }
    return
}
/*
func (o *Http) Start() {
	channel := make(chan bool);
	stop := false;
	isTls := false;
	var err string;
	
	tlsConfig := &tls.Config{};
	tlsConfig.Certificates = []tls.Certificate{};
	o.router = make(map[int]*mux.Router);
	counter := 0;
	for _, v := range o.Server {
		
		
		o.HTTPService = &http.Server{Addr: ":" + o.HttpPort, Handler: o.router[counter]};
		o.HTTPSService = &http.Server{Addr: ":" + o.HttpsPort, Handler: o.router[counter], TLSConfig: tlsConfig};
		endpoints := make(map[string]bool,0);
		if strings.Trim(v.ServerName," ") ==""{
			v.ServerName = "localhost";
		}
		_, ok := o.Path[v.ServerName];
		if  ok {
			err = "Server name is not shoud be the same";
			stop = true;
			break;
		}
		//replaceServerName := strings.Replace(v.ServerName,"*","{subdomain}",-1);
		v.subrouter = o.router[counter].MatcherFunc(func(hh *http.Request, bb *mux.RouteMatch)bool{
			fmt.Println(v.ServerName);
			return false;
		}).Subrouter();
		//v.subrouter = o.router[counter].Host(v.ServerName).Subrouter();
		
		if v.EnableHttps {
			isTls = true;
		}
		if v.EnableHttps{
			tmp, certErr := tls.LoadX509KeyPair("etc/http/certs/"+v.Cert,"etc/http/certs/"+v.Key);
			if certErr !=nil{
				fmt.Println(certErr);
				err = "Error on certificate. "+v.Cert;
				stop = true;
			}
			
			tlsConfig.Certificates = append(tlsConfig.Certificates, tmp);
		}
	}
	fmt.Println("http server start");
	go o.listenHttp(channel, *o.HTTPService);
	if isTls {
		go o.listenHttps(channel,*o.HTTPSService);
	}
	tlsConfig.BuildNameToCertificate();
	for !stop{
		x := <-channel;
		if(x == true){
			stop = true;
			o.HTTPService.Shutdown(nil);
		}
	}
	if(err != ""){
		fmt.Println(err);
	}
	
}
func(o *Http) MatcherFunc(hh *http.Request, bb *mux.RouteMatch)(r bool){
	
	return r;
}

func (o *pathConfig) ServeHTTP(w http.ResponseWriter,r *http.Request){
	upgrade := false
	protocol := `http://`;
	secureProtocol := "https://";
	if r.TLS != nil {
		protocol = `https://`;
	}
	if upgrade {
		secureProtocol = "wss://";
		switch protocol{
		case "http://":
			protocol = "ws://";
		case "https://":
			protocol = "wss://";
		}
	}
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusNoContent);
		return;
	}
    for _, header := range r.Header["Upgrade"] {
        if header == "websocket" {
            upgrade = true
            break
        }
    }

	hostSplit,_ := url.Parse(protocol+r.Host);
	refererSplit,_ := url.Parse(r.Referer());
	switch o.Mode{
	case "gql","websocket","file":
		cookie,_ := r.Cookie("NUEVE_SESSION");
		var cookieValue []byte;
		if cookie != nil {
			cookieValue = []byte(cookie.Value);
		}
		//fmt.Println(string(cookieValue))
		SessionStart(w,r,&cookieValue,"NUEVE_SESSION")
	}
	httpsURI, _,_ := net.SplitHostPort(hostSplit.Host);
	
	if o.httpsPort != "443" && o.enableHttps && o.redirect{
		httpsURI += ":"+o.httpsPort;
	}
	
	if o.redirect && o.enableHttps && r.TLS == nil {
		http.Redirect(w,r,secureProtocol+httpsURI+r.RequestURI,301);
		return;
	}
	isAllow := o.isAllowOrigin( hostSplit, refererSplit , o.Mode);
	if !isAllow{
		w.WriteHeader(http.StatusUnauthorized);
		return;
	}
	upgrader.CheckOrigin = func(r *http.Request) bool { 
		
		return isAllow;
	}
	
	if refererSplit.Host != "" {
		w.Header().Set("Access-Control-Allow-Origin", protocol+refererSplit.Host);
		w.Header().Set("Access-Control-Allow-Credentials", "true");
	}
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE");
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization");
	//w.Header().Set("Access-Control-Max-Age", "86400");
	
	if o.OnBegin != nil {
		o.OnBegin(hostSplit, refererSplit, o.Mode );
	}
	switch o.Mode{
	case "file":
		file, fErr := os.Open(o.Path+r.RequestURI);
		fileStat,_ := file.Stat();
		
		if fErr != nil || fileStat.IsDir(){
			file, fErr = os.Open(o.Path+r.RequestURI+"/"+o.FileDefault);
			fileStat,_ = file.Stat();
		}
		
		if o.Rewrite && fErr != nil{
			file, fErr = os.Open(o.Path+"/"+o.RewriteTo);
			fileStat,_ = file.Stat();
		}
		
		if fErr != nil || fileStat.IsDir(){
			w.WriteHeader(http.StatusNotFound);
			fmt.Fprint(w,"file not found, archivo no se encuentra");

			return;
		}
		
		http.ServeContent(w,r,o.Path+"/"+r.RequestURI,time.Time{},file);
		break;
	case "gql":
		w.Header().Set("Content-Type", "application/json; charset=UTF-8");
		
		//por favor, revisa que o.serverName exista, si no existe entonces devuelvele un dedito
		rx := o.gqlRender[o.serverName].GQLRender(w,r);
		fmt.Fprint(w,rx);
		break;
	case "websocket":
		headers := http.Header{};
		headers.Set("Sec-WebSocket-Protocol", "graphql-transport-ws");
		headers.Set("Sec-WebSocket-Version", "13");
		headers.Set("Content-Type", "application/json; charset=UTF-8");
		headers.Set("Access-Control-Allow-Origin", protocol+refererSplit.Host);

		id := uuid.New().String();
		WsIds[id] = make(chan bool,1);
		var upgraderError error;
		WsChannels[id], upgraderError = upgrader.Upgrade(w, r, headers)
		if upgraderError != nil{

			fmt.Println(upgraderError);
			select{
				case WsIds[id] <- false:
			}
			close(WsIds[id]);
			delete(WsIds,id);
			delete(WsChannels,id);
			fmt.Println(upgraderError);
		}
		defer WsChannels[id].Close();
		while:=true;
		for while{
			mt, message, err := WsChannels[id].ReadMessage();
			if err != nil {
				select{
				case WsIds[id] <- false:
				}
				close(WsIds[id]);
				delete(WsIds,id);
				delete(WsChannels,id);
				while = false;
				break;
			}
			go o.WebSocketMessage(mt, message, id);
			
		}
	default:
		w.WriteHeader(http.StatusExpectationFailed);
		fmt.Fprint(w,"Mode "+o.Mode+" not exists.");
	}

	if o.OnFinish != nil {
		o.OnFinish();
	}
	return;
}
func (o *pathConfig) WebSocketMessage(mt int, message []byte, id string ){
	
	o.gqlRender[o.serverName].GQLRenderSubscription(mt,message,id);
}
func (o *pathConfig) isAllowOrigin( host *url.URL, referer *url.URL, mode string) bool {
	isAllow :=true;

	if o.CheckOrigin != nil{
		isAllow = o.CheckOrigin( host, referer, mode );
	}

	return isAllow;
}
func WriteWebsocketMessage(mt int , id string,message []byte){
	if WsChannels[id] != nil{
		WsChannels[id].WriteMessage(mt,message);
	}
}
*/