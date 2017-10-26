package sentence

import (
	"strings"

	"github.com/renstrom/fuzzysearch/fuzzy"
)

//go:generate evt2gogen -t Sentence

type Sentence []string

func New(s string) Sentence {
	return Sentence(strings.Split(s, " "))
}

func (s Sentence) Len() (l int) {
	for _, w := range s {
		l += len(w)
	}
	return
}

func (s Sentence) Contains(word string) bool {
	for _, w := range s {
		if w == word {
			return true
		}
	}
	return false
}

func (s Sentence) Remove(word string) (ns Sentence){
	for i, w := range s {
		if w == word {
			if i == 0 {
				ns = s[1:]
			} else if i == s.Len()-1 {
				ns = s[:i-1]
			} else {
				ns = append(s[:i-1], s[i+1:]...)
			}
		}
	}
	return
}

func (s Sentence) Similarity(to Sentence) (score float64) {
	for _, word := range s {
		res := fuzzy.RankFind(word, []string(to))
		// if len(to) <= i {
		//         break
		// }
		// toword := []rune(to[i])
		//
		// var wordscore float64
		// for j, r := range word {
		//         if len(toword) <= j {
		//                 break
		//         }
		//         if r == toword[j] {
		//                 wordscore += 1
		//         }
		// }
		//
		// wordlen := len(word)
		// if len(word) < len(toword) {
		//         wordlen = len(toword)
		// }
		if len(res) > 0 {
			score += float64(res[len(res)-1].Distance)
		} else {
			score += float64(len(word))
		}

	}
	slen := s.Len()
	diff := slen - to.Len()
	if slen < to.Len() {
		slen = to.Len()
		diff *= -1
	}
	score += float64(diff)
	score = float64(slen) - score
	score /= float64(slen)
	return
}
