package gss

import (
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/kinds"
)

type SubscriptionIDByConnectionID map[string]string
type QueryArgsMap map[string]string

type SubscribeFilter interface {
	RegisterConnectionIDFromDocument(connID string, subID string, doc *ast.Document, variables map[string]interface{})
	RemoveSubscriptionIDFromConnectionID(connID, subID string)
	RemoveConnectionIDFromChannels(connID string)
	GetChannelRegisteredConnectionIDs(channel string) SubscriptionIDByConnectionID
}

type ChannelSerializer interface {
	Serialize(field string, args QueryArgsMap) string
}

type channelSerializerFunc func(field string, args QueryArgsMap) string

func (f channelSerializerFunc) Serialize(field string, args QueryArgsMap) string {
	return f(field, args)
}

func getNewChannelSerializerFunc() channelSerializerFunc {
	return func(field string, args QueryArgsMap) string {
		sargs := []string{}
		for k := range args {
			sargs = append(sargs, k)
		}
		sort.Slice(sargs, func(i, j int) bool {
			return sargs[i] <= sargs[j]
		})
		strList := []string{field}
		for i := range sargs {
			strList = append(strList, args[sargs[i]])
		}
		return strings.Join(strList, ":")
	}
}

type subscribeFilter struct {
	SubscribeFilter
	Serializer            ChannelSerializer
	ConnectionIDByChannel map[string]*sync.Map
	ChannelByConnectionID map[string]*sync.Map
}

func NewSubscribeFilter() *subscribeFilter {
	return &subscribeFilter{
		Serializer:            getNewChannelSerializerFunc(),
		ConnectionIDByChannel: map[string]*sync.Map{},
		ChannelByConnectionID: map[string]*sync.Map{},
	}
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

func ifToStr(d interface{}) string {
	if v, ok := d.(string); ok {
		return v
	}
	if v, ok := d.(int); ok {
		return strconv.Itoa(v)
	}
	return ""
}

func getArgKeyValueFromAstValue(variables map[string]interface{}, arg *ast.Argument) (string, string, bool) {
	var k, v string
	val := arg.Value
	if val.GetKind() == kinds.Variable {
		n := val.GetValue().(*ast.Name)
		k = n.Value
		vv, ok := variables[n.Value]
		if !ok {
			return "", "", false
		}
		v = ifToStr(vv)
	} else {
		k = arg.Name.Value
		vv := arg.Value.GetValue().(string)
		v = ifToStr(vv)
	}
	if v != "" {
		return k, v, true
	}
	return "", "", false
}

func channelsForSelectionSets(variables map[string]interface{}, sets []*ast.SelectionSet) map[string]map[string]string {
	nameList := map[string]map[string]string{}
	for _, set := range sets {
		if len(set.Selections) < 1 {
			continue
		}
		field := set.Selections[0].(*ast.Field)
		args := map[string]string{}
		for _, arg := range field.Arguments {
			if k, v, ok := getArgKeyValueFromAstValue(variables, arg); ok {
				args[k] = v
			}
		}
		if len(args) > 0 {
			nameList[field.Name.Value] = args
		}
	}
	return nameList
}

func (f *subscribeFilter) RegisterConnectionIDFromDocument(connID string, subID string, doc *ast.Document, variables map[string]interface{}) {
	defs := operationDefinitionsWithOperation(doc, "subscription")
	sets := selectionSetsForOperationDefinitions(defs)
	for field, args := range channelsForSelectionSets(variables, sets) {
		ch := f.Serializer.Serialize(field, args)
		if m, ok := f.ConnectionIDByChannel[ch]; ok {
			m.Store(connID, subID)
		} else {
			m := &sync.Map{}
			m.Store(connID, subID)
			f.ConnectionIDByChannel[ch] = m
		}
		if m, ok := f.ChannelByConnectionID[connID]; ok {
			m.Store(ch, subID)
		} else {
			m := &sync.Map{}
			m.Store(ch, subID)
			f.ChannelByConnectionID[connID] = m
		}
	}
}

func (f *subscribeFilter) RemoveConnectionIDFromChannels(connID string) {
	channels := []string{}
	if m1, ok := f.ChannelByConnectionID[connID]; ok {
		m1.Range(func(k, v interface{}) bool {
			channels = append(channels, k.(string))
			return true
		})
		delete(f.ChannelByConnectionID, connID)
	}
	for _, ch := range channels {
		if m, ok := f.ConnectionIDByChannel[ch]; ok {
			m.Delete(connID)
		}
	}
}

func (f *subscribeFilter) RemoveSubscriptionIDFromConnectionID(connID, subID string) {
	var ch string
	m1, ok := f.ChannelByConnectionID[connID]
	if !ok {
		return
	}
	m1.Range(func(k, v interface{}) bool {
		if v.(string) == subID {
			ch = k.(string)
		}
		return ch == ""
	})
	if ch == "" {
		return
	}
	m1.Delete(ch)
	m2, ok := f.ConnectionIDByChannel[ch]
	if !ok {
		return
	}
	if s, ok := m2.Load(connID); ok && s.(string) == subID {
		m2.Delete(connID)
	}
}

func (f *subscribeFilter) GetChannelRegisteredConnectionIDs(channel string) SubscriptionIDByConnectionID {
	founds := SubscriptionIDByConnectionID{}
	if m, ok := f.ConnectionIDByChannel[channel]; ok {
		m.Range(func(k, v interface{}) bool {
			founds[k.(string)] = v.(string)
			return true
		})
	}
	return founds
}
