package gss

import (
	"io"
	"testing"

	"github.com/functionalfoundry/graphqlws"
	"github.com/graphql-go/graphql"
)

type connForTest struct {
	graphqlws.Connection
	id              string
	user            interface{}
	ReceivedOpID    string
	ReceivedPayload *graphqlws.DataMessagePayload
	ReceivedError   error
}

func (c *connForTest) ID() string {
	return c.id
}

func (c *connForTest) User() interface{} {
	return c.user
}

func (c *connForTest) SendData(opID string, d *graphqlws.DataMessagePayload) {
	c.ReceivedOpID = opID
	c.ReceivedPayload = d
}

func (c *connForTest) SendError(e error) {
	c.ReceivedError = e
}

func TestNewSubscribeService(t *testing.T) {
	user := map[string]string{}
	user["foo"] = "world"

	schema := buildSchema()

	// Query
	query := `
		{
			hello
		}
	`
	f := func(conn *graphqlws.Connection, reqData *RequestData) bool {
		return true
	}
	s := NewSubscribeService(schema, f)

	conn1 := &connForTest{
		id:   "hoge",
		user: user,
	}

	sub1 := &graphqlws.Subscription{
		ID:            "foo",
		Query:         query,
		Variables:     map[string]interface{}{},
		OperationName: "",
		Fields:        []string{},
		Connection:    conn1,
		SendData:      func(d *graphqlws.DataMessagePayload) {},
	}

	errs := s.AddSubscription(conn1, sub1)
	if len(errs) > 0 {
		t.Error("errors should not exists")
	}

	subs := s.Subscriptions()
	if _, exists := subs[conn1]; !exists {
		t.Error("conn1 should exists")
	}

	s.RemoveSubscription(conn1, sub1)
	subs = s.Subscriptions()
	if _, exists := subs[conn1]; exists {
		t.Error("conn1 should not exists")
	}

	s.AddSubscription(conn1, sub1)
	s.RemoveSubscriptions(conn1)

	subs = s.Subscriptions()
	if _, exists := subs[conn1]; exists {
		t.Error("conn1 should not exists")
	}

	s.NewSubscriptionHandler(func(tokenstring string) (interface{}, error) {
		return user, nil
	})
}

type testUser struct {
	ID            string
	JoinedChannel string
	Payloads      []*graphqlws.DataMessagePayload
}

func (u *testUser) GetSendData() func(*graphqlws.DataMessagePayload) {
	return func(d *graphqlws.DataMessagePayload) {
		u.Payloads = append(u.Payloads, d)
	}
}

type testSampleComment struct {
	ID      string `json:"id"`
	Content string `json:"content"`
}

type testNotification struct {
	Content string `json:"content"`
}

type testTmpleView struct {
	io.Writer
	Captured []byte
}

func (v *testTmpleView) Write(p []byte) (int, error) {
	v.Captured = append(v.Captured, p...)
	return len(p), nil
}

func TestSubscribeServicePublish(t *testing.T) {
	user1 := &testUser{ID: "test1", JoinedChannel: "foo"}
	user2 := &testUser{ID: "test2", JoinedChannel: "bar"}
	user3 := &testUser{ID: "test3", JoinedChannel: "baz"}
	user4 := &testUser{ID: "test4", JoinedChannel: "foo"}
	user5 := &testUser{ID: "test5", JoinedChannel: "foo"}

	schema := buildSchema()

	subService := NewSubscribeService(schema, func(conn *graphqlws.Connection, reqData *RequestData) bool {
		user := (*conn).User().(testUser)
		for _, userID := range reqData.Users {
			if user.ID == userID {
				return true
			}
		}
		return false
	})

	// Query
	query := `
subscription mySubscribe($commentId: ID!) {
	newComment(id: $commentId) {
		id content
	}
	notification {
		content
	}
}
`

	for _, user := range []*testUser{user1, user2, user3, user4, user5} {
		sub := &graphqlws.Subscription{
			ID:    user.ID + "-sub-" + user.JoinedChannel,
			Query: query,
			Variables: map[string]interface{}{
				"commentId": user.JoinedChannel,
			},
			Connection: &connForTest{
				id:   user.ID + "-conn",
				user: user,
			},
			SendData: user.GetSendData(),
		}
		subService.AddSubscription(sub.Connection, sub)
	}

	clearPayloads := func() {
		for _, user := range []*testUser{user1, user2, user3, user4, user5} {
			user.Payloads = []*graphqlws.DataMessagePayload{}
		}
	}

	subService.Publish(&RequestData{
		Channel: "newComment:foo",
		Payload: testSampleComment{ID: "id1", Content: "TestSend1"},
	})
	for _, user := range []*testUser{user1, user4, user5} {
		if len(user.Payloads) != 1 {
			t.Error("user.Payloads count should be 1")
		}
		d := user.Payloads[0].Data.(map[string]interface{})
		if d["newComment"] == nil {
			t.Error("newComment should exists")
		}
		if _, exists := d["notification"]; exists {
			t.Error("notification shoud not exists")
		}

	}
	for _, user := range []*testUser{user2, user3} {
		if len(user.Payloads) != 0 {
			t.Error("user.Payloads count should be 0")
		}
	}

	clearPayloads()

	subService.Publish(&RequestData{
		Channel: "notification",
		Payload: testNotification{Content: "hogefuga"},
	})

	for _, user := range []*testUser{user1, user2, user3, user4, user5} {
		if len(user.Payloads) != 1 {
			t.Error("user.Payloads count should be 1")
		}
	}

}

func buildSchema() *graphql.Schema {
	rootQuery := graphql.ObjectConfig{
		Name: "RootQuery",
		Fields: graphql.Fields{
			"hello": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return "world", nil
				},
			},
		},
	}
	rootSubscription := graphql.ObjectConfig{
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
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.ID),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return p.Context.Value(GraphQLContextKey("newComment")), nil
				},
			},
			"notification": &graphql.Field{
				Type: graphql.NewObject(graphql.ObjectConfig{
					Name: "Notification",
					Fields: graphql.Fields{
						"content": &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
					},
				}),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return p.Context.Value(GraphQLContextKey("notification")), nil
				},
			},
		},
	}
	schemaConfig := graphql.SchemaConfig{
		Query:        graphql.NewObject(rootQuery),
		Subscription: graphql.NewObject(rootSubscription),
	}
	schema, _ := graphql.NewSchema(schemaConfig)
	return &schema
}
