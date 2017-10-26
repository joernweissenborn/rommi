package main

import (
	"encoding/json"
	"io"
	"time"
)

type weather struct {
	Time        time.Time
	Temp        float64
	TempMin     float64
	TempMax     float64
	Pressure    float64
	Humidity    float64
	Weather     string
	WeatherDesc string
	Cloudiness  float64
	WindSpeed   float64
	WindDir     float64
	Rain3h      float64
	Snow3h      float64
}

type encodedWeather map[string]interface{}

func decodeJSON(r io.Reader) (m encodedWeather, err error) {
	dec := json.NewDecoder(r)
	err = dec.Decode(&m)
	return
}
