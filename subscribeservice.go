package gss

import (
	"context"
	"sync"

	"github.com/functionalfoundry/graphqlws"
	"github.com/graphql-go/graphql"
)

type GraphQLContextKey string
type CanSendToUserFunc func(conn *graphqlws.Connection, reqData *RequestData) bool

type SubscribeService struct {
	graphqlws.SubscriptionManager
	Schema        *graphql.Schema
	Pool          graphqlws.SubscriptionManager
	Filter        SubscribeFilter
	notifyChan    chan *RequestData
	canSendToUser CanSendToUserFunc
}

func NewSubscribeService(schema *graphql.Schema, handleCount uint, c CanSendToUserFunc) *SubscribeService {
	return &SubscribeService{
		Schema:        schema,
		Pool:          graphqlws.NewSubscriptionManager(schema),
		Filter:        NewSubscribeFilter(),
		notifyChan:    make(chan *RequestData, handleCount),
		canSendToUser: c,
	}
}

func (s *SubscribeService) AddSubscription(conn graphqlws.Connection, sub *graphqlws.Subscription) []error {
	errs := s.Pool.AddSubscription(conn, sub)
	if errs != nil {
		return errs
	}

	s.Filter.ReplaceFieldsFromDocument(sub)

	return nil
}

func (s *SubscribeService) RemoveSubscription(conn graphqlws.Connection, sub *graphqlws.Subscription) {
	s.Pool.RemoveSubscription(conn, sub)
}

func (s *SubscribeService) RemoveSubscriptions(conn graphqlws.Connection) {
	s.Pool.RemoveSubscriptions(conn)
}

func (s *SubscribeService) Subscriptions() graphqlws.Subscriptions {
	return s.Pool.Subscriptions()
}

func (s *SubscribeService) Publish(reqData *RequestData) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, GraphQLContextKey("payload"), reqData.Payload)
	for conn, subsByID := range s.Pool.Subscriptions() {
		if len(reqData.Users) > 0 && !s.canSendToUser(&conn, reqData) {
			continue
		}
		for _, sub := range subsByID {
			if sub.MatchesField(reqData.Channel) {
				res := graphql.Do(graphql.Params{
					Schema:         *s.Schema, // The GraphQL schema
					RequestString:  sub.Query,
					VariableValues: sub.Variables,
					OperationName:  sub.OperationName,
					Context:        ctx,
				})
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

func (s *SubscribeService) GetNotifierChan() chan *RequestData {
	return s.notifyChan
}

func (s *SubscribeService) Start(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case data := <-s.notifyChan:
				go s.Publish(data)
			}
		}
	}()
}
