package cumulocity

import (
	"time"
)

type Id string
type Name string
type Owner string

type Device struct {
	Id            Id
	Name          Name
	Owner         Owner
	CreationTime  time.Time
	ChildDevices  []Id
	ParentDevices []Id
}
