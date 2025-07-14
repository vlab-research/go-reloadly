package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var operatorListCmd = &cobra.Command{
	Use:   "list [country]",
	Short: "List all operators for a given country",
	Long:  "List all available operators for a given country code. This shows basic information about each operator including ID, name, and supported services.",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("requires exactly one argument: country code")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		country := args[0]

		svc, err := LoadTopupsService(cmd)
		if err != nil {
			return err
		}

		operators, err := svc.Topups().OperatorsByCountry(country)
		if err != nil {
			return err
		}

		if len(operators) == 0 {
			fmt.Printf("No operators found for country: %s\n", country)
			return nil
		}

		fmt.Printf("Operators for %s:\n", country)
		fmt.Printf("%-8s %-30s %-15s %-14s %-10s %-10s %-20s\n", "ID", "Name", "Country", "Denomination", "Bundle", "Data", "FixedAmounts")
		fmt.Printf("%-8s %-30s %-15s %-14s %-10s %-10s %-20s\n", "---", "----", "-------", "-------------", "------", "----", "------------")

		for _, op := range operators {
			fixed := ""
			if op.DenominationType == "FIXED" && len(op.FixedAmounts) > 0 {
				fixed = formatFixedAmounts(op.FixedAmounts, 5)
			}
			fmt.Printf("%-8d %-30s %-15s %-14s %-10t %-10t %-20s\n",
				op.OperatorID,
				truncateString(op.Name, 28),
				op.Country.Name,
				op.DenominationType,
				op.Bundle,
				op.Data,
				fixed)
		}

		fmt.Printf("\nTotal operators: %d\n", len(operators))
		return nil
	},
}

// Helper function to truncate strings for table display
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// Helper to format fixed amounts, truncating if too many
func formatFixedAmounts(amounts []float64, max int) string {
	if len(amounts) == 0 {
		return ""
	}
	end := len(amounts)
	truncated := false
	if max > 0 && len(amounts) > max {
		end = max
		truncated = true
	}
	str := ""
	for i, amt := range amounts[:end] {
		if i > 0 {
			str += ", "
		}
		str += fmt.Sprintf("%.2f", amt)
	}
	if truncated {
		str += ", ..."
	}
	return str
}

func init() {
	operatorsCmd.AddCommand(operatorListCmd)
}
