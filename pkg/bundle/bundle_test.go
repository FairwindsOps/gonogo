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
		file    []string
		want    *BundleConfig
		wantErr bool
	}{
		{
			name: "bundle1",
			file: []string{"testdata/bundle_read_check.yaml"},
			want: &BundleConfig{
				Addons: []*Bundle{
					{
						Name:                  "metrics-server",
						Versions:              Versions{"5.10.2", "5.10.14"},
						Notes:                 "A text field with general notes",
						Source:                Source{"metrics-server", "https://charts.bitnami.com/bitnami"},
						Warnings:              []string{"warning 1", "warning 2"},
						CompatibleK8sVersions: K8sVersions{"1.18", "1.20"},
						NecessaryAPIVersions:  []string{"apps/v1", "v1"},
						ValuesSchema:          "",
						OpaChecks:             []string{"Check One", "Check Two"},
					},
				},
			},
			wantErr: false,
		},
		{
			name:    "file is not valid syntax",
			file:    []string{"testdata/invalid_bundle.yaml"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "file does not exist",
			file:    []string{"farglebargle"},
			want:    nil,
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

func TestReadConfigDefaultBundles(t *testing.T) {
	got, err := ReadConfig(nil)
	assert.NoError(t, err)

	var found bool
	for _, addon := range got.Addons {
		if addon.Name != "argocd" {
			continue
		}

		found = true
		assert.Equal(t, "9.4.15", addon.Versions.Start)
		assert.Equal(t, "9.5.6", addon.Versions.End)
		assert.Equal(t, "argo-cd", addon.Source.Chart)
		assert.Equal(t, "https://argoproj.github.io/argo-helm", addon.Source.Repository)
		assert.Equal(t, "1.25", addon.CompatibleK8sVersions.Min)
		assert.Contains(t, addon.NecessaryAPIVersions, "apiextensions.k8s.io/v1")
	}

	assert.True(t, found, "expected embedded argocd bundle")
}
