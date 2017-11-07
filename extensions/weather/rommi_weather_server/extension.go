package main

import (
	"rommi/lib/extension"
	"rommi/lib/language/sentence"
)

var ext = extension.Extension{
	Name:       "Weather Service",
	Descriptor: descriptor,
	Actions: []extension.Action{
		tellCurWeather,
	},
}

var tellCurWeather = extension.Action{
	Function: "TellCurrentWeather",
	Name:     "Tell The Curren Weather",
	Sentences: []sentence.Sentence{
		sentence.New("wie ist wetter"),
		sentence.New("wie ist wetter im moment"),
		sentence.New("sage mir das wetter"),
	},
	Parameter: []byte{},
}
