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

func (c *Client) sendBatch(protocol *Protocol, payload string) error {
	if err := protocol.Send(payload); err != nil {
		log.Errorf("action: send_batch | result: fail | client_id: %v | error: %v", c.config.ID, err)
		return err
	}
	response, err := protocol.RecvResponse()
	if err != nil {
		log.Errorf("action: receive_response | result: fail | client_id: %v | error: %v", c.config.ID, err)
		return err
	}
	if response == "OK" {
		log.Infof("action: send_batch | result: success | client_id: %v", c.config.ID)
	} else {
		log.Errorf("action: send_batch | result: fail | client_id: %v | response: %v", c.config.ID, response)
	}
	return nil
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
	protocol := NewProtocol(c.conn)
	batcher := NewBatcher(c.config.BatchMaxAmount)
	reader := NewCSVReader(c.config.FilePath)
	err = reader.Open()
	if err != nil {
		log.Errorf("action: read_file | result: fail | client_id: %v | error: %v", c.config.ID, err)
		c.closeConnection()
		return
	}
	defer reader.Close()
	csvBet := NewCSVBet(reader)

	for {
		bet, finished, err := csvBet.NextBet(c.config.ID)
		if finished {
			break
		}
		if err != nil {
			log.Errorf("action: read_bet | result: fail | client_id: %v | error: %v", c.config.ID, err)
			c.closeConnection()
			return
		}
		if batcher.CanAddBet(bet) {
			batcher.AddBet(bet)
		} else {
			if err := c.sendBatch(protocol, batcher.GetBatch()); err != nil {
				c.closeConnection()
				return
			}
			batcher.Clean(bet)
		}
	}

	if !batcher.IsEmpty() {
		if err := c.sendBatch(protocol, batcher.GetBatch()); err != nil {
			c.closeConnection()
			return
		}
	}

	if err := protocol.Send("END"); err != nil {
		log.Errorf("action: send_finish_flag | result: fail | client_id: %v | error: %v", c.config.ID, err)
		c.closeConnection()
		return
	}

	response, err := protocol.RecvResponse()
	if err != nil {
		log.Errorf("action: receive_response | result: fail | client_id: %v | error: %v", c.config.ID, err)
		c.closeConnection()
		return
	}
	log.Infof("action: receive_response | result: success | client_id: %v | response: %v", c.config.ID, response)
	c.closeConnection()
}
