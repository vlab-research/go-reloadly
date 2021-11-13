package cmd

import (
	"github.com/spf13/cobra"
)

var topupsCmd = &cobra.Command{
	Use:   "topups",
	Short: "Make airtime recharge to mobile numbers",
	Long:  "Make airtime recharge to mobile numbers",
}

func init() {
	rootCmd.AddCommand(topupsCmd)
}
