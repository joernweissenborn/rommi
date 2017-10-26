package extension

import (
	"bytes"
	"errors"
	"rommi/brain/service"
	"time"

	"github.com/ThingiverseIO/thingiverseio/config"
	tvio "github.com/ThingiverseIO/thingiverseio/core"
	"github.com/ThingiverseIO/thingiverseio/descriptor"
	"github.com/joernweissenborn/eventual2go"
	"github.com/ugorji/go/codec"
)

//go:generate evt2gogen -t Extension

type Extension struct {
	Name       string
	Descriptor string
	Actions    []Action

	input tvio.InputCore
	gone  *eventual2go.Completer
}

func (e *Extension) Activate() (err error) {
	d, err := descriptor.Parse(e.Descriptor)
	if err != nil {
		return
	}
	cfg := config.Configure()
	cfg.Debug = true
	tracker, provider := tvio.DefaultBackends()
	e.input, err = tvio.NewInputCore(d, cfg, tracker, provider...)
	if err != nil {
		return
	}
	f := e.input.ConnectedObservable().NextChange()
	e.input.Run()
	if !f.WaitUntilTimeout(5 * time.Second) {
		e.input.Shutdown()
		err = errors.New("Could not connect to service")
	} else {
		e.gone = eventual2go.NewCompleter()
		e.input.ConnectedObservable().NextChange().Then(e.onDisconnect)
	}
	return
}
func (e *Extension) Execute(a service.Action)  { a.(Action).execute(e.input) }
func (e *Extension) Gone() *eventual2go.Future { return e.gone.Future() }
func (e *Extension) GetName() string           { return e.Name }
func (e *Extension) GetActions() (actions []service.Action) {
	for _, action := range e.Actions {
		actions = append(actions, action)
	}
	return
}

func (e *Extension) onDisconnect(bool)bool {
	e.gone.Complete(e)
	e.input.Shutdown()
	return false
}

func Encode(e Extension) (data []byte, err error) {
	var h codec.MsgpackHandle
	var buf bytes.Buffer
	enc := codec.NewEncoder(&buf, &h)
	err = enc.Encode(e)
	data = buf.Bytes()
	return
}

func Decode(data []byte) (e *Extension, err error) {
	e = &Extension{}
	var h codec.MsgpackHandle
	buf := bytes.NewBuffer(data)
	dec := codec.NewDecoder(buf, &h)
	err = dec.Decode(e)
	return
}
