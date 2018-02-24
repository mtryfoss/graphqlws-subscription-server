package main

import (
	"context"
	"flag"
	"log"
	"sync"

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

	server := gss.NewServer(conf.Server.Port)
	handler := gss_handler.NewHandler(listener)

	receiver := gss.NewReceiver(conf.Server.MaxHandlerCount)
	receiver.Start(ctx, wg, listener)

	authCallback := AuthenticateCallback(conf.Auth.SecretKey)

	server.RegisterHandle("/subscription", handler.NewWebsocketHandler(authCallback))
	server.RegisterHandle("/notify", handler.NewNotifyHandler(receiver.GetNotifierChan()))

	server.Start(ctx, wg)

	wg.Wait()
}
