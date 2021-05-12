package api

import (
	"errors"

	"github.com/unifiedpush/go_dbus_connector/dbus"
	"github.com/unifiedpush/go_dbus_connector/store"
)

var client *dbus.Client
var dataStore *store.Storage

func Initialize(name string) {
	if client != nil {
		client.Close()
	}
	client = dbus.NewClient(name)
	err := client.InitializeDefaultConnection()
	if err != nil {
		//TODO
	}

}

//GetDistributor returns current selected distributor or empty string
func GetDistributor() string {
	return dataStore.Distributor
}

//GetDistributors lists all distributors that are available to register with
// (note the difference from GetDistributor singular)
func GetDistributors() ([]string, error) {
	return client.ListDistributors()
}

//GetInstance returns instance from token (for internal use) or empty string if not found
func GetInstance(token string) string {
	for i, j := range dataStore.Instances {
		if token == j.Token.String() {
			return i
		}
	}

	return ""
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
		if valid := inStringSlice(s, id); !valid {
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

func inStringSlice(a []string, b string) bool {
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
