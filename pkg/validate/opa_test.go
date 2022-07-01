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

package validate

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplitResourcePath(t *testing.T) {
	tests := []struct {
		name    string
		args    string
		want    []string
		wantErr bool
	}{
		{
			name:    "test for group version and resource",
			args:    "apps/v1/deployments",
			want:    []string{"apps/", "v1/", "deployments"},
			wantErr: false,
		},
		{
			name:    "test version and resource with blank group pass",
			args:    "v1/secrets",
			want:    []string{"", "v1/", "secrets"},
			wantErr: false,
		},
		{
			name:    "test for version and resource with blank group error",
			args:    "v1/pods",
			want:    []string{"v1/", "pods"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			splitResourcePath(tt.args)
			if tt.wantErr {
				err := fmt.Errorf("path not split properly")
				assert.Error(t, err)
			} else {
				assert.NoError(t, nil)
			}
		})
	}
}
