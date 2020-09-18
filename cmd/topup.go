/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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
	"errors"
	"fmt"

	// "github.com/nandanrao/go-reloadly/reloadly"
	"strconv"

	"github.com/nandanrao/go-reloadly/reloadly"
	"github.com/spf13/cobra"
)

// TopupCmd represents the Topup command
var TopupCmd = &cobra.Command{
	Use:   "topup",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 3 {
			return errors.New("requires at least mobile and amount and country")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		mobile := args[0]
		amount, err := strconv.ParseFloat(args[1], 64)
		country := args[2]

		svc, err := LoadService(cmd)
		if err != nil {
			return err
		}
		fmt.Println("Authorized with Reloadly")

		operatorName, err := cmd.Flags().GetString("operator")
		if err != nil {
			return err
		}

		var op *reloadly.Operator
		if operatorName != "" {
			op, err = svc.SearchOperator(country, operatorName)
		} else {
			op, err = svc.OperatorsAutoDetect(mobile, country)

		}
		if err != nil {
			fmt.Println(err)
			return nil
		}

		fmt.Println(fmt.Sprintf("Using operator: %v", op.Name))

		res, err := svc.TopupBySuggestedAmount("cli", mobile, op, amount, 0.0)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		fmt.Printf("Topup response: %v", res)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(TopupCmd)

	// TopupCmd.Flags().Float64P("amount", "a", 0.0, "amount for topup")
	TopupCmd.Flags().Float64P("tolerance", "t", 0.0, "tolerance for topup")
	TopupCmd.Flags().String("operator", "", "operator")
}
