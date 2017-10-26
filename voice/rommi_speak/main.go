package main

import (
	"fmt"
	"os"

	"github.com/ThingiverseIO/thingiverseio"
)

const (
	descriptor string = "function Speak(Text string)"
)

var ()

type speakRequest struct {
	Text string
}

func main() {
	input, err := thingiverseio.NewInput(descriptor)
	if err != nil {
		panic(err)
	}
	input.Run()
	fmt.Println("Speaking:", os.Args[1])
	r, _ := input.Call("Speak", speakRequest{os.Args[1]})
	<-r.AsChan()
}
