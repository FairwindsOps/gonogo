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

	"k8s.io/klog"
)

var (
	group    string
	version  string
	resource string
)

// getClusterManifests gets manifests from the cluster not included in helm release
func (m *match) getClusterManifests() ([]map[string]interface{}, error) {
	var manifests []map[string]interface{}
	resources := m.Bundle.Resources

	namespaces := m.Helm.GetNamespaces()

	for _, namespace := range namespaces.Items {
		ns := namespace.Name
		for _, r := range resources {
			splitResourcePath(r)
			objs, err := m.Helm.GetClusterObjects(group, version, resource, ns)
			if err != nil {
				klog.Error()
				continue
			}
			for _, i := range objs {
				manifests = append(manifests, i.Object)
			}
		}
	}
	return manifests, nil
}

// RunOPAChecks evaluates rego defined in bundle spec against helm charts and cluster objects and returns an error
func (m *match) runOPAChecks() error {
	if len(m.Bundle.OpaChecks) < 1 {
		return nil
	}

	manifests, err := splitYAML([]byte(m.Release.Manifest))
	if err != nil {
		return err
	}

	clusterManifests, err := m.getClusterManifests()
	if err != nil {
		return err
	}

	manifests = append(manifests, clusterManifests...)

	for _, o := range m.Bundle.OpaChecks {
		for _, y := range manifests {
			m.addActionItem(o, y)
		}
	}

	return nil
}

// addActionItem runs rego against manifest using passed in opa check from bundle and appends to actionItems
func (m *match) addActionItem(o string, y map[string]interface{}) {
	client := helm.NewHelm()

	r, err := rego.RunRegoForItemV2(context.TODO(), o, y, client.Kube, nil)
	if err != nil {
		klog.Error(err)
	}

	for _, l := range r {
		b, err := yaml.Marshal(l)
		if err != nil {
			klog.Error(err)
			continue
		}

		var i *ActionItem

		err = yaml.Unmarshal(b, &i)
		if err != nil {
			klog.Error(err)
			continue
		}
		if i.ResourceKind == "" {
			i.ResourceKind = y["kind"].(string)
		}
		if i.ResourceName == "" {
			i.ResourceName = y["metadata"].(map[string]interface{})["name"].(string)
		}
		if i.ResourceNamespace == "" {
			i.ResourceNamespace = y["metadata"].(map[string]interface{})["namespace"].(string)
		}
		m.AddonOutput.ActionItems = append(m.AddonOutput.ActionItems, i)
	}
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

// splitResourcePath takes resource string defined in bundle and splits into separate strings to be passed to apiserver so that we can dynamically look up objects
func splitResourcePath(path string) {
	rs := strings.SplitAfter(path, "/")

	if len(rs) == 3 {
		group = rs[0]
		version = rs[1]
		resource = rs[2]
	} else {
		group = ""
		version = rs[0]
		resource = rs[1]
	}

}
