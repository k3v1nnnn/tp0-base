package common

func sendBatch(c *Client, protocol *Protocol, payload string) error {
	if err := protocol.Send(payload); err != nil {
		log.Errorf("action: send_batch | result: fail | client_id: %v | error: %v", c.config.ID, err)
		return err
	}
	response, err := protocol.RecvResponse()
	if err != nil {
		log.Errorf("action: receive_response | result: fail | client_id: %v | error: %v", c.config.ID, err)
		return err
	}
	if response != "OK" {
		log.Errorf("action: send_batch | result: fail | client_id: %v | response: %v", c.config.ID, response)
	}
	return nil
}

func runBetsHandler(c *Client, exitChan <-chan struct{}) bool {
	err := c.createClientSocket(exitChan)
	if err != nil {
		return false
	}
	protocol := NewProtocol(c.conn)
	batcher := NewBatcher(c.config.BatchMaxAmount)
	reader := NewCSVReader(c.config.FilePath)
	err = reader.Open()
	if err != nil {
		log.Errorf("action: read_file | result: fail | client_id: %v | error: %v", c.config.ID, err)
		c.closeConnection()
		return false
	}
	defer reader.Close()
	csvBet := NewCSVBet(reader)

	for {
		select {
		case <-exitChan:
			log.Infof("action: exit | result: success | client_id: %v", c.config.ID)
			c.closeConnection()
			return false
		default:
		}
		bet, finished, err := csvBet.NextBet(c.config.ID)
		if finished {
			break
		}
		if err != nil {
			log.Errorf("action: read_bet | result: fail | client_id: %v | error: %v", c.config.ID, err)
			c.closeConnection()
			return false
		}
		if batcher.CanAddBet(bet) {
			batcher.AddBet(bet)
		} else {
			if err := sendBatch(c, protocol, batcher.GetBatch()); err != nil {
				c.closeConnection()
				return false
			}
			batcher.Clean(bet)
		}
	}

	if !batcher.IsEmpty() {
		if err := sendBatch(c, protocol, batcher.GetBatch()); err != nil {
			c.closeConnection()
			return false
		}
	}

	if err := protocol.Send("END"); err != nil {
		log.Errorf("action: send_end | result: fail | client_id: %v | error: %v", c.config.ID, err)
		c.closeConnection()
		return false
	}

	response, err := protocol.RecvResponse()
	if err != nil {
		log.Errorf("action: receive_end_response | result: fail | client_id: %v | error: %v", c.config.ID, err)
		c.closeConnection()
		return false
	}
	if response == "OK" {
		log.Infof("action: receive_end_response | result: success | client_id: %v", c.config.ID)
	} else {
		log.Errorf("action: receive_end_response | result: fail | client_id: %v | response: %v", c.config.ID, response)
	}
	c.closeConnection()
	return true
}
