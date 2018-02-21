package graphqlws_subscription_server

import (
	"bytes"
	"encoding/json"
	"net/http"

	gss "github.com/taiyoh/graphqlws-subscription-server"
)

type ChannelNotification struct {
	notifyChan chan gss.ChannelRequestData
}

type UsersNotification struct {
	notifyChan chan gss.UserRequestData
}

func (r *ChannelNotification) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	data := gss.ChannelRequestData{}
	if contentType := req.Header.Get("Content-type"); contentType != "application/json" {
		w.WriteHeader(400)
		w.Write([]byte("Content-type requires application/json"))
		return
	}
	bufbody := new(bytes.Buffer)
	bufbody.ReadFrom(req.Body)
	err := json.Unmarshal(bufbody.Bytes(), &data)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("OK"))
		return
	}
	r.notifyChan <- data
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func (r *UsersNotification) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	data := gss.UserRequestData{}
	if contentType := req.Header.Get("Content-type"); contentType != "application/json" {
		w.WriteHeader(400)
		w.Write([]byte("Content-type requires application/json"))
		return
	}
	bufbody := new(bytes.Buffer)
	bufbody.ReadFrom(req.Body)
	err := json.Unmarshal(bufbody.Bytes(), &data)
	if err != nil {
		w.WriteHeader(400)
		return
	}
	r.notifyChan <- data
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func (h *Handler) NewNotifyChannelHandler(ch chan gss.ChannelRequestData) http.Handler {
	return &ChannelNotification{notifyChan: ch}
}

func (h *Handler) NewNotifyUsersHandler(ch chan gss.UserRequestData) http.Handler {
	return &UsersNotification{notifyChan: ch}
}
