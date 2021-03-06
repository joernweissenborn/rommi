
/*
 * generated by event_generator
 *
 * DO NOT EDIT
 */

package audio

import "github.com/joernweissenborn/eventual2go"



type AudioCompleter struct {
	*eventual2go.Completer
}

func NewAudioCompleter() *AudioCompleter {
	return &AudioCompleter{eventual2go.NewCompleter()}
}

func (c *AudioCompleter) Complete(d Audio) {
	c.Completer.Complete(d)
}

func (c *AudioCompleter) Future() *AudioFuture {
	return &AudioFuture{c.Completer.Future()}
}

type AudioFuture struct {
	*eventual2go.Future
}

func (f *AudioFuture) Result() Audio {
	return f.Future.Result().(Audio)
}

type AudioCompletionHandler func(Audio) Audio

func (ch AudioCompletionHandler) toCompletionHandler() eventual2go.CompletionHandler {
	return func(d eventual2go.Data) eventual2go.Data {
		return ch(d.(Audio))
	}
}

func (f *AudioFuture) Then(ch AudioCompletionHandler) *AudioFuture {
	return &AudioFuture{f.Future.Then(ch.toCompletionHandler())}
}

func (f *AudioFuture) AsChan() chan Audio {
	c := make(chan Audio, 1)
	cmpl := func(d chan Audio) AudioCompletionHandler {
		return func(e Audio) Audio {
			d <- e
			close(d)
			return e
		}
	}
	ecmpl := func(d chan Audio) eventual2go.ErrorHandler {
		return func(error) (eventual2go.Data, error) {
			close(d)
			return nil, nil
		}
	}
	f.Then(cmpl(c))
	f.Err(ecmpl(c))
	return c
}

type AudioStreamController struct {
	*eventual2go.StreamController
}

func NewAudioStreamController() *AudioStreamController {
	return &AudioStreamController{eventual2go.NewStreamController()}
}

func (sc *AudioStreamController) Add(d Audio) {
	sc.StreamController.Add(d)
}

func (sc *AudioStreamController) Join(s *AudioStream) {
	sc.StreamController.Join(s.Stream)
}

func (sc *AudioStreamController) JoinFuture(f *AudioFuture) {
	sc.StreamController.JoinFuture(f.Future)
}

func (sc *AudioStreamController) Stream() *AudioStream {
	return &AudioStream{sc.StreamController.Stream()}
}

type AudioStream struct {
	*eventual2go.Stream
}

type AudioSubscriber func(Audio)

func (l AudioSubscriber) toSubscriber() eventual2go.Subscriber {
	return func(d eventual2go.Data) { l(d.(Audio)) }
}

func (s *AudioStream) Listen(ss AudioSubscriber) *eventual2go.Completer {
	return s.Stream.Listen(ss.toSubscriber())
}

func (s *AudioStream) ListenNonBlocking(ss AudioSubscriber) *eventual2go.Completer {
	return s.Stream.ListenNonBlocking(ss.toSubscriber())
}

type AudioFilter func(Audio) bool

func (f AudioFilter) toFilter() eventual2go.Filter {
	return func(d eventual2go.Data) bool { return f(d.(Audio)) }
}

func toAudioFilterArray(f ...AudioFilter) (filter []eventual2go.Filter){

	filter = make([]eventual2go.Filter, len(f))
	for i, el := range f {
		filter[i] = el.toFilter()
	}
	return
}

func (s *AudioStream) Where(f ...AudioFilter) *AudioStream {
	return &AudioStream{s.Stream.Where(toAudioFilterArray(f...)...)}
}

func (s *AudioStream) WhereNot(f ...AudioFilter) *AudioStream {
	return &AudioStream{s.Stream.WhereNot(toAudioFilterArray(f...)...)}
}

func (s *AudioStream) TransformWhere(t eventual2go.Transformer, f ...AudioFilter) *eventual2go.Stream {
	return s.Stream.TransformWhere(t, toAudioFilterArray(f...)...)
}

func (s *AudioStream) Split(f AudioFilter) (*AudioStream, *AudioStream)  {
	return s.Where(f), s.WhereNot(f)
}

func (s *AudioStream) First() *AudioFuture {
	return &AudioFuture{s.Stream.First()}
}

func (s *AudioStream) FirstWhere(f... AudioFilter) *AudioFuture {
	return &AudioFuture{s.Stream.FirstWhere(toAudioFilterArray(f...)...)}
}

func (s *AudioStream) FirstWhereNot(f ...AudioFilter) *AudioFuture {
	return &AudioFuture{s.Stream.FirstWhereNot(toAudioFilterArray(f...)...)}
}

func (s *AudioStream) AsChan() (c chan Audio, stop *eventual2go.Completer) {
	c = make(chan Audio)
	stop = s.Listen(pipeToAudioChan(c))
	stop.Future().Then(closeAudioChan(c))
	return
}

func pipeToAudioChan(c chan Audio) AudioSubscriber {
	return func(d Audio) {
		c <- d
	}
}

func closeAudioChan(c chan Audio) eventual2go.CompletionHandler {
	return func(d eventual2go.Data) eventual2go.Data {
		close(c)
		return nil
	}
}

type AudioCollector struct {
	*eventual2go.Collector
}

func NewAudioCollector() *AudioCollector {
	return &AudioCollector{eventual2go.NewCollector()}
}

func (c *AudioCollector) Add(d Audio) {
	c.Collector.Add(d)
}

func (c *AudioCollector) AddFuture(f *AudioFuture) {
	c.Collector.Add(f.Future)
}

func (c *AudioCollector) AddStream(s *AudioStream) {
	c.Collector.AddStream(s.Stream)
}

func (c *AudioCollector) Get() Audio {
	return c.Collector.Get().(Audio)
}

func (c *AudioCollector) Preview() Audio {
	return c.Collector.Preview().(Audio)
}

type AudioObservable struct {
	*eventual2go.Observable
}

func NewAudioObservable (value Audio) (o *AudioObservable) {
	return &AudioObservable{eventual2go.NewObservable(value)}
}

func (o *AudioObservable) Value() Audio {
	return o.Observable.Value().(Audio)
}

func (o *AudioObservable) Change(value Audio) {
	o.Observable.Change(value)
}

func (o *AudioObservable) OnChange(s AudioSubscriber) (cancel *eventual2go.Completer) {
	return o.Observable.OnChange(s.toSubscriber())
}

func (o *AudioObservable) Stream() (*AudioStream) {
	return &AudioStream{o.Observable.Stream()}
}


func (o *AudioObservable) AsChan() (c chan Audio, cancel *eventual2go.Completer) {
	return o.Stream().AsChan()
}

func (o *AudioObservable) NextChange() (f *AudioFuture) {
	return o.Stream().First()
}
