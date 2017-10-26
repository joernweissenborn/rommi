package voice

import (
	"fmt"

	"github.com/ThingiverseIO/thingiverseio"
)

const descriptor string = `
function Speak(Text string)()
property InSpeech: Is bool
`

type SpeakRequest struct {
	Text string
}

type InSpeech struct {
	Is bool
}

type Voice struct {
	*thingiverseio.Input
	inSpeech *thingiverseio.PropertyObservable
}

func New() (v *Voice, err error) {
	i, err := thingiverseio.NewInput(descriptor)
	if err != nil {
		return
	}
	v = &Voice{
		Input: i}
	return
}

func (v *Voice) Speak(text string) {
	v.Trigger("Speak", SpeakRequest{text})
}

func (v *Voice) Speakf(format string, values ...interface{}) {
	v.Trigger("Speak", SpeakRequest{fmt.Sprintf(format, values...)})
}

func (v *Voice) ObserveInSpeech() {
	v.StartObservation("InSpeech")
}

func (v *Voice) InSpeechObservable() *thingiverseio.PropertyObservable {
	if v.inSpeech == nil {
		v.inSpeech, _, _ = v.GetPropertyObservable("InSpeech")
	}
	return v.inSpeech
}
