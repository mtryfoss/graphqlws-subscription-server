package gss

import (
	"sort"
	"strings"

	"github.com/functionalfoundry/graphqlws"
	"github.com/graphql-go/graphql/language/ast"
)

type SubscribeFilter interface {
	ReplaceFieldsFromDocument(subscription *graphqlws.Subscription)
}

type subscribeFilter struct {
	SubscribeFilter
}

func NewSubscribeFilter() *subscribeFilter {
	return &subscribeFilter{}
}

type astArgs struct {
	Key string
	Val string
}

func operationDefinitionsWithOperation(
	doc *ast.Document,
	op string,
) []*ast.OperationDefinition {
	defs := []*ast.OperationDefinition{}
	for _, node := range doc.Definitions {
		if node.GetKind() == "OperationDefinition" {
			if def, ok := node.(*ast.OperationDefinition); ok {
				if def.Operation == op {
					defs = append(defs, def)
				}
			}
		}
	}
	return defs
}

func selectionSetsForOperationDefinitions(
	defs []*ast.OperationDefinition,
) []*ast.SelectionSet {
	sets := []*ast.SelectionSet{}
	for _, def := range defs {
		if set := def.GetSelectionSet(); set != nil {
			sets = append(sets, set)
		}
	}
	return sets
}

func nameForSelectionSet(set *ast.SelectionSet) (string, bool) {
	if len(set.Selections) >= 1 {
		if field, ok := set.Selections[0].(*ast.Field); ok {
			args := []astArgs{}
			for _, arg := range field.Arguments {
				args = append(args, astArgs{
					Key: arg.Name.Value,
					Val: arg.Value.GetValue().(string),
				})
			}
			sort.Slice(args, func(i, j int) bool {
				return args[i].Key <= args[j].Key
			})
			joinedArgs := []string{field.Name.Value}
			for _, a := range args {
				joinedArgs = append(joinedArgs, a.Val)
			}
			return strings.Join(joinedArgs, ":"), true
		}
	}
	return "", false
}

func namesForSelectionSets(sets []*ast.SelectionSet) []string {
	names := []string{}
	for _, set := range sets {
		if name, ok := nameForSelectionSet(set); ok {
			names = append(names, name)
		}
	}
	return names
}

func (f *subscribeFilter) ReplaceFieldsFromDocument(subscription *graphqlws.Subscription) {
	defs := operationDefinitionsWithOperation(subscription.Document, "subscription")
	sets := selectionSetsForOperationDefinitions(defs)
	fields := namesForSelectionSets(sets)
	subscription.Fields = fields
}
