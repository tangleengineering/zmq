package zmq

import (
	"time"

	"github.com/pebbe/zmq4"
)

type Client struct {
	address       string
	socket        *zmq4.Socket
	subscriptions map[MessageType]chan Message
	stopChan      chan struct{}
}

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

func (c *Client) Connect() error {
	go c.handleMessages()
	return c.connect()
}
func (c *Client) connect() error {
	return c.socket.Connect(c.address)
}

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
