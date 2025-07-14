package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var operatorInfoCmd = &cobra.Command{
	Use:   "info [operator-id]",
	Short: "Get basic information about an operator by ID",
	Long:  "Get basic information about an operator by its ID. This includes name, country, supported amounts, and other details.",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("requires exactly one argument: operator ID")
		}
		_, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("operator ID must be a valid integer")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		operatorID, _ := strconv.ParseInt(args[0], 10, 64)

		svc, err := LoadTopupsService(cmd)
		if err != nil {
			return err
		}

		operator, err := svc.Topups().GetOperatorByID(operatorID)
		if err != nil {
			return err
		}

		fmt.Printf("Operator Information:\n")
		fmt.Printf("  ID: %d\n", operator.OperatorID)
		fmt.Printf("  Name: %s\n", operator.Name)
		fmt.Printf("  Country: %s (%s)\n", operator.Country.Name, operator.Country.IsoName)
		fmt.Printf("  Bundle: %t\n", operator.Bundle)
		fmt.Printf("  Data: %t\n", operator.Data)
		fmt.Printf("  Pin: %t\n", operator.Pin)
		fmt.Printf("  Supports Local Amounts: %t\n", operator.SupportsLocalAmounts)
		fmt.Printf("  Denomination Type: %s\n", operator.DenominationType)
		fmt.Printf("  Sender Currency: %s (%s)\n", operator.SenderCurrencyCode, operator.SenderCurrencySymbol)
		fmt.Printf("  Destination Currency: %s (%s)\n", operator.DestinationCurrencyCode, operator.DestinationCurrencySymbol)
		fmt.Printf("  Commission: %.2f%%\n", operator.Commission)
		fmt.Printf("  International Discount: %.2f%%\n", operator.InternationalDiscount)
		fmt.Printf("  Local Discount: %.2f%%\n", operator.LocalDiscount)
		fmt.Printf("  Most Popular Amount: %.2f\n", operator.MostPopularAmount)
		fmt.Printf("  Most Popular Local Amount: %.2f\n", operator.MostPopularLocalAmount)
		fmt.Printf("  Min Amount: %.2f\n", operator.MinAmount)
		fmt.Printf("  Max Amount: %.2f\n", operator.MaxAmount)
		fmt.Printf("  Local Min Amount: %.2f\n", operator.LocalMinAmount)
		fmt.Printf("  Local Max Amount: %.2f\n", operator.LocalMaxAmount)

		if operator.Fx.Rate > 0 {
			fmt.Printf("  FX Rate: %.4f %s\n", operator.Fx.Rate, operator.Fx.CurrencyCode)
		}

		if len(operator.FixedAmounts) > 0 {
			fmt.Printf("  Fixed Amounts: %v\n", operator.FixedAmounts)
		}

		if len(operator.LocalFixedAmounts) > 0 {
			fmt.Printf("  Local Fixed Amounts: %v\n", operator.LocalFixedAmounts)
		}

		if len(operator.SuggestedAmounts) > 0 {
			fmt.Printf("  Suggested Amounts: %v\n", operator.SuggestedAmounts)
		}

		if len(operator.SuggestedAmountsMap) > 0 {
			fmt.Printf("  Suggested Amounts Map:\n")
			for _, amount := range operator.SuggestedAmountsMap {
				fmt.Printf("    Pay: %.2f, Sent: %.2f\n", amount.Pay, amount.Sent)
			}
		}

		if len(operator.LogoUrls) > 0 {
			fmt.Printf("  Logo URLs: %v\n", operator.LogoUrls)
		}

		return nil
	},
}

func init() {
	operatorsCmd.AddCommand(operatorInfoCmd)
}
