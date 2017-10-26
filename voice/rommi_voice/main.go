package main

import (
	"os"

	"github.com/ThingiverseIO/logger"
	"github.com/ThingiverseIO/thingiverseio"
	flite "github.com/gen2brain/flite-go"
)

const descriptor string = `
function Speak(Text string)()
property InSpeech: Is bool
`

var (
	log    = logger.New("Rommi Voice Server")
	output thingiverseio.Output
)

type speakRequest struct {
	Text string
}

type inSpeech struct {
	Is bool
}

func main() {
	log.Init("Starting Up")
	output, err := thingiverseio.NewOutput(descriptor)
	if err != nil {
		log.Error("error starting output: ", err)
		os.Exit(1)
	}
	// The valid names are "awb", "kal16", "kal", "rms" and "slt"
	voice, err := flite.VoiceSelect("slt")
	if err != nil {
		log.Error("error loading voice: ", err)
		os.Exit(1)
	}
	var speakreq speakRequest
	requests, _ := output.Requests().AsChan()

	output.SetProperty("InSpeech", inSpeech{false})
	output.Run()

	for req := range requests {
		req.Decode(&speakreq)
		output.SetProperty("InSpeech", inSpeech{true})
		log.Info("Playing text: ", speakreq.Text)
		flite.TextToSpeech(speakreq.Text, voice, "play")
		output.Reply(req, nil)
		log.Success("Done")
		output.SetProperty("InSpeech", inSpeech{false})
	}
}
