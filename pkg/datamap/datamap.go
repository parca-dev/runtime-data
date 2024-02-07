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

const (
	OpOffset Operation = iota
	OpSize
)

type DataMap struct {
	Structs []*Struct
}

type Struct struct {
	StructName string
	Op         Operation

	// Value is only required for SizeOf operation.
	Value reflect.Value

	// Only required for OffsetOf operation.
	Fields []*Field
}

type Field struct {
	Name  string
	Value reflect.Value
}

// For testing.
func (s Struct) fieldNames() []string {
	names := make([]string, len(s.Fields))
	for i, f := range s.Fields {
		names[i] = f.Name
	}
	return names
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

	structs, err := generateFromStructType(st, sv.Elem())
	if err != nil {
		return nil, fmt.Errorf("failed to generate query from struct type: %w", err)
	}
	dm := DataMap{
		Structs: structs,
	}
	return &dm, nil
}

func generateFromStructType(st reflect.Type, sv reflect.Value) ([]*Struct, error) {
	var (
		structName = func(st reflect.Type) string {
			name := st.Name()
			if name == "" {
				// Anonymous struct, use the field name.
				name = st.String()
			}
			return name
		}
		groupByStructAndOp = make(map[string]*Struct)
		groupKey           = func(name string, op Operation) string {
			return fmt.Sprintf("%s-%d", name, op)
		}
		addActionByField = func(structName string, op Operation, fields ...*Field) {
			key := groupKey(structName, op)
			if sm, exists := groupByStructAndOp[key]; exists {
				sm.Fields = append(sm.Fields, fields...)
				return
			}

			groupByStructAndOp[key] = &Struct{
				StructName: structName,
				Op:         op,
				Fields:     fields,
			}
		}
		addActionByStruct = func(structName string, op Operation, value reflect.Value) {
			key := groupKey(structName, op)
			if _, exists := groupByStructAndOp[key]; exists {
				return
			}

			groupByStructAndOp[key] = &Struct{
				StructName: structName,
				Op:         op,
				Value:      value,
			}
		}

		name = structName(st)
	)
	for i := 0; i < st.NumField(); i++ {
		field := st.Field(i)
		fieldValue := sv.Field(i)

		var (
			tagValue string
			ok       bool
			op       Operation
		)
		tagValue, ok = field.Tag.Lookup(tagOffsetOf)
		if ok {
			op = OpOffset
		}
		if tagValue == "" {
			tagValue, ok = field.Tag.Lookup(tagSizeOf)
			if ok {
				op = OpSize
			}
		}

		if tagValue == "" || tagValue == "-" {
			continue
		}

		switch op {
		case OpOffset:
			if strings.Contains(tagValue, ".") {
				parts := strings.Split(tagValue, ".")
				if len(parts) != 2 {
					return nil, fmt.Errorf("invalid offset tag, only one-level nesting is supported: %s", tagValue)
				}
				addActionByField(parts[0], op, &Field{Name: parts[1], Value: fieldValue})
				continue
			}

			if field.Type.Kind() == reflect.Struct {
				nestedStructs, err := generateFromStructType(field.Type, fieldValue)
				if err != nil {
					return nil, fmt.Errorf("failed to generate query from struct type: %w", err)
				}

				for _, ns := range nestedStructs {
					key := ns.StructName
					if strings.HasPrefix(key, "struct ") {
						addActionByField(tagValue, op, ns.Fields...)
					} else {
						addActionByField(key, op, ns.Fields...)
					}
				}
				continue
			}

			// Struct tag given, we will use the struct tag.
			addActionByField(name, op, &Field{Name: tagValue, Value: fieldValue})
		case OpSize:
			if fieldValue.Kind() != reflect.Int64 {
				return nil, fmt.Errorf("size tag must be an int64, got %s", fieldValue.Kind())
			}
			addActionByStruct(tagValue, op, fieldValue)
		}
	}

	return maps.Values(groupByStructAndOp), nil
}
