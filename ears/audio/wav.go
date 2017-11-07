package audio

import (
	"bytes"
	"unsafe"

	"github.com/joernweissenborn/eventual2go"
	wave "github.com/zenwerk/go-wave"
)

type addAudioToWavEvent struct{}

type waveBuf struct{ *bytes.Buffer }

func (*waveBuf) Close() error { return nil }

type Wav struct {
	err error
	buf *waveBuf
	w   *wave.Writer
	r   *eventual2go.Reactor
}

func NewWav(audio *AudioStream) (w *Wav) {
	w = &Wav{
		buf: &waveBuf{&bytes.Buffer{}},
		r:   eventual2go.NewReactor(),
	}
	w.r.React(addAudioToWavEvent{}, w.addAudio)
	w.r.AddStream(addAudioToWavEvent{}, audio.Stream)
	return
}

func (w *Wav) Close() (err error) {
	w.r.Shutdown(nil)
	w.r.ShutdownFuture().WaitUntilComplete()
	if err == nil {
		err = w.w.Close()
	}
	return w.err
}

func (w *Wav) Data() (data []byte) {
	return w.buf.Bytes()
}

func (w *Wav) addAudio(d eventual2go.Data) {
	if w.err != nil {
		return
	}
	audio := d.(Audio)
	if w.w == nil {
		param := wave.WriterParam{
			Out:           w.buf,
			Channel:       audio.Channels(),
			SampleRate:    audio.Rate(),
			BitsPerSample: 16,
		}

		w.w, w.err = wave.NewWriter(param)
		if w.err != nil {
			return
		}
	}
	_, w.err = w.w.WriteSample16(audio.Samples())
	if w.err != nil {
		return
	}
}

func cast(data []int16) []uint16 {
	out := *(*[]uint16)(unsafe.Pointer(&data))
	for i := range data {
		out[i] += 0x7FFF
	}
	return out
}
