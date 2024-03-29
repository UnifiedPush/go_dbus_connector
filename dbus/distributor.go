package dbus

import (
	"github.com/godbus/dbus/v5"
	"unifiedpush.org/go/dbus_connector/definitions"
)

func NewDistributor(object dbus.BusObject) *Distributor {
	return &Distributor{
		object: object,
	}
}

type Distributor struct {
	object dbus.BusObject
}

func (d *Distributor) Register(name, token, description string) (definitions.RegisterStatus, string) {
	var status, reason string
	err := d.object.Call(definitions.DistributorInterface+".Register", dbus.Flags(0), name, token, description).Store(&status, &reason)
	if err != nil {
		return definitions.RegisterStatusFailedRequest, ""
	}

	registerStatus, ok := definitions.RegisterStatusMap[status]
	if !ok {
		return definitions.RegisterStatusFailedRequest, ""
	}

	return registerStatus, reason
}

func (d *Distributor) Unregister(token string) (err error) {
	err = d.object.Call(definitions.DistributorInterface+".Unregister", dbus.FlagNoReplyExpected, token).Err
	return
}
