package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type OpenMeteoResponse struct {
	gorm.Model
	Latitude              float64      `json:"latitude"`
	Longitude             float64      `json:"longitude"`
	Generationtime_ms     float64      `json:"generationtime_ms"`
	Utc_offset_seconds    int          `json:"utc_offset_seconds"`
	Timezone              string       `json:"timezone"`
	Timezone_abbreviation string       `json:"timezone_abbrevation"`
	Elevation             float64      `json:"elevation"`
	Hourly_units          Hourly_units `json:"hourly_units"`
	Hourly                Hourly       `json:"hourly"`
}
type Hourly_units struct {
	Time           string `json:"time"`
	Temperature_3m string `json:"temperature_3m"`
}

type Hourly struct {
	//Id             uint      `json:"primarykey"`
	Time           []string  `json:"time"`
	Temperature_3m []float64 `json:"temperature_3m"`
}

var dbm *sql.DB

func connectDB() {
	//username := "root"
	//password := "root"
	//host := "localhost"
	//database := "weather"

	//dsn := username + ":" + password + "@tcp(" + host + ":3306)/" + database + "?charset=utf8mb4&parseTime=True&loc=Local"
	//db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/weather")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("database connected successfully...")
	dbm = db
}

func createTable() {
	query := `create table weather_info (
		    wid int primary key auto_increment,
			temperature_3m  float ,
			created_at datetime default CURRENT_TIMESTAMP
			 )`

	_, err := dbm.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(" table created....")
}

func main() {
	connectDB()
	createTable()

	var r OpenMeteoResponse
	resp, err := http.Get("https://api.open-meteo.com/v1/forecast?latitude=11.7117117&longitude=79.3271609&timezone=IST&hourly=temperature_2m&hourly=relativehumidity_2m&hourly=windspeed_10m&hourly=winddirection_10m&hourly=pressure_msl&hourly=soil_temperature_6cm&hourly=visibility&hourly=rain")
	//err = dbm.create(&hinfo).Error
	if err != nil {
		fmt.Println("error while inserting data:", err)
		return
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&r)
	if err != nil {
		return
	}
	fmt.Println(r)
	return

}
