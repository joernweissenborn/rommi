package main

import (
	"encoding/json"
	"io"
)

type forecast struct {
	City     string
	Country  string
	Forecast []weather
}

func (f forecast) current() (w weather) { return f.Forecast[0] }

func (f forecast) toJSON(w io.Writer) {
	enc := json.NewEncoder(w)
	enc.Encode(f)
}
