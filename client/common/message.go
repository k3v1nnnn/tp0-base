package common

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"strings"
)

const BetConfirmationMessage = "a"
const BetMaxLength = 70
const ConfigMaxLength = 8

type Message struct {
	betSeparator string
	separator string
	maxLength int
	filler string
}

func NewMessage() *Message {
	message := &Message{
		betSeparator: ",",
		separator: "#",
		maxLength: BetMaxLength, //bytes
		filler: "@",
	}
	return message
}

func (c *Message) serializeBet(bet *Bet) []byte {
	_info := bet.information()
	info := "b" + c.separator + strings.Join(_info, c.separator)
	missingBytes := c.maxLength - len(info)
	if missingBytes < 0 {
		info = ""
	} else {
		info = info + strings.Repeat(c.filler, missingBytes)
	}
	return []byte(info)
}

func (c *Message) deserializeBetConfirmation(buffer []byte) string {
	return string(buffer)
}

func (c *Message) serializeConfig(batchSize int, id string) []byte {
	info := fmt.Sprintf("bi#%d#%v", batchSize, id)
	missingBytes := ConfigMaxLength - len(info)
	info = info + strings.Repeat(c.filler, missingBytes)
	return []byte(info)
}

func (c *Message) serializeBets(bets []Bet, last bool, batchSize int) []byte {
	var _info []string
	for _, bet := range bets {
		_bet := bet.information()
		_betJoin := strings.Join(_bet, c.separator)
		_betJoin = strings.Replace(_betJoin, " ", "_", -1)
		_info = append(_info, _betJoin)
	}
	lastBetBatch := "bc"
	if last {
		lastBetBatch = "bf"
	}
	info := lastBetBatch + c.separator + strings.Join(_info, c.betSeparator)
	missingBytes := (batchSize * c.maxLength) - len(info)
	if missingBytes < 0 {
		log.Infof("action: bet_sent | result: in_progress | amount: %v | bytes: %v",
			len(bets), missingBytes)
		info = ""
	} else {
		info = info + strings.Repeat(c.filler, missingBytes)
	}
	return []byte(info)
}