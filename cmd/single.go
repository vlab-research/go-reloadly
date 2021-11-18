package cmd

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/vlab-research/go-reloadly/reloadly"
)

var singleCmd = &cobra.Command{
	Use:   "single",
	Short: "Make airtime recharge to a single mobile number",
	Long:  "Make airtime recharge to a single mobile number",
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
