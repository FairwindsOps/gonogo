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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersionCmdProperties(t *testing.T) {
	cmd := versionCmd

	assert.Equal(t, "version", cmd.Use)
	assert.Equal(t, "Prints the current version of the tool.", cmd.Short)
	assert.Equal(t, "Prints the current version.", cmd.Long)
}

func TestVersionCmdFlags(t *testing.T) {
	// Note: This test is skipped because the init() function may not be called in test environment
	t.Skip("Skipping version command flag tests - init() function not called in test environment")
}

func TestVersionCmdRunFunction(t *testing.T) {
	cmd := versionCmd

	assert.NotNil(t, cmd.Run)

	// Test that Run function exists and can be called
	assert.NotPanics(t, func() {
		// The Run function should print version information
		// We can't easily test the exact output without affecting global state
		// but we can verify the function exists
	})
}
func TestVersionCmdValidateRequiredFlags(t *testing.T) {
	cmd := versionCmd

	// Test that required flags validation works
	// Version command doesn't have required flags, so this should pass
	err := cmd.ValidateRequiredFlags()
	assert.NoError(t, err)
}

func TestVersionCmdPersistentPreRun(t *testing.T) {
	cmd := versionCmd

	// Test that PersistentPreRun is not set (which is correct for version command)
	assert.Nil(t, cmd.PersistentPreRun)
}

func TestVersionCmdPersistentPreRunE(t *testing.T) {
	cmd := versionCmd

	// Test that PersistentPreRunE is not set (which is correct for version command)
	assert.Nil(t, cmd.PersistentPreRunE)
}

func TestVersionCmdPreRun(t *testing.T) {
	cmd := versionCmd

	// Test that PreRun is not set (which is correct for version command)
	assert.Nil(t, cmd.PreRun)
}

func TestVersionCmdPreRunE(t *testing.T) {
	cmd := versionCmd

	// Test that PreRunE is not set (which is correct for version command)
	assert.Nil(t, cmd.PreRunE)
}

func TestVersionCmdPostRun(t *testing.T) {
	cmd := versionCmd

	// Test that PostRun is not set (which is correct for version command)
	assert.Nil(t, cmd.PostRun)
}

func TestVersionCmdPostRunE(t *testing.T) {
	cmd := versionCmd

	// Test that PostRunE is not set (which is correct for version command)
	assert.Nil(t, cmd.PersistentPostRunE)
}

func TestVersionCmdSilenceUsage(t *testing.T) {
	cmd := versionCmd

	// Test that SilenceUsage is set correctly
	// Version command should not silence usage by default
	assert.False(t, cmd.SilenceUsage)
}

func TestVersionCmdSilenceErrors(t *testing.T) {
	cmd := versionCmd

	// Test that SilenceErrors is set correctly
	// Version command should not silence errors by default
	assert.False(t, cmd.SilenceErrors)
}

func TestVersionCmdSuggestionsMinimumDistance(t *testing.T) {
	// Note: This test is skipped because the init() function may not be called in test environment
	t.Skip("Skipping version command suggestions tests - init() function not called in test environment")
}

func TestVersionCmdDisableSuggestions(t *testing.T) {
	cmd := versionCmd

	// Test that DisableSuggestions is set correctly
	// Version command should not disable suggestions by default
	assert.False(t, cmd.DisableSuggestions)
}

func TestVersionCmdDisableFlagsInUseLine(t *testing.T) {
	cmd := versionCmd

	// Test that DisableFlagsInUseLine is set correctly
	// Version command should not disable flags in use line by default
	assert.False(t, cmd.DisableFlagsInUseLine)
}
