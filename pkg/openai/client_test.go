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
package openai

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// MockClient implements the Client interface for testing
type MockClient struct {
	Response *UpgradeAnalysisResponse
	Error    error
}

func (m *MockClient) AnalyzeUpgrade(input UpgradeAnalysisInput) (*UpgradeAnalysisResponse, error) {
	if m.Error != nil {
		return nil, m.Error
	}
	return m.Response, nil
}

func TestNewClient(t *testing.T) {
	client := NewClient("test-api-key")
	
	openaiClient, ok := client.(*OpenAIClient)
	assert.True(t, ok)
	assert.Equal(t, "test-api-key", openaiClient.APIKey)
	assert.Equal(t, "https://api.openai.com/v1", openaiClient.BaseURL)
	assert.Equal(t, "gpt-4o-mini", openaiClient.Model)
}

func TestNewClientWithModel(t *testing.T) {
	client := NewClientWithModel("test-api-key", "gpt-4")
	
	openaiClient, ok := client.(*OpenAIClient)
	assert.True(t, ok)
	assert.Equal(t, "test-api-key", openaiClient.APIKey)
	assert.Equal(t, "https://api.openai.com/v1", openaiClient.BaseURL)
	assert.Equal(t, "gpt-4", openaiClient.Model)
}

func TestBuildUpgradeAnalysisPrompt(t *testing.T) {
	input := UpgradeAnalysisInput{
		AppName:             "nginx-ingress",
		ClusterVersion:      "1.28.0",
		CurrentChartVersion: "4.7.1",
		DesiredChartVersion: "4.8.0",
		RepoURL:             "https://kubernetes.github.io/ingress-nginx",
	}

	prompt := buildUpgradeAnalysisPrompt(input)
	
	// Check that the prompt contains all the expected information
	assert.Contains(t, prompt, "1.28.0")
	assert.Contains(t, prompt, "nginx-ingress")
	assert.Contains(t, prompt, "4.7.1")
	assert.Contains(t, prompt, "4.8.0")
	assert.Contains(t, prompt, "https://kubernetes.github.io/ingress-nginx")
	assert.Contains(t, prompt, "breaking changes")
	assert.Contains(t, prompt, "CRD changes")
}

func TestMockClient(t *testing.T) {
	mockResponse := &UpgradeAnalysisResponse{
		Analysis: "Test analysis response",
		BreakingChanges: []string{"Breaking change 1"},
		Considerations:  []string{"Consideration 1"},
		Recommendations: []string{"Recommendation 1"},
	}

	mockClient := &MockClient{
		Response: mockResponse,
	}

	input := UpgradeAnalysisInput{
		AppName:             "test-app",
		ClusterVersion:      "1.28.0",
		CurrentChartVersion: "1.0.0",
		DesiredChartVersion: "2.0.0",
		RepoURL:             "https://example.com/charts",
	}

	response, err := mockClient.AnalyzeUpgrade(input)
	
	assert.NoError(t, err)
	assert.Equal(t, mockResponse, response)
}