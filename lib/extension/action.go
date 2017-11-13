package extension

import (
	"rommi/lib/language/sentence"

	"github.com/ThingiverseIO/thingiverseio/core"
	"github.com/ThingiverseIO/thingiverseio/message"
)

type Action struct {
	Function  string
	Name      string
	Sentences []sentence.Sentence
	Parameter []byte
}

func (a Action) GetName() string                   { return a.Name }
func (a Action) GetSentences() []sentence.Sentence { return a.Sentences }

func (a Action) execute(i core.InputCore) (err error) {
	_, _, _, err = i.Request(a.Function, message.TRIGGER, a.Parameter)
	return
}
