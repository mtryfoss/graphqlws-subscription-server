package main

import (
	"context"
	"flag"
	"log"

	gss "github.com/taiyoh/graphqlws-subscription-server"
	gss_handler "github.com/taiyoh/graphqlws-subscription-server/handler"
)

func main() {
	var confPath string
	flag.StringVar(&confPath, "config", "config.toml", "config path")
	flag.Parse()

	conf, err := gss.NewConf(confPath)
	if err != nil {
		log.Fatalln("conf load error")
	}

	listener := gss.NewListener(conf.Auth.DummyUserID)

	schema, err := LoadSchema(listener)
	if err != nil {
		log.Fatalln("GraphQL schema is invalid")
	}

	ctx := context.Background()

	server := gss.NewServer(conf.Server)
	handler := gss_handler.NewHandler(&schema)

	server.RegisterHandle("/subscription", handler.NewWebsocketHandler(AuthenticateCallback(conf.Auth.SecretKey)))
	server.RegisterHandle("/notify_channel", handler.NewNotifyChannelHandler())
	server.RegisterHandle("/notify_users", handler.NewNotifyUsersHandler())

	server.Start(ctx)
}
