package cmd

import (
	"os"
	"rommi/modules/brain"
	"rommi/lib/train"

	"github.com/spf13/cobra"
)

func init() {
	trainCmd.AddCommand(trainSentenceCmd)
	trainSentenceCmd.AddCommand(trainSentenceUpdateCmd)
	trainSentenceCmd.AddCommand(trainSentenceShowCmd)
}

var trainSentenceCmd = &cobra.Command{
	Use:   "sentence",
	Short: "Create, update and modify rommi's training sentences.",
}

var trainSentenceUpdateCmd = &cobra.Command{
	Use:   "update [MODELPATH]",
	Short: "Updates the sentences with rommi's currently known sentences.",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		b, err := brain.New()
		if err != nil {
			return
		}
		b.Run()
		cmd.Println("Getting Rommi's Sentences")
		sentences := b.GetSentences()
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
		cmd.Println("Updating Sentences")
		err = t.AddSentence(sentences...)
		if err != nil {
			return
		}
		cmd.Println("Success")
		return
	},
}

var trainSentenceShowCmd = &cobra.Command{
	Use:   "show [MODELPATH]",
	Short: "Show all sentences in the training database.",
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
		cmd.Println("")

		cmd.Println("Sentences in database")
		cmd.Println("=====================")
		cmd.Println("")
		for _, id := range t.GetAllIds() {
			s,_:=t.GetSentence(id)
			cmd.Printf("%s: %s\n", id.FullString(), s)
		}
		return
	},
}
