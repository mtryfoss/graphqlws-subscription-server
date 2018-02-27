package gss

import (
	"testing"

	"github.com/functionalfoundry/graphqlws"
	"github.com/graphql-go/graphql/language/parser"
)

func TestChannelManagerSubscribeAndUnsubscribe(t *testing.T) {

	user := map[string]string{}
	user["foo"] = "world"

	// Query
	query := `
		subscription mySubscribe($id: ID!, $aaa: String!) {
			hello(id: $id, aaa: $aaa) {
				foo
				bar
			}
		}
	`

	conn1 := &connForTest{
		id:   "hoge",
		user: user,
	}

	sub1 := &graphqlws.Subscription{
		ID:    "foo",
		Query: query,
		Variables: map[string]interface{}{
			"id":  1,
			"aaa": "fuu",
		},
		OperationName: "",
		Connection:    conn1,
		SendData:      func(d *graphqlws.DataMessagePayload) {},
	}

	document, _ := parser.Parse(parser.ParseParams{
		Source: query,
	})

	sub1.Document = document

	f := NewSubscribeFilter()
	f.ReplaceFieldsFromDocument(sub1)

	if len(sub1.Fields) != 1 {
		t.Error("subscription.Fields count should be 1")
	}
	if sub1.Fields[0] != "hello:fuu:1" {
		t.Error("subscription.Fields[0] should hello:fuu:1")
	}
}
