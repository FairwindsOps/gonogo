package validate

import (
	_ "context"
	"fmt"

	_ "github.com/fairwindsops/insights-plugins/opa/pkg/rego"
)

func (m *match) RunOPAChecks() error {
	//  RunRegoForItem(ctx context.Context, regoStr string, params map[string]interface{}, obj map[string]interface{}, dataFn KubeDataFunction, insightsInfo *InsightsInfo)
	// ([]interface{}, error)

	fmt.Println(m.Bundle.OpaChecks)

	// unmarshal, err :=

	// r, err := rego.RunRegoForItem(context.TODO(), somerego, nil, obj?,  )

	return nil
}

// func (m *match) AddActionItems(item *ActionItem) error {

// }
