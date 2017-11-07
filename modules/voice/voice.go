package voice

import (
	"fmt"

	"github.com/ThingiverseIO/thingiverseio"
)

const Descriptor string = `
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
	i, err := thingiverseio.NewInput(Descriptor)
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
	v.Speak(fmt.Sprintf(format, values...))
}

func (v *Voice) SpeakAndWait(text string) {
	r, _ := v.Call("Speak", SpeakRequest{text})
	<-r.AsChan()
}

func (v *Voice) SpeakfAndWait(format string, values ...interface{}) {
	v.SpeakAndWait(fmt.Sprintf(format, values...))
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
