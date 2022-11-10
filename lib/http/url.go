package http

import (
	"log"
	"net"
	"net/http"
	"net/url"
)

func (o *URL) Split(r *http.Request) {
	port := ""
	var splitError error
	isURI := false
	if _, err := url.ParseRequestURI("http://" + r.Host); err == nil {
		isURI = true
		o.Host, o.Port, splitError = net.SplitHostPort(r.Host)
	} else {
		log.Println("18", err.Error())
		_, _, splitError = net.SplitHostPort(r.Host)
	}

	if splitError != nil && isURI {
		tmpPort := ":80"
		if r.TLS != nil {
			tmpPort = ":443"
		}
		o.Host, o.Port, splitError = net.SplitHostPort(r.Host + tmpPort)
	}
	if splitError != nil {
		log.Println("29", "splitError", splitError.Error())
	}
	o.RequestURI = r.RequestURI
	o.TLS = false
	isURI = false
	o.Scheme = "http"

	for _, header := range r.Header["Upgrade"] {
		if header == "websocket" {
			o.Scheme = "ws"
			break
		}
	}
	if r.TLS != nil {
		o.TLS = true
		switch o.Scheme {
		case "http":
			o.Scheme = "https"
			break
		case "ws":
			o.Scheme = "wss"
			break
		}
	}

	if o.Port != "" {
		port = ":" + o.Port
	}
	o.Method = r.Method
	o.Referer = r.Referer()
	o.URL = o.Scheme + "://" + o.Host + port
	o.Origin.URL = r.Header.Get("Origin")

	if o.Origin.URL != "" {
		origin, err := url.ParseRequestURI(o.Origin.URL)
		o.Origin.Scheme = origin.Scheme
		if err == nil {
			isURI = true
			o.Origin.Host, o.Origin.Port, splitError = net.SplitHostPort(origin.Host)
		} else {
			log.Println("70", origin.Host, err.Error())
			_, _, splitError = net.SplitHostPort(r.Host)
		}
		//o.Origin.Host, o.Origin.Port, splitError = net.SplitHostPort(origin.Host)
		tmpPort := ":80"
		if r.TLS != nil {
			tmpPort = ":443"
		}
		if splitError != nil && isURI {

			o.Origin.Host, o.Origin.Port, splitError = net.SplitHostPort(origin.Host + tmpPort)
		}
		if splitError != nil {
			log.Println("83", splitError.Error(), origin.Host+tmpPort)
		}
	}
}
