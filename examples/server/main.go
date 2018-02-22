package main

import (
	"context"
	"flag"
	"log"
	"sync"

	"github.com/graphql-go/graphql"
	gss "github.com/taiyoh/graphqlws-subscription-server"
	gss_handler "github.com/taiyoh/graphqlws-subscription-server/handler"
)

func main() {
	var confPath string
	flag.StringVar(&confPath, "config", "config.toml", "config path")
	flag.Parse()

	conf, err := NewConf(confPath)
	if err != nil {
		log.Fatalln("conf load error")
	}

	listener := gss.NewListener()

	fields := graphql.Fields{}
	types := []gss.GraphQLType{NewComment()}
	for _, t := range types {
		fields[t.FieldName()] = t.GetField(listener)
	}
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Subscription: graphql.NewObject(
			graphql.ObjectConfig{
				Name:   "RootSubscription",
				Fields: fields,
			},
		),
	})
	if err != nil {
		log.Fatalln("GraphQL schema is invalid")
	}

	listener.BuildManager(&schema)

	ctx := context.Background()
	wg := &sync.WaitGroup{}

	server := gss.NewServer(conf.Server.Port)
	handler := gss_handler.NewHandler(listener)

	receiver := gss.NewReceiver(conf.Server.MaxHandlerCount)
	receiver.Start(ctx, wg, listener)

	authCallback := AuthenticateCallback(conf.Auth.SecretKey)

	server.RegisterHandle("/subscription", handler.NewWebsocketHandler(authCallback))
	server.RegisterHandle("/notify_channel", handler.NewNotifyChannelHandler(receiver.GetChannelNotifierChan()))
	server.RegisterHandle("/notify_users", handler.NewNotifyUsersHandler(receiver.GetUserNotifierChan()))

	server.Start(ctx, wg)

	wg.Wait()
}
