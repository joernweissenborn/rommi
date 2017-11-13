package audio

import "time"

//go:generate evt2gogen -t Audio

type Audio interface {
	Channels() int
	Samples() []int16
	Size() int
	Rate() int
}

func Duration(a Audio) time.Duration {
	return time.Duration(float64(a.Size())/float64(a.Rate())) * time.Second
}
