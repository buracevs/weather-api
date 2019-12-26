package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/buracevs/weather-api/cmd/weather/models"
	"github.com/gorilla/mux"
)

type HandlerForHttp struct {
	ds models.DataSource
}

func NewHandlerForHttp(src models.DataSource) *HandlerForHttp {
	return &HandlerForHttp{
		ds: src,
	}
}

// SaveToDataBase godoc
// @Summary Creates new weather record in database
// @Success 204 "no body returned"
// @Router /{id}/add-data [post]
// @Param id path integer true "Sensor id"
// @Param data body models.ArduinoSensorData false "Sensor data"
// @Failure 500 {string} string "Server Error"
func (h *HandlerForHttp) SaveToDataBase(writer http.ResponseWriter, request *http.Request) {
	var parsedData models.ArduinoSensorData
	var deviceId int
	var err error

	vars := mux.Vars(request)
	if _, ok := vars["id"]; ok {
		deviceId, err = strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			log.Println("error", err)
			os.Exit(1)
		}
	}

	err = json.NewDecoder(request.Body).Decode(&parsedData)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		log.Println("error", err)
		os.Exit(1)
	}

	result, err := h.ds.Save(deviceId, parsedData)
	if err != nil {
		log.Fatal("err", err)
	}

	log.Println("result is ", result)

	writer.WriteHeader(http.StatusNoContent)
	writer.Write([]byte(""))
}

// GetDataRange godoc
// @Summary Returns data from database
// @Success 200 {array} models.ArduinoSensorData "Array of temperature/pressure data"
// @Router /get-data/{id} [get]
// @Param id path integer true "Sensor id"
// @Failure 500 {string} string "Server Error"
func (h *HandlerForHttp) GetDataRange(writer http.ResponseWriter, request *http.Request) {

	vars := mux.Vars(request)
	if _, ok := vars["id"]; ok {

		sensorId, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			log.Println("error", err)
			os.Exit(1)
		}

		weatherDt, err := h.ds.GetSensorsData(sensorId)
		if err != nil {
			log.Fatal("err", err)
		}
		data := models.WeatherData{
			Page: 1,
			Data: weatherDt,
		}

		bytes, err := json.Marshal(data)
		if err != nil {
			log.Fatal("err", err)
		}
		writer.Header().Set("Content-Type", "application/json; charset=UTF-8")
		writer.WriteHeader(http.StatusOK)
		writer.Write(bytes)
	}

}

/*func saveData(ds models.DataSource, sensorId int, data models.ArduinoSensorData) {
	saved, err := ds.Save(sensorId, data)
	if err != nil {

	}
}*/
