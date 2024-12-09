package data

import(
	"time"
)

type WeatherData struct {
	Temperature float64
	Pressure    float64
	RelPressure float64
	Voltage     float64
	Humidity    float64
	ReceivedAt  time.Time
}

func NewWeatherData(Temperature float64,
	Pressure float64,
	RelPressure float64,
	Voltage float64,
	Humidity float64,
	ReceivedAt time.Time,
	) WeatherData {
	return WeatherData{
		Temperature: Temperature,
		Pressure:    Pressure,
		RelPressure: RelPressure,
		Voltage:     Voltage,
		Humidity:    Humidity,
		ReceivedAt:  ReceivedAt,
	}
}
