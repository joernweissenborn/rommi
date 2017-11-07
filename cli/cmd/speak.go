package cmd

import (
	"fmt"
	"rommi/voice"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(speakCmd)
}

var speakCmd = &cobra.Command{
	Use:   "speak SENTENCE",
	Short: "Makes Rommi speak a sentence.",
	RunE: func(_ *cobra.Command, args []string) error {
		v, err := voice.New()
		if err != nil {
			return err
		}
		v.Run()
		s := strings.Join(args, " ")
		fmt.Println("Speaking:", s)
		v.SpeakAndWait(s)
		return nil
	},
}
