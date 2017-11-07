package ears

import (
	"rommi/lib/language/sentence"
	"rommi/lib/language/wordlist"

	"github.com/ThingiverseIO/thingiverseio"
	"github.com/joernweissenborn/eventual2go"
)

const Descriptor = `
function StartListen()()
function StopListen()()
function StartRecording()(OK bool)
function StopRecording()(OK bool)
function GetRecordedWav()(Wav bin)
function SetWordList(Words []string)
property Listenting: Is bool
stream Utterances: Sentence string
`

type StartListenRequest struct{}
type StopListenRequest struct{}

type StartRecordingRequest struct{}
type StartRecordingReply struct{ OK bool }

type StopRecordingRequest struct{}
type StopRecordingReply struct{ OK bool }

type GetRecordedWavRequest struct{}
type GetRecordedWavReply struct{ Wav []byte }

type SetWordListRequest struct{ Words []string }
type SetWordListReply struct{}

type Listening struct{ Is bool }
type Utterance struct{ Sentence string }

type Ears struct {
	*thingiverseio.Input
}

func New() (v *Ears, err error) {
	i, err := thingiverseio.NewInput(Descriptor)
	if err != nil {
		return
	}
	v = &Ears{i}
	return
}

func (e *Ears) SetWordList(wl wordlist.WordList) {
	var swl SetWordListRequest
	for word := range wl {
		swl.Words = append(swl.Words, word)
	}
	e.Trigger("SetWordList", swl)
}

func (e *Ears) SetWordListAndWait(wl wordlist.WordList) {
	var swl SetWordListRequest
	for word := range wl {
		swl.Words = append(swl.Words, word)
	}
	e.Call("SetWordList", swl)
}

func (e *Ears) StartListen() { e.Trigger("StartListen", nil) }
func (e *Ears) StopListen()  { e.Trigger("StopListen", nil) }

func (e *Ears) StartRecording() (ok bool) {
	r, _ := e.Call("StartRecording", nil)
	var rep StartRecordingReply
	(<-r.AsChan()).Decode(&rep)
	return rep.OK
}

func (e *Ears) StopRecording() (ok bool) {
	r, _ := e.Call("StopRecording", nil)
	var rep StopRecordingReply
	(<-r.AsChan()).Decode(&rep)
	return rep.OK
}

func (e *Ears) GetRecordedWav() (wav []byte) {
	r, _ := e.Call("GetRecordedWav", nil)
	var rep GetRecordedWavReply
	(<-r.AsChan()).Decode(&rep)
	return rep.Wav
}

func (e *Ears) Utterances() *sentence.SentenceStream {
	s, _ := e.GetStream("Utterances")
	toSentence := func(d eventual2go.Data) eventual2go.Data {
		var utt Utterance
		d.(thingiverseio.StreamEvent).Value(&utt)
		return sentence.New(utt.Sentence)
	}
	return &sentence.SentenceStream{s.Transform(toSentence)}
}
