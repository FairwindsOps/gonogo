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
package helm

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// ChartVersion represents a version of a helm chart
type ChartVersion struct {
	Version    string   `json:"version"`
	AppVersion string   `json:"appVersion"`
	URLs       []string `json:"urls"`
	Digest     string   `json:"digest"`
	Created    string   `json:"created"`
	Removed    bool     `json:"removed,omitempty"`
}

// ChartEntry represents a chart in the repository index
type ChartEntry struct {
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	Home        string                  `json:"home"`
	Sources     []string                `json:"sources"`
	Maintainers []map[string]string     `json:"maintainers"`
	Icon        string                  `json:"icon"`
	APIVersion  string                  `json:"apiVersion"`
	Condition   string                  `json:"condition"`
	Tags        string                  `json:"tags"`
	Deprecated  bool                    `json:"deprecated"`
	Annotations map[string]string       `json:"annotations"`
	KubeVersion string                  `json:"kubeVersion"`
	Versions    map[string]ChartVersion `json:"versions"`
}

// RepositoryIndex represents the index.yaml of a helm repository
type RepositoryIndex struct {
	APIVersion string                  `json:"apiVersion"`
	Generated  string                  `json:"generated"`
	Entries    map[string][]ChartEntry `json:"entries"`
}

// ValidateChartVersion checks if a specific version exists in a helm repository
func (h *Helm) ValidateChartVersion(repoURL, chartName, version string) (*ChartVersion, error) {
	// Fetch the repository index
	index, err := h.fetchRepositoryIndex(repoURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch repository index: %w", err)
	}

	// Find the chart
	entries, exists := index.Entries[chartName]
	if !exists {
		return nil, fmt.Errorf("chart '%s' not found in repository", chartName)
	}

	// Find the specific version
	for _, entry := range entries {
		for ver, chartVersion := range entry.Versions {
			if ver == version {
				return &chartVersion, nil
			}
		}
	}

	return nil, fmt.Errorf("version '%s' not found for chart '%s' in repository", version, chartName)
}

// GetAvailableVersions returns all available versions for a chart in a repository
func (h *Helm) GetAvailableVersions(repoURL, chartName string) ([]string, error) {
	// Fetch the repository index
	index, err := h.fetchRepositoryIndex(repoURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch repository index: %w", err)
	}

	// Find the chart
	entries, exists := index.Entries[chartName]
	if !exists {
		return nil, fmt.Errorf("chart '%s' not found in repository", chartName)
	}

	var versions []string
	for _, entry := range entries {
		for version := range entry.Versions {
			versions = append(versions, version)
		}
	}

	return versions, nil
}

// fetchRepositoryIndex fetches the index.yaml from a helm repository
func (h *Helm) fetchRepositoryIndex(repoURL string) (*RepositoryIndex, error) {
	// Ensure the URL ends with /index.yaml
	if !strings.HasSuffix(repoURL, "/index.yaml") {
		if strings.HasSuffix(repoURL, "/") {
			repoURL += "index.yaml"
		} else {
			repoURL += "/index.yaml"
		}
	}

	// Make HTTP request
	resp, err := http.Get(repoURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch repository index: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch repository index: HTTP %d", resp.StatusCode)
	}

	// Read and parse the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read repository index: %w", err)
	}

	var index RepositoryIndex
	if err := json.Unmarshal(body, &index); err != nil {
		return nil, fmt.Errorf("failed to parse repository index: %w", err)
	}

	return &index, nil
}
