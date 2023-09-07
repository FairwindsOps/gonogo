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
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"k8s.io/klog"

	"github.com/fairwindsops/gonogo/pkg/helm"
	"github.com/fairwindsops/gonogo/pkg/validate"
)

var bundleFile []string

func init() {
	rootCmd.AddCommand(checkCmd)
	checkCmd.PersistentFlags().StringSliceVarP(&bundleFile, "bundle", "b", []string{}, "bundle file(s) to use")
}

var checkCmd = &cobra.Command{
	Use:     "check [path to Bundle config file]",
	Short:   "Check for Helm releases that can be updated",
	Long:    `Check for Helm releases that can be updated`,
	PreRunE: validateArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 0 {
			for _, fileName := range bundleFile {
				bundleFile = append(bundleFile, fileName)
			}
		}

		config := &validate.Config{
			Helm:   helm.NewHelm(),
			Bundle: bundleFile,
		}

		out, err := config.Validate()
		if err != nil {
			klog.Error(err)
		}
		fmt.Println(out)
	},
}

func validateArgs(cmd *cobra.Command, args []string) error {
	if len(args) != 0 {
		_, err := os.Stat(args[0])

		if os.IsNotExist(err) {
			return fmt.Errorf("bundle file %s does not exist", args[0])
		}
	}
	return nil
}
