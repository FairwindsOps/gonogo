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
	"log"
	"os"

	"github.com/fairwindsops/hall-monitor/pkg/bundle"
	"github.com/spf13/cobra"
)

// readSpecFileCmd represents the readSpecFile command
var readSpecFileCmd = &cobra.Command{
	Use:     "read [file to process]",
	Short:   "Parse file",
	Long:    `Provide file that adheres to the bundle spec for parsing`,
	PreRunE: validateArgs,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Reading spec %s\n", args[0])
		config, err := bundle.ReadConfig(args[0])
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%v\n", config)
	},
}

func validateArgs(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("you must specify a spec file")
	}

	_, err := os.Stat(args[0])
	if os.IsNotExist(err) {
		return fmt.Errorf("spec file %s does not exist", args[0])
	}
	return err
}

func init() {
	rootCmd.AddCommand(readSpecFileCmd)
}
