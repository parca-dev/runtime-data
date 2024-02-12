//nolint:unused
package datamap

import (
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestGenerateQuery(t *testing.T) {
	tests := []struct {
		name      string
		mapStruct any
		want      []*RouteNode
		wantErr   bool
	}{
		{
			name: "no tags",
			mapStruct: &struct {
				a int
				b int
			}{},
			wantErr: true,
		},
		{
			name: "sizeof",
			mapStruct: &struct {
				a int `sizeof:"A"`
				b int `sizeof:"B"`
			}{},
			want: []*RouteNode{
				{
					Type: "A",
					Extractors: []*Extractor{
						{
							Source: "A",
							Op:     OpSizeOf,
						},
					},
				},
				{
					Type: "B",
					Extractors: []*Extractor{
						{
							Source: "B",
							Op:     OpSizeOf,
						},
					},
				},
			},
		},
		{
			name: "offsetof tags without at least one struct type",
			mapStruct: &struct {
				a int `offsetof:"a"`
				b int `offsetof:"b"`
			}{},
			wantErr: true,
		},
		{
			name: "simple",
			mapStruct: &struct {
				gopherA int `offsetof:"simple.a"`
				b       int `offsetof:"simple.b"`
			}{},
			want: []*RouteNode{
				{
					Type: "simple",
					Extractors: []*Extractor{
						{
							Source: "a",
							Op:     OpOffsetOf,
						},
						{
							Source: "b",
							Op:     OpOffsetOf,
						},
					},
				},
			},
		},
		{
			name: "ignore",
			mapStruct: &struct {
				a int `offsetof:"-"`
				b int `offsetof:"_MyStruct.b"`
				c int `offsetof:""`
			}{},
			want: []*RouteNode{
				{
					Type: "_MyStruct",
					Extractors: []*Extractor{
						{
							Source: "b",
							Op:     OpOffsetOf,
						},
					},
				},
			},
		},
		{
			name:      "not a struct pointer",
			mapStruct: 1,
			wantErr:   true,
		},
		{
			name: "simple with paths",
			mapStruct: &struct {
				a int `offsetof:"_StructTypeILookFor.field_one"`
				b int `offsetof:"_AnotherStructTypeILookFor.field_two"`
				c int `sizeof:"_AnotherStructTypeILookFor.field_three"`
			}{},
			want: []*RouteNode{
				{
					Type: "_StructTypeILookFor",
					Extractors: []*Extractor{
						{
							Source: "field_one",
							Op:     OpOffsetOf,
						},
					},
				},
				{
					Type: "_AnotherStructTypeILookFor",
					Extractors: []*Extractor{
						{
							Source: "field_two",
							Op:     OpOffsetOf,
						},
						{
							Source: "field_three",
							Op:     OpSizeOf,
						},
					},
				},
			},
		},
		{
			name: "nested target",
			mapStruct: &struct {
				a int `offsetof:"_StructTypeILookFor.field_one"`
				b int `offsetof:"_AnotherStructTypeILookFor.nested_struct.field_two"`
				c int `offsetof:"_AnotherStructTypeILookFor.nested_struct.field_three"`
			}{},
			want: []*RouteNode{
				{
					Type: "_StructTypeILookFor",
					Extractors: []*Extractor{
						{
							Source: "field_one",
							Op:     OpOffsetOf,
						},
					},
				},
				{
					Type:       "_AnotherStructTypeILookFor",
					Extractors: nil,
					Next: &RouteNode{
						Type: "nested_struct",
						Extractors: []*Extractor{
							{
								Source: "field_two",
								Op:     OpOffsetOf,
							},
							{
								Source: "field_three",
								Op:     OpOffsetOf,
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.mapStruct)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("Failed case: %s, reason: %s", tt.name, err)
				}
				return
			}

			sort.Slice(got.Routes, func(i, j int) bool {
				return got.Routes[i].Type < got.Routes[j].Type
			})
			sort.Slice(tt.want, func(i, j int) bool {
				return tt.want[i].Type < tt.want[j].Type
			})
			diff := cmp.Diff(
				tt.want, got.Routes,
				cmpopts.IgnoreUnexported(RouteNode{}, Extractor{}),
			)
			if diff != "" {
				t.Errorf("Failed case: %s, mismatch (-want +got): %s", tt.name, diff)
			}
		})
	}
}
