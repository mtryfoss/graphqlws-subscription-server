package main

import (
	"context"
	"flag"
	"log"
	"sync"

	"github.com/functionalfoundry/graphqlws"
	gss "github.com/taiyoh/graphqlws-subscription-server"
)

func main() {
	var confPath string
	flag.StringVar(&confPath, "config", "config.toml", "config path")
	flag.Parse()

	conf, err := NewConf(confPath)
	if err != nil {
		log.Fatalln("conf load error")
	}

	subChan := make(chan *gss.SubscribeEvent, conf.Server.MaxHandlerCount)
	unsubChan := make(chan *gss.UnsubscribeEvent, conf.Server.MaxHandlerCount)

	types := []gss.GraphQLType{NewComment(subChan, unsubChan)}
	schema, err := gss.BuildSchema(types)
	if err != nil {
		log.Fatalln("GraphQL schema is invalid")
	}

	listener := gss.NewListener(schema)

	ctx := context.Background()
	wg := &sync.WaitGroup{}

	receiver := gss.NewReceiver(conf.Server.MaxHandlerCount)
	receiver.Start(ctx, wg, listener)

	authCallback := AuthenticateCallback(conf.Auth.SecretKey)

	server := gss.NewServer(conf.Server.Port)

	server.RegisterHandle("/subscription", graphqlws.NewHandler(graphqlws.HandlerConfig{
		SubscriptionManager: listener,
		Authenticate:        authCallback,
	}))

	server.RegisterHandle("/notify", gss.NewNotifyHandler(receiver.GetNotifierChan()))

	server.Start(ctx, wg)

	wg.Wait()
}
