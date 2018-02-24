package main

import (
	"context"
	"flag"
	"log"
	"sync"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/functionalfoundry/graphqlws"
	"github.com/graphql-go/graphql"
	toml "github.com/pelletier/go-toml"

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

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Subscription: graphql.NewObject(
			graphql.ObjectConfig{
				Name: "RootSubscription",
				Fields: graphql.Fields{
					"newComment": &graphql.Field{
						Type: graphql.NewObject(graphql.ObjectConfig{
							Name: "Comment",
							Fields: graphql.Fields{
								"id":      &graphql.Field{Type: graphql.NewNonNull(graphql.ID)},
								"content": &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
							},
						}),
						Args: graphql.FieldConfigArgument{
							"roomId": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.ID)},
						},
						Resolve: gss.BuildResolve(NewComment(subChan, unsubChan)),
					},
				},
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

//
// <!-- Authenticate section start
//

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

//
// Authenticate section end -->
//

//
// <!-- Config section start
//

type Conf struct {
	Server ServerConf `toml:"server"`
	Auth   AuthConf   `toml:"auth"`
}

type ServerConf struct {
	Port            uint `toml:"port"`
	MaxHandlerCount uint `toml:"max_handler_count"`
}

type AuthConf struct {
	SecretKey string `toml:"secret_key"`
}

func NewConf(path string) (*Conf, error) {
	config, err := toml.LoadFile(path)
	if err != nil {
		return nil, err
	}

	conf := &Conf{}
	config.Unmarshal(conf)

	return conf, nil
}

//
// Config section end -->
//

//
// <!-- Resolver section start
//

type SampleComment struct {
	id      string `json:"id"`
	content string `json:"content"`
}

func newDummyResponse() *SampleComment {
	return &SampleComment{id: "id", content: "content"}
}

type Comment struct {
	gss.GraphQLResolve
	subscribeChan   chan *gss.SubscribeEvent
	unsubscribeChan chan *gss.UnsubscribeEvent
}

func NewComment(subChan chan *gss.SubscribeEvent, unsubChan chan *gss.UnsubscribeEvent) *Comment {
	return &Comment{subscribeChan: subChan, unsubscribeChan: unsubChan}
}

func (c *Comment) OnPayload(payload interface{}, p graphql.ResolveParams) (interface{}, error) {
	comment := payload.(SampleComment)
	return comment, nil
}

func (c *Comment) OnSubscribe(p graphql.ResolveParams) (interface{}, error) {
	user := p.Context.Value("user").(ConnectedUser)
	connID := p.Context.Value("connID").(string)
	channelName := "newComment:" + p.Args["roomId"].(string)
	c.subscribeChan <- gss.NewSubscribeEvent(channelName, connID, user.Name())
	return newDummyResponse(), nil
}

func (c *Comment) OnUnsubscribe(p graphql.ResolveParams) (interface{}, error) {
	user := p.Context.Value("user").(ConnectedUser)
	connID := p.Context.Value("connID").(string)
	c.unsubscribeChan <- gss.NewUnsubscribeEvent(connID, user.Name())
	return newDummyResponse(), nil
}

//
// Resolver section end -->
//
