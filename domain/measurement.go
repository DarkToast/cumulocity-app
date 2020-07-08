package domain

type Temperature float64
type Humidity float64
type AirPressure float64
type DeviceId int

type Measurement struct {
	DeviceId    DeviceId
	Temperature Temperature
	Humidity    Humidity
	AirPressure AirPressure
}
