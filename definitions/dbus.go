package definitions

import (
	"io/fs"
	"os"
	"path/filepath"
)

const (
	DistributorPrefix    = "org.unifiedpush.Distributor"
	DistributorPath      = "/org/unifiedpush/Distributor"
	DistributorInterface = "org.unifiedpush.Distributor1"
)

const (
	ConnectorPath      = "/org/unifiedpush/Connector"
	ConnectorInterface = "org.unifiedpush.Connector1"

	ConnectorBackgroundArgument = "UNIFIEDPUSH_DBUS_BACKGROUND_ACTIVATION"
)

const (
	ConnectorPerm fs.FileMode = 0600
)

//storagePaths provides a basic cache for StoragePath
var storagePaths = map[string]string{}

//StoragePath appName only recommends using something that can be a filename for now
func StoragePath(appName string) string {
	if a, ok := storagePaths[appName]; ok {
		return a
	}

	basedir := os.Getenv("XDG_CONFIG_HOME")
	if len(basedir) == 0 {
		basedir = os.Getenv("HOME")
		if len(basedir) == 0 {
			panic(nil)
			//TODO
		}
		basedir = filepath.Join(basedir, ".config")
	}
	basedir = filepath.Join(basedir, "unifiedpush", "connectors")
	err := os.MkdirAll(basedir, 0700)
	if err != nil {
	}
	finalFilename := filepath.Join(basedir, appName+".json")
	storagePaths[appName] = finalFilename
	return finalFilename
}

type RegisterStatus int

const (
	RegisterStatusNewEndpoint RegisterStatus = iota
	RegisterStatusRefused
	RegisterStatusFailed
	RegisterStatusFailedRequest = 99
)

var RegisterStatusMap = map[string]RegisterStatus{
	"NEW_ENDPOINT":         RegisterStatusNewEndpoint,
	"REGISTRATION_REFUSED": RegisterStatusRefused,
	"REGISTRATION_FAILED":  RegisterStatusFailed,
}
