package audio

import (
	"bytes"

	"github.com/joernweissenborn/eventual2go"
	wave "github.com/joernweissenborn/go-wave"
)

type WavAudio struct {
	channels, size, rate int
	samples              []int16
}

func ReadWav(path string) (a Audio, err error) {
	r, err := wave.NewReader(path)
	if err != nil {
		return
	}
	wa := &WavAudio{
		channels: int(r.FmtChunk.Data.Channel),
		size:     int(r.NumSamples),
		rate:     int(r.FmtChunk.Data.SamplesPerSec),
	}
	var sample []int
	wa.samples = make([]int16, wa.size)
	for i := 0; i < wa.size; i++ {
		sample, err = r.ReadSampleInt()
		if err != nil {
			return
		}
		wa.samples[i] = int16(sample[0])
	}
	a = wa
	return

}

func (w *WavAudio) Channels() int    { return w.channels }
func (w *WavAudio) Samples() []int16 { return w.samples }
func (w *WavAudio) Size() int        { return w.size }
func (w *WavAudio) Rate() int        { return w.rate }

type addAudioToWavEvent struct{}

type waveBuf struct{ *bytes.Buffer }

func (*waveBuf) Close() error { return nil }

type WavRecorder struct {
	err error
	buf *waveBuf
	w   *wave.Writer
	r   *eventual2go.Reactor
}

func NewWavRecorder(audio *AudioStream) (w *WavRecorder) {
	w = &WavRecorder{
		buf: &waveBuf{&bytes.Buffer{}},
		r:   eventual2go.NewReactor(),
	}
	w.r.React(addAudioToWavEvent{}, w.addAudio)
	w.r.AddStream(addAudioToWavEvent{}, audio.Stream)
	return
}

func (w *WavRecorder) Close() (err error) {
	w.r.Shutdown(nil)
	w.r.ShutdownFuture().WaitUntilComplete()
	if err == nil {
		err = w.w.Close()
	}
	return w.err
}

func (w *WavRecorder) Data() (data []byte) {
	return w.buf.Bytes()
}

func (w *WavRecorder) addAudio(d eventual2go.Data) {
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
	for _, s := range audio.Samples() {
		_, w.err = w.w.WriteSample16([]int16{s})
		if w.err != nil {
			return
		}
	}
}
