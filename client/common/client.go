package common

import (
	"net"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("log")

// ClientConfig Configuration used by the client
type ClientConfig struct {
	ID             string
	ServerAddress  string
	FilePath       string
	BatchMaxAmount int
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

func (c *Client) closeConnection() {
	c.conn.Close()
	log.Infof("action: close_connection | result: success | client_id: %v", c.config.ID)
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
	log.Infof("action: connect | result: success | client_id: %v", c.config.ID)
	return nil
}

func (c *Client) StartClient(exitChan <-chan struct{}) {
	select {
	case <-exitChan:
		log.Infof("action: exit | result: success | client_id: %v", c.config.ID)
		return
	default:
	}
	if !runBetsHandler(c, exitChan) {
		return
	}
	runWinnersHandler(c, exitChan)
}
