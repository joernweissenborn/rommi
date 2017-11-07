package pa

import (
	"errors"
	"rommi/ears/audio"
	"unsafe"

	"github.com/joernweissenborn/eventual2go"
	"github.com/xlab/portaudio-go/portaudio"
)

type PaAudio struct {
	data                              unsafe.Pointer
	sampleCount, sampleRate, channels int
}

func (a *PaAudio) Samples() []int16 { return (*(*[1 << 24]int16)(a.data))[:a.sampleCount] }
func (a *PaAudio) Rate() int        { return a.sampleRate }
func (a *PaAudio) Size() int        { return a.sampleCount }
func (a *PaAudio) Channels() int    { return a.channels }

type PortAudio struct {
	inStream               *audio.AudioStreamController
	paInStream             *portaudio.Stream
	sampleRate, sampleSize int
	channels               int
}

func New(sampleRate, sampleSize, channels int) (pa *PortAudio, err error) {
	code := portaudio.Initialize()
	err = paError(code)
	if err != nil {
		return
	}

	pa = &PortAudio{
		sampleRate: sampleRate,
		sampleSize: sampleSize,
		channels:   channels,
	}
	return
}

func (a *PortAudio) OpenRecordStream() (in *audio.AudioStream, err error) {
	a.inStream = audio.NewAudioStreamController()
	in = a.inStream.Stream()
	var stream *portaudio.Stream
	code := portaudio.OpenDefaultStream(&stream,
		int32(a.channels),
		0,
		portaudio.PaInt16,
		float64(a.sampleRate),
		uint(a.sampleSize),
		a.inCallback,
		nil)
	err = paError(code)
	a.paInStream = stream
	return
}

func (pa *PortAudio) inCallback(input unsafe.Pointer, _ unsafe.Pointer, sampleCount uint,
	_ *portaudio.StreamCallbackTimeInfo, _ portaudio.StreamCallbackFlags, _ unsafe.Pointer) int32 {
	pa.inStream.Add(&PaAudio{
		data:        input,
		sampleCount: int(sampleCount),
		sampleRate:  pa.sampleRate,
		channels:    pa.channels,
	})
	return int32(portaudio.PaContinue)
}

func (pa *PortAudio) StartRecording() (err error) {
	code := portaudio.StartStream(pa.paInStream)
	err = paError(code)
	return
}

func (pa *PortAudio) StopRecording() (err error) {
	code := portaudio.StopStream(pa.paInStream)
	err = paError(code)
	return
}
func (pa *PortAudio) Shutdown(eventual2go.Data) (err error) {
	portaudio.StopStream(pa.paInStream)
	portaudio.CloseStream(pa.paInStream)
	portaudio.Terminate()
	return

}

func paError(code portaudio.Error) (err error) {
	if portaudio.ErrorCode(code) != portaudio.PaNoError {
		err = errors.New(portaudio.GetErrorText(code))
	}
	return
}
