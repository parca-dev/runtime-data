// Copyright 2022-2024 The Parca Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package buildid

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestFromFile(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "go binary",
			args: args{
				path: "./testdata/readelf-sections",
			},
			want: "8HZi_313fFZIwx9R85S5/pagPyamQ7GjRRvxkDrCh/VF65lKUDP8KhNqvmQ31J/Iv_9XZ3HkWjhOW0faRQX",
		},
		{
			name: "rust binary",
			args: args{
				path: "./testdata/rust",
			},
			want: "ea8a38018312ad155fa70e471d4e0039ff9971c6",
		},
		{
			name: "rust binary build with bazel",
			args: args{
				path: "./testdata/bazel-rust",
			},
			want: "983bd888c60ead8e",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := os.Open(tt.args.path)
			if err != nil {
				t.Fatal(err)
			}
			got, err := FromFile(f)
			if (err != nil) != tt.wantErr {
				t.Errorf("FromFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("FromFile() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
