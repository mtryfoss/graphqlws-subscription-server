package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/functionalfoundry/graphqlws"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
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

	schema, err := getSchema()
	if err != nil {
		log.Fatalln("GraphQL schema is invalid: ", err)
	}

	canSendToUser := func(conn *graphqlws.Connection, reqData *gss.RequestData) bool {
		user := (*conn).User().(ConnectedUser)
		for _, userName := range reqData.Users {
			if userName == user.Name() {
				return true
			}
		}
		return false
	}

	subService := gss.NewSubscribeService(schema, canSendToUser)

	ctx := context.Background()
	wg := &sync.WaitGroup{}

	notifyChan := make(chan *gss.RequestData, conf.Server.MaxHandlerCount)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case data := <-notifyChan:
				go subService.Publish(data)
			}
		}
	}()

	authCallback := AuthenticateCallback(conf.Auth.SecretKey)

	mux := http.NewServeMux()

	d, _ := os.Getwd()
	mux.Handle("/", http.FileServer(http.Dir(filepath.Join(d, "build"))))

	mux.Handle("/graphql", handler.New(&handler.Config{
		Schema:   schema,
		Pretty:   true,
		GraphiQL: true,
	}))
	mux.Handle("/subscription", subService.NewSubscriptionHandler(authCallback))
	mux.Handle("/notify", gss.NewNotifyHandler(notifyChan))

	server := &http.Server{
		Addr:    ":" + strconv.Itoa(int(conf.Server.Port)),
		Handler: mux,
	}
	startServer(server, ctx, wg)

	wg.Wait()
}

//
// <!-- Authenticate section start
//

type ConnectedUser struct {
	ID string
}

func (u ConnectedUser) Name() string {
	return u.ID
}

func AuthenticateCallback(secretkey string) graphqlws.AuthenticateFunc {
	return func(tokenstring string) (interface{}, error) {
		return ConnectedUser{ID: tokenstring}, nil
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
// <!-- Schema section start
//

func getSchema() (*graphql.Schema, error) {
	resolve := func(p graphql.ResolveParams) (interface{}, error) {
		payload := p.Context.Value(gss.GraphQLContextKey("payload"))
		if payload == nil {
			return nil, errors.New("payload not found")
		}
		return payload, nil
	}
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(
			graphql.ObjectConfig{
				Name: "RootQuery",
				Fields: graphql.Fields{
					"hello": &graphql.Field{
						Type: graphql.NewNonNull(graphql.String),
						Resolve: func(p graphql.ResolveParams) (interface{}, error) {
							return "world", nil
						},
					},
				},
			},
		),
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
						Resolve: resolve,
					},
					"notification": &graphql.Field{
						Type: graphql.NewObject(graphql.ObjectConfig{
							Name: "Notification",
							Fields: graphql.Fields{
								"id":      &graphql.Field{Type: graphql.NewNonNull(graphql.ID)},
								"content": &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
							},
						}),
						Resolve: resolve,
					},
				},
			},
		),
	})
	return &schema, err
}

//
// Schema section end -->
//

//
// <!-- Server section start
//

func startServer(srv *http.Server, ctx context.Context, wg *sync.WaitGroup) {
	log.Println("Starting subscription server on " + srv.Addr)

	wg.Add(2)
	go func() {
		defer wg.Done()
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				if err := srv.Shutdown(ctx); err != nil {
					log.Fatal(err)
				}
				return
			}
		}
	}()
}

//
// Server section end -->
//
