//nolint:unused
package datamap

import (
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGenerateQuery(t *testing.T) {
	type simple struct {
		gopherA int    `offsetof:"a"`
		b       string `offsetof:"b"`
	}
	type simpleWithPaths struct {
		a int `offsetof:"_StructTypeILookFor.field_one"`
		b int `offsetof:"_AnotherStructTypeILookFor.field_two"`
		c int `offsetof:"_AnotherStructTypeILookFor.field_three"`
	}
	type ignore struct {
		a int `offsetof:"-"`
		b int `offsetof:"b"`
		c int `offsetof:""`
	}

	type _want struct {
		name   string
		op     Operation
		fields []string
	}

	tests := []struct {
		name      string
		mapStruct any
		want      []_want
		wantErr   bool
	}{
		{
			name: "no tags",
			mapStruct: &struct {
				a int
				b string
			}{},
			want: []_want{},
		},
		{
			name:      "simple",
			mapStruct: &simple{},
			want: []_want{
				{
					name:   "simple",
					op:     OpOffset,
					fields: []string{"a", "b"},
				},
			},
		},
		{
			name:      "simple with paths",
			mapStruct: &simpleWithPaths{},
			want: []_want{
				{
					name:   "_StructTypeILookFor",
					op:     OpOffset,
					fields: []string{"field_one"},
				},
				{
					name:   "_AnotherStructTypeILookFor",
					fields: []string{"field_two", "field_three"},
				},
			},
		},
		{
			name: "nested",
			mapStruct: &struct {
				a      int `offsetof:"_StructTypeILookFor.field_one"`
				astilf struct {
					b int `offsetof:"field_two"`
					c int `offsetof:"field_three"`
				} `offsetof:"_AnotherStructTypeILookFor"`
			}{},
			want: []_want{
				{
					name:   "_StructTypeILookFor",
					op:     OpOffset,
					fields: []string{"field_one"},
				},
				{
					name:   "_AnotherStructTypeILookFor",
					op:     OpOffset,
					fields: []string{"field_two", "field_three"},
				},
			},
		},
		{
			name: "deeply nested",
			mapStruct: &struct {
				a int `offsetof:"_StructTypeILookFor.field_one"`
				b struct {
					c int `offsetof:"field_two"`
					d struct {
						e int `offsetof:"field_three"`
					} `offsetof:"_YetAnotherStructTypeILookFor"`
				} `offsetof:"_AnotherStructTypeILookFor"`
			}{},
			want: []_want{
				{
					name:   "_StructTypeILookFor",
					op:     OpOffset,
					fields: []string{"field_one"},
				},
				{
					name:   "_AnotherStructTypeILookFor",
					op:     OpOffset,
					fields: []string{"field_two"},
				},
				{
					name:   "_YetAnotherStructTypeILookFor",
					op:     OpOffset,
					fields: []string{"field_three"},
				},
			},
		},
		{
			name:      "ignore",
			mapStruct: &ignore{},
			want: []_want{
				{
					name:   "ignore",
					fields: []string{"b"},
				},
			},
		},
		{
			name: "sizeof",
			mapStruct: &struct {
				a int `sizeof:"A"`
				b int `sizeof:"B"`
			}{},
			want: []_want{
				{
					name:   "A",
					op:     OpSize,
					fields: []string{},
				},
				{
					name:   "B",
					op:     OpSize,
					fields: []string{},
				},
			},
		},
		{
			name:      "not a struct pointer",
			mapStruct: 1,
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.mapStruct)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("Failed case: %s, reason: %s", tt.name, "unexpected error")
				}
				return
			}
			if len(got.Structs) != len(tt.want) {
				t.Errorf("Failed case: %s, reason: %s", tt.name, "length mismatch")
				t.Errorf("Structs: got = %v, want %v", got.Structs, tt.want)
				return
			}
			gotStructs := make([]_want, len(got.Structs))
			for i, s := range got.Structs {
				gotStructs[i] = _want{
					name:   s.StructName,
					op:     s.Op,
					fields: sort.StringSlice(s.fieldNames()),
				}
			}
			sort.Slice(gotStructs, func(i, j int) bool {
				return gotStructs[i].name < gotStructs[j].name
			})
			sort.Slice(tt.want, func(i, j int) bool {
				return tt.want[i].name < tt.want[j].name
			})
			if diff := cmp.Diff(gotStructs, tt.want, cmp.AllowUnexported(Struct{}, _want{})); diff != "" {
				t.Errorf("Failed case: %s, reason: %s", tt.name, diff)
			}
		})
	}
}
