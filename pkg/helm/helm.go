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
	"context"
	"encoding/json"
	"fmt"

	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/releaseutil"
	helmstoragev3 "helm.sh/helm/v3/pkg/storage"
	driverv3 "helm.sh/helm/v3/pkg/storage/driver"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog"
)

// Helm represents all current releases that we can find in the cluster
type Helm struct {
	Releases []*release.Release
	Kube     *kubernetes.Clientset
	Dynamic  *DynamicClientInstance
}

// NewHelm returns a basic helm struct
func NewHelm() *Helm {
	return &Helm{
		Kube:    getKubeClient(),
		Dynamic: getDynamicInstance(),
	}
}

// GetReleasesVersionThree retrieves helm 3 releases from Secrets
func (h *Helm) GetReleasesVersionThree() error {
	hs := driverv3.NewSecrets(h.Kube.CoreV1().Secrets(""))
	helmClient := helmstoragev3.Init(hs)
	namespaces := h.GetNamespaces()

	releases, err := helmClient.ListDeployed()
	if err != nil {
		return err
	}
	for _, namespace := range namespaces.Items {
		ns := namespace.Name

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

// GetNamespaces retrieves a list of namespaces for a cluster
func (h *Helm) GetNamespaces() *v1.NamespaceList {
	ns, err := h.Kube.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		klog.Error(err)
	}
	return ns
}

// GetClusterObjects returns a list of unstructured.Unstructured objects
func (h *Helm) GetClusterObjects(group string, version string, resource string, namespace string) ([]unstructured.Unstructured, error) {
	resourceId := schema.GroupVersionResource{
		Group:    group,
		Version:  version,
		Resource: resource,
	}
	list, err := h.Dynamic.Client.Resource(resourceId).Namespace(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		klog.Error(err)
	}

	return list.Items, nil
}
