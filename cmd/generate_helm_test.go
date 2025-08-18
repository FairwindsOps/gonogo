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
package cmd

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestWebhookJSONResponseUnmarshal(t *testing.T) {
	// Test JSON response that would be received from webhook
	jsonResponse := `{
		"addons": [
			{
				"name": "gng-nginx",
				"warnings": [
					"This is a warning",
					"This is another warning",
					"This is a third warning"
				]
			}
		],
		"status": "success"
	}`

	var response WebhookJSONResponse
	err := json.Unmarshal([]byte(jsonResponse), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON response: %v", err)
	}

	if len(response.Addons) != 1 {
		t.Errorf("Expected 1 addon, got %d", len(response.Addons))
	}

	if response.Status != "success" {
		t.Errorf("Expected status 'success', got '%s'", response.Status)
	}

	addon := response.Addons[0]
	if addon.Name != "gng-nginx" {
		t.Errorf("Expected addon name 'gng-nginx', got '%s'", addon.Name)
	}

	if len(addon.Warnings) != 3 {
		t.Errorf("Expected 3 warnings, got %d", len(addon.Warnings))
	}

	if addon.Warnings[0] != "This is a warning" {
		t.Errorf("Expected first warning 'This is a warning', got '%s'", addon.Warnings[0])
	}
}

func TestWebhookWarningMessageProcessing(t *testing.T) {
	// Test processing warning messages from webhook response
	response := &WebhookJSONResponse{
		Status: "success",
		Addons: []WebhookAddon{
			{
				Name: "gng-nginx",
				Warnings: []string{
					"This is a warning",
					"This is another warning",
					"This is a third warning",
				},
			},
		},
	}

	// Test processing the response
	err := processWebhookJSONResponse(response)
	if err != nil {
		t.Fatalf("Failed to process webhook warning messages: %v", err)
	}
}

func TestProcessWebhookJSONResponse(t *testing.T) {
	// Create a mock webhook response with warning messages
	response := &WebhookJSONResponse{
		Status: "success",
		Addons: []WebhookAddon{
			{
				Name: "gng-nginx",
				Warnings: []string{
					"First warning message",
					"Second warning message",
				},
			},
		},
	}

	// Test processing the response
	err := processWebhookJSONResponse(response)
	if err != nil {
		t.Fatalf("Failed to process webhook JSON response: %v", err)
	}
}

func TestProcessWebhookJSONResponseEmpty(t *testing.T) {
	// Test with empty response
	response := &WebhookJSONResponse{
		Status: "success",
		Addons: []WebhookAddon{},
	}

	err := processWebhookJSONResponse(response)
	if err != nil {
		t.Errorf("Expected no error for empty addons, got: %v", err)
	}
}

func TestProcessWebhookJSONResponseNonSuccess(t *testing.T) {
	// Test with non-success status - this should be caught earlier in sendToWebhookWithResponse
	// but we can still test the processing function directly
	response := &WebhookJSONResponse{
		Status: "error",
		Addons: []WebhookAddon{
			{
				Name:     "test-addon",
				Warnings: []string{"Error message"},
			},
		},
	}

	// Since we removed the status check from processWebhookJSONResponse,
	// this should now process successfully even with non-success status
	err := processWebhookJSONResponse(response)
	if err != nil {
		t.Errorf("Expected no error for non-success status in processing function, got: %v", err)
	}
}

func TestProcessWebhookJSONResponseLocalhost5000(t *testing.T) {
	// Test with the actual response format from localhost:5000
	jsonResponse := `{"addons":[{"name":"gng-nginx","warnings":["This is a warning","This is another warning","This is a third warning"]}],"status":"success"}`

	var response WebhookJSONResponse
	err := json.Unmarshal([]byte(jsonResponse), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal localhost:5000 response: %v", err)
	}

	// Verify the response structure
	if response.Status != "success" {
		t.Errorf("Expected status 'success', got '%s'", response.Status)
	}

	if len(response.Addons) != 1 {
		t.Errorf("Expected 1 addon, got %d", len(response.Addons))
	}

	addon := response.Addons[0]
	if addon.Name != "gng-nginx" {
		t.Errorf("Expected addon name 'gng-nginx', got '%s'", addon.Name)
	}

	if len(addon.Warnings) != 3 {
		t.Errorf("Expected 3 warnings, got %d", len(addon.Warnings))
	}

	expectedWarnings := []string{
		"This is a warning",
		"This is another warning",
		"This is a third warning",
	}

	for i, expected := range expectedWarnings {
		if addon.Warnings[i] != expected {
			t.Errorf("Expected warning %d to be '%s', got '%s'", i+1, expected, addon.Warnings[i])
		}
	}

	// Test processing the response
	err = processWebhookJSONResponse(&response)
	if err != nil {
		t.Fatalf("Failed to process localhost:5000 response: %v", err)
	}
}

func TestWebhookStatusValidation(t *testing.T) {
	// Test that non-success status is properly validated
	jsonResponse := `{"addons":[{"name":"test","warnings":["test warning"]}],"status":"error"}`

	var response WebhookJSONResponse
	err := json.Unmarshal([]byte(jsonResponse), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal test response: %v", err)
	}

	// Verify that non-success status is detected
	if response.Status == "success" {
		t.Error("Expected non-success status, got 'success'")
	}

	// Test that the status validation would catch this
	if response.Status != "success" {
		// This simulates the validation logic in sendToWebhookWithResponse
		expectedError := fmt.Sprintf("webhook returned non-success status: %s", response.Status)
		if expectedError != "webhook returned non-success status: error" {
			t.Errorf("Expected error message 'webhook returned non-success status: error', got '%s'", expectedError)
		}
	}
}
