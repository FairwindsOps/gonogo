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
	"bytes"
	"context"
	"io"
	"strings"

	"github.com/fairwindsops/hall-monitor/pkg/helm"
	"github.com/fairwindsops/insights-plugins/plugins/opa/pkg/rego"
	"gopkg.in/yaml.v3"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog"
)

var clientset = helm.GetConfigInstance()
var dynamicClient = helm.GetDynamicInstance().Client

var (
	group    string
	version  string
	resource string
)

// RunOPAChecks evaluates rego defined in bundle spec against helm charts and cluster objects and returns an error
func (m *match) RunOPAChecks() error {
	if len(m.Bundle.OpaChecks) < 1 {
		return nil
	}

	manifests, err := splitYAML([]byte(m.Release.Manifest))
	if err != nil {
		klog.Error(err)
		return nil
	}

	resources := m.Bundle.Resources

	namespaces := helm.GetNamespaces()

	for _, namespace := range namespaces.Items {
		ns := namespace.Name
		for _, r := range resources {
			rs := strings.SplitAfter(r, "/")

			if len(rs) == 3 {
				group = rs[0]
				version = rs[1]
				resource = rs[2]
			} else {
				group = ""
				version = rs[0]
				resource = rs[1]
			}

			objs, err := GetClusterObjects(dynamicClient, context.TODO(), group, version, resource, ns)
			if err != nil {
				klog.Error()
				continue
			}
			for _, i := range objs {
				manifests = append(manifests, i.Object)
			}
		}
	}

	for _, o := range m.Bundle.OpaChecks {
		for _, y := range manifests {
			r, err := rego.RunRegoForItemV2(context.TODO(), o, y, clientset, nil)
			if err != nil {
				klog.Error(err)
				continue
			}
			for _, l := range r {
				b, err := yaml.Marshal(l)
				if err != nil {
					klog.Error(err)
					continue
				}
				var item *ActionItem
				err = yaml.Unmarshal(b, &item)
				if err != nil {
					klog.Error(err)
					continue
				}
				if item.ResourceKind == "" {
					item.ResourceKind = y["kind"].(string)
				}
				if item.ResourceName == "" {
					item.ResourceName = y["metadata"].(map[string]interface{})["name"].(string)
				}
				if item.ResourceNamespace == "" {
					item.ResourceNamespace = y["metadata"].(map[string]interface{})["namespace"].(string)
				}
				m.AddonOutput.ActionItems = append(m.AddonOutput.ActionItems, item)
			}

		}
	}

	return nil
}

// splitYAML takes a list of Helm manifests and splits them into separate files
func splitYAML(objects []byte) ([]map[string]interface{}, error) {

	dec := yaml.NewDecoder(bytes.NewReader(objects))

	var output []map[string]interface{}

	for {
		var value map[string]interface{}
		err := dec.Decode(&value)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		output = append(output, value)
	}
	return output, nil
}

// GetClusterObjects returns a list of unstructured.Unstructured objects
func GetClusterObjects(dynamic dynamic.Interface, ctx context.Context, group string, version string, resource string, namespace string) ([]unstructured.Unstructured, error) {
	resourceId := schema.GroupVersionResource{
		Group:    group,
		Version:  version,
		Resource: resource,
	}
	list, err := dynamic.Resource(resourceId).Namespace(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		klog.Error(err)
	}

	return list.Items, nil
}
