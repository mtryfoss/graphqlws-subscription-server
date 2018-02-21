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

type TypeResolver interface {
	func OnPayload(interface{}, context.Context) (interface{}, error)
	func OnSubscribe(context.Context, *gss.Listener) (interface{}, error)
	func OnUnsubscribe(context.Context, *gss.Listener) (interface{}, error)
	func GetField(*gss.Listener) (string, *graphql.Field)
}

type typeResolver struct {
	TypeResolver
}

func (r *typeResolver) GetResolve(listener, *gss.Listener) func(p graphql.ResolveParams) (interface{}, error) {
	return func(p graphql.ResolveParams) (interface{}, error) {
		if payload := p.Context.Value("payload"); payload != nil { // payload exists
			return r.OnPayload(payload, ctx)
		}
		if s := p.Context.Value("onSubscribe"); s != nil { // AddSubscription called
			return r.OnSubscribe(ctx, listener)
		}
		if s := p.Context.Value("onUnsubscribe"); s != nil { // removeSubscription called
			return r.OnUnubscribe(ctx, listener)
		}
		return nil, errors.New("no payload exists")
	}
}

type CommentResolver struct {
	typeResolver
}

type NewCommentResolver() *CommentResolver {
	return &CommentResolver{}
}

func (r *CommentResolver) GetField(listener *gss.Listener) (string, *graphql.Field) {
	t := graphql.ObjectConfig{
		Name: "Comment",
		Fields: graphql.Fields{
			"id":      &graphql.Field{Type: graphql.NewNonNull(graphql.ID)},
			"content": &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		},
	}
	field := &graphql.Field{
		Type: graphql.NewObject(t),
		Args: graphql.FieldConfigArgument{
			"roomId": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.ID)},
		},
		Resolve: r.GetResolve(listener),
	}
	return "newComment", field
}

func (r *CommentResolver) OnPayload(payload interface{}, ctx context.Context) (interface{}, error) {
	comment := payload.(SampleComment)
	return comment, nil
}

func (r *CommentResolver) OnSubscribe(ctx context.Context, listener *gss.Listener) (interface{}, error) {
	user := p.Context.Value("user").(ConnectedUser)
	connID := string(p.Context.Value("connID"))
	channelName := "newComment:" + string(p.Args["roomId"])
	listener.Subscribe(channelName, connID, user.Name())
	dummyComment := &SampleComment{"ping", "ping"}
	return dummyComment, nil
}

func (r *CommentResolver) OnUnsubscribe(ctx context.Context, listener *gss.Listener) (interface{}, error) {
	user := p.Context.Value("user").(ConnectedUser)
	connID := string(p.Context.Value("connID"))
	listener.Unsubscribe(connID, user.Name())
	dummyComment := &SampleComment{"ping", "ping"}
	return dummyComment, nil
}
