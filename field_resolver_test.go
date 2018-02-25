package gss

import (
	"context"
	"testing"

	"github.com/graphql-go/graphql"
)

type GraphQLResolveTest1 struct {
	GraphQLResolve
	OnPayloadCalled     bool
	OnSubscribeCalled   bool
	OnUnsubscribeCalled bool
}

func (t *GraphQLResolveTest1) OnPayload(payload interface{}, p graphql.ResolveParams) (interface{}, error) {
	t.OnPayloadCalled = true
	return payload, nil
}

func (t *GraphQLResolveTest1) OnSubscribe(p graphql.ResolveParams) (interface{}, error) {
	t.OnSubscribeCalled = true
	return p.Context.Value(GraphQLContextKey("onSubscribe")), nil
}

func (t *GraphQLResolveTest1) OnUnsubscribe(p graphql.ResolveParams) (interface{}, error) {
	t.OnUnsubscribeCalled = true
	return p.Context.Value(GraphQLContextKey("onUnsubscribe")), nil
}

type TestBuildResolveCase struct {
	Label               string
	ContextKey          interface{}
	ContextVal          interface{}
	OnPayloadCalled     bool
	OnSubscribeCalled   bool
	OnUnsubscribeCalled bool
}

func TestBuildResolve(t *testing.T) {
	ctx := context.Background()
	p := graphql.ResolveParams{}
	seed := &GraphQLResolveTest1{}
	resolve := BuildResolve(seed)

	p.Context = ctx
	res, err := resolve(p)
	if res != nil {
		t.Error("response shoud not exists")
	}
	if err == nil {
		t.Error("error should exists")
	}
	if err.Error() != "no payload exists" {
		t.Error("unexpected error message")
	}

	cases := []TestBuildResolveCase{
		TestBuildResolveCase{
			Label:               "onSubscribe",
			ContextKey:          GraphQLContextKey("onSubscribe"),
			ContextVal:          "a",
			OnPayloadCalled:     false,
			OnSubscribeCalled:   true,
			OnUnsubscribeCalled: false,
		},
		TestBuildResolveCase{
			Label:               "onUnsubscribe",
			ContextKey:          GraphQLContextKey("onUnsubscribe"),
			ContextVal:          "b",
			OnPayloadCalled:     false,
			OnSubscribeCalled:   false,
			OnUnsubscribeCalled: true,
		},
		TestBuildResolveCase{
			Label:               "onPayload",
			ContextKey:          GraphQLContextKey("payload"),
			ContextVal:          "foobar",
			OnPayloadCalled:     true,
			OnSubscribeCalled:   false,
			OnUnsubscribeCalled: false,
		},
	}

	for _, testCase := range cases {
		p := graphql.ResolveParams{
			Context: context.WithValue(context.Background(), testCase.ContextKey, testCase.ContextVal),
		}
		seed.OnPayloadCalled = false
		seed.OnSubscribeCalled = false
		seed.OnUnsubscribeCalled = false
		ret, err := resolve(p)
		if err != nil {
			t.Error(testCase.Label + ": error should not exists. msg: " + err.Error())
		}
		if ret == nil {
			t.Error(testCase.Label + ": return value should exists")
		}
		if ret.(string) != testCase.ContextVal.(string) {
			t.Error(testCase.Label + ": return value should same as '" + testCase.ContextVal.(string) + "'")
		}
		if seed.OnPayloadCalled != testCase.OnPayloadCalled {
			v := "false"
			if testCase.OnPayloadCalled {
				v = "true"
			}
			t.Error(testCase.Label + ": OnPayloadCalled should " + v)
		}
		if seed.OnSubscribeCalled != testCase.OnSubscribeCalled {
			v := "false"
			if testCase.OnSubscribeCalled {
				v = "true"
			}
			t.Error(testCase.Label + ": OnSubscriptionCalled should " + v)
		}
		if seed.OnUnsubscribeCalled != testCase.OnUnsubscribeCalled {
			v := "false"
			if testCase.OnUnsubscribeCalled {
				v = "true"
			}
			t.Error(testCase.Label + ": OnUnsubscribeCalled should " + v)
		}
	}

}
