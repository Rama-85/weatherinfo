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

type weather_info struct {
	gorm.Model
	Id             uint      `json:"primarykey"`
	Time           string    `json:"time"`
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
	fmt.Println("weather_info table created....")
}

func main() {
	connectDB()
	createTable()

	var winfo weather_info
	resp, err := http.Get("https://api.open-meteo.com/v1/forecast?latitude=11.7117117&longitude=79.3271609&timezone=IST&hourly=temperature_2m&hourly=relativehumidity_2m&hourly=windspeed_10m&hourly=winddirection_10m&hourly=pressure_msl&hourly=soil_temperature_6cm&hourly=visibility&hourly=rain")
	//err = dbm.create(&winfo).Error
	if err != nil {
		fmt.Println("error while inserting data:", err)
		return
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&winfo)
	if err != nil {
		return
	}
	fmt.Println(winfo)
	return

}
