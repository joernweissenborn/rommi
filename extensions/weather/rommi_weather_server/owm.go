package main

import (
	"time"
)

func forecastFromOpenWeatherMap(m encodedWeather) (fc forecast) {
	var ws []weather

	for _, wi := range m["list"].([]interface{}) {
		w := wi.(map[string]interface{})
		t := time.Unix(int64(w["dt"].(float64)), 0)

		var winddir, rain, snow float64

		if r, ok := w["wind"]; ok {
			if h, ok := r.(map[string]interface{})["dir"]; ok {
				winddir = h.(float64)
			}
		}

		if r, ok := w["rain"]; ok {
			if h, ok := r.(map[string]interface{})["3h"]; ok {
				rain = h.(float64)
			}
		}

		if r, ok := w["snow"]; ok {
			if h, ok := r.(map[string]interface{})["3h"]; ok {
				snow = h.(float64)
			}
		}

		ws = append(ws, weather{
			Time:        t,
			Temp:        w["main"].(map[string]interface{})["temp"].(float64),
			TempMin:     w["main"].(map[string]interface{})["temp_min"].(float64),
			TempMax:     w["main"].(map[string]interface{})["temp_max"].(float64),
			Pressure:    w["main"].(map[string]interface{})["pressure"].(float64),
			Humidity:    w["main"].(map[string]interface{})["humidity"].(float64),
			Weather:     w["weather"].([]interface{})[0].(map[string]interface{})["main"].(string),
			WeatherDesc: w["weather"].([]interface{})[0].(map[string]interface{})["description"].(string),
			Cloudiness:  w["clouds"].(map[string]interface{})["all"].(float64),
			WindSpeed:   w["wind"].(map[string]interface{})["speed"].(float64),
			WindDir:     winddir,
			Rain3h:      rain,
			Snow3h:      snow,
		})
	}

	fc = forecast{
		City:     m["city"].(map[string]interface{})["name"].(string),
		Country:  m["city"].(map[string]interface{})["country"].(string),
		Forecast: ws,
	}
	return
}
