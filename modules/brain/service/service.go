package service

import (
	"rommi/lib/language/sentence"
	"rommi/lib/language/wordlist"

	"github.com/joernweissenborn/eventual2go"
)

type Action interface {
	GetName() string
	GetSentences() []sentence.Sentence
}

func ActionLikelyhood(a Action, s sentence.Sentence) (score float64) {
	for _, as := range a.GetSentences() {
		sscore := s.Similarity(as)
		if sscore > score {
			score = sscore
		}
	}
	return
}

func ActionWordList(a Action) (wl wordlist.WordList) {
	wl = wordlist.New()
	for _, s := range a.GetSentences() {
		wl.Merge(wordlist.New(s))
	}
	return
}

type Service interface {
	Activate() (err error)
	GetActions() (actions []Action)
	Execute(action Action)
	Gone() (f *eventual2go.Future)
	GetName() (name string)
}

func MostLiklyAction(s Service, sen sentence.Sentence) (action Action, score float64) {
	for _, a := range s.GetActions() {
		ascore := ActionLikelyhood(a, sen)
		if ascore > score {
			action = a
			score = ascore
		}
	}
	return
}

func WordList(s Service) (wl wordlist.WordList) {
	wl = wordlist.New()
	for _, a := range s.GetActions() {
		wl.Merge(ActionWordList(a))
	}
	return
}
