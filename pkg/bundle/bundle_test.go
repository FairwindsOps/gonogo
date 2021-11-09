/*
Copyright © 2021 FairwindsOps Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package bundle

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadConfig(t *testing.T) {
	tests := []struct {
		name    string
		file    string
		want    *BundleConfig
		wantErr bool
	}{
		{
			name: "bundle1",
			file: "testdata/bundle1.yaml",
			want: &BundleConfig{
				Addons: []Bundle{
					{
						Name:                  "external-dns",
						Versions:              Versions{"1.1.0", "1.2.0"},
						Notes:                 "A text field with general notes",
						Source:                Source{"external-dns", "https://charts.bitnami.com/bitnami"},
						Warnings:              []string{"warning 1", "warning 2"},
						CompatibleK8SVersions: []string{"1.20", "1.18"},
						NecessaryAPIVersions:  []string{"apps/v1", "core/v1"},
						ValuesSchema:          "",
						OpaChecks:             []string{"Type From the existing OPA Package", "OPACustomCheck", "Handle resources"},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "file is not valid syntax",
			file: "testdata/invalid_bundle.yaml",
			want: &BundleConfig{},
			wantErr: true,
		},
		{
			name:    "file does not exist",
			file:    "farglebargle",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadConfig(tt.file)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.EqualValues(t, got, tt.want)
			}
		})
	}
}
