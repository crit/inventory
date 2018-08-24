package cmd

import (
	"os"

	"github.com/crit/inventory/internal/inventory"
	"github.com/spf13/cobra"
)

// reportCmd represents the report command
var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Report on current levels of each item.",
	Long:  `Report on current levels of each item.`,
	Run: func(cmd *cobra.Command, args []string) {
		inventory.Report(os.Stdout)
	},
}

func init() {
	RootCmd.AddCommand(reportCmd)
}
