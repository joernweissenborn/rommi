package cmd

import (
	"os"
	"rommi/lib/train"
	"rommi/modules/brain"

	"github.com/spf13/cobra"
)

func init() {
	trainCmd.AddCommand(trainTriggerWordCmd)
	trainTriggerWordCmd.AddCommand(trainTriggerWordUpdateCmd)
	trainTriggerWordCmd.AddCommand(trainTriggerWordShowCmd)
}

var trainTriggerWordCmd = &cobra.Command{
	Use:   "triggerword",
	Short: "Create, update and modify rommi's training sentences.",
}

var trainTriggerWordUpdateCmd = &cobra.Command{
	Use:   "update [MODELPATH]",
	Short: "Updates the triggerword with rommi's current triggerword.",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		b, err := brain.New()
		if err != nil {
			return
		}
		b.Run()
		cmd.Println("Getting Rommi's triggerword")
		triggerword := b.GetTriggerWord()
		cmd.Println("Success")

		path, err := os.Getwd()
		if err != nil {
			return
		}
		if len(args) != 0 {
			path = args[0]
		}
		cmd.Println("Opening Model Directory:", path)
		t, err := train.Open(path)
		if err != nil {
			return
		}
		cmd.Println("Success")
		cmd.Println("Updating triggerword")
		err = t.SetTriggerWord(triggerword)
		if err != nil {
			return
		}
		cmd.Println("Success")
		return
	},
}

var trainTriggerWordShowCmd = &cobra.Command{
	Use:   "show [MODELPATH]",
	Short: "Show the triggerword in the training database.",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		path, err := os.Getwd()
		if err != nil {
			return
		}
		if len(args) != 0 {
			path = args[0]
		}
		cmd.Println("Opening Model Directory:", path)
		t, err := train.Open(path)
		if err != nil {
			return
		}
		cmd.Println("Success")
		word := t.GetTriggerWord()
		cmd.Println("Triggerword in database")
		cmd.Println("=======================")
		cmd.Println("")
		cmd.Println(word)
		return
	},
}
