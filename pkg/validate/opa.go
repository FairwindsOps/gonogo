package validate

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/fairwindsops/hall-monitor/pkg/helm"
	"github.com/fairwindsops/insights-plugins/opa/pkg/rego"
	"gopkg.in/yaml.v3"
	"k8s.io/klog"
)

func (m *match) RunOPAChecks() error {
	if len(m.Bundle.OpaChecks) < 1 {
		return nil
	}

	manifests, err := splitYAML([]byte(m.Release.Manifest))
	if err != nil {
		klog.Error(err)
		return err
	}

	client := helm.NewHelm("")

	for _, o := range m.Bundle.OpaChecks {
		for _, y := range manifests {
			r, err := rego.RunRegoForItem(context.TODO(), o, nil, y, client, nil)
			if err != nil {
				klog.Error(err)
				continue
			}
			fmt.Println("value of r is", r)

		}
	}

	return nil
}

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
