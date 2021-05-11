package main

import (
	"fmt"

	"github.com/unifiedpush/go_dbus_connector/dbus"
)

func main() {
	c := dbus.NewClient()
	c.InitializeDefaultConnection()
	defer c.Close()

	fmt.Println(c.PickDistributor("org.unifiedpush.Distributor.gotify"))
}
