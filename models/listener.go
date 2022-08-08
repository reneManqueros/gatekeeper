package models

import (
	"encoding/base64"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
)

type Listener struct {
	Address             string   `json:"address"`
	AuthenticatedTokens []string `json:"authenticated_tokens"`
	Destination         string   `json:"destination"`
	Host                string   `json:"host"`
}

func (l *Listener) Serve() {
	err := http.ListenAndServe(l.Address, http.HandlerFunc(l.handleRequest))
	if err != nil {
		log.Println(err)
	}
}

func (l *Listener) AuthRequired(wr http.ResponseWriter, req *http.Request) {
	wr.Header().Set("WWW-Authenticate", `Basic`)
	http.Error(wr, "Unauthorized", http.StatusUnauthorized)
}

func (l *Listener) handleRequest(wr http.ResponseWriter, req *http.Request) {
	isAuth := false

	if value, ok := req.Header["Authorization"]; ok && len(value) > 0 {
		authHeader := value[0]
		authHeader = strings.TrimPrefix(authHeader, "Basic ")
		data, err := base64.StdEncoding.DecodeString(authHeader)
		if err != nil {
			log.Println("error:", err)
		}
		for _, userToken := range l.AuthenticatedTokens {
			if string(data) == userToken {
				isAuth = true
				break
			}
		}
	}

	if isAuth == false {
		l.AuthRequired(wr, req)
		return
	}

	client := &http.Client{}
	req.RequestURI = ""

	delHopHeaders(req.Header)

	if clientIP, _, err := net.SplitHostPort(req.RemoteAddr); err == nil {
		appendHostToXForwardHeader(req.Header, clientIP)
	}
	req.URL.Scheme = "http"
	req.Host = l.Destination
	req.URL.Host = req.Host
	resp, err := client.Do(req)
	if err != nil {
		http.Error(wr, "Server Error", http.StatusInternalServerError)
		log.Println("ServeHTTP:", err)
	}
	defer resp.Body.Close()

	delHopHeaders(resp.Header)
	copyHeader(wr.Header(), resp.Header)
	wr.WriteHeader(resp.StatusCode)
	_, _ = io.Copy(wr, resp.Body)
}
