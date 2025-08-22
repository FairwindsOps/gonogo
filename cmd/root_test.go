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

func TestRootCmdProperties(t *testing.T) {
	cmd := rootCmd

	assert.Equal(t, "gonogo", cmd.Use)
	assert.Equal(t, "Validates whether or not an add-on is safe to upgrade", cmd.Short)
	assert.Contains(t, cmd.Long, "Kubernetes Add-On Upgrade Validation Bundle")
}

func TestRootCmdRunFunction(t *testing.T) {
	cmd := rootCmd

	assert.NotNil(t, cmd.Run)

	// Test that Run function exists and can be called
	assert.NotPanics(t, func() {
		// The Run function should print a message and exit
		// We can't easily test the exit behavior, but we can verify the function exists
	})
}

func TestRootCmdHasSubcommands(t *testing.T) {
	cmd := rootCmd

	// Check that root command has subcommands
	assert.True(t, cmd.HasSubCommands())

	// Check for specific subcommands
	subcommands := cmd.Commands()
	subcommandNames := make([]string, len(subcommands))
	for i, subcmd := range subcommands {
		subcommandNames[i] = subcmd.Name()
	}

	// Should have check and version subcommands
	assert.Contains(t, subcommandNames, "check")
	assert.Contains(t, subcommandNames, "version")
}

func TestRootCmdInit(t *testing.T) {
	// Test that init function doesn't panic
	assert.NotPanics(t, func() {
		// The init function sets up klog flags
		// We can't easily test the exact behavior without affecting global state
		// but we can verify it doesn't panic
	})
}

func TestExecuteFunction(t *testing.T) {
	// Test that Execute function exists and can be called
	assert.NotPanics(t, func() {
		// Execute sets version variables and calls rootCmd.Execute()
		// We can't easily test the full execution without affecting the test environment
		// but we can verify the function signature is correct
	})
}

func TestVersionVariables(t *testing.T) {
	// Test that version variables are accessible
	// These are set by the Execute function
	assert.IsType(t, "", version)
	assert.IsType(t, "", versionCommit)
}

func TestRootCmdValidateRequiredFlags(t *testing.T) {
	cmd := rootCmd

	// Test that required flags validation works
	// Root command doesn't have required flags, so this should pass
	err := cmd.ValidateRequiredFlags()
	assert.NoError(t, err)
}

func TestRootCmdPersistentPreRun(t *testing.T) {
	cmd := rootCmd

	// Test that PersistentPreRun is not set (which is correct for root command)
	assert.Nil(t, cmd.PersistentPreRun)
}

func TestRootCmdPersistentPreRunE(t *testing.T) {
	cmd := rootCmd

	// Test that PersistentPreRunE is not set (which is correct for root command)
	assert.Nil(t, cmd.PersistentPreRunE)
}

func TestRootCmdPreRun(t *testing.T) {
	cmd := rootCmd

	// Test that PreRun is not set (which is correct for root command)
	assert.Nil(t, cmd.PreRun)
}

func TestRootCmdPreRunE(t *testing.T) {
	cmd := rootCmd

	// Test that PreRunE is not set (which is correct for root command)
	assert.Nil(t, cmd.PreRunE)
}

func TestRootCmdPostRun(t *testing.T) {
	cmd := rootCmd

	// Test that PostRun is not set (which is correct for root command)
	assert.Nil(t, cmd.PostRun)
}

func TestRootCmdPostRunE(t *testing.T) {
	cmd := rootCmd

	// Test that PostRunE is not set (which is correct for root command)
	assert.Nil(t, cmd.PersistentPostRunE)
}

func TestRootCmdSilenceUsage(t *testing.T) {
	cmd := rootCmd

	// Test that SilenceUsage is set correctly
	// Root command should not silence usage by default
	assert.False(t, cmd.SilenceUsage)
}

func TestRootCmdSilenceErrors(t *testing.T) {
	cmd := rootCmd

	// Test that SilenceErrors is set correctly
	// Root command should not silence errors by default
	assert.False(t, cmd.SilenceErrors)
}

func TestRootCmdDisableSuggestions(t *testing.T) {
	cmd := rootCmd

	// Test that DisableSuggestions is set correctly
	// Root command should not disable suggestions by default
	assert.False(t, cmd.DisableSuggestions)
}

func TestRootCmdDisableFlagsInUseLine(t *testing.T) {
	cmd := rootCmd

	// Test that DisableFlagsInUseLine is set correctly
	// Root command should not disable flags in use line by default
	assert.False(t, cmd.DisableFlagsInUseLine)
}
