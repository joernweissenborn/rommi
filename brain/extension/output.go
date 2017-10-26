package extension

import (
	"github.com/ThingiverseIO/thingiverseio"
	"github.com/joernweissenborn/eventual2go"
)

const desc = `
function RegisterExtension(Extension bin)
`

type RegisterExtension struct {
	Extension []byte
}

type ExtensionOutput struct {
	*thingiverseio.Output
}

func NewOutput() (eo ExtensionOutput, err error) {
	o, err := thingiverseio.NewOutput(desc)
	if err != nil {
		return
	}
	eo = ExtensionOutput{o}
	return
}

func (eo ExtensionOutput) Extensions() (e *ExtensionStream) {
	return &ExtensionStream{eo.Requests().Where(isRegExt).TransformConditional(eo.toExtension)}
}

func isRegExt(r *thingiverseio.Request) bool { return r.Function == "RegisterExtension" }

func (eo ExtensionOutput) toExtension(in eventual2go.Data) (out eventual2go.Data, ok bool) {
	req := in.(*thingiverseio.Request)
	eo.Reply(req, nil)
	var re RegisterExtension
	err := req.Decode(&re)
	if err != nil {
		return
	}
	ext, err := Decode(re.Extension)
	if err != nil {
		return
	}
	out = ext
	ok = true
	return
}
