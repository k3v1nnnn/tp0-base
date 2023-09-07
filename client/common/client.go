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
	FilePath string
	BatchSize int
}

// Client Entity that encapsulates how
type Client struct {
	config ClientConfig
	conn   *Connection
	bet *Bet
	file *File
}

// NewClient Initializes a new client receiving the configuration
// as a parameter
func NewClient(config ClientConfig, bet *Bet) *Client {
	client := &Client{
		config: config,
		bet: bet,
		conn: newConnection(config.ID, config.ServerAddress),
		file: NewFile(config.FilePath),
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
			c.file.Close()
			c.conn.end()
			log.Infof("action: graceful_exit | result: success | signal: %v", signalReceived.String())
			break loop
		default:
		}

		c.conn.start()

		// Send Config
		wasConfigSent := c.conn.sendConfig(c.config.BatchSize, c.config.ID)
		if !wasConfigSent {
			log.Fatalf("action: config_sent | result: fail")
			c.conn.end()
			break loop
		}
		wasConfigConfirm := c.conn.readConfirmation()
		if !wasConfigConfirm {
			log.Fatalf("action: config_confirmation | result: fail")
			c.conn.end()
			break loop
		}

		c.file.Open()
		var bets []Bet
		log.Infof("action: bet_sent | result: in_progress")
		for {
			line := c.file.ReadLine()
			if line != "" {
				if len(bets) < c.config.BatchSize {
					bets = append(bets, c.file.getBet(c.config.ID, line))
				} else {
					wasBetSent := c.conn.sendBetBatch(bets, false, c.config.BatchSize)
					if !wasBetSent {
						log.Fatalf("action: bet_batch_sent | result: fail | bytes: %v",
							len(bets))
					}
					wasBetConfirm := c.conn.readConfirmation()
					if !wasBetConfirm {
						log.Fatalf("action: bet_batch_confirm | result: fail | bytes: %v",
							len(bets))
					}
					bets = []Bet{c.file.getBet(c.config.ID, line)}
				}
			} else {
				break
			}
		}
		if len(bets) > 0 {
			wasBetSent := c.conn.sendBetBatch(bets, true, c.config.BatchSize)
			if !wasBetSent {
				log.Fatalf("action: bet_batch_sent | result: fail | bytes: %v",
					len(bets))
			}
		}
		c.file.Close()
		wasBetConfirm := c.conn.readConfirmation()
		if !wasBetConfirm {
			log.Fatalf("action: bet_batch_confirm | result: fail | bytes: %v",
				len(bets))
		}
		log.Infof("action: bet_sent | result: success")
		c.conn.end()
		//Winners Request
		attempt := 1
		for {
			c.conn.start()
			wasRequestSent := c.conn.sendRequest(c.config.ID)
			if !wasRequestSent {
				log.Fatalf("action: winner_request_sent | result: fail")
			}
			hasResponse, winners := c.conn.readResponse()
			if !hasResponse {
				log.Fatalf("action: winner_response | result: fail")
			}
			if winners != "" {
				log.Infof("action: request_winners | result: success | amount_winners: %v",
					winners)
				c.conn.end()
				break
			}
			c.conn.end()
			log.Infof("action: request_winners | result: in_progress | sleep: %v",
				attempt)
			time.Sleep(c.config.LoopPeriod)
			attempt = attempt + 1
		}
		break loop
	}
	log.Infof("action: loop_finished | result: success | client_id: %v", c.config.ID)
}
