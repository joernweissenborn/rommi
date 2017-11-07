package cmd

import (
	"fmt"
	"rommi/brain"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(tellCmd)
}

var tellCmd = &cobra.Command{
	Use:   "tell COMMAND",
	Short: "Tells Rommi a command.",
	RunE: func(_ *cobra.Command, args []string) error {
		b, err := brain.New()
		if err != nil {
			return err
		}
		b.Run()
		cmd := strings.Join(args, " ")
		fmt.Println("Telling:", cmd)
		b.TellCommandAndWait(cmd)
		return nil
	},
}
