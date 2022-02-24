package validate

import (
	"context"

	"github.com/fairwindsops/insights-plugins/opa/pkg/rego"
)

func (m *match) RunOPAChecks() error {
	//  RunRegoForItem(ctx context.Context, regoStr string, params map[string]interface{}, obj map[string]interface{}, dataFn KubeDataFunction, insightsInfo *InsightsInfo)
	// ([]interface{}, error)

	r, err := rego.RunRegoForItem(context.TODO(), )

	return nil
}
