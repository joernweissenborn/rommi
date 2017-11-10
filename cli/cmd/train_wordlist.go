package cmd

import (
	"os"
	"rommi/lib/train"

	"github.com/spf13/cobra"
)

func init() {
	trainCmd.AddCommand(trainWordListCmd)
}

var trainWordListCmd = &cobra.Command{
	Use:   "words [MODELPATH]",
	Short: "View rommi's training words.",
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
		wl := t.GetWordList()
		cmd.Println("Success")
		cmd.Println("Words in database")
		cmd.Println("=====================")
		cmd.Println("")
		for w := range wl {
			cmd.Printf("%s\n", w)
		}
		return
	},
}
