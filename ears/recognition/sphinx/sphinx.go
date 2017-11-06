package sphinx

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"rommi/ears/audio"
	"strings"

	"github.com/joernweissenborn/eventual2go/typedevents"
	"github.com/xlab/pocketsphinx-go/sphinx"
)

type Recognizer struct {
	*sphinx.Decoder
	cfg         *Config
	utterances  *typedevents.StringStreamController
	inUtterance bool
}

type Config struct {
	Hmm, Dict, Lm, G2PModel string
	WorkDir                 string
	SampleRate              float32
}

func (cfg Config) sphinxCfg() *sphinx.Config {
	return sphinx.NewConfig(
		sphinx.HMMDirOption(cfg.Hmm),
		sphinx.DictFileOption(cfg.Dict),
		sphinx.LMFileOption(cfg.Lm),
		sphinx.SampleRateOption(cfg.SampleRate),
	)
}

func New(cfg *Config) (sr *Recognizer, err error) {
	_, err = os.Stat(cfg.Dict)
	if err != nil {
		if err = ioutil.WriteFile(cfg.Dict, []byte{}, os.ModePerm); err != nil {
			return
		}
	}
	dec, err := sphinx.NewDecoder(cfg.sphinxCfg())
	if err != nil {
		return
	}
	sr = &Recognizer{
		Decoder:    dec,
		cfg:        cfg,
		utterances: typedevents.NewStringStreamController(),
	}
	return
}

func (sr *Recognizer) SetWordList(words []string) (err error) {

	// only take action if we have to, regenerating the dict takes 5min on rpi
	curWL, err := sr.readWordList()
	if err != nil {
		return
	}
	newWL := wordlist.FromStringSlice(words)

	if curWL.ContainsAll(newWL) && len(curWL) == len(newWL) {
		return
	}

	sr.inUtterance = false
	sr.EndUtt()
	defer sr.StartUtt()

	// Word Create file
	path := sr.wordListPath()
	ioutil.WriteFile(path, []byte(strings.Join(words, "\n")), os.ModePerm)

	// Call g2p
	cmd := exec.Command("g2p-seq2seq",
		fmt.Sprintf("--decode=%s", path),
		fmt.Sprintf("--model=%s", sr.cfg.G2PModel),
		fmt.Sprintf("--output=%s", sr.cfg.Dict),
	)

	err := cmd.Run()
	if err != nil {
		return
	}

	// Reconfigure
	sr.Decoder.Reconfigure(sr.cfg.sphinxCfg())
}

func (sr *Recognizer) wordListPath() (path string) {
	path := filepath.Join(sr.cfg.WorkDir, "wordlist")
	return
}

func (sr *Recognizer) readWordList() (worldList wordlist.Wordlist, err error) {
	f, err := os.Open(sr.wordListPath())
	if err != nil {
		return
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	wordlist = wordlist{}
	for scanner.Scan() {
		wordlist.AddString(scanner.Text())
	}
	err = scanner.Err()
	return
}

func (sr *Recognizer) Recognize(in *audio.AudioStream) (utterances *typedevents.StringStream) {
	sr.StartUtt()
	in.Listen(sr.recognize)
	return sr.utterances.Stream()
}

func (sr *Recognizer) recognize(a audio.Audio) {
	sr.ProcessRaw(a.Samples(), false, false)
	if sr.IsInSpeech() {
		sr.inUtterance = true
	} else if sr.inUtterance {
		sr.inUtterance = false
		sr.EndUtt()
		hyp, _ := sr.Hypothesis()
		sentence := fmt.Sprint(hyp)
		if len(sentence) > 0 {
			sr.utterances.Add(sentence)
		}
		sr.StartUtt()
	}
}
