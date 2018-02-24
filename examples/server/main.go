package main

import (
	"context"
	"flag"
	"log"
	"sync"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/functionalfoundry/graphqlws"
	"github.com/graphql-go/graphql"
	gss "github.com/taiyoh/graphqlws-subscription-server"
)

type ConnectedUser struct {
	jwt.StandardClaims
}

func (u ConnectedUser) Name() string {
	return u.Subject
}

func AuthenticateCallback(secretkey string) graphqlws.AuthenticateFunc {
	return func(tokenstring string) (interface{}, error) {
		user := ConnectedUser{}
		_, err := jwt.ParseWithClaims(tokenstring, &user, func(token *jwt.Token) (interface{}, error) {
			return []byte(secretkey), nil
		})
		return user, err
	}
}

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
	fields := graphql.Fields{}
	for _, t := range types {
		fields[t.FieldName()] = gss.BuildField(t)
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

	listener := gss.NewListener(&schema, conf.Server.MaxHandlerCount, subChan, unsubChan)

	ctx := context.Background()
	wg := &sync.WaitGroup{}

	listener.Start(ctx, wg)

	server := gss.NewServer(conf.Server.Port)

	server.RegisterHandle("/subscription", graphqlws.NewHandler(graphqlws.HandlerConfig{
		SubscriptionManager: listener,
		Authenticate:        AuthenticateCallback(conf.Auth.SecretKey),
	}))

	server.RegisterHandle("/notify", gss.NewNotifyHandler(listener.GetNotifierChan()))

	server.Start(ctx, wg)

	wg.Wait()
}
