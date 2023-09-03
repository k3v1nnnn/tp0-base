package common

import (
	"net"

	log "github.com/sirupsen/logrus"
)

type Connection struct {
	id string
	address string
	protocol string
	confirmationLength int
	conn  net.Conn
}

func newConnection(id string, address string) *Connection{
	connection := &Connection{
		id: id,
		address: address,
		protocol: "tcp",
	}
	return connection
}

func (c *Connection) start() net.Conn {
	conn, err := net.Dial(c.protocol, c.address)
	if err != nil {
		log.Fatalf(
			"action: connect | result: fail | client_id: %v | error: %v",
			c.id, err,
		)
	}
	c.conn = conn
	return c.conn
}

func (c *Connection) send(bet *Bet) bool{
	message := NewMessage()
	serializedMessage := message.serializeBet(bet)
	read, err := c.conn.Write(serializedMessage)
	if err != nil {
		log.Fatalf(
			"action: send | result: fail | client_id: %v | error: %v",
			c.id, err,
		)
	}
	return read == len(serializedMessage)
}

func (c *Connection) read() (string, bool) {
	message := NewMessage()
	buffer := make([]byte, c.confirmationLength)
	write, err := c.conn.Read(buffer)
	if err != nil {
		log.Fatalf(
			"action: read | result: fail | client_id: %v | error: %v",
			c.id, err,
		)
	}
	return message.deserializeBetConfirmation(buffer), write == c.confirmationLength
}

func (c *Connection) end() {
	err := c.conn.Close()
	if err != nil {
		log.Fatalf(
			"action: disconnect | result: fail | client_id: %v | error: %v",
			c.id, err,
		)
	}
}