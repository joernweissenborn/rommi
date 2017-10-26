package wordlist

import (
	"rommi/brain/language/sentence"
	"strings"
)

type WordList map[string]int

func New(sentences ...sentence.Sentence) (wl WordList) {
	wl = WordList{}
	wl.Add(sentences...)
	return
}

func (wl WordList) Add(sentences ...sentence.Sentence) {
	for _, s := range sentences {
		for _, word := range s {
			word = strings.ToLower(word)
			wl[word]++
		}
	}
}

func (wl WordList) Merge(t WordList) {
	for w, s := range t {
		if _, ok := wl[w]; ok {
			wl[w] += s
		} else {
			wl[w] = s
		}
	}
}

func (wl WordList) Strings() (s []string) {
	for w := range wl {
		s = append(s, w)
	}
	return
}
