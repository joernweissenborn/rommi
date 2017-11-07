package main

import (
	"os"
	"os/signal"
	"rommi/lib/extension"
	"rommi/modules/voice"

	"github.com/ThingiverseIO/logger"
	"github.com/ThingiverseIO/thingiverseio"
)

const descriptor = `
function TellChuckQuote()()
`

var log = logger.New("Rommi Chuck Quotes")

func main() {
	log.Init("Starting up")

	output, err := thingiverseio.NewOutput(descriptor)
	if err != nil {
		log.Error("Error creating output: ", err)
		os.Exit(1)
		return
	}

	v, err := voice.New()
	if err != nil {
		log.Error("Error creating voice: ", err)
		os.Exit(1)
		return
	}
	v.Run()

	rep := func(r *thingiverseio.Request) {
		q := getRandomQuote()
		log.Info("Telling Quote: ", q)
		v.Speak(q)
		output.Reply(r, nil)
	}

	isTCQ := func(r *thingiverseio.Request) bool { return r.Function == "TellChuckQuote" }
	output.Requests().Where(isTCQ).Listen(rep)

	output.Run()

	extension.New(ext)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
