package gss

import (
	"context"
	"sync"

	"github.com/functionalfoundry/graphqlws"
	"github.com/graphql-go/graphql"
)

type ListenerContextKey string

type SubscribeService struct {
	graphqlws.SubscriptionManager
	pool       graphqlws.SubscriptionManager
	calculator SubscribeCalculator
	filter     SubscribeFilter
	notifyChan chan *RequestData
	subChan    chan *SubscribeEvent
	unsubChan  chan *UnsubscribeEvent
}

func NewSubscribeService(schema *graphql.Schema, handleCount uint, subChan chan *SubscribeEvent, unsubChan chan *UnsubscribeEvent) *SubscribeService {
	return &SubscribeService{
		pool:       graphqlws.NewSubscriptionManager(schema),
		filter:     NewSubscribeFilter(),
		calculator: NewSubscribeCalculator(schema),
		notifyChan: make(chan *RequestData, handleCount),
		subChan:    subChan,
		unsubChan:  unsubChan,
	}
}

func (s *SubscribeService) AddSubscription(conn graphqlws.Connection, sub *graphqlws.Subscription) []error {
	rctx := NewResolveContext(conn.ID(), sub.ID, conn.User(), "onSubscribe", true)
	s.calculator.DoGraphQL(rctx, sub.Query, sub.Variables, sub.OperationName)
	errs := s.pool.AddSubscription(conn, sub)
	if errs != nil {
		return errs
	}

	return nil
}

func (s *SubscribeService) RemoveSubscription(conn graphqlws.Connection, sub *graphqlws.Subscription) {
	rctx := NewResolveContext(conn.ID(), sub.ID, conn.User(), "onUnsubscribe", true)
	s.calculator.DoGraphQL(rctx, sub.Query, sub.Variables, sub.OperationName)
	s.pool.RemoveSubscription(conn, sub)
}

func (s *SubscribeService) RemoveSubscriptions(conn graphqlws.Connection) {
	for _, sub := range s.Subscriptions()[conn] {
		rctx := NewResolveContext(conn.ID(), sub.ID, conn.User(), "onUnsubscribe", true)
		s.calculator.DoGraphQL(rctx, sub.Query, sub.Variables, sub.OperationName)
	}
	s.pool.RemoveSubscriptions(conn)
}

func (s *SubscribeService) Subscriptions() graphqlws.Subscriptions {
	return s.pool.Subscriptions()
}

func (s *SubscribeService) Publish(connIds ConnIDBySubscriptionID, payload interface{}) {
	for conn, subsById := range s.pool.Subscriptions() {
		for subID, sub := range subsById {
			if _, exists := connIds[subID]; exists {
				rctx := NewResolveContext(conn.ID(), subID, conn.User(), "payload", payload)
				res := s.calculator.DoGraphQL(rctx, sub.Query, sub.Variables, sub.OperationName)
				d := &graphqlws.DataMessagePayload{
					Data: res.Data,
				}
				if res.HasErrors() {
					d.Errors = graphqlws.ErrorsFromGraphQLErrors(res.Errors)
				}
				sub.SendData(d)
			}
		}
	}
}

func (s *SubscribeService) SubscribeFilter() SubscribeFilter {
	return s.filter
}

func (s *SubscribeService) SubscribeCalculator() SubscribeCalculator {
	return s.calculator
}

func (s *SubscribeService) GetNotifierChan() chan *RequestData {
	return s.notifyChan
}

func (s *SubscribeService) Start(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	filter := s.SubscribeFilter()
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case data := <-s.notifyChan:
				var connIds ConnIDBySubscriptionID
				if len(data.Users) > 0 {
					connIds = filter.GetUserSubscriptionIDs(data.Channel, data.Users)
				} else {
					connIds = filter.GetChannelSubscriptionIDs(data.Channel)
				}
				if len(connIds) > 0 {
					go s.Publish(connIds, data.Payload)
				}
			case data := <-s.subChan:
				s.filter.Subscribe(data.Channel, data.ConnID, data.SubscriptionID, data.User)
			case data := <-s.unsubChan:
				s.filter.Unsubscribe(data.ConnID, data.User)
			}
		}
	}()
}
