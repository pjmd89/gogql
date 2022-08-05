package http

import (
	"net"
	"net/http"
	"net/url"
)

func(o * URL) Split(r *http.Request){
	port := "";
	o.Host,o.Port,_ = net.SplitHostPort(r.Host);
	o.RequestURI = r.RequestURI;
	o.TLS = false;
	o.Scheme = "http";

	for _, header := range r.Header["Upgrade"] {
        if header == "websocket" {
            o.Scheme = "ws"
            break
        }
    }
	if r.TLS != nil{
		o.TLS = true;
		switch o.Scheme{
		case "http":
			o.Scheme = "https";
			break;
		case "ws":
			o.Scheme = "wss";
			break;
		}
	}
	
	if o.Port != ""{
		port = ":"+o.Port;
	}
	o.Method = r.Method;
	o.Referer = r.Referer();
	o.URL = o.Scheme+"://"+o.Host+port;
	o.Origin.URL = r.Header.Get("Origin");

	if o.Origin.URL != "" {
		origin,_ := url.Parse(o.Origin.URL)
		o.Origin.Scheme = origin.Scheme
		o.Origin.Host,o.Origin.Port,_ = net.SplitHostPort(origin.Host);
	}
}