package domain

type Temperature float32
type Humidity float32
type DeviceId int

type Measurement struct {
	DeviceId    DeviceId
	Temperature Temperature
	Humidity    Humidity
}
