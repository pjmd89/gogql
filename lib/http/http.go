package http

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/pjmd89/gogql/lib"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)
type GQLDefault struct{
}

func(o *GQLDefault) GQLRender(w http.ResponseWriter,r *http.Request){
	
}
func(o *GQLDefault) GetServerName() string{
	return "localhost";
}
func Init(gql... Gql) *Http {
	mapGQL := make(map[string]Gql);
	for _,v := range gql{
		mapGQL[v.GetServerName()] = v;
	}
	o := &Http{HttpPort: "8080", HttpsPort: "8443", Path: make(map[string]string),gql: mapGQL};
	lib.GetJson("http/http.json", &o);
	
	if  o.HttpPort == "" {
		o.HttpPort = "8080";
	}
	if  o.HttpsPort == "" {
		o.HttpsPort = "8443";
	}
	return o;
}
func (o *Http) Start() {
	channel := make(chan bool);
	stop := false;
	isTls := false;
	var err string;
	o.router = mux.NewRouter();
	o.router.Use(handlers.CompressHandler);
	o.router.NotFoundHandler = o;
	o.router.MethodNotAllowedHandler = o;
	tlsConfig := &tls.Config{};
	tlsConfig.Certificates = []tls.Certificate{};
	o.HTTPService = &http.Server{Addr: ":" + o.HttpPort, Handler: o.router};
	o.HTTPSService = &http.Server{Addr: ":" + o.HttpsPort, Handler: o.router, TLSConfig: tlsConfig};
	
	for _, v := range o.Server {
		if strings.Trim(v.ServerName," ") ==""{
			v.ServerName = "localhost";
		}
		_, ok := o.Path[v.ServerName];
		if  ok {
			err = "Server name is not shoud be the same";
			stop = true;
			break;
		}
		v.subrouter = o.router.Host(v.ServerName).Subrouter();
		if v.EnableHttps {
			isTls = true;
		}
		if v.EnableHttps{
			tmp, certErr := tls.LoadX509KeyPair("etc/http/certs/"+v.ServerName+"/"+v.Cert,"etc/http/certs/"+v.ServerName+"/"+v.Key);
			if certErr !=nil{
				fmt.Println(certErr);
				err = "Error on certificate. "+v.Cert;
				stop = true;
			}
			
			tlsConfig.Certificates = append(tlsConfig.Certificates, tmp);
		}
		for _, vv := range v.Path {
			vv.Url = v.ServerName;
			if vv.Path == "" {
				vv.Path = "htdocs";
			}
			if strings.Trim(vv.FileDefault, " ") == ""{
				vv.FileDefault = "index.html";
			}
			if strings.Trim(vv.RewriteTo, " ") == ""{
				vv.RewriteTo = "index.html";
			}
			vv.redirect = v.Redirect;
			vv.enableHttps = v.EnableHttps;
			vv.httpsPort = o.HttpsPort;
			vv.serverName = v.ServerName;
			if vv.Mode == "gql"{
				x:= v.subrouter.Methods("POST", "OPTIONS").Path(vv.Endpoint);
				vv.gqlRender = o.gql;
				x.Handler(vv)
			} else{				
				x:= v.subrouter.Methods("GET", "OPTIONS").PathPrefix(vv.Endpoint);
				x.Handler(http.StripPrefix(vv.Endpoint,vv))
				
			}
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
func (o *Http) listenHttp(channel chan bool, handler http.Server){
	channel <- false;
	err := handler.ListenAndServe();
	if(err != nil){
		fmt.Println("http server start error: " + err.Error());
		channel <- true;
	}
	
}
func(o *Http) listenHttps(channel chan bool, handler http.Server){
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
func(o *Http)ServeHTTP(w http.ResponseWriter,r *http.Request){
	w.WriteHeader(http.StatusNotFound);
	fmt.Fprint(w,"file not found, archivo no se encuentra");
}
func (o *pathConfig) ServeHTTP(w http.ResponseWriter,r *http.Request){
	//hostSplit := strings.Split(r.Host, ":");
	httpsURI := o.Url;
	protocol := `http://`;
	if o.httpsPort != "443" && o.enableHttps && o.redirect{
		httpsURI += ":"+o.httpsPort;
	} 
	if r.TLS != nil {
		protocol = `https://`;
	}
	if o.redirect && o.enableHttps && r.TLS == nil {
		http.Redirect(w,r,"https://"+httpsURI+r.RequestURI,301);
		return;
	}
	if strings.Trim(o.AllowOrigin, " ") != "" {
		w.Header().Set("Access-Control-Allow-Origin", protocol+o.AllowOrigin);
		w.Header().Set("Access-Control-Allow-Credentials", "true");
	}
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE");
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization");
	//w.Header().Set("Access-Control-Max-Age", "86400");
	if r.Method == "OPTIONS" {
		return;
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
			fmt.Println(r.Cookie("NUEVE_SESSION"));
			cookie,_ := r.Cookie("NUEVE_SESSION");
			var cookieValue []byte;
			if cookie != nil {
				cookieValue = []byte(cookie.Value);
			}
			SessionStart(w,r,&cookieValue,"NUEVE_SESSION")
			rx := o.gqlRender[o.serverName].GQLRender(w,r);
			fmt.Fprint(w,rx);
			break;
		default:
			w.WriteHeader(http.StatusExpectationFailed);
			fmt.Fprint(w,"Mode "+o.Mode+" not exists.");
		}
	return;
}
