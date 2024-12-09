package server

import (
	"encoding/json"
	"io"
	"log"
	"luis/wetterserver/data"
	"luis/wetterserver/database"
	"net/http"
)

type HttpServer struct {
	database database.Database
}

func CreateNewHttpServer(db database.Database) HttpServer {
	var server = HttpServer{
		database: db,
	}
	server.StartServer()
	return server
}

func (server *HttpServer) StartServer() {
	http.HandleFunc("/weather_data", server.weatherHandler)
	http.ListenAndServe(":8080", nil)
}

func (server *HttpServer) weatherHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		server.outboundHandler(w, r)
	case "POST":
		server.inboundHandler(w, r)
	default:
		{
			w.WriteHeader(http.StatusMethodNotAllowed)
			log.Println("unallowed method used on /weather_data:", r.Method)
			return
		}
	}
}

func (server *HttpServer) outboundHandler(w http.ResponseWriter, r *http.Request) {
	data, err := server.database.GetLastDataset()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("error loading data:", err)
		return
	}
	json, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("error handling json data:", err)
		return
	}
	w.Write(json)
}

func (server *HttpServer) inboundHandler(w http.ResponseWriter, r *http.Request) {
	raw_data, error := io.ReadAll(r.Body)
	if error != nil || len(raw_data) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("body is nil or invalid")
		return
	}
	if !json.Valid(raw_data) {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("body is not a valid json")
		return
	}
	var parsed data.WeatherData
	error = json.Unmarshal(raw_data, &parsed)
	if error != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("json is invalid")
		return
	}
	db_error := server.database.InsertWeatherData(&parsed)
	if db_error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("error storing data:", db_error)
		return
	}

}
