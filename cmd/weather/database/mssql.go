package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/buracevs/weather-api/cmd/weather/models"
	_ "github.com/denisenkom/go-mssqldb"
)

type MssqlDao struct {
	dataSource string
	userId     string
	password   string
}

func MakeMssqlDao(dataSource, userId, password string) MssqlDao {
	return MssqlDao{
		dataSource: dataSource,
		userId:     userId,
		password:   password,
	}
}

//Save saves sensor data into mssql database
func (sqlDao MssqlDao) Save(sensorId int, data models.ArduinoSensorData) (bool, error) {
	connectionString := fmt.Sprintf("data source=%s;user id=%s;password=%s; database=Weather", sqlDao.dataSource, sqlDao.userId, sqlDao.password)
	conn, err := sql.Open("sqlserver", connectionString)
	defer conn.Close()
	if err != nil {
		log.Fatal("Db error", err)
	}

	res, err := conn.Exec("insert into SensorData ([SensorId],[TempC],[Pressure]) values (@p1,@p2,@p3)", sensorId, data.Temperature, data.Pressure)
	if err != nil {
		log.Fatal("Db exec error", err)
	}
	ineserted, err := res.RowsAffected()
	if err != nil {
		log.Fatal("Db insert error", err)
	}
	if ineserted == 1 {
		return true, err
	}
	return false, err
}

//GetSensorsData returns set of data of a custom sensor
func (sqlDao MssqlDao) GetSensorsData(sensorId int) ([]models.ArduinoSensorData, error) {

	connectionString := fmt.Sprintf("data source=%s;user id=%s;password=%s; database=Weather", sqlDao.dataSource, sqlDao.userId, sqlDao.password)
	conn, err := sql.Open("sqlserver", connectionString)
	defer conn.Close()
	if err != nil {
		log.Fatal("Db error", err)
	}
	log.Println("sensorid", sensorId)
	weatherData := make([]models.ArduinoSensorData, 1)
	rows, err := conn.Query("select [TempC], [Pressure] from SensorData  where [SensorId]=@SensorId", sql.Named("SensorId", sensorId))
	log.Println(rows)
	if err != nil {
		log.Fatal("Db error", err)
	}

	for rows.Next() {
		response := models.ArduinoSensorData{}
		err = rows.Scan(&response.Temperature, &response.Pressure)
		if err != nil {
			log.Fatal("Db error", err)
		}
		weatherData = append(weatherData, response)
	}
	log.Println(weatherData, sensorId)
	return weatherData, err
}
