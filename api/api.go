package api

import (
	"errors"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/unifiedpush/go_dbus_connector/dbus"
	"github.com/unifiedpush/go_dbus_connector/definitions"
	"github.com/unifiedpush/go_dbus_connector/store"
)

var client *dbus.Client
var dataStore *store.Storage

type connector struct {
	c dbus.ConnectorHandler
	t chan time.Time
}

func (ch connector) Message(a, b, c string) {
	if ch.t != nil {
		ch.t <- time.Now()
	}
	ch.c.Message(a, b, c)
}
func (ch connector) NewEndpoint(a, b string) {
	if ch.t != nil {
		ch.t <- time.Now()
	}
	ch.c.NewEndpoint(a, b)
}

func (ch connector) Unregistered(a string) {
	if ch.t != nil {
		ch.t <- time.Now()
	}
	removeToken(a)
	ch.c.Unregistered(a)
}

//Initializes the bus and object
func Initialize(name string, handler dbus.ConnectorHandler) {
	if len(name) == 0 {
		//TODO err
	}
	if client != nil {
		client.Close()
	}
	client = dbus.NewClient(name)
	err := client.InitializeDefaultConnection()
	if err != nil {
		//TODO
	}
	err = client.StartHandling(dbus.NewConnector(connector{c: handler}))
	if err != nil {
		//TODO
	}
	dataStore = store.NewStorage(name)
}

//InitializeAndCheck is a convienience method that handles initialization and checking whether app started in the background.
// The background check checks whether the argument UNIFIEDPUSH_DBUS_BACKGROUND_ACTIVATION is input.
// Listens for 3 seconds after the last message and then stops.
func InitializeAndCheck(name string, handler dbus.ConnectorHandler) {
	if !containsString(os.Args, definitions.ConnectorBackgroundArgument) {
		Initialize(name, handler)
		return
	}
	lastCallTime := make(chan time.Time)

	//TODO might result in double running through connector,
	// should only be a problem for unregister but inefficient find better architecture
	Initialize(name, connector{c: handler, t: lastCallTime})

	for {
		//if another message arrives or 3 seconds happen whichever first
		select {
		case <-lastCallTime:
			continue
			//TODO make time adjustable?
		case <-time.After(3 * time.Second):
			break
		}
	}
	os.Exit(0)
}

//Actual UP methods

// TryRegister registers a new instance.
// value of instance can be empty string for the default instance
// registration endpoint is returned through the callback if method is successful
func Register(instance string) (registerStatus definitions.RegisterStatus, registrationFailReason string, err error) {
	if len(GetDistributor()) == 0 {
		err = errors.New("No distributor selected")
		return
	}

	err = saveNewToken(instance)
	if err != nil {
		return
	}
	status, reason := client.PickDistributor(GetDistributor()).Register(dataStore.AppName, getToken(instance))
	if status == definitions.RegisterStatusFailed || status == definitions.RegisterStatusRefused {
		err = removeToken(instance)
	}
	return status, reason, err
}

//TryUnregister attempts unregister, results are returned through callback
func TryUnregister(instance string) {
	client.PickDistributor(GetDistributor()).Unregister(getToken(instance))
}

//Distributor things

//GetDistributor returns current selected distributor or empty string
func GetDistributor() string {
	return dataStore.Distributor
}

//GetDistributors lists all distributors that are available to register with
// (note the difference from GetDistributor singular)
func GetDistributors() ([]string, error) {
	return client.ListDistributors()
}

//SaveDistributor saves the distributor preference to use for future registrations
// valid values are picked from the list returned by GetDistributors
// all instances (registered to the previous distributor) have to unregister before running this
func SaveDistributor(id string) error {
	if err := storeIsEmpty(); err != nil {
		return err
	}

	//TODO @S1m should I force this check of ensuring input is a valid distrib
	if s, err := GetDistributors(); err == nil {
		if valid := containsString(s, id); !valid {
			return errors.New("Not an ID of a valid distributor")
		}
	}

	dataStore.Distributor = id
	return dataStore.Commit()
}

// RemoveDistributor removes the currently set distributor
// all instances (registered to the previous distributor) have to unregister before running this
func RemoveDistributor() error {
	if err := storeIsEmpty(); err != nil {
		return err
	}
	dataStore.Distributor = ""
	return dataStore.Commit()
}

//Token things

//getToken returns token for instance or empty string if instance doesn't exist
func getToken(instance string) string {
	a, ok := dataStore.Instances[instance]
	if !ok {
		return ""
	}
	return a.Token
}

func saveNewToken(instance string) error {
	token, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	dataStore.Instances[instance] = store.Instance{token.String()}
	return dataStore.Commit()
}

func removeToken(instance string) error {
	delete(dataStore.Instances, instance)
	return dataStore.Commit()
}

//getInstance returns instance from token (for internal use) or empty string if not found
func getInstance(token string) string {
	for i, j := range dataStore.Instances {
		if token == j.Token {
			return i
		}
	}

	return ""
}

//utils

func containsString(a []string, b string) bool {
	for _, i := range a {
		if b == i {
			return true
		}
	}
	return false
}
func storeIsEmpty() error {
	if len(dataStore.Instances) != 0 {
		return errors.New("Instances are not unregistered") //TODO define error properly
	}
	return nil
}
