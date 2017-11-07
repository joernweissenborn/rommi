package core

import (
	"path/filepath"
	"rommi/lib/audio"
	"rommi/lib/audio/pa"
	"rommi/lib/recognition"
	"rommi/lib/recognition/sphinx"
	"rommi/modules/ears"
	"rommi/modules/voice"

	"github.com/ThingiverseIO/logger"
	"github.com/ThingiverseIO/thingiverseio"
	"github.com/ThingiverseIO/thingiverseio/message"
	"github.com/joernweissenborn/eventual2go"
	"github.com/joernweissenborn/eventual2go/typedevents"
)

var (
	log        = logger.New("Rommis Ears")
	channels   = 1
	sampleRate = 16000
	sampleSize = 512
)

type core struct {
	*eventual2go.Reactor
	listening     *typedevents.BoolObservable
	recognizer    recognition.Recognizer
	recorder      audio.Recorder
	recordedAudio *audio.AudioStream
	recordedWav   *audio.Wav
	recording     bool
	output        *thingiverseio.Output
	shutdown      *eventual2go.Shutdown
}

func Start(path string) (err error) {
	log.Init("Starting Up")

	c := &core{
		Reactor:   eventual2go.NewReactor(),
		listening: typedevents.NewBoolObservable(false),
		shutdown:  eventual2go.NewShutdown(),
	}

	log.Init("Initializing Audio")

	c.recorder, err = pa.New(sampleRate, sampleSize, channels)
	if err != nil {
		log.Error("Error Opening Audio: ", err)
		return
	}

	c.shutdown.Register(c.recorder.(eventual2go.Shutdowner))

	c.recordedAudio, err = c.recorder.OpenRecordStream()
	if err != nil {
		log.Error("Error Opening Record Stream: ", err)
		return
	}

	log.Success("Done")

	log.Init("Initializing Recognizer")

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

	log.Init("Initializing Voice")
	v, err := voice.New()
	if err != nil {
		log.Error("Error Starting Voice: ", err)
		return
	}
	v.ObserveInSpeech()
	c.AddObservable(voice.InSpeech{}, v.InSpeechObservable().Observable)
	c.React(voice.InSpeech{}, c.onSpeechChange)
	v.Run()
	log.Success("Done")

	log.Init("Initializing Output")
	c.output, err = thingiverseio.NewOutput(ears.Descriptor)
	if err != nil {
		log.Error("Error Starting Output: ", err)
		return
	}
	setListening := func(listen bool) { c.output.SetProperty("Listenting", ears.Listening{Is: listen}) }
	c.listening.OnChange(setListening)
	c.AddStream(ears.SetWordListReply{}, c.output.RequestsWhereFunction("SetWordList").Stream)
	c.React(ears.SetWordListReply{}, c.onSetWordList)

	c.AddStream(ears.StartRecordingRequest{}, c.output.RequestsWhereFunction("StartRecording").Stream)
	c.React(ears.StartRecordingRequest{}, c.onStartRecording)

	c.AddStream(ears.StopRecordingRequest{}, c.output.RequestsWhereFunction("StopRecording").Stream)
	c.React(ears.StopRecordingRequest{}, c.onStopRecording)

	c.AddStream(ears.GetRecordedWavRequest{}, c.output.RequestsWhereFunction("GetRecordedWav").Stream)
	c.React(ears.GetRecordedWavRequest{}, c.onGetRecordedWav)

	listening := func(audio.Audio) bool { return c.listening.Value() }
	c.recognizer.Recognize(c.recordedAudio.Where(listening)).Listen(streamUtterance(c.output))

	log.Init("Start Recording")
	err = c.startRecording()
	if err != nil {
		log.Error("Error Start Recording: ", err)
		return
	}
	log.Init("Start Listenting")
	c.startListening()

	log.Success("Done")

	return
}

func (c *core) startRecording() (err error) {
	if c.recording {
		return
	}
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
	if !c.recording {
		return
	}
	log.Info("Stop Recording")
	err = c.recorder.StopRecording()
	if err != nil {
		log.Error("Error Stop Recording: ", err)
		return
	}
	c.recording = false
	return
}

func (c *core) startListening() (err error) {
	log.Info("Start Listening")
	c.listening.Change(true)
	return
}

func (c *core) stopListening() (err error) {
	log.Info("Stop Listening")
	c.listening.Change(false)
	return
}

func (c *core) onSetWordList(d eventual2go.Data) {
	req := d.(*message.Request)
	if c.recording {
		c.stopRecording()
		defer c.startRecording()
	}
	var wordlist ears.SetWordListRequest
	req.Decode(&wordlist)
	log.Infof("Got New World List: %s", wordlist.Words)
	if err := c.recognizer.SetWordList(wordlist.Words); err != nil {
		log.Error("Error Setting Word List: ", err)
	} else {
		log.Success("Done")
	}
	c.output.Reply(req, nil)
}

func (c *core) onStartRecording(d eventual2go.Data) {
	req := d.(*message.Request)
	c.stopListening()
	c.stopRecording()
	c.recordedWav = audio.NewWav(c.recordedAudio)
	log.Info("Start Recording Wav")
	c.output.Reply(req, ears.StartRecordingReply{true})
	c.startRecording()
}

func (c *core) onStopRecording(d eventual2go.Data) {
	req := d.(*message.Request)
	c.stopRecording()
	err := c.recordedWav.Close()
	log.Info("Stop Recording Wav")
	if err != nil {
		log.Error("Error Recording Wav: ", err)
	}
	c.output.Reply(req, ears.StopRecordingReply{err == nil})
	c.startRecording()
	c.startListening()
}

func (c *core) onGetRecordedWav(d eventual2go.Data) {
	req := d.(*message.Request)
	log.Info("Sending Recording Wav, length is ", len(c.recordedWav.Data()))
	c.output.Reply(req, ears.GetRecordedWavReply{c.recordedWav.Data()})
}

func (c *core) onSpeechChange(d eventual2go.Data) {
	var s voice.InSpeech
	d.(thingiverseio.Property).Value(&s)
	if s.Is {
		log.Info("In Speech Detected, Stopped Listening")
		c.stopListening()
	} else {
		log.Info("In Speech Ended, Started Listening")
		c.startListening()
	}
}

func streamUtterance(output *thingiverseio.Output) typedevents.StringSubscriber {
	return func(utt string) {
		log.Infof("Got Utterance: %s", utt)
		output.AddStream("Utterances", ears.Utterance{utt})
	}
}
