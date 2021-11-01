/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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

	"github.com/spf13/cobra"
)

// readSpecFileCmd represents the readSpecFile command
var readSpecFileCmd = &cobra.Command{
	Use:   "read [file to process]",
	Short: "Parse file",
	Long: `Provide file that adheres to the bundle spec for parsing`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Reading spec file")


	},
}

var (
	filename string
)

func init() {
	rootCmd.AddCommand(readSpecFileCmd)
	readSpecFileCmd.Flags().StringVarP(&filename, "filename", "f", "", "A bundle spec file to parse")
}
 