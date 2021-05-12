package dbus

import (
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
		c.conn.Close()
		//TODO: does this need to be error handled
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
	//TODO check interface impl maybe in future
	//node, err := introspect.Call(object)
	//if err != nil {
	//	//TODO
	//}

	//good := false
	//for _, i := range node.Interfaces {
	//	if i.Name == definitions.DistributorInterface {
	//		good = true
	//		break
	//	}
	//}

	//if !good {
	//	return nil
	//}
	return NewDistributor(object)
}

//StartHandling exports the connector interface and requests the app's name on the bus
func (c *Client) StartHandling(connector Connector) error {
	err := c.conn.Export(connector, definitions.ConnectorPath, definitions.ConnectorInterface)
	if err != nil {
		//TODO
	}

	//TODO introspect?
	//	n = introspect.Node{
	//	}
	//
	//	err = c.conn.Export(introspect.NewIntrospectable(n), definitions.ConnectorPath, "org.freedesktop.DBus.Introspectable")

	name, err := c.conn.RequestName(c.name, dbus.NameFlagDoNotQueue)
	if err != nil || name != dbus.RequestNameReplyPrimaryOwner {
		//TODO
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
