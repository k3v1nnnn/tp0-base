package common

import (
	"net"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("log")

// ClientConfig Configuration used by the client
type ClientConfig struct {
	ID            string
	ServerAddress string
	FirstName     string
	LastName      string
	Document      string
	Birthdate     string
	Number        string
}

// Client Entity that encapsulates how
type Client struct {
	config ClientConfig
	conn   net.Conn
}

// NewClient Initializes a new client receiving the configuration
// as a parameter
func NewClient(config ClientConfig) *Client {
	client := &Client{
		config: config,
	}
	return client
}

// CreateClientSocket Initializes client socket. In case of
// failure, error is printed in stdout/stderr and exit 1
// is returned
func (c *Client) createClientSocket(exitChan <-chan struct{}) error {
	conn, err := net.Dial("tcp", c.config.ServerAddress)
	if err != nil {
		select {
			case <-exitChan:
				log.Infof("action: close_connection | result: success | client_id: %v", c.config.ID)
			default:
				log.Criticalf("action: connect | result: fail | client_id: %v | error: %v", c.config.ID, err)
		}
		return err
	}
	c.conn = conn
	return nil
}

func (c *Client) StartClient(exitChan <-chan struct{}) {
	select {
		case <-exitChan:
			log.Infof("action: exit | result: success | client_id: %v", c.config.ID)
			return
		default:
	}
	err := c.createClientSocket(exitChan)
	if err != nil {
		return
	}

	bet := Bet{
		AgencyID:  c.config.ID,
		FirstName: c.config.FirstName,
		LastName:  c.config.LastName,
		Document:  c.config.Document,
		Birthdate: c.config.Birthdate,
		Number:    c.config.Number,
	}

	protocol := NewProtocol(c.conn)
	payload := SerializeBet(bet)

	if err := protocol.Send(payload); err != nil {
		log.Errorf("action: send_bet | result: fail | client_id: %v | error: %v",
			c.config.ID, err)
		c.conn.Close()
		return
	}

	response, err := protocol.RecvResponse()
	if err != nil {
		log.Errorf("action: receive_response | result: fail | client_id: %v | error: %v",
			c.config.ID, err)
		c.conn.Close()
		log.Infof("action: close_connection | result: success | client_id: %v", c.config.ID)
		return
	}
	log.Infof("action: receive_response | result: success | client_id: %v | response: %v",
		c.config.ID, response)

	if response == "OK" {
		log.Infof("action: apuesta_enviada | result: success | dni: %v | numero: %v",
			c.config.Document, c.config.Number)
	} else {
		log.Errorf("action: apuesta_enviada | result: fail | client_id: %v | response: %v",
			c.config.ID, response)
	}

	c.conn.Close()
	log.Infof("action: close_connection | result: success | client_id: %v", c.config.ID)
}
