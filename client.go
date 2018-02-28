package zmq

import (
	"errors"
	"time"

	"github.com/pebbe/zmq4"
)

// Client represents a connection to an IOTA node that supports ZeroMQ.
type Client struct {
	address       string
	socket        *zmq4.Socket
	subscriptions map[MessageType]chan Message
	stopChan      chan struct{}
}

// NewClient returns a new client used for connecting to an IOTA node's ZeroMQ server
// located at the supplied address.
func NewClient(address string) (*Client, error) {
	c := Client{
		address:       address,
		subscriptions: make(map[MessageType]chan Message),
		stopChan:      make(chan struct{}),
	}

	socket, err := zmq4.NewSocket(zmq4.SUB)
	if err != nil {
		return nil, err
	}

	err = socket.SetRcvtimeo(5 * time.Second)
	if err != nil {
		panic(err)
	}

	c.socket = socket
	return &c, nil
}

// Connect connects to the previously supplied address and handles messages.
// Any previously subscribed messages will continue to function after
// reconnecting.
func (c *Client) Connect() error {
	go c.handleMessages()
	return c.connect()
}

func (c *Client) connect() error {
	return c.socket.Connect(c.address)
}

// Disconnect disconnects from the supplied address and stops the message
// handler.
func (c *Client) Disconnect() error {
	c.stopChan <- struct{}{}
	return c.disconnect()
}

func (c *Client) disconnect() error {
	return c.socket.Disconnect(c.address)
}

// Subscribe returns a channel through which messages of the requested type
// will be passed when received from the node.
func (c *Client) Subscribe(msg MessageType) (chan Message, error) {
	var ch chan Message
	var ok bool
	if ch, ok = c.subscriptions[msg]; !ok {
		ch = make(chan Message)
		c.subscriptions[msg] = ch
	}
	err := c.socket.SetSubscribe(msgTypes[msg])
	return ch, err
}

// Unsubscribe closes the channel through which messages of the requested type
// are being passed when received from the node.
func (c *Client) Unsubscribe(msg MessageType) error {
	if _, ok := c.subscriptions[msg]; !ok {
		return errors.New("No subscription exists for this message type")
	}
	close(c.subscriptions[msg])
	delete(c.subscriptions, msg)
	return c.socket.SetUnsubscribe(msgTypes[msg])
}
