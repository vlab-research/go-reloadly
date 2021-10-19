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
nnlimitations under the License.
*/
package cmd

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/nandanrao/go-reloadly/reloadly"
	"github.com/spf13/cobra"
)

// operatorsCmd represents the operators command
var operatorsCmd = &cobra.Command{
	Use:   "operators",
	Short: "Check operator support for a given amount",
	Long:  `Check operator support for a given amount.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 3 {
			return errors.New("requires country and amount and tolerance")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {

		country := args[0]
		target, err := strconv.ParseFloat(args[1], 64)
		if err != nil {
			return err
		}

		tol, err := strconv.ParseFloat(args[2], 64)
		if err != nil {
			return err
		}

		svc, err := LoadService(cmd)

		ops, err := svc.OperatorsByCountry(country)

		if err != nil {
			return err
		}

		for _, op := range ops {
			amt, err := reloadly.GetSuggestedAmount(&op, target, tol)
			if err != nil {
				fmt.Printf("failed with %v \n", op.Name)
				continue
			}
			fmt.Println(op.Name)
			fmt.Println(amt)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(operatorsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// operatorsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// operatorsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
