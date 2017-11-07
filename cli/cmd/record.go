package cmd

import (
	"bufio"
	"io/ioutil"
	"os"
	"rommi/modules/ears"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(recordCmd)
}

var recordCmd = &cobra.Command{
	Use:   "record FILENAME",
	Short: "Records audio with rommi's ears and saves it as WAVE file",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return cmd.Usage()
		}
		e, err := ears.New()
		if err != nil {
			return err
		}
		e.Run()
		if !e.StartRecording() {
			cmd.Println("Failed to start Recording")
			return nil
		}
		cmd.Println("Started Recording, press return to stop")
		consoleReader := bufio.NewReaderSize(os.Stdin, 1)
		consoleReader.ReadByte()
		if !e.StopRecording() {
			cmd.Println("An Error Happened during Recording")
			return nil
		}
		cmd.Println("Stopped Recording")

		wav := e.GetRecordedWav()
		cmd.Println("Lenght of recorded WAVE is:", len(wav))
		cmd.Println("Writing File to", args[0])
		err = ioutil.WriteFile(args[0], wav, 0644)
		if err != nil {
			cmd.Println("Error Writing File", err)
		} else {
			cmd.Println("Success")
		}
		return nil
	},
}
