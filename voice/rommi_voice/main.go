package main

import (
	"os"
	"rommi/voice"

	"github.com/ThingiverseIO/logger"
	"github.com/ThingiverseIO/thingiverseio"
	flite "github.com/gen2brain/flite-go"
)

func main() {
	log := logger.New("Rommi Voice Server")
	log.Init("Starting Up")
	output, err := thingiverseio.NewOutput(voice.Descriptor)
	if err != nil {
		log.Error("error starting output: ", err)
		os.Exit(1)
	}
	// The valid names are "awb", "kal16", "kal", "rms" and "slt"
	v, err := flite.VoiceSelect("slt")
	if err != nil {
		log.Error("error loading voice: ", err)
		os.Exit(1)
	}
	var speakreq voice.SpeakRequest
	requests, _ := output.Requests().AsChan()

	output.SetProperty("InSpeech", voice.InSpeech{false})
	output.Run()

	for req := range requests {
		req.Decode(&speakreq)
		output.SetProperty("InSpeech", voice.InSpeech{true})
		log.Info("Playing text: ", speakreq.Text)
		flite.TextToSpeech(speakreq.Text, v, "play")
		output.Reply(req, nil)
		log.Success("Done")
		output.SetProperty("InSpeech", voice.InSpeech{false})
	}
}
