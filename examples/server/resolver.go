package main

import (
	"context"
	"errors"

	"github.com/graphql-go/graphql"
	gss "github.com/taiyoh/graphqlws-subscription-server"
)

type SampleComment struct {
	id      string `json:"id"`
	content string `json:"content"`
}

type Comment struct {
	GraphQLTypeImpl	
	fieldName: string
}

type NewComment() *CommentResolver {
	return &Comment{fieldName: "newComment"}
}

func (c *Comment) GetField(listener *gss.Listener) *graphql.Field {
	t := graphql.ObjectConfig{
		Name: "Comment",
		Fields: graphql.Fields{
			"id":      &graphql.Field{Type: graphql.NewNonNull(graphql.ID)},
			"content": &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		},
	}
	return &graphql.Field{
		Type: graphql.NewObject(t),
		Args: graphql.FieldConfigArgument{
			"roomId": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.ID)},
		},
		Resolve: r.GetResolve(listener),
	}
}

func (c *Comment) OnPayload(payload interface{}, ctx context.Context) (interface{}, error) {
	comment := payload.(SampleComment)
	return comment, nil
}

func (c *Comment) OnSubscribe(ctx context.Context, listener *gss.Listener) (interface{}, error) {
	user := p.Context.Value("user").(ConnectedUser)
	connID := string(p.Context.Value("connID"))
	channelName := r.fieldName + ":" + string(p.Args["roomId"])
	listener.Subscribe(channelName, connID, user.Name())
	dummyComment := &SampleComment{"ping", "ping"}
	return dummyComment, nil
}

func (c *Comment) OnUnsubscribe(ctx context.Context, listener *gss.Listener) (interface{}, error) {
	user := p.Context.Value("user").(ConnectedUser)
	connID := string(p.Context.Value("connID"))
	listener.Unsubscribe(connID, user.Name())
	dummyComment := &SampleComment{"ping", "ping"}
	return dummyComment, nil
}
