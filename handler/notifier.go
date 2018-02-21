package graphqlws_subscription_server

import (
	"bytes"
	"encoding/json"
	"net/http"

	gss "github.com/taiyoh/graphqlws-subscription-server"
)

type ChannelNotification struct {
	notifyChan chan gss.ChannelRequestPayload
}

type UsersNotification struct {
	notifyChan chan gss.UserRequestPayload
}

func (r *ChannelNotification) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	payload := gss.ChannelRequestPayload{}
	bufbody := new(bytes.Buffer)
	bufbody.ReadFrom(req.Body)
	err := json.Unmarshal(bufbody.Bytes(), &payload)
	if err != nil {
		w.WriteHeader(400)
		return
	}
	r.notifyChan <- payload
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func (r *UsersNotification) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	payload := gss.UserRequestPayload{}
	bufbody := new(bytes.Buffer)
	bufbody.ReadFrom(req.Body)
	err := json.Unmarshal(bufbody.Bytes(), &payload)
	if err != nil {
		w.WriteHeader(400)
		return
	}
	r.notifyChan <- payload
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func (h *Handler) NewNotifyChannelHandler(ch chan gss.ChannelRequestPayload) http.Handler {
	return &ChannelNotification{notifyChan: ch}
}

func (h *Handler) NewNotifyUsersHandler(ch chan gss.UserRequestPayload) http.Handler {
	return &UsersNotification{notifyChan: ch}
}
