package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var operatorGetCmd = &cobra.Command{
	Use:   "get [country] [operator-name]",
	Short: "Get basic information about an operator by country and name",
	Long:  "Get basic information about an operator by country code and operator name. This includes all operator details like supported amounts, fees, and configuration.",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return fmt.Errorf("requires exactly two arguments: country code and operator name")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		country := args[0]
		operatorName := args[1]

		svc, err := LoadTopupsService(cmd)
		if err != nil {
			return err
		}

		operator, err := svc.Topups().SearchOperator(country, operatorName)
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

		if operator.MostPopularAmount != nil {
			fmt.Printf("  Most Popular Amount: %.2f\n", *operator.MostPopularAmount)
		} else {
			fmt.Printf("  Most Popular Amount: <not set>\n")
		}

		if operator.MostPopularLocalAmount != nil {
			fmt.Printf("  Most Popular Local Amount: %.2f\n", *operator.MostPopularLocalAmount)
		} else {
			fmt.Printf("  Most Popular Local Amount: <not set>\n")
		}

		if operator.MinAmount != nil {
			fmt.Printf("  Min Amount: %.2f\n", *operator.MinAmount)
		} else {
			fmt.Printf("  Min Amount: <not set>\n")
		}

		if operator.MaxAmount != nil {
			fmt.Printf("  Max Amount: %.2f\n", *operator.MaxAmount)
		} else {
			fmt.Printf("  Max Amount: <not set>\n")
		}

		if operator.LocalMinAmount != nil {
			fmt.Printf("  Local Min Amount: %.2f\n", *operator.LocalMinAmount)
		} else {
			fmt.Printf("  Local Min Amount: <not set>\n")
		}

		if operator.LocalMaxAmount != nil {
			fmt.Printf("  Local Max Amount: %.2f\n", *operator.LocalMaxAmount)
		} else {
			fmt.Printf("  Local Max Amount: <not set>\n")
		}

		if operator.Fx.Rate > 0 {
			fmt.Printf("  FX Rate: %.4f %s\n", operator.Fx.Rate, operator.Fx.CurrencyCode)
		}

		if len(operator.GetFixedAmounts()) > 0 {
			fmt.Printf("  Fixed Amounts: %v\n", operator.GetFixedAmounts())
		}

		if len(operator.GetLocalFixedAmounts()) > 0 {
			fmt.Printf("  Local Fixed Amounts: %v\n", operator.GetLocalFixedAmounts())
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
	operatorsCmd.AddCommand(operatorGetCmd)
}
