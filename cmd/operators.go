package cmd

import (
	"github.com/spf13/cobra"
)

var operatorsCmd = &cobra.Command{
	Use:   "operators",
	Short: "Make airtime recharge to mobile numbers",
	Long:  "Make airtime recharge to mobile numbers",
}

func init() {
	topupsCmd.AddCommand(operatorsCmd)
}
