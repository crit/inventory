package cmd

import (
	"errors"
	"fmt"
	"os/user"

	"github.com/crit/inventory/internal/inventory"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add to an inventory items count by name.",
	Long:  `Add to an inventory items count by name.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args)%2 != 0 {
			return errors.New("add needs count/name pairs (example: inventory add 1 sd-card)")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		var errs []error

		usr, err := user.Current()

		if err != nil {
			usr.Username = "unknown"
		}

		for i := 0; i < len(args); i = i + 2 {
			change, name := args[i], args[i+1]
			if err := inventory.NewEntry(change, name, usr.Username); err != nil {
				errs = append(errs, fmt.Errorf("%s %s error: %v", change, name, err))
			} else {
				fmt.Printf("added %s %s\n", change, name)
			}
		}

		if len(errs) != 0 {
			fmt.Printf("  encountered errors: %q", errs)
		}
	},
}

func init() {
	RootCmd.AddCommand(addCmd)
}
