package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(trainCmd)
}

var trainCmd = &cobra.Command{
	Use:   "train",
	Short: "Train Rommi.",
}
