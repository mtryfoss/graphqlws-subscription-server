package graphqlws_subscription_server

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	gss "github.com/taiyoh/graphqlws-subscription-server"
)

type Notification struct {
	notifyChan chan *gss.RequestData
}

func readJSONContent(req *http.Request) (*gss.RequestData, error) {
	if contentType := req.Header.Get("Content-Type"); contentType != "application/json" {
		return nil, errors.New("Content-type requires application/json")
	}
	bufbody := new(bytes.Buffer)
	if _, err := bufbody.ReadFrom(req.Body); err != nil {
		return nil, err
	}
	data := &gss.RequestData{}
	if err := json.Unmarshal(bufbody.Bytes(), data); err != nil {
		return nil, errors.New("cannot parse invalid JSON request data.")
	}
	if err := data.Validate(); err != nil {
		return nil, err
	}
	return data, nil
}

func (r *Notification) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	data, err := readJSONContent(req)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprint(w, err.Error())
		return
	}
	r.notifyChan <- data
	fmt.Fprint(w, "OK")
}

func (h *Handler) NewNotifyHandler(ch chan *gss.RequestData) http.Handler {
	return &Notification{notifyChan: ch}
}
