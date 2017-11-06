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

func FromStringSlice(words []string) (wl WordList) {
	wl = WordList{}
	wl.AddString(words...)
	return
}

func (wl WordList) AddString(s ...string) {
	for _, word := range s {
		word = strings.ToLower(word)
		wl[word]++
	}
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

func (wl WordList) ContainsAll(t WordList) (contains bool) {
	contains = len(t) == 0
	for word := range t {
		_, contains = wl[word]
		if !contains {
			return
		}
	}
	return
}
