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
	"strconv"

	"github.com/dop251/goja"
)

// Client represents an OpenAI API client
type Client interface {
	// AnalyzeUpgrade analyzes the upgrade path between two helm chart versions
	AnalyzeUpgrade(input UpgradeAnalysisInput) (*UpgradeAnalysisResponse, error)
}

// GojaOpenAIClient implements the Client interface using goja JavaScript runtime
type GojaOpenAIClient struct {
	APIKey string
	Model  string
	vm     *goja.Runtime
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
	Analysis        string   `json:"analysis"`
	BreakingChanges []string `json:"breaking_changes,omitempty"`
	Considerations  []string `json:"considerations,omitempty"`
	Recommendations []string `json:"recommendations,omitempty"`
}

// NewClient creates a new OpenAI client using goja JavaScript runtime
func NewClient(apiKey string) Client {
	return &GojaOpenAIClient{
		APIKey: apiKey,
		Model:  "gpt-4o-search-preview", // Default to web search enabled model
		vm:     goja.New(),
	}
}

// NewClientWithModel creates a new OpenAI client with a specific model using goja
func NewClientWithModel(apiKey, model string) Client {
	return &GojaOpenAIClient{
		APIKey: apiKey,
		Model:  model,
		vm:     goja.New(),
	}
}

