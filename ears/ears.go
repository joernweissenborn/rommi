package ears

import (
	"rommi/brain/language/sentence"
	"rommi/brain/language/wordlist"

	"github.com/ThingiverseIO/thingiverseio"
	"github.com/joernweissenborn/eventual2go"
)

const descriptor = `
function SetWordList(Words []string)
stream Utterances: Sentence string
`

type SetWordList struct {
	Words []string
}

type Utterances struct {
	Sentence string
}

type Ears struct {
	*thingiverseio.Input
}

func New() (v *Ears, err error) {
	i, err := thingiverseio.NewInput(descriptor)
	if err != nil {
		return
	}
	v = &Ears{i}
	return
}

func (e *Ears) SetWordList(wl wordlist.WordList) {
	var swl SetWordList
	for word := range wl {
		swl.Words = append(swl.Words, word)
	}
	e.Call("SetWordList", swl)
}

func (e *Ears) Utterances() *sentence.SentenceStream {
	s, _ := e.GetStream("Utterances")
	toSentence := func(d eventual2go.Data) eventual2go.Data {
		var utt Utterances
		d.(thingiverseio.StreamEvent).Value(&utt)
		return sentence.New(utt.Sentence)
	}
	return &sentence.SentenceStream{s.Transform(toSentence)}
}
