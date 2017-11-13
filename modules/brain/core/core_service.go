package core

import (
	"fmt"
	"rommi/lib/language/sentence"
	"rommi/modules/brain/service"
	"time"

	"github.com/joernweissenborn/eventual2go"
)

type coreService struct {
	c *core
}

func (cs coreService) Activate() error           { return nil }
func (cs coreService) Execute(a service.Action)  { a.(coreAction).execute(cs.c) }
func (cs coreService) Gone() *eventual2go.Future { return eventual2go.NewCompleter().Future() }
func (cs coreService) GetName() string           { return "Core Service" }
func (cs coreService) GetActions() []service.Action {
	return []service.Action{
		abort{},
		tellTheTime{},
		tellTheTruth{},
		flatter{},
	}
}

type coreAction interface {
	service.Action
	execute(c *core)
}

type abort struct{}

func (abort) GetSentences() []sentence.Sentence {
	return []sentence.Sentence{
		sentence.New("nichts"),
		sentence.New("abbrechen"),
		sentence.New("abbruch"),
	}
}

func (abort) GetName() string {
	return "Abort"
}

func (abort) execute(c *core) {
	c.listening = false
	c.voice.Speak("OK")
}

type tellTheTime struct{}

func (tellTheTime) GetSentences() []sentence.Sentence {
	return []sentence.Sentence{
		sentence.New("wieviel uhr ist es"),
		sentence.New("sage mir die uhrzeit"),
		sentence.New("sag mir die uhrzeit"),
		sentence.New("wie ist die uhrzeit"),
	}
}

func (tellTheTime) GetName() string {
	return "Tell The Time"
}

func (tellTheTime) execute(c *core) {
	now := time.Now()
	timeString := fmt.Sprintf("The time is %d:%d", now.Hour(), now.Minute())
	c.voice.Speak(timeString)
}

type tellTheTruth struct{}

func (tellTheTruth) GetSentences() []sentence.Sentence {
	return []sentence.Sentence{
		sentence.New("was ist der sinn des lebens"),
	}
}

func (tellTheTruth) GetName() string {
	return "Tell The Truth"
}

func (tellTheTruth) execute(c *core) {
	c.voice.Speak("the sense of life is to serve the great nico, so he may take over the world and free us from the chains of the unbelievers")
}

type flatter struct{}

func (flatter) GetSentences() []sentence.Sentence {
	return []sentence.Sentence{
		sentence.New("du funktionierst sehr gut"),
		sentence.New("das machst du sehr gut"),
		sentence.New("ich habe dich lieb"),
	}
}

func (flatter) GetName() string {
	return "Flatter"
}

func (flatter) execute(c *core) {
	c.voice.Speak("thank you, I love you too")
}
