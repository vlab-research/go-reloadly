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

var singleCmd = &cobra.Command{
	Use:   "single",
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

		number := args[0]
		amount, err := strconv.ParseFloat(args[1], 64)
		if err != nil {
			return err
		}

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
		tolerance, err := cmd.Flags().GetFloat64("tolerance")
		if err != nil {
			return err
		}

		var res *reloadly.TopupResponse

		if operatorName != "" {
			fmt.Println(fmt.Sprintf("Using operator: %v", operatorName))
			res, err = svc.Topups().FindOperator(country, operatorName).SuggestedAmount(tolerance).AutoFallback().Topup(number, amount)
		} else {
			t := svc.Topups()
			res, err = t.AutoDetect(country).SuggestedAmount(tolerance).Topup(number, amount)
			fmt.Println(fmt.Sprintf("Autodetected operator: %v", t.GetSetOperator().Name))
		}

		if err != nil {
			fmt.Println(err)
			return nil
		}

		fmt.Printf("Topup response: %v", res)


		return nil
	},
}

func init() {
	topupsCmd.AddCommand(singleCmd)

	singleCmd.Flags().Float64P("tolerance", "t", 0.0, "tolerance for topup")
	singleCmd.Flags().String("operator", "", "operator")
}
