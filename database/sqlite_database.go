package database

import (
	"database/sql"
	"log"
	"luis/wetterserver/data"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type SqliteDatabase struct {
	currentConnection *sql.DB
}

func CreateNewSqliteDatabaseConnection() SqliteDatabase {
	var db = SqliteDatabase{}
	db.StartDatabase()
	return db
}

func (db *SqliteDatabase) StartDatabase() {
	db.createDatabase()
}

func (db *SqliteDatabase) openDatabase() error {
	if db.currentConnection != nil {
		log.Println("reusing connection")
		return nil
	}
	conn, err := sql.Open("sqlite3", "/db/weather.db?parseTime=true")
	if err != nil {
		log.Println("error opening or creating db:", err)
		return err
	}
	db.currentConnection = conn
	log.Println("successfully created db")
	return nil
}

func (db *SqliteDatabase) createDatabase() {
	err := db.openDatabase()
	if err != nil {
		log.Fatal("Error creating Database file: ", err)
		return
	}

	_, err = db.currentConnection.Exec("CREATE TABLE if not exists data (id INTEGER PRIMARY KEY AUTOINCREMENT, temperature REAL NOT NULL,pressure REAL NOT NULL, relative_pressure REAL NOT NULL, humidity REAL NOT NULL, voltage REAL NOT NULL, created_at DATETIME NOT NULL)")
	if err != nil {
		log.Fatal("Error creating Database table: ", err)
		return
	}
}

func (db *SqliteDatabase) InsertWeatherData(data *data.WeatherData) error {
	err := db.openDatabase()
	if err != nil {
		log.Println("unable to write data to db:", err)
		return err
	}
	sql := "INSERT INTO data (temperature, pressure,relative_pressure,voltage,humidity,created_at) values (?,?,?,?,?,?)"
	stmt, err := db.currentConnection.Prepare(sql)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(
		data.Temperature,
		data.Pressure,
		data.RelPressure,
		data.Voltage,
		data.Humidity,
		time.Now().Format(time.RFC3339),
	)
	if err != nil {
		return err
	}

	return nil
}

func (db *SqliteDatabase) GetLastDataset() (*data.WeatherData, error) {
	err := db.openDatabase()
	if err != nil {
		log.Println("unable to read data from db:", err)
		return nil, err
	}
	rows, err := db.currentConnection.Query("SELECT * FROM data ORDER BY created_at DESC LIMIT 1;")
	if err != nil {
		log.Println("unable to perfrom query from db:", err)
		return nil, err
	}
	defer rows.Close()
	var returnData data.WeatherData
	for rows.Next() {
		var temperature float64
		var pressure float64
		var relative_pressure float64
		var voltage float64
		var humidity float64
		var id int
		var created_at time.Time
		err = rows.Scan(&id, &temperature, &pressure, &relative_pressure, &humidity, &voltage, &created_at)
		if err != nil {
			log.Println("malformed data query:", err)
			return nil, err
		}
		returnData = data.NewWeatherData(temperature, pressure, relative_pressure, voltage, humidity,created_at)
	}

	return &returnData, nil
}
