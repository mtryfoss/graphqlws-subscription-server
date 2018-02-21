package main

import (
	"github.com/graphql-go/graphql"
	gss "github.com/taiyoh/graphqlws-subscription-server"
)

type SampleComment struct {
	id      string `json:"id"`
	content string `json:"content"`
}

type Comment struct {
	gss.GraphQLTypeImpl
	fieldName string
}

func NewComment() *Comment {
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
		Resolve: c.GetResolve(listener),
	}
}

func (c *Comment) OnPayload(payload interface{}, p graphql.ResolveParams) (interface{}, error) {
	comment := payload.(SampleComment)
	return comment, nil
}

func (c *Comment) OnSubscribe(p graphql.ResolveParams, listener *gss.Listener) (interface{}, error) {
	user := p.Context.Value("user").(ConnectedUser)
	connID := p.Context.Value("connID").(string)
	channelName := c.FieldName() + ":" + p.Args["roomId"].(string)
	listener.Subscribe(channelName, connID, user.Name())
	dummyComment := &SampleComment{"ping", "ping"}
	return dummyComment, nil
}

func (c *Comment) OnUnsubscribe(p graphql.ResolveParams, listener *gss.Listener) (interface{}, error) {
	user := p.Context.Value("user").(ConnectedUser)
	connID := p.Context.Value("connID").(string)
	listener.Unsubscribe(connID, user.Name())
	dummyComment := &SampleComment{"ping", "ping"}
	return dummyComment, nil
}
