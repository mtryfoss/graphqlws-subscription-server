package graphqlws_subscription_server

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	gss "github.com/taiyoh/graphqlws-subscription-server"
)

type Notification struct {
	notifyChan chan *gss.RequestData
}

type NotificationResponse struct {
	Success bool     `json:"success"`
	Errors  []string `json:"errors"`
}

func successResponse() []byte {
	r := &NotificationResponse{Success: true}
	b, _ := json.Marshal(r)
	return b
}

func failResponse(errs []string) []byte {
	r := &NotificationResponse{Success: false, Errors: errs}
	b, _ := json.Marshal(r)
	return b
}

func readJSONContent(req *http.Request) (*gss.RequestData, error) {
	if contentType := req.Header.Get("Content-Type"); contentType != "application/json" {
		return nil, errors.New("Content-Type requires application/json")
	}
	bufbody := new(bytes.Buffer)
	if _, err := bufbody.ReadFrom(req.Body); err != nil {
		return nil, err
	}
	data, err := gss.NewRequestDataFromBytes(bufbody.Bytes())
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r *Notification) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	w.Header().Set("Content-Type", "application/json")
	data, err := readJSONContent(req)
	if err != nil {
		w.WriteHeader(400)
		w.Write(failResponse([]string{err.Error()}))
		return
	}
	r.notifyChan <- data
	w.Write(successResponse())
}

func (h *Handler) NewNotifyHandler(ch chan *gss.RequestData) http.Handler {
	return &Notification{notifyChan: ch}
}
