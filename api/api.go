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

// ErrInstanceNotUnregistered informs if instances are not unregistered when executing a distributor change method
var ErrInstanceNotUnregistered = errors.New("Instance isn't unregistered yet")

var client *dbus.Client

// maybe mutex the globals?
var dataStore *store.Storage

type connector struct {
	c dbus.ConnectorHandler
	t chan time.Time
}

func (ch connector) Message(token, msg, id string) {
	if ch.t != nil {
		ch.t <- time.Now()
	}
	instance := getInstance(token)
	ch.c.Message(instance, msg, id)
}

func (ch connector) NewEndpoint(token, endpoint string) {
	if ch.t != nil {
		ch.t <- time.Now()
	}

	instance := getInstance(token)
	ch.c.NewEndpoint(instance, endpoint)
}

func (ch connector) Unregistered(token string) {
	if ch.t != nil {
		ch.t <- time.Now()
	}
	instance := getInstance(token)
	// FIXME how should I handle this
	removeToken(token) //nolint:errcheck
	ch.c.Unregistered(instance)
}

// Initializes the bus and object
func Initialize(name string, handler dbus.ConnectorHandler) error {
	if len(name) == 0 {
		return errors.New("invalid name")
	}
	if client != nil {
		client.Close()
	}
	client = dbus.NewClient(name)
	err := client.InitializeDefaultConnection()
	if err != nil {
		return errors.New("DBus Error")
	}

	// if the handler passed in is already of type connector (from InitializeAndCheck), don't wrap it in another connector. If its not then wrap with connector
	var conn connector
	var ok bool
	if conn, ok = handler.(connector); !ok {
		conn = connector{c: handler}
	}

	err = client.StartHandling(dbus.NewConnector(conn))
	if err != nil {
		return errors.New("DBus Error")
	}
	dataStore = store.NewStorage(name)
	if dataStore == nil {
		return errors.New("Storage Err")
	}
	return nil
}

// InitializeAndCheck is a convienience method that handles initialization and checking whether app started in the background.
// The background check checks whether the argument UNIFIEDPUSH_DBUS_BACKGROUND_ACTIVATION is input.
// Listens for 3 seconds after the last message and then exits.
// if running in the background this panics on error
func InitializeAndCheck(name string, handler dbus.ConnectorHandler) error {
	if !containsString(os.Args, definitions.ConnectorBackgroundArgument) {
		return Initialize(name, handler)
	}
	lastCallTime := make(chan time.Time)

	err := Initialize(name, connector{c: handler, t: lastCallTime})
	if err != nil {
		panic(err)
		// panic bc in bg listener
	}

	func() {
		for {
			// if another message arrives or 3 seconds happen whichever first
			select {
			case <-lastCallTime:
				continue
			case <-time.After(3 * time.Second):
				return
			}
		}
	}()
	client.Close()
	os.Exit(0)
	return nil
}

// Actual UP methods

// TryRegister registers a new instance.
// value of instance can be empty string for the default instance
// registration endpoint is returned through the callback if method is successful
func Register(instance string) (registerStatus definitions.RegisterStatus, registrationFailReason string, err error) {
	if len(GetDistributor()) == 0 {
		err = errors.New("No distributor selected")
		return
	}

	if len(getToken(instance)) == 0 {
		err = saveNewToken(instance)
		if err != nil {
			return
		}
	}
	status, reason := client.PickDistributor(GetDistributor()).Register(dataStore.AppName, getToken(instance))
	if status == definitions.RegisterStatusFailed || status == definitions.RegisterStatusRefused {
		err = removeToken(instance)
	}
	return status, reason, err
}

// TryUnregister attempts unregister, results are returned through callback
// any error returned is before unregister requested from dbus
func TryUnregister(instance string) error {
	return client.PickDistributor(GetDistributor()).Unregister(getToken(instance))
}

// Distributor things

// GetDistributor returns current selected distributor or empty string
func GetDistributor() string {
	return dataStore.Distributor
}

// GetDistributors lists all distributors that are available to register with
// (note the difference from GetDistributor singular)
func GetDistributors() ([]string, error) {
	return client.ListDistributors()
}

// SaveDistributor saves the distributor preference to use for future registrations
// valid values are picked from the list returned by GetDistributors
// all instances (registered to the previous distributor) have to unregister before running this
func SaveDistributor(id string) error {
	if err := storeIsEmpty(); err != nil {
		return err
	}

	// checks if valid distrib
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

// Token things

// getToken returns token for instance or empty string if instance doesn't exist
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
	dataStore.Instances[instance] = store.Instance{Token: token.String()}
	return dataStore.Commit()
}

func removeToken(instance string) error {
	delete(dataStore.Instances, instance)
	return dataStore.Commit()
}

// getInstance returns instance from token (for internal use) or empty string if not found
func getInstance(token string) string {
	for i, j := range dataStore.Instances {
		if token == j.Token {
			return i
		}
	}

	return ""
}

// utils

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
		return ErrInstanceNotUnregistered
	}
	return nil
}
