package gss

import (
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

	// Schema (via https://github.com/graphql-go/graphql )
	fields := graphql.Fields{
		"hello": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				v, _ := p.Context.Value(GraphQLContextKey("user")).(map[string]string)["foo"]
				return v, nil
			},
		},
	}
	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
	schema, _ := graphql.NewSchema(schemaConfig)

	// Query
	query := `
		{
			hello
		}
	`
	f := func(conn *graphqlws.Connection, reqData *RequestData) bool {
		return true
	}
	s := NewSubscribeService(&schema, 10, f)
	s.GetNotifierChan()

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
	s.NewNotifyHandler()
}
