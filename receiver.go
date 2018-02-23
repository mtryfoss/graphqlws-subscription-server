package gss

import (
	"context"
	"encoding/json"
	"errors"
	"sync"

	"github.com/functionalfoundry/graphqlws"
)

type RequestData struct {
	Users   []string    `json:"users"`
	Channel string      `json:"channel"`
	Payload interface{} `json:"payload"`
}

func (d *RequestData) Validate() error {
	if d.Channel == "" {
		return errors.New("require channel")
	}
	if d.Payload == nil {
		return errors.New("require payload")
	}
	return nil
}

func NewRequestDataFromBytes(b []byte) (*RequestData, error) {
	data := &RequestData{}
	if err := json.Unmarshal(b, data); err != nil {
		return nil, errors.New("cannot parse invalid JSON request data")
	}
	if err := data.Validate(); err != nil {
		return nil, err
	}
	return data, nil
}

type Receiver struct {
	notifyChan chan *RequestData
}

func NewReceiver(handleCount uint) *Receiver {
	return &Receiver{
		notifyChan: make(chan *RequestData, handleCount),
	}
}

func (r *Receiver) GetNotifierChan() chan *RequestData {
	return r.notifyChan
}

func (r *Receiver) Start(ctx context.Context, wg *sync.WaitGroup, l *Listener) {
	sendData := func(subscriptions graphqlws.Subscriptions, payload interface{}) {
		for conn := range subscriptions {
			for _, subscription := range subscriptions[conn] {
				res := l.DoGraphQL(BuildCtx("payload", payload, conn), subscription)
				d := &graphqlws.DataMessagePayload{
					Data: res.Data,
				}
				if res.HasErrors() {
					d.Errors = graphqlws.ErrorsFromGraphQLErrors(res.Errors)
				}
				subscription.SendData(d)
			}
		}
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case data := <-r.notifyChan:
				if len(data.Users) > 0 {
					sendData(l.GetUserSubscriptions(data.Channel, data.Users), data.Payload)
				} else {
					sendData(l.GetChannelSubscriptions(data.Channel), data.Payload)
				}
			}
		}
	}()
}
