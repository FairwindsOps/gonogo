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
	"encoding/json"
	"fmt"

	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/releaseutil"
	helmstoragev3 "helm.sh/helm/v3/pkg/storage"
	driverv3 "helm.sh/helm/v3/pkg/storage/driver"
)

// Helm represents all current releases that we can find in the cluster
type Helm struct {
	Releases  []*release.Release
	Kube      *kube
	Namespace string
}


// NewHelm returns a basic helm struct
func NewHelm(namespace string) *Helm {
	return &Helm{
		Kube:      GetConfigInstance(),
		Namespace: namespace,
	}
}

// GetReleasesVersionThree retrieves helm 3 releases from Secrets
func (h *Helm) GetReleasesVersionThree() error {
	hs := driverv3.NewSecrets(h.Kube.Client.CoreV1().Secrets(h.Namespace))
	helmClient := helmstoragev3.Init(hs)
	namespaces := GetNamespaces()

	releases, err := helmClient.ListDeployed()
	if err != nil {
		return err
	}
	for _, namespace := range namespaces.Items {
		ns := namespace.Name
		if h.Namespace != "" && ns != h.Namespace {
			continue
		}
		filteredReleases := h.deployedReleasesPerNamespace(ns, releases)
		for _, r := range filteredReleases {
			rel, err := helmToRelease(r)
			if err != nil {
				return fmt.Errorf("error converting helm r '%s/%s' to internal object\n   %w", r.Namespace, r.Name, err)
			}
			h.Releases = append(h.Releases, rel)
		}
	}
	return nil
}

func (h *Helm) deployedReleasesPerNamespace(namespace string, releases []*release.Release) []*release.Release {
	return releaseutil.All(deployed, relNamespace(namespace)).Filter(releases)
}

func deployed(rls *release.Release) bool {
	return rls.Info.Status == release.StatusDeployed
}

func relNamespace(ns string) releaseutil.FilterFunc {
	return func(rls *release.Release) bool {
		return rls.Namespace == ns
	}
}

func helmToRelease(helmRelease interface{}) (*release.Release, error) {
	jsonRel, err := json.Marshal(helmRelease)
	if err != nil {
		return nil, fmt.Errorf("error marshaling release: %s", err.Error())
	}
	return marshalToRelease(jsonRel)
}

// marshalToRelease marshals release data into the local Release type so we have a common type regardless of helm version
func marshalToRelease(jsonRel []byte) (*release.Release, error) {
	var ret = new(release.Release)
	err := json.Unmarshal(jsonRel, ret)
	return ret, err
}
