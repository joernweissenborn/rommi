package core

import (
	"path/filepath"
	"rommi/ears/audio"
	"rommi/ears/audio/pa"
	"rommi/ears/recognition"
	"rommi/ears/recognition/sphinx"
	"rommi/voice"

	"github.com/ThingiverseIO/logger"
	"github.com/ThingiverseIO/thingiverseio"
	"github.com/ThingiverseIO/thingiverseio/message"
	"github.com/joernweissenborn/eventual2go"
	"github.com/joernweissenborn/eventual2go/typedevents"
)

const descriptor = `
function SetWordList(Words []string)
stream Utterances: Sentence string
`

type setWordList struct {
	Words []string
}

type utterance struct {
	Sentence string
}

type setWordListEvent struct{}

var (
	log        = logger.New("Rommis Ears")
	channels   = 1
	sampleRate = 16000
	sampleSize = 512
)

type core struct {
	*eventual2go.Reactor
	output     *thingiverseio.Output
	recognizer recognition.Recognizer
	recorder   audio.Recorder
	recording  bool
	shutdown   *eventual2go.Shutdown
}

func Start(path string) (err error) {
	log.Init("Starting Up")

	c := &core{
		Reactor:  eventual2go.NewReactor(),
		shutdown: eventual2go.NewShutdown(),
	}

	log.Init("Initializing Audio")

	c.recorder, err = pa.New(sampleRate, sampleSize, channels)
	if err != nil {
		log.Error("Error Opening Audio: ", err)
		return
	}

	c.shutdown.Register(c.recorder.(eventual2go.Shutdowner))

	audio, err := c.recorder.OpenRecordStream()
	if err != nil {
		log.Error("Error Opening Record Stream: ", err)
		return
	}

	log.Success("Done")

	log.Info("Initializing Recognizer")

	hmm := filepath.Join(path, "hmm")
	dict := filepath.Join(path, "dict")
	lm := filepath.Join(path, "lm.bin")
	g2p := filepath.Join(path, "g2pmodel")

	log.Initf("Path is '%s'", path)
	log.Initf("HMM Path is '%s'", hmm)
	log.Initf("Dict Path is '%s'", dict)
	log.Initf("LM Path is '%s'", lm)
	log.Initf("G2P Model Path is '%s'", g2p)

	cfg := &sphinx.Config{
		Hmm:        hmm,
		Dict:       dict,
		Lm:         lm,
		G2PModel:   g2p,
		WorkDir:    path,
		SampleRate: float32(sampleRate),
	}
	c.recognizer, err = sphinx.New(cfg)
	if err != nil {
		log.Error("Error Opening Recognizer: ", err)
		return
	}

	log.Success("Done")

	log.Info("Initializing Voice")
	v, err := voice.New()
	if err != nil {
		log.Error("Error Starting Voice: ", err)
		return
	}
	v.ObserveInSpeech()
	c.AddObservable(voice.InSpeech{}, v.InSpeechObservable().Observable)
	c.React(voice.InSpeech{}, c.onSpechChange)
	v.Run()
	log.Success("Done")

	log.Info("Initializing Output")
	c.output, err = thingiverseio.NewOutput(descriptor)
	if err != nil {
		log.Error("Error Starting Output: ", err)
		return
	}
	isSetWL := func(r *message.Request) bool { return r.Function == "SetWordList" }
	c.AddStream(setWordListEvent{}, c.output.Requests().Where(isSetWL).Stream)
	c.React(setWordListEvent{}, c.onSetWordList)
	c.recognizer.Recognize(audio).Listen(streamUtterance(c.output))

	log.Info("Start Recording")
	err = c.startRecording()
	if err != nil {
		log.Error("Error Start Recording: ", err)
		return
	}

	log.Success("Done")

	return
}

func (c *core) startRecording() (err error) {
	log.Info("Start Recording")
	err = c.recorder.StartRecording()
	if err != nil {
		log.Error("Error Start Recording: ", err)
		return
	}
	c.recording = true
	return
}

func (c *core) stopRecording() (err error) {
	log.Info("Stop Recording")
	err = c.recorder.StopRecording()
	if err != nil {
		log.Error("Error Stop Recording: ", err)
		return
	}
	c.recording = false
	return
}

func (c *core) onSetWordList(d eventual2go.Data) {
	req := d.(*message.Request)
	if c.recording {
		c.stopRecording()
		defer c.startRecording()
	}
	var wordlist setWordList
	req.Decode(&wordlist)
	log.Infof("Got New World List: %s", wordlist.Words)
	if err := c.recognizer.SetWordList(wordlist.Words); err != nil {
		log.Error("Error Setting Word List: ", err)
	} else {
		log.Success("Done")
	}
	c.output.Reply(req, nil)
}

func (c *core) onSpechChange(d eventual2go.Data) {
	var s voice.InSpeech
	d.(thingiverseio.Property).Value(&s)
	if s.Is {
		log.Info("In Speech Detected, Stopped Recording")
		c.stopRecording()
	} else {
		log.Info("In Speech Ended, Started Recording")
		c.startRecording()
	}
}

func streamUtterance(output *thingiverseio.Output) typedevents.StringSubscriber {
	return func(utt string) {
		log.Infof("Got Utterance: %s", utt)
		output.AddStream("Utterances", utterance{utt})
	}
}
