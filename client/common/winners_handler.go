package common

import "time"

func runWinnersHandler(c *Client, exitChan <-chan struct{}) {
	retries := 10
	for i := 0; i < retries; i++ {
		err := c.createClientSocket(exitChan)
		if err != nil {
			return
		}
		protocol := NewProtocol(c.conn)

		if err := protocol.Send(SerializeWinnerFlag(c.config.ID)); err != nil {
			log.Errorf("action: send_winners_query | result: fail | client_id: %v | error: %v", c.config.ID, err)
			c.closeConnection()
			return
		}
		log.Infof("action: send_winners_query | result: success | client_id: %v", c.config.ID)

		response, err := protocol.RecvResponse()
		if err != nil {
			log.Errorf("action: receive_winners_response | result: fail | client_id: %v | error: %v", c.config.ID, err)
			c.closeConnection()
			return
		}

		if response != "WAIT" {
			var winners []string
			if response != "" {
				winners = UnserializeWinners(response)
			}
			log.Infof("action: consulta_ganadores | result: success | cant_ganadores: %v", len(winners))
			c.closeConnection()
			return
		}
		c.closeConnection()
		if i < retries-1 {
			time.Sleep(time.Duration(1<<i) * time.Second)
		}
	}
	log.Errorf("action: consulta_ganadores | result: fail | client_id: %v | error: max retries exceeded", c.config.ID)
}
