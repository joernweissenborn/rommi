package brain

import (
	"fmt"

	"github.com/ThingiverseIO/thingiverseio"
)

const Descriptor string = `
function TellCommand(Cmd string)
`

type TellCommandRequest struct {
	Cmd string
}

type Brain struct {
	*thingiverseio.Input
}

func New() (b *Brain, err error) {
	i, err := thingiverseio.NewInput(Descriptor)
	if err != nil {
		return
	}
	b = &Brain{
		Input: i}
	return
}

func (b *Brain) TellCommand(text string) {
	b.Trigger("TellCommand", TellCommandRequest{text})
}

func (b *Brain) TellCommandf(format string, values ...interface{}) {
	b.TellCommand(fmt.Sprintf(format, values...))
}

func (b *Brain) TellCommandAndWait(text string) {
	r, _ := b.Call("TellCommand", TellCommandRequest{text})
	<-r.AsChan()
}

func (b *Brain) SpeakfAndWait(format string, values ...interface{}) {
	b.TellCommandAndWait(fmt.Sprintf(format, values...))
}
