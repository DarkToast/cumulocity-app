package cumulocity

import (
	"time"
)

type DeviceId string
type Name string
type Owner string

type Device struct {
	Id             DeviceId
	Name           Name
	Owner          Owner
	Created        time.Time
	ChildDeviceIds []DeviceId
	ParentDeviceId *DeviceId
}
