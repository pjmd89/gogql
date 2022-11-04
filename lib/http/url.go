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
	o.Host, o.Port, splitError = net.SplitHostPort(r.Host)
	if splitError != nil {
		tmpPort := ":80"
		if r.TLS != nil {
			tmpPort = ":443"
		}
		o.Host, o.Port, splitError = net.SplitHostPort(r.Host + tmpPort)
		if splitError != nil {
			log.Println(splitError.Error())
		}
	}
	o.RequestURI = r.RequestURI
	o.TLS = false
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
		origin, _ := url.Parse(o.Origin.URL)
		o.Origin.Scheme = origin.Scheme
		o.Origin.Host, o.Origin.Port, splitError = net.SplitHostPort(origin.Host)
		if splitError != nil {
			tmpPort := ":80"
			if r.TLS != nil {
				tmpPort = ":443"
			}
			o.Origin.Host, o.Origin.Port, splitError = net.SplitHostPort(origin.Host + tmpPort)
			if splitError != nil {
				log.Println(splitError.Error())
			}
		}
	}
}
