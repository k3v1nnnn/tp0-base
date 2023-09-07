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
	responseLength int
}

func newConnection(id string, address string) *Connection{
	connection := &Connection{
		id: id,
		address: address,
		protocol: "tcp",
		confirmationLength: 1,
		responseLength:70,
	}
	return connection
}

func (c *Connection) start() {
	conn, err := net.Dial(c.protocol, c.address)
	if err != nil {
		log.Fatalf(
			"action: connect | result: fail | client_id: %v | error: %v",
			c.id, err,
		)
	}
	c.conn = conn
}

func (c *Connection) sendBet(bet *Bet) bool{
	message := NewMessage()
	serializedMessage := message.serializeBet(bet)
	write, err := c.conn.Write(serializedMessage)
	if err != nil {
		log.Fatalf(
			"action: bet_sent | result: fail | dni: %v | numero: %v | error: %v",
			bet.dni, bet.number, err)
	}
	return write == len(serializedMessage)
}

func (c *Connection) sendBetBatch(bets []Bet, last bool, maxBatchSize int) bool{
	message := NewMessage()
	serializedMessage := message.serializeBets(bets, last, maxBatchSize)
	write, err := c.conn.Write(serializedMessage)
	if err != nil {
		log.Fatalf(
			"action: bet_batch_sent | result: fail | amount: %v | error: %v",
			len(bets), err)
	}
	return write == len(serializedMessage)
}

func (c *Connection) sendConfig(maxBatchSize int, id string) bool{
	message := NewMessage()
	serializedMessage := message.serializeConfig(maxBatchSize, id)
	write, err := c.conn.Write(serializedMessage)
	if err != nil {
		log.Fatalf(
			"action: config_sent | result: fail | error: %v",
			err)
	}
	return write == len(serializedMessage)
}

func (c *Connection) readConfirmation() bool {
	message := NewMessage()
	buffer := make([]byte, c.confirmationLength)
	read, err := c.conn.Read(buffer)
	if err != nil {
		log.Fatalf(
			"action: bet_sent | result: fail | client_id: %v | error: %v",
			c.id, err,
		)
	}
	confirmation:= message.deserializeBetConfirmation(buffer)
	return confirmation == BetConfirmationMessage && read == len(buffer)
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

func (c *Connection) sendRequest(id string) bool {
	message := NewMessage()
	serializedMessage := message.serializeRequest(id)
	write, err := c.conn.Write(serializedMessage)
	if err != nil {
		log.Fatalf(
			"action: request_sent | result: fail | error: %v",
			err)
	}
	return write == len(serializedMessage)
}

func (c *Connection) readResponse() (bool, int) {
	message := NewMessage()
	buffer := make([]byte, c.responseLength)
	read, err := c.conn.Read(buffer)
	if err != nil {
		log.Fatalf(
			"action: response_read | result: fail | client_id: %v | error: %v",
			c.id, err,
		)
	}
	return read == len(buffer), message.deserializeResponse(buffer)
}