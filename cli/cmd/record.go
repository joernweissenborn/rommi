package cmd

import (
	"errors"
	"io/ioutil"
	"rommi/modules/ears"

	"github.com/ThingiverseIO/console"
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
		wav, err := record()
		if err != nil {
			return err
		}
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

func record() (wav []byte, err error) {
	e, err := ears.New()
	if err != nil {
		return
	}
	e.Run()
	console.Println("Waiting for ears.")
	e.WaitUntilConnected()
	if console.AskEnterOrAbort("Ready To Record, press enter to start or 'q' enter to abort", "q") {
		err = errors.New("Aborted")
		return
	}
	if !e.StartRecording() {
		err = errors.New("Failed to start Recording")
		return
	}
	console.AskEnter("Started Recording, press enter to stop")
	if !e.StopRecording() {
		err = errors.New("An Error Happened during Recording")
		return
	}
	console.Println("Stopped Recording")
	wav = e.GetRecordedWav()
	return
}
