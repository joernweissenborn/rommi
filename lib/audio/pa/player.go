package pa

import (
	"rommi/lib/audio"
	"sync"
	"unsafe"

	"github.com/xlab/portaudio-go/portaudio"
)

var (
	samples = []int16{}
	wg      = &sync.WaitGroup{}
)

type PaAudioPlayer struct {
}

func NewPlayer() (p *PaAudioPlayer, err error) {
	code := portaudio.Initialize()
	err = paError(code)
	if err != nil {
		return
	}
	p = &PaAudioPlayer{}
	return
}

func (p *PaAudioPlayer) Close() (err error) {
	code := portaudio.Terminate()
	err = paError(code)
	return
}

func (p *PaAudioPlayer) Play(a audio.Audio) (err error) {

	var stream *portaudio.Stream
	code := portaudio.OpenDefaultStream(
		&stream,
		0,
		int32(a.Channels()),
		portaudio.PaInt16,
		float64(a.Rate()),
		4096,
		callback,
		nil)
	err = paError(code)
	if err != nil {
		return
	}
	samples = a.Samples()
	wg.Add(1)
	code = portaudio.StartStream(stream)
	err = paError(code)
	if err != nil {
		return
	}
	wg.Wait()

	code = portaudio.StopStream(stream)
	err = paError(code)
	if err != nil {
		return
	}
	return
}

func callback(_ unsafe.Pointer, output unsafe.Pointer, sampleCount uint,
	_ *portaudio.StreamCallbackTimeInfo, _ portaudio.StreamCallbackFlags, _ unsafe.Pointer) int32 {

	const (
		statusContinue = int32(portaudio.PaContinue)
		statusComplete = int32(portaudio.PaComplete)
		statusAbort    = int32(portaudio.PaAbort)
	)
	if len(samples) == 0 {
		wg.Done()
		return statusComplete
	}
	if sampleCount > uint(len(samples)) {
		sampleCount = uint(len(samples))
	}
	out := (*(*[1 << 24]int16)(output))[:sampleCount]
	for i := 0; i < int(sampleCount); i++ {
		out[i] = samples[i]
	}
	samples = samples[sampleCount:]
	return statusContinue
}
