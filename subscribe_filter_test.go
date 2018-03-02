package gss

import (
	"testing"

	"github.com/functionalfoundry/graphqlws"
	"github.com/graphql-go/graphql/language/parser"
)

func TestSubscriptionFilter(t *testing.T) {
	user := map[string]string{}
	user["foo"] = "world"

	conn1 := &connForTest{
		id:   "hoge",
		user: user,
	}

	sub1 := &graphqlws.Subscription{
		ID: "foo",
		Query: `
		subscription f($id:ID!, $aaa:String!){
			hello(id: $id, aaa: $aaa) {
				foo
				bar
			}
		}
		`,
		Variables: map[string]interface{}{
			"id":  1,
			"aaa": "fuu",
		},
		OperationName: "",
		Connection:    conn1,
		SendData:      func(d *graphqlws.DataMessagePayload) {},
	}

	sub2 := &graphqlws.Subscription{
		ID: "bar",
		Query: `
		subscription {
			hoge(bbb: "ccc") {
				fuga
				piyo
			}
		}
		`,
		Variables:     map[string]interface{}{},
		OperationName: "",
		Connection:    conn1,
		SendData:      func(d *graphqlws.DataMessagePayload) {},
	}

	doc1, _ := parser.Parse(parser.ParseParams{
		Source: sub1.Query,
	})
	doc2, _ := parser.Parse(parser.ParseParams{
		Source: sub2.Query,
	})

	sub1.Document = doc1
	sub2.Document = doc2
	f := NewSubscribeFilter()

	t.Run("register and get", func(t *testing.T) {
		f.RegisterConnectionIDFromDocument(conn1.ID(), sub1.ID, sub1.Document, sub1.Variables)
		f.RegisterConnectionIDFromDocument(conn1.ID(), sub2.ID, sub2.Document, sub2.Variables)

		if len(f.ChannelByConnectionID) != 1 {
			t.Error("registered conn.ID count should be 1")
		}
		subcount := 0
		f.ChannelByConnectionID[conn1.ID()].Range(func(k, v interface{}) bool {
			subcount++
			return true
		})
		if subcount != 2 {
			t.Error("registered subscription count should be 2")
		}

		idmap := f.GetChannelRegisteredConnectionIDs("hello:fuu:1")
		if len(idmap) != 1 {
			t.Error("idmap keys count should be 1")
		}
		if subID, ok := idmap[conn1.ID()]; !ok {
			t.Error("idmap key has conn1.ID")
		} else if subID != sub1.ID {
			t.Error("subscription.id should be sub1.ID")
		}
		idmap = f.GetChannelRegisteredConnectionIDs("hoge:ccc")
		if len(idmap) != 1 {
			t.Error("idmap keys count should be 1")
		}
		if subID, ok := idmap[conn1.ID()]; !ok {
			t.Error("idmap key has conn1.ID")
		} else if subID != sub2.ID {
			t.Error("subscription.id should be sub2.ID")
		}
	})

	t.Run("remove", func(t *testing.T) {
		f.RemoveSubscriptionIDFromConnectionID(conn1.ID(), sub2.ID)
		idmap := f.GetChannelRegisteredConnectionIDs("hoge:ccc")
		if len(idmap) != 0 {
			t.Error("idmap keys count should be 0")
		}
		f.RemoveConnectionIDFromChannels(conn1.ID())
		idmap = f.GetChannelRegisteredConnectionIDs("hello:fuu:1")
		if len(idmap) != 0 {
			t.Error("idmap keys count should be 0")
		}
	})
}
