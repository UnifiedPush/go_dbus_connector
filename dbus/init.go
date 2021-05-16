package dbus

import (
	"errors"
	"strings"

	"github.com/godbus/dbus/v5"
	"github.com/unifiedpush/go_dbus_connector/definitions"
)

type Client struct {
	conn *dbus.Conn
	name string
}

func NewClient(appName string) *Client {
	return &Client{name: appName}
}

func (c *Client) InitializeDefaultConnection() error {
	var err error
	c.conn, err = dbus.ConnectSessionBus()
	return err
}

func (c *Client) Close() {
	if c.conn != nil {
		err := c.conn.Close()
		if err != nil {
			panic(err) // huge error if ever can't close
		}
	}
}

func (c *Client) ListDistributors() (distributors []string, err error) {
	var s []string
	err = c.conn.BusObject().Call("org.freedesktop.DBus.ListNames", 0).Store(&s)
	if err != nil {
		return
	}

	for _, name := range s {
		if strings.HasPrefix(name, definitions.DistributorPrefix) {
			distributors = append(distributors, name)
		}
	}
	return
}

func (c *Client) PickDistributor(dist string) *Distributor {
	if !strings.HasPrefix(dist, definitions.DistributorPrefix) {
		return nil
	}

	object := c.conn.Object(dist, definitions.DistributorPath)
	return NewDistributor(object)
}

// StartHandling exports the connector interface and requests the app's name on the bus
func (c *Client) StartHandling(connector Connector) error {
	err := c.conn.Export(connector, definitions.ConnectorPath, definitions.ConnectorInterface)
	if err != nil {
		return err
	}

	name, err := c.conn.RequestName(c.name, dbus.NameFlagDoNotQueue)
	if err != nil {
		return err
	}
	if name != dbus.RequestNameReplyPrimaryOwner {
		return errors.New("Cannot request name on dbus")
	}

	return nil
}

func NewConnector(handler ConnectorHandler) Connector {
	return Connector{
		h: handler,
	}
}

type ConnectorHandler interface {
	Message(token, message, msgID string)
	NewEndpoint(token, endpoint string)
	Unregistered(token string)
}

type Connector struct {
	h ConnectorHandler
}

func (c Connector) Message(token, message, msgID string) *dbus.Error {
	c.h.Message(token, message, msgID)
	return nil
}

func (c Connector) NewEndpoint(token, endpoint string) *dbus.Error {
	c.h.NewEndpoint(token, endpoint)
	return nil
}

func (c Connector) Unregistered(token string) *dbus.Error {
	c.h.Unregistered(token)
	return nil
}
