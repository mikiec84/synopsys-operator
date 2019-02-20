// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package synopsysctl

import (
	"fmt"

	"github.com/spf13/cobra"
)

// describeCmd represents the describe command
var describeCmd = &cobra.Command{
	Use:   "describe",
	Short: "Print a detailed description of the selected resource",
	Args: func(cmd *cobra.Command, args []string) error {
		numArgs := 1
		if len(args) < numArgs {
			return fmt.Errorf("Not enough arguments")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Describing Non-Synopsys Resource")
		kubeCmdArgs := append([]string{"describe"}, args...)
		out, err := RunKubeCmd(kubeCmdArgs...)
		if err != nil {
			fmt.Printf("Error Describing the Resource with KubeCmd: %s\n", err)
		}
		fmt.Printf("%+v\n", out)
	},
}

var describeBlackduckCmd = &cobra.Command{
	Use:   "blackduck",
	Short: "Describe an instance of Blackduck",
	Run: func(cmd *cobra.Command, args []string) {
		// Read Commandline Parameters
		namespace = args[0]

		fmt.Println("Describing Blackduck")
		out, err := RunKubeCmd("describe", "blackduck", namespace, "-n", namespace)
		if err != nil {
			fmt.Printf("Error Describing the Blackduck with KubeCmd: %s\n", err)
		}
		fmt.Printf("%+v\n", out)
	},
}

var describeOpsSightCmd = &cobra.Command{
	Use:   "opssight",
	Short: "Describe an instance of OpsSight",
	Run: func(cmd *cobra.Command, args []string) {
		// Read Commandline Parameters
		namespace = args[0]

		fmt.Println("Describing OpsSight")
		out, err := RunKubeCmd("describe", "opssight", namespace, "-n", namespace)
		if err != nil {
			fmt.Printf("Error Describing the OpsSight with KubeCmd: %s\n", err)
		}
		fmt.Printf("%+v\n", out)
	},
}

var describeAlertCmd = &cobra.Command{
	Use:   "alert",
	Short: "Describe an instance of Alert",
	Run: func(cmd *cobra.Command, args []string) {
		// Read Commandline Parameters
		namespace = args[0]

		fmt.Println("Describing Alert")
		out, err := RunKubeCmd("describe", "alert", namespace, "-n", namespace)
		if err != nil {
			fmt.Printf("Error Describing the Alert with KubeCmd: %s\n", err)
		}
		fmt.Printf("%+v\n", out)
	},
}

func init() {
	describeCmd.DisableFlagParsing = true
	rootCmd.AddCommand(describeCmd)
	describeCmd.AddCommand(describeBlackduckCmd)
	describeCmd.AddCommand(describeOpsSightCmd)
	describeCmd.AddCommand(describeAlertCmd)
}
