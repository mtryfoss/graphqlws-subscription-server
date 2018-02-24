package main

import (
	"github.com/graphql-go/graphql"
	gss "github.com/taiyoh/graphqlws-subscription-server"
)

type SampleComment struct {
	id      string `json:"id"`
	content string `json:"content"`
}

func newDummyResponse() *SampleComment {
	return &SampleComment{id: "id", content: "content"}
}

type Comment struct {
	gss.GraphQLType
	fieldName string
}

func NewComment() *Comment {
	return &Comment{fieldName: "newComment"}
}

func (c *Comment) FieldName() string {
	return c.fieldName
}

func (c *Comment) GetType() graphql.ObjectConfig {
	return graphql.ObjectConfig{
		Name: "Comment",
		Fields: graphql.Fields{
			"id":      &graphql.Field{Type: graphql.NewNonNull(graphql.ID)},
			"content": &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		},
	}
}

func (c *Comment) GetArgs() map[string]*graphql.ArgumentConfig {
	return map[string]*graphql.ArgumentConfig{
		"roomId": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.ID)},
	}
}

func (c *Comment) OnPayload(payload interface{}, p graphql.ResolveParams) (interface{}, error) {
	comment := payload.(SampleComment)
	return comment, nil
}

func (c *Comment) OnSubscribe(listener *gss.Listener, p graphql.ResolveParams) (interface{}, error) {
	user := p.Context.Value("user").(ConnectedUser)
	connID := p.Context.Value("connID").(string)
	channelName := c.FieldName() + ":" + p.Args["roomId"].(string)
	listener.Subscribe(channelName, connID, user.Name())
	return newDummyResponse(), nil
}

func (c *Comment) OnUnsubscribe(listener *gss.Listener, p graphql.ResolveParams) (interface{}, error) {
	user := p.Context.Value("user").(ConnectedUser)
	connID := p.Context.Value("connID").(string)
	listener.Unsubscribe(connID, user.Name())
	return newDummyResponse(), nil
}
