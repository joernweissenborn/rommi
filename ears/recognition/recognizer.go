package recognition

import (
	"rommi/ears/audio"

	"github.com/joernweissenborn/eventual2go/typedevents"
)

type Recognizer interface {
	Recognize(in *audio.AudioStream) (utterances *typedevents.StringStream)
	SetWordList(words []string)
}
