package store

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/unifiedpush/go_dbus_connector/definitions"
)

type Instance struct {
	Token string
}

func NewStorage(appName string) *Storage {
	var st Storage
	b, err := os.ReadFile(definitions.StoragePath(appName))
	if errors.Is(err, os.ErrNotExist) {
		return &Storage{
			AppName:   appName,
			Instances: map[string]Instance{},
		}
	}
	err = json.Unmarshal(b, &st)
	if err != nil {
		return nil
	}
	return &st
}

type Storage struct {
	AppName     string
	Distributor string
	//map key is instance name
	Instances map[string]Instance
}

func (s *Storage) Commit() error {
	b, err := json.Marshal(s)
	if err != nil {
		return err
	}
	err = os.WriteFile(definitions.StoragePath(s.AppName), b, definitions.ConnectorPerm)
	return err
}
