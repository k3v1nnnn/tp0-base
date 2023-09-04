package common

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
)

// ClientConfig Configuration used by the client
type ClientConfig struct {
	ID            string
	ServerAddress string
	LoopLapse     time.Duration
	LoopPeriod    time.Duration
}

// Client Entity that encapsulates how
type Client struct {
	config ClientConfig
	conn   *Connection
	bet *Bet
}

// NewClient Initializes a new client receiving the configuration
// as a parameter
func NewClient(config ClientConfig, bet *Bet) *Client {
	client := &Client{
		config: config,
		bet: bet,
		conn: newConnection(config.ID, config.ServerAddress),
	}
	return client
}

// StartClientLoop Send messages to the client until some time threshold is met
func (c *Client) StartClientLoop() {
	signalHandler := make(chan os.Signal, 1)
	signal.Notify(signalHandler, syscall.SIGTERM)
loop:
	// Send messages if the loopLapse threshold has not been surpassed
	for timeout := time.After(c.config.LoopLapse); ; {
		select {
		case <-timeout:
	        log.Infof("action: timeout_detected | result: success | client_id: %v",
                c.config.ID,
            )
			break loop
		case signalReceived := <-signalHandler:
			c.conn.end()
			log.Infof("action: graceful_exit | result: success | signal: %v", signalReceived.String())
			break loop
		default:
		}

		c.conn.start()
		wasBetSent := c.conn.sendBet(c.bet)
		if wasBetSent {
			log.Infof("action: bet_sent | result: in_progress | dni: %v | numero: %v",
				c.bet.dni, c.bet.number)
		}
		wasBetConfirm := c.conn.readConfirmation()
		if wasBetConfirm {
			log.Infof("action: bet_sent | result: success | dni: %v | numero: %v",
				c.bet.dni, c.bet.number)
			c.conn.end()
			break loop
		}
		time.Sleep(c.config.LoopPeriod)
	}
	log.Infof("action: loop_finished | result: success | client_id: %v", c.config.ID)
}
