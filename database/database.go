package database

import "luis/wetterserver/data"

type Database interface {
	StartDatabase()
	InsertWeatherData(data *data.WeatherData) error
	GetLastDataset() (*data.WeatherData, error)
}
