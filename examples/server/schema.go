package main

import (
	"errors"

	"github.com/graphql-go/graphql"
	gss "github.com/taiyoh/graphqlws-subscription-server"
)

type SampleComment struct {
	id      string `json:"id"`
	content string `json:"content"`
}

// this is sample code.
func LoadSchema(listener *gss.Listener) (graphql.Schema, error) {
	commentType := graphql.ObjectConfig{
		Name: "Comment",
		Fields: graphql.Fields{
			"id":      &graphql.Field{Type: graphql.NewNonNull(graphql.ID)},
			"content": &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		},
	}
	fields := graphql.Fields{
		"newComment": &graphql.Field{
			Type: graphql.NewObject(commentType),
			Args: graphql.FieldConfigArgument{
				"roomId": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.ID)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if payload := p.Context.Value("payload"); payload != nil { // payload exists
					comment := payload.(SampleComment)
					return comment, nil
				}
				if s := p.Context.Value("onSubscribe"); s != nil { // AddSubscription called
					user := p.Context.Value("user").(ConnectedUser)
					connID := string(p.Context.Value("connID"))
					channelName := "newComment:" + string(p.Args["roomId"])
					listener.Subscribe(channelName, connID, user.Name())
					dummyComment := &SampleComment{"ping", "ping"}
					return dummyComment, nil
				}
				if s := p.Context.Value("onUnsubscribe"); s != nil { // removeSubscription called
					user := p.Context.Value("user").(ConnectedUser)
					connID := string(p.Context.Value("connID"))
					listener.Unsubscribe(connID, user.Name())
					dummyComment := &SampleComment{"ping", "ping"}
					return dummyComment, nil
				}
				return nil, errors.New("no payload exists")
			},
		},
	}
	rootSubscription := graphql.ObjectConfig{Name: "RootSubscription", Fields: fields}
	schemaConfig := graphql.SchemaConfig{Subscription: graphql.NewObject(rootSubscription)}
	return graphql.NewSchema(schemaConfig)
}
