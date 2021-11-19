// Copyright 2021 FairwindsOps, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License

package helm

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"helm.sh/helm/v3/pkg/release"
)


func Test_helmToRelease(t *testing.T) {
	tests := []struct {
		name        string
		helmRelease interface{}
		want        *release.Release
		wantErr     bool
		errMsg      string
	}{
		{
			name:        "test err in json.Marshal",
			helmRelease: map[string]interface{}{"foo": make(chan int)},
			want:        nil,
			wantErr:     true,
			errMsg:      "error marshaling release: json: unsupported type: chan int",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := helmToRelease(tt.helmRelease)
			if tt.wantErr {
				assert.EqualError(t, err, tt.errMsg)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
