package main

import (
	"fmt"
	"strings"

	"github.com/gen2brain/beeep"
	up "unifiedpush.org/go/dbus_connector/api"
	"unifiedpush.org/go/dbus_connector/definitions"
)

var Endpoint string

type NotificationHandler struct{}

func (n NotificationHandler) Message(instance, message, id string) {
	fmt.Println("new message received")
	// this message can be in whatever format you like, in this case the title and message body are two strings seperated by a '-'
	parts := strings.Split(message, "-")

	title := "No Title Provided"
	if len(parts) > 1 {
		title = parts[1]
	}

	err := beeep.Notify(title, parts[0], "")
	if err != nil {
		panic(err)
	}
}

func (n NotificationHandler) NewEndpoint(instance, endpoint string) {
	// the endpoint should be sent to whatever server your app is using
	Endpoint = endpoint
	fmt.Println("New endpoint received", Endpoint)
}

func (n NotificationHandler) Unregistered(instance string) {
	Endpoint = ""
	fmt.Println("endpoint unregistered", Endpoint)
}

func main() {
	connector := NotificationHandler{}
	up.InitializeAndCheck("cc.malhotra.karmanyaah.testapp.golibrary", "Golang Example App", connector)

	if len(up.GetDistributor()) == 0 { // not picked distributor yet
		pickDist()
	}
	// run this for each instance on each application startup to get the most up-to-date info
	result, reason, err := up.Register("")
	if err != nil {
		panic(err)
	}
	switch result {
	case definitions.RegisterStatusFailed:
		fmt.Println("registration failed because", reason)
		return
	case definitions.RegisterStatusRefused:
		fmt.Println("Registration refused", reason)
		return
	default:
		fmt.Println("will receive registration soon", reason)
	}

	// do whatever your app does
	fmt.Println("app waiting now")
	<-make(chan struct{})
}

func pickDist() {
	dist, err := up.GetDistributors()
	if err != nil {
		panic(err)
	}

	var distributor string

	if len(dist) == 0 {
		fmt.Printf("No distributor so can't be push notifications, exiting")
	} else if len(dist) == 1 {
		distributor = dist[0]
		fmt.Println("Picking the only distributor available", distributor)
	} else {
		fmt.Println("avalible distributors")
		for i, j := range dist {
			fmt.Println(i, j)
		}
		fmt.Print("Pick one distributor by number  ")
		var num int
		fmt.Scanln(&num)
		distributor = dist[num]
		fmt.Println("Picked distributor", distributor)
	}
	err = up.SaveDistributor(distributor)
	if err != nil {
		panic(err)
	}
}
