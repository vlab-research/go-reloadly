package cmd

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/vlab-research/go-reloadly/reloadly"
)

// operatorsCmd represents the operators command
var operatorsCmd = &cobra.Command{
	Use:   "operators",
	Short: "Check operator support for a given amount.",
	Long:  "Check operator support for a given amount.",
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
		if err != nil {
			return err
		}

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
}
