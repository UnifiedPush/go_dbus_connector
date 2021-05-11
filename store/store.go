package store

import (
	"encoding/json"
	"io/ioutil"

	"github.com/google/uuid"
	"github.com/unifiedpush/go_dbus_connector/definitions"
)

type Storage struct {
	AppName     string
	Distributor string
	Instances   map[string]struct { //map key is instance
		Token uuid.UUID
	}
}

func (s *Storage) Commit() error {
	b, err := json.Marshal(s)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(definitions.StoragePath(s.AppName), b, definitions.ConnectorPerm)
	return err
}
