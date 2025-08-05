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
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Client represents an OpenAI API client
type Client interface {
	// AnalyzeUpgrade analyzes the upgrade path between two helm chart versions
	AnalyzeUpgrade(input UpgradeAnalysisInput) (*UpgradeAnalysisResponse, error)
}

// OpenAIClient implements the Client interface
type OpenAIClient struct {
	APIKey  string
	BaseURL string
	Model   string
}

// UpgradeAnalysisInput contains the information needed for upgrade analysis
type UpgradeAnalysisInput struct {
	AppName             string `json:"app_name"`
	ClusterVersion      string `json:"cluster_version"`
	CurrentChartVersion string `json:"current_chart_version"`
	DesiredChartVersion string `json:"desired_chart_version"`
	RepoURL             string `json:"repo_url"`
}

// UpgradeAnalysisResponse contains the analysis from OpenAI
type UpgradeAnalysisResponse struct {
	Analysis        string `json:"analysis"`
	BreakingChanges []string `json:"breaking_changes,omitempty"`
	Considerations  []string `json:"considerations,omitempty"`
	Recommendations []string `json:"recommendations,omitempty"`
}

// OpenAIRequest represents the request structure for OpenAI API
type OpenAIRequest struct {
	Model       string        `json:"model"`
	Messages    []ChatMessage `json:"messages"`
	MaxTokens   int           `json:"max_tokens,omitempty"`
	Temperature float64       `json:"temperature,omitempty"`
}

// ChatMessage represents a single message in the chat
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// OpenAIResponse represents the response from OpenAI API
type OpenAIResponse struct {
	Choices []struct {
		Message ChatMessage `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"error,omitempty"`
}

// NewClient creates a new OpenAI client
func NewClient(apiKey string) Client {
	return &OpenAIClient{
		APIKey:  apiKey,
		BaseURL: "https://api.openai.com/v1",
		Model:   "gpt-o3", // Default to a cost-effective model
	}
}

// NewClientWithModel creates a new OpenAI client with a specific model
func NewClientWithModel(apiKey, model string) Client {
	return &OpenAIClient{
		APIKey:  apiKey,
		BaseURL: "https://api.openai.com/v1",
		Model:   model,
	}
}

// AnalyzeUpgrade analyzes the upgrade path between two helm chart versions
func (c *OpenAIClient) AnalyzeUpgrade(input UpgradeAnalysisInput) (*UpgradeAnalysisResponse, error) {
	prompt := buildUpgradeAnalysisPrompt(input)
	
	request := OpenAIRequest{
		Model: c.Model,
		Messages: []ChatMessage{
			{
				Role:    "system",
				Content: "You are a Kubernetes and Helm expert. Analyze helm chart upgrades and provide detailed insights about potential breaking changes, considerations, and recommendations. Focus on practical, actionable advice.",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
		MaxTokens:   2000,
		Temperature: 0.1, // Low temperature for more consistent, factual responses
	}

	response, err := c.makeAPICall(request)
	if err != nil {
		return nil, fmt.Errorf("failed to call OpenAI API: %w", err)
	}

	if len(response.Choices) == 0 {
		return nil, fmt.Errorf("no response choices received from OpenAI")
	}

	analysis := response.Choices[0].Message.Content
	
	return &UpgradeAnalysisResponse{
		Analysis: analysis,
	}, nil
}

// buildUpgradeAnalysisPrompt creates the prompt for upgrade analysis
func buildUpgradeAnalysisPrompt(input UpgradeAnalysisInput) string {
	return fmt.Sprintf(`I have a Kubernetes cluster running version %s. My installation of %s is currently on helm chart %s. I want to update to %s. Can you analyze and summarize the release notes between the two chart versions, paying attention for any breaking changes? Also highlight any considerations like CRD changes between the two versions that could cause issues.

Application Details:
- App Name: %s
- Cluster Version: %s  
- Current Chart Version: %s
- Desired Chart Version: %s
- Repository URL: %s

Please provide:
1. A summary of major changes between these versions
2. Any breaking changes that require action
3. CRD changes or considerations
4. Recommended upgrade steps or precautions
5. Any version compatibility issues with the cluster

Please be specific and actionable in your recommendations.

If you do not know the answer to a question, say "I don't know".

Give me a verbose analysis of the release notes between the two chart versions. I need you to be very specific and detailed in your analysis including information about new features, bug fixes, and any other information.`,
		input.ClusterVersion,
		input.AppName,
		input.CurrentChartVersion,
		input.DesiredChartVersion,
		input.AppName,
		input.ClusterVersion,
		input.CurrentChartVersion,
		input.DesiredChartVersion,
		input.RepoURL,
	)
}

// makeAPICall makes the actual HTTP request to OpenAI API
func (c *OpenAIClient) makeAPICall(request OpenAIRequest) (*OpenAIResponse, error) {
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", c.BaseURL+"/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.APIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var response OpenAIResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if response.Error != nil {
		return nil, fmt.Errorf("OpenAI API error: %s", response.Error.Message)
	}

	return &response, nil
}