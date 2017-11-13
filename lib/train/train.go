package train

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"rommi/lib/audio"
	"rommi/lib/language/sentence"
	"rommi/lib/language/wordlist"
	"strings"

	"github.com/ThingiverseIO/uuid"
)

type AudioDB map[string]map[string][]uuid.UUID

func (db AudioDB) HasAudio(speaker, condition string, id uuid.UUID) (has bool) {
	for _, uid := range db[speaker][condition] {
		has = uid == id
		if has {
			return
		}
	}
	return
}

type SentenceDB map[string]uuid.UUID

func (db SentenceDB) GetById(id uuid.UUID) (s string, ok bool) {
	for sen, sid := range db {
		ok = sid == id
		if ok {
			s = sen
			return
		}
	}
	return
}

type Train struct {
	path string
	adb  AudioDB
	sdb  SentenceDB
}

func Open(path string) (t *Train, err error) {
	fi, err := os.Stat(path)
	if err != nil {
		err = os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return
		}
	} else {
		if !fi.IsDir() {
			err = errors.New("Modelpath must by a directory")
			return
		}
	}
	t = &Train{
		path: path,
	}
	err = t.openAudioDB()
	if err != nil {
		return
	}
	err = t.openSentenceDB()
	return
}

func (t *Train) sentencesPath() string { return filepath.Join(t.path, "sentences.json") }
func (t *Train) audioPath() string     { return filepath.Join(t.path, "audiodb") }

// AudioDB Layout
//
// root
// |
// === Speaker1
// |   |
// |   == Condition1
// |   |
// |   == Condition2
// |
// === Speaker1

func (t *Train) openAudioDB() (err error) {
	t.adb = AudioDB{}
	walk := func(path string, info os.FileInfo, err error) error {
		rel, err := filepath.Rel(t.audioPath(), path)
		if err != nil {
			return err
		}
		split := strings.Split(rel, string(filepath.Separator))
		if len(split) == 3 {
			speaker, condition, wav := split[0], split[1], split[2]
			if !strings.HasSuffix(wav, ".wav") {
				return nil
			}
			id := uuid.UUID(strings.TrimSuffix(wav, ".wav"))
			if _, ok := t.adb[speaker]; !ok {
				t.adb[speaker] = map[string][]uuid.UUID{}
			}
			t.adb[speaker][condition] = append(t.adb[speaker][condition], id)
		}
		return nil
	}
	err = filepath.Walk(t.audioPath(), walk)
	if err != nil {
		return
	}
	return
}

func (t *Train) GetSpeakers() (speakers []string) {
	for speaker := range t.adb {
		speakers = append(speakers, speaker)
	}
	return
}

func (t *Train) GetConditions(speaker string) (conditions []string) {
	for condition := range t.adb[speaker] {
		conditions = append(conditions, condition)
	}
	return
}

func (t *Train) AvailableAudio(speaker, condition string) (id []uuid.UUID) {
	return t.adb[speaker][condition]
}

func (t *Train) HasAudio(speaker, condition string, id uuid.UUID) (has bool) {
	return t.adb.HasAudio(speaker, condition, id)
}

func (t *Train) SaveAudio(speaker, condition string, id uuid.UUID, wav []byte) (err error) {
	dir := filepath.Join(t.audioPath(), speaker, condition)
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return
	}
	path := filepath.Join(dir, fmt.Sprintf("%s.wav", id.FullString()))
	err = ioutil.WriteFile(path, wav, 0644)
	if err != nil {
		return
	}
	if !t.HasAudio(speaker, condition, id) {
		t.adb[speaker][condition] = append(t.adb[speaker][condition], id)
	}
	return
}

func (t *Train) GetAudio(speaker, condition string, id uuid.UUID) (a audio.Audio, err error) {
	dir := filepath.Join(t.audioPath(), speaker, condition)
	path := filepath.Join(dir, fmt.Sprintf("%s.wav", id.FullString()))
	a, err = audio.ReadWav(path)
	return
}

func (t *Train) GetWordList() (wl wordlist.WordList) {
	wl = wordlist.WordList{}
	for s := range t.sdb {
		wl.Add(sentence.New(s))
	}
	return
}

func (t *Train) GetAllIds() (ids []uuid.UUID) {
	for _, id := range t.sdb {
		if id != uuid.UUID("trigger") {
			ids = append(ids, id)
		}
	}
	return
}

func (t *Train) GetTriggerWord() (word string) {
	word, _ = t.sdb.GetById(uuid.UUID("trigger"))
	return
}

func (t *Train) SetTriggerWord(word string) (err error) {
	t.sdb[word] = uuid.UUID("trigger")
	err = t.writeSentenceDB()
	return
}

func (t *Train) GetSentence(id uuid.UUID) (sentence string, ok bool) {
	sentence, ok = t.sdb.GetById(id)
	return
}

func (t *Train) AddSentence(sentence ...string) (err error) {
	var changed bool
	for _, s := range sentence {
		if _, ok := t.sdb[s]; !ok {
			changed = true
			t.sdb[s] = uuid.New()
		}
	}

	if changed {
		err = t.writeSentenceDB()
	}
	return
}

func (t *Train) openSentenceDB() (err error) {
	t.sdb = SentenceDB{}
	data, err := ioutil.ReadFile(t.sentencesPath())
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	err = json.Unmarshal(data, &t.sdb)
	return
}

func (t *Train) writeSentenceDB() (err error) {
	data, err := json.MarshalIndent(t.sdb, "", "\t")
	if err != nil {
		return
	}
	err = ioutil.WriteFile(t.sentencesPath(), data, os.ModePerm)
	return
}
