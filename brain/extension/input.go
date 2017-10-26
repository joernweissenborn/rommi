package extension

import (
	"github.com/ThingiverseIO/thingiverseio"
)

type ExtensionInput struct {
	*thingiverseio.Input
	ext []byte
}

func New(e Extension) (ei ExtensionInput, err error) {
	ext, err := Encode(e)
	if err != nil {
		return
	}

	i, err := thingiverseio.NewInput(desc)
	if err != nil {
		return
	}
	ei = ExtensionInput{
		Input: i,
		ext:   ext,
	}
	ei.ConnectedObservable().OnChange(ei.register)
	ei.Run()
	return
}

func (ei ExtensionInput) register(connected bool) {
	if !connected {
		return
	}
	re := RegisterExtension{ei.ext}
	ei.TriggerAll("RegisterExtension", re)
}
