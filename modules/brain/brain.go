package brain

import (
	"fmt"

	"github.com/ThingiverseIO/thingiverseio"
)

const Descriptor string = `
function TellCommand(Cmd string)
function GetTriggerWord() (TriggerWord string)
function GetServices() (Services []string)
function GetServiceActions(Service string) (Actions []string)
function GetActionSentences(Service string, Action string) (Sentences []string)
`

type TellCommandRequest struct {
	Cmd string
}

type GetTriggerWordRequest struct{}

type GetTriggerWordReply struct {
	TriggerWord string
}

type GetServicesRequest struct{}

type GetServicesReply struct {
	Services []string
}

type GetServiceActionsRequest struct {
	Service string
}

type GetServiceActionsReply struct {
	Actions []string
}

type GetActionSentencesRequest struct {
	Service string
	Action  string
}

type GetActionSentencesReply struct {
	Sentences []string
}

type Brain struct {
	*thingiverseio.Input
}

func New() (b *Brain, err error) {
	i, err := thingiverseio.NewInput(Descriptor)
	if err != nil {
		return
	}
	b = &Brain{
		Input: i}
	return
}

func (b *Brain) TellCommand(text string) {
	b.Trigger("TellCommand", TellCommandRequest{text})
}

func (b *Brain) TellCommandf(format string, values ...interface{}) {
	b.TellCommand(fmt.Sprintf(format, values...))
}

func (b *Brain) TellCommandAndWait(text string) {
	r, _ := b.Call("TellCommand", TellCommandRequest{text})
	<-r.AsChan()
}

func (b *Brain) TellCommandfAndWait(format string, values ...interface{}) {
	b.TellCommandAndWait(fmt.Sprintf(format, values...))
}

func (b *Brain) GetTriggerWord() (triggerWord string) {
	r, _ := b.Call("GetTriggerWord", nil)
	var res GetTriggerWordReply
	(<-r.AsChan()).Decode(&res)
	triggerWord = res.TriggerWord
	return
}

func (b *Brain) GetServices() (services []string) {
	r, _ := b.Call("GetServices", nil)
	var res GetServicesReply
	(<-r.AsChan()).Decode(&res)
	services = res.Services
	return
}

func (b *Brain) GetServiceActions(service string) (actions []string) {
	r, _ := b.Call("GetServiceActions", GetServiceActionsRequest{Service: service})
	var res GetServiceActionsReply
	(<-r.AsChan()).Decode(&res)
	actions = res.Actions
	return
}

func (b *Brain) GetSentences() (sentences []string) {
	for _, service := range b.GetServices() {
		ss := b.GetServiceSentences(service)
		sentences = append(sentences, ss...)
	}
	return
}

func (b *Brain) GetServiceSentences(service string) (sentences []string) {
	actions := b.GetServiceActions(service)
	for _, action := range actions {
		ss := b.GetActionSentences(service, action)
		sentences = append(sentences, ss...)
	}
	return
}

func (b *Brain) GetActionSentences(service string, action string) (sentences []string) {
	r, _ := b.Call("GetActionSentences", GetActionSentencesRequest{Service: service, Action: action})
	var res GetActionSentencesReply
	(<-r.AsChan()).Decode(&res)
	sentences = res.Sentences
	return
}
