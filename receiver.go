package graphqlws_subscription_server

import (
	"context"
	"sync"

	"github.com/functionalfoundry/graphqlws"
)

type ChannelRequestData struct {
	Channel string      `json:"channel"`
	payload interface{} `json:"payload"`
}

type UserRequestData struct {
	Users   []string    `json:"users"`
	payload interface{} `json:"payload"`
}

func (p *ChannelRequestData) Payload() interface{} {
	return p.payload
}

func (p *UserRequestData) Payload() interface{} {
	return p.payload
}

type NotifyRequestData interface {
	Payload() interface{}
}

type Receiver struct {
	notifyChannelChan chan ChannelRequestData
	notifyUserChan    chan UserRequestData
}

func NewReceiver(handleCount uint) *Receiver {
	return &Receiver{
		notifyChannelChan: make(chan ChannelRequestData, handleCount),
		notifyUserChan:    make(chan UserRequestData, handleCount),
	}
}

func (r *Receiver) GetChannelNotifierChan() chan ChannelRequestData {
	return r.notifyChannelChan
}

func (r *Receiver) GetUserNotifierChan() chan UserRequestData {
	return r.notifyUserChan
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
			case data := <-r.notifyChannelChan:
				sendData(l.GetChannelSubscriptions(data.Channel), data.Payload())
			case data := <-r.notifyUserChan:
				sendData(l.GetUserSubscriptions(data.Users), data.Payload())
			}
		}
	}()
}
