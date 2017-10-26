package core

import (
	"rommi/brain/extension"
	"rommi/brain/language/sentence"
	"rommi/brain/language/wordlist"
	"rommi/brain/service"
	"rommi/ears"
	"rommi/voice"
	"time"

	"github.com/ThingiverseIO/logger"
	"github.com/joernweissenborn/eventual2go"
)

var (
	log                  = logger.New("Rommis Brain")
	likelynessThreshhold = 0.6
	triggerWord          = "computer"
)

type core struct {
	*eventual2go.Reactor
	services  map[string]service.Service
	ears      *ears.Ears
	listening bool
	voice     *voice.Voice
	wordlist  wordlist.WordList
}

func Start() (err error) {
	log.Init("Initializing")
	c := &core{
		Reactor:  eventual2go.NewReactor(),
		services: map[string]service.Service{},
	}

	log.Init("Initializing Voice")
	c.voice, err = voice.New()
	if err != nil {
		log.Error(err)
		return
	}
	c.voice.ConnectedObservable().OnChange(c.onVoiceConnChanged)
	conn := c.voice.ConnectedObservable().NextChange()
	c.voice.Run()
	conn.WaitUntilTimeout(5 * time.Second)

	log.Init("Initializing Ears")
	c.ears, err = ears.New()
	if err != nil {
		log.Error(err)
		return
	}
	c.React(earsConnectionEvent{}, c.onEarsConnChanged)
	c.AddObservable(earsConnectionEvent{}, c.ears.ConnectedObservable().Observable)
	c.React(utteranceEvent{}, c.onUtterance)
	c.AddStream(utteranceEvent{}, c.ears.Utterances().Stream)
	c.ears.StartConsume("Utterances")
	c.ears.Run()

	c.React(stopListenToCmdEvent{}, c.onStopListenToCmd)

	c.React(registerServiceEvent{}, c.onRegisterService)
	c.React(removeServiceEvent{}, c.onRemoveService)
	c.Fire(registerServiceEvent{}, coreService{c})

	log.Init("Initializing Extension Output")
	eo, err := extension.NewOutput()
	if err != nil {
		log.Error(err)
		return
	}
	c.AddStream(registerServiceEvent{}, eo.Extensions().Stream)
	eo.Run()

	log.Init("Done")
	return
}

func (c *core) updateWordList(eventual2go.Data) {
	log.Info("Rebuilding Wordlist")
	c.wordlist = wordlist.WordList{triggerWord: 1}
	for _, srv := range c.services {
		c.wordlist.Merge(service.WordList(srv))
	}
	log.Debug("New Wordlist is: ", c.wordlist)
	if c.ears.Connected() {
		log.Info("Sending Wordlist To Ear Server")
		c.ears.SetWordList(c.wordlist)
	}
	log.Success("Done")

}

func (c *core) onRegisterService(d eventual2go.Data) {
	srv := d.(service.Service)
	if _, ok := c.services[srv.GetName()]; !ok {
		log.Infof("Service Arrived: %s", srv.GetName())
		c.voice.Speakf("Service Arrived: %s", srv.GetName())

		if err := srv.Activate(); err != nil {
			log.Errorf("Error Activating Service: %s", err)
			c.voice.Speakf("Error Activating Service: %s", err)
			return
		}
		c.services[srv.GetName()] = srv
		c.AddFuture(removeServiceEvent{}, srv.Gone())
		log.Success("Service Activated")
		c.voice.Speak("Service Activated")
		c.updateWordList(nil)
	}
}

func (c *core) onRemoveService(d eventual2go.Data) {
	srv := d.(service.Service)
	if _, ok := c.services[srv.GetName()]; ok {
		log.Infof("Service Gone: %s", srv.GetName())
		c.voice.Speakf("Service Gone: %s", srv.GetName())

		delete(c.services, srv.GetName())
		log.Success("Service Removed")
		c.voice.Speak("Service Removed")
		c.updateWordList(nil)
	}
}

func (c *core) onUtterance(d eventual2go.Data) {
	s := d.(sentence.Sentence)
	log.Info("Got Utterance: ", s)

	if s.Contains(triggerWord) {
		s = s.Remove(triggerWord)
		ok := c.checkForAction(s)
		if !ok {
			c.startListenToCmd()
			c.FireIn(stopListenToCmdEvent{}, nil, 10*time.Second)
			c.voice.Speak("Awaiting Command")
		}
	} else if c.listening {
		c.listening = !c.checkForAction(s)
	}
}

func (c *core) checkForAction(s sentence.Sentence) (ok bool) {
	var score float64
	var srv service.Service
	var action service.Action
	for _, ssrv := range c.services {
		saction, sscore := service.MostLiklyAction(ssrv, s)
		if sscore > score {
			score = sscore
			srv = ssrv
			action = saction
		}
	}
	ok = score >= likelynessThreshhold
	if ok {
		log.Infof("Most likely action is '%s' of service '%s' with a score of '%f'", action.GetName(), srv.GetName(), score)
		log.Infof("Executing Action")
		srv.Execute(action)
	} else {
		log.Info("No Likly Action Reacognized")
	}
	return
}

func (c *core) onEarsConnChanged(d eventual2go.Data) {
	connected := d.(bool)
	if connected {
		log.Success("Ears Connected")
		c.voice.Speak("Ears service connected. I am able to hear you!")
		log.Info("Sending Wordlist To Ear Server")
		c.ears.SetWordList(c.wordlist)
	} else {
		log.Error("Ears Disconnected")
		c.voice.Speak("Ears service disconnected. I am not able to hear you anymore!")
	}
}

func (c *core) onVoiceConnChanged(connected bool) {
	if connected {
		log.Success("Voice Connected")
		c.voice.Speak("Voice service connected. I am able to speak to you!")
	} else {
		log.Error("Voice Disconnected")
	}
}

func (c *core) startListenToCmd() {
	log.Info("Start listening for commands")
	c.listening = true
}

func (c *core) onStopListenToCmd(eventual2go.Data) {
	log.Info("Stop listening for commands")
	c.listening = false
}

func (c *core) callForServices() {
	for _ = range time.Tick(10 * time.Millisecond) {
		// c.output.Emit("CallForServices", nil, nil)
	}
}
