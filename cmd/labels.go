/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

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
// cmd package
package cmd

import (
	"github.com/donofden/tale-gmail/pkg/talegmail/labels"
	"github.com/spf13/cobra"
)

// labelsCmd represents the labels command
var labelsCmd = &cobra.Command{
	Use:   "labels",
	Short: "Get Labels used in Gmail Account",
	Long:  `Get Labels used in Gmail Account as a list`,
	Run: func(cmd *cobra.Command, args []string) {
		labels.ListLabels()
	},
}

func init() {
	rootCmd.AddCommand(labelsCmd)
}
