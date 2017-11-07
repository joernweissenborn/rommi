package train

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"rommi/brain/language/sentence"
	"rommi/brain/language/wordlist"

	"github.com/ThingiverseIO/thingiverseio/uuid"
)

type SentenceDB map[string]uuid.UUID

type Train struct {
	path string
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
	return
}

func (t *Train) GetWordList() (wl wordlist.WordList, err error) {
	wl = wordlist.WordList{}
	db, err := t.OpenSentenceDB()
	if err != nil {
		return
	}
	for s := range db {
		wl.Add(sentence.New(s))
	}
	return
}

func (t *Train) AddSentence(sentence ...string) (err error) {
	db, err := t.OpenSentenceDB()
	if err != nil {
		return
	}

	var changed bool
	for _, s := range sentence {
		if _, ok := db[s]; !ok {
			changed = true
			db[s] = uuid.New()
		}
	}

	if changed {
		err = t.writeSentenceDB(db)
	}
	return
}

func (t *Train) OpenSentenceDB() (sentences SentenceDB, err error) {
	sentences = SentenceDB{}
	data, err := ioutil.ReadFile(t.sentencesPath())
	if err != nil {
		if os.IsNotExist(err) {
			err = nil
		}
		return
	}
	err = json.Unmarshal(data, &sentences)
	return
}

func (t *Train) writeSentenceDB(sentences SentenceDB) (err error) {
	data, err := json.MarshalIndent(sentences, "", "\t")
	if err != nil {
		return
	}
	err = ioutil.WriteFile(t.sentencesPath(), data, os.ModePerm)
	return
}

func (t *Train) sentencesPath() string {
	return filepath.Join(t.path, "sentences.json")
}