// AnalyzeUpgrade analyzes the upgrade path between two helm chart versions using goja and OpenAI JavaScript SDK
func (c *GojaOpenAIClient) AnalyzeUpgrade(input UpgradeAnalysisInput) (*UpgradeAnalysisResponse, error) {
	// Initialize the JavaScript environment
	if err := c.setupJavaScriptEnvironment(); err != nil {
		return nil, fmt.Errorf("failed to setup JavaScript environment: %w", err)
	}

	// Convert input to JSON for JavaScript
	inputJSON, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal input: %w", err)
	}

	// Build the prompt
	prompt := buildUpgradeAnalysisPrompt(input)

	// Use a synchronous JavaScript approach instead of async/await
	// since goja doesn't have a built-in event loop
	script := fmt.Sprintf(`
		(function() {
			try {
				var input = JSON.parse(%s);
				var apiKey = %s;
				var prompt = %s;
				
				// Create OpenAI instance
				var openai = new OpenAI({ apiKey: apiKey });
				
				// Prepare the request with web search capabilities
				var requestBody = JSON.stringify({
					model: 'gpt-4o-search-preview',
					messages: [
						{
							role: 'system',
							content: 'You are a Kubernetes and Helm expert. Analyze helm chart upgrades and provide detailed insights about potential breaking changes, considerations, and recommendations. Focus on practical, actionable advice.'
						},
						{
							role: 'user',
							content: prompt
						}
					],
					web_search_options: {
						search_context_size: 'high'
					},
					max_tokens: 2000,
					temperature: 0.1
				});
				
				// Make the HTTP request using our Go-backed fetch
				var response = syncFetch('https://api.openai.com/v1/chat/completions', {
					method: 'POST',
					headers: {
						'Content-Type': 'application/json',
						'Authorization': 'Bearer ' + apiKey
					},
					body: requestBody
				});
				
				if (!response.ok) {
					throw new Error('OpenAI API request failed: ' + response.status + ' ' + response.statusText);
				}
				
				var result = response.json();
				
				if (result.error) {
					throw new Error('OpenAI API error: ' + result.error.message);
				}
				
				if (!result.choices || result.choices.length === 0) {
					throw new Error('No response choices received from OpenAI');
				}
				
				var analysis = result.choices[0].message.content;
				
				return JSON.stringify({
					analysis: analysis,
					breaking_changes: [],
					considerations: [],
					recommendations: []
				});
			} catch (error) {
				throw error;
			}
		})();
		`,
		strconv.Quote(string(inputJSON)),
		strconv.Quote(c.APIKey),
		strconv.Quote(prompt),
	)

	result, err := c.vm.RunString(script)
	if err != nil {
		return nil, fmt.Errorf("failed to execute JavaScript: %w", err)
	}

	// Parse the result
	resultJSON := result.String()
	var response UpgradeAnalysisResponse
	if err := json.Unmarshal([]byte(resultJSON), &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

// setupJavaScriptEnvironment sets up the JavaScript environment with OpenAI SDK and necessary functions
func (c *GojaOpenAIClient) setupJavaScriptEnvironment() error {
	// Set up console for logging
	c.vm.Set("console", map[string]interface{}{
		"log": func(args ...interface{}) {
			// In production, you might want to use a proper logger
			for _, arg := range args {
				fmt.Print(arg, " ")
			}
			fmt.Println()
		},
		"error": func(args ...interface{}) {
			fmt.Print("ERROR: ")
			for _, arg := range args {
				fmt.Print(arg, " ")
			}
			fmt.Println()
		},
	})

	// Set up both async and sync fetch implementations for HTTP requests
	c.vm.Set("fetch", c.createFetchFunction())
	c.vm.Set("syncFetch", c.createSyncFetchFunction())

	// Load a minimal OpenAI SDK implementation in JavaScript
	openaiSDK := `
		// Minimal OpenAI SDK implementation
		class OpenAI {
			constructor(config) {
				this.apiKey = config.apiKey;
				this.baseURL = config.baseURL || 'https://api.openai.com/v1';
			}

			async createChatCompletion(params) {
				const response = await fetch(this.baseURL + '/chat/completions', {
					method: 'POST',
					headers: {
						'Content-Type': 'application/json',
						'Authorization': 'Bearer ' + this.apiKey
					},
					body: JSON.stringify(params)
				});

				if (!response.ok) {
					throw new Error('OpenAI API request failed: ' + response.status + ' ' + response.statusText);
				}

				return await response.json();
			}
		}

		// Global function to analyze upgrade
		async function analyzeUpgrade(inputJSON, apiKey, prompt) {
			try {
				const input = JSON.parse(inputJSON);
				const openai = new OpenAI({ apiKey: apiKey });

				const completion = await openai.createChatCompletion({
					model: '` + c.Model + `',
					messages: [
						{
							role: 'system',
							content: 'You are a Kubernetes and Helm expert. Analyze helm chart upgrades and provide detailed insights about potential breaking changes, considerations, and recommendations. Focus on practical, actionable advice.'
						},
						{
							role: 'user',
							content: prompt
						}
					],
					max_tokens: 2000,
					temperature: 0.1
				});

				const analysis = completion.choices[0].message.content;

				// Return the response in the expected format
				return JSON.stringify({
					analysis: analysis,
					breaking_changes: [],
					considerations: [],
					recommendations: []
				});
			} catch (error) {
				console.error('Error in analyzeUpgrade:', error);
				throw error;
			}
		}
	`

	// Load the OpenAI SDK into the JavaScript runtime
	_, err := c.vm.RunString(openaiSDK)
	if err != nil {
		return fmt.Errorf("failed to load OpenAI SDK: %w", err)
	}

	return nil
}

// createFetchFunction creates a fetch function that can be used in JavaScript
func (c *GojaOpenAIClient) createFetchFunction() func(string, map[string]interface{}) *goja.Promise {
	return func(url string, options map[string]interface{}) *goja.Promise {
		// Create a new promise
		promise, resolve, reject := c.vm.NewPromise()

		// Execute the HTTP request in a goroutine
		go func() {
			defer func() {
				if r := recover(); r != nil {
					reject(c.vm.ToValue(fmt.Sprintf("fetch panic: %v", r)))
				}
			}()

			// Extract options
			method := "GET"
			var body []byte
			headers := make(map[string]string)

			if options != nil {
				if m, ok := options["method"].(string); ok {
					method = m
				}
				if b, ok := options["body"].(string); ok {
					body = []byte(b)
				}
				if h, ok := options["headers"].(map[string]interface{}); ok {
					for k, v := range h {
						if str, ok := v.(string); ok {
							headers[k] = str
						}
					}
				}
			}

			// Create HTTP request
			var req *http.Request
			var err error

			if body != nil {
				req, err = http.NewRequest(method, url, bytes.NewReader(body))
			} else {
				req, err = http.NewRequest(method, url, nil)
			}

			if err != nil {
				reject(c.vm.ToValue(fmt.Sprintf("failed to create request: %v", err)))
				return
			}

			// Set headers
			for k, v := range headers {
				req.Header.Set(k, v)
			}

			// Make the request
			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				reject(c.vm.ToValue(fmt.Sprintf("fetch failed: %v", err)))
				return
			}
			defer resp.Body.Close()

			// Read response body
			respBody, err := io.ReadAll(resp.Body)
			if err != nil {
				reject(c.vm.ToValue(fmt.Sprintf("failed to read response: %v", err)))
				return
			}

			// Create response object
			response := map[string]interface{}{
				"ok":         resp.StatusCode >= 200 && resp.StatusCode < 300,
				"status":     resp.StatusCode,
				"statusText": resp.Status,
				"json": func() *goja.Promise {
					jsonPromise, jsonResolve, jsonReject := c.vm.NewPromise()

					go func() {
						var jsonData interface{}
						if err := json.Unmarshal(respBody, &jsonData); err != nil {
							jsonReject(c.vm.ToValue(fmt.Sprintf("failed to parse JSON: %v", err)))
							return
						}
						jsonResolve(c.vm.ToValue(jsonData))
					}()

					return jsonPromise
				},
				"text": func() *goja.Promise {
					textPromise, textResolve, _ := c.vm.NewPromise()
					go func() {
						textResolve(c.vm.ToValue(string(respBody)))
					}()
					return textPromise
				},
			}

			resolve(c.vm.ToValue(response))
		}()

		return promise
	}
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

// createSyncFetchFunction creates a synchronous fetch function for JavaScript
func (c *GojaOpenAIClient) createSyncFetchFunction() func(string, map[string]interface{}) map[string]interface{} {
	return func(url string, options map[string]interface{}) map[string]interface{} {
		// Extract options
		method := "GET"
		var body []byte
		headers := make(map[string]string)

		if options != nil {
			if m, ok := options["method"].(string); ok {
				method = m
			}
			if b, ok := options["body"].(string); ok {
				body = []byte(b)
			}
			if h, ok := options["headers"].(map[string]interface{}); ok {
				for k, v := range h {
					if str, ok := v.(string); ok {
						headers[k] = str
					}
				}
			}
		}

		// Create HTTP request
		var req *http.Request
		var err error

		if body != nil {
			req, err = http.NewRequest(method, url, bytes.NewReader(body))
		} else {
			req, err = http.NewRequest(method, url, nil)
		}

		if err != nil {
			return map[string]interface{}{
				"ok":     false,
				"status": 0,
				"error":  err.Error(),
			}
		}

		// Set headers
		for k, v := range headers {
			req.Header.Set(k, v)
		}

		// Make the request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return map[string]interface{}{
				"ok":     false,
				"status": 0,
				"error":  err.Error(),
			}
		}
		defer resp.Body.Close()

		// Read response body
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return map[string]interface{}{
				"ok":     false,
				"status": resp.StatusCode,
				"error":  err.Error(),
			}
		}

		// Parse JSON response
		var jsonData interface{}
		jsonErr := json.Unmarshal(respBody, &jsonData)

		// Create response object
		response := map[string]interface{}{
			"ok":         resp.StatusCode >= 200 && resp.StatusCode < 300,
			"status":     resp.StatusCode,
			"statusText": resp.Status,
			"json": func() interface{} {
				if jsonErr != nil {
					return map[string]interface{}{
						"error": jsonErr.Error(),
					}
				}
				return jsonData
			},
			"text": func() string {
				return string(respBody)
			},
		}

		return response
	}
}
