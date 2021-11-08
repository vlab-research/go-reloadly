package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// topupsCmd represents the topups command
var topupsCmd = &cobra.Command{
	Use:   "topups",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("topups called")
	},
}

func init() {
	rootCmd.AddCommand(topupsCmd)
}
