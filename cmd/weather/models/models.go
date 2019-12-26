package models

//ArduinoSensorData receives temperature in Celcius
//and pressure in mmHg
type ArduinoSensorData struct {
	Temperature float32 `json:"temperature,omitempty" example:"22.1"`
	Pressure    float32 `json:"pressure,omitempty" example:"736.15"`
}

//WeatherData return type for api with pagination information
type WeatherData struct {
	Page int
	Data []ArduinoSensorData
}

//DataSource for changing data source
type DataSource interface {
	Save(sensorID int, data ArduinoSensorData) (bool, error)
	GetSensorsData(sensorID int) ([]ArduinoSensorData, error)
}
