package main

import (
	"rommi/lib/extension"
	"rommi/lib/language/sentence"
)

var ext = extension.Extension{
	Name:       "Chuck Norris Facts",
	Descriptor: descriptor,
	Actions: []extension.Action{
		tellChuckQuote,
	},
}

var tellChuckQuote = extension.Action{
	Function: "TellChuckQuote",
	Name:     "Tell Chuck Norris Qote",
	Sentences: []sentence.Sentence{
		sentence.New("erzaehle einen witz"),
	},
	Parameter: []byte{},
}
