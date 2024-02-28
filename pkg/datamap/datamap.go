package datamap

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"golang.org/x/exp/maps"
)

const (
	tagOffsetOf = "offsetof"
	tagSizeOf   = "sizeof"
)

type Operation int

func (o Operation) String() string {
	switch o {
	case OpOffsetOf:
		return "OffsetOf"
	case OpSizeOf:
		return "SizeOf"
	default:
		return "Unknown"
	}
}

func (o Operation) minimumRequiredRouteLength() int {
	switch o {
	case OpOffsetOf:
		return 2
	case OpSizeOf:
		return 1
	default:
		return -1
	}
}

const (
	OpOffsetOf Operation = iota
	OpSizeOf
)

type DataMap struct {
	Routes []*RouteNode
}

type Extractor struct {
	Source string

	Op          Operation
	targetValue *reflect.Value
}

func (d *Extractor) Set(value int64) error {
	if !d.targetValue.CanSet() {
		return fmt.Errorf("field from struct %s is not settable", d.targetValue.Type().Name())
	}
	switch d.targetValue.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		d.targetValue.SetInt(value)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		d.targetValue.SetUint(uint64(value))
	default:
		return fmt.Errorf("field from struct %s is not of type int or uint, type: %s", d.targetValue.Type().Name(), d.targetValue.Kind())
	}
	return nil
}

type RouteNode struct {
	prev *RouteNode
	Next *RouteNode

	Type       string
	Extractors []*Extractor
}

func (rn *RouteNode) IsLeaf() bool {
	return rn.Next == nil
}

func (rn *RouteNode) Leaf() *RouteNode {
	curr := rn
	for curr.Next != nil {
		curr = curr.Next
	}
	return curr
}

func (rn *RouteNode) path() []*RouteNode {
	var (
		path = []*RouteNode{rn}
		curr = rn
	)
	for curr.Next != nil {
		curr = curr.Next
		path = append(path, curr)
	}
	return path
}

func (rn *RouteNode) Key() string {
	var (
		parts = rn.path()
		key   = ""
	)
	for _, p := range parts {
		key += p.Type + "."
	}
	return key
}

func newRouteFromTagValue(path string) *RouteNode {
	parts := strings.Split(path, ".")
	var (
		head = &RouteNode{Type: parts[0]}
		curr = head
	)
	for _, p := range parts[1:] {
		curr.Next = &RouteNode{Type: p, prev: curr}
		curr = curr.Next
	}
	return head
}

// New generates a DataMap from the given struct.
// Given argument must be a pointer to a struct.
// The struct fields must be tagged with `offsetof` or `sizeof` tags.
// The pointer is needed to be able to set the fields.
func New(layoutMap any) (*DataMap, error) {
	if layoutMap == nil {
		return nil, errors.New("layoutMap must not be nil")
	}
	if reflect.TypeOf(layoutMap).Kind() != reflect.Ptr {
		return nil, errors.New("layoutMap must be a pointer to a struct")
	}

	// We have to preserve a mutable value of the struct to be able to set the fields.
	sv := reflect.ValueOf(layoutMap)
	st := sv.Elem().Type()
	if st.Kind() != reflect.Struct {
		return nil, fmt.Errorf("layoutMap must be a struct %s, got %s", st.Name(), st.Kind())
	}

	routes, err := readRoutesFromMapStruct(st, sv.Elem())
	if err != nil {
		return nil, fmt.Errorf("failed to generate query from struct type: %w", err)
	}
	dm := DataMap{
		Routes: routes,
	}
	return &dm, nil
}

// readRoutesFromMapStruct reads the routes from the given struct.
// Op(Offset|Size) tags are used to determine the operation.
// The struct fields must be tagged with `offsetof` or `sizeof` tags.
// The pointer is needed to be able to set the fields.
// e.g.:
//
//	sizeof(Type)
//	sizeof(StructType.(StructType)*.Field)
//	offsetof(StructType.Field)
//	offsetof(StructType.(StructType)*.Field)
func readRoutesFromMapStruct(st reflect.Type, sv reflect.Value) ([]*RouteNode, error) {
	var (
		groupBy = make(map[string]*RouteNode)
		add     = func(path string, name string, op Operation, fieldValue reflect.Value) {
			if r, exists := groupBy[path]; exists {
				r.Leaf().Extractors = append(r.Leaf().Extractors, &Extractor{
					Source: name, Op: op, targetValue: &fieldValue,
				},
				)
				return
			}
			route := newRouteFromTagValue(path)
			route.Leaf().Extractors = []*Extractor{
				{Source: name, Op: op, targetValue: &fieldValue},
			}
			groupBy[path] = route
		}
	)
	for i := 0; i < st.NumField(); i++ {
		field := st.Field(i)
		fieldValue := sv.Field(i)

		if !fieldValue.CanSet() && !sv.CanSet() {
			return nil, fmt.Errorf("field %s is not settable", field.Name)
		}
		if !isIntType(field.Type) {
			return nil, fmt.Errorf("field %s is not of type int or uint, type: %s", field.Name, field.Type.Kind())
		}

		var (
			tagValue string
			ok       bool
			op       Operation
		)
		tagValue, ok = field.Tag.Lookup(tagOffsetOf)
		if ok {
			op = OpOffsetOf
		}
		if tagValue == "" {
			tagValue, ok = field.Tag.Lookup(tagSizeOf)
			if ok {
				op = OpSizeOf
			}
		}
		if tagValue == "" || tagValue == "-" {
			continue
		}

		parts := strings.Split(tagValue, ".")
		if len(parts) < op.minimumRequiredRouteLength() {
			return nil, fmt.Errorf("field %s: invalid tag value: %s", field.Name, tagValue)
		}

		// Separate the field name from the path.
		var (
			path      string
			fieldName string
		)
		if len(parts) == 1 {
			path = tagValue
			fieldName = tagValue
		} else {
			path = strings.Join(parts[:len(parts)-1], ".")
			fieldName = parts[len(parts)-1]
		}
		add(path, fieldName, op, fieldValue)
	}

	if len(groupBy) == 0 {
		return nil, errors.New("no fields found with offsetof or sizeof tag")
	}
	return maps.Values(groupBy), nil
}

func isIntType(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return true
	default:
		return false
	}
}
