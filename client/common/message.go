package common

import "strings"

const BetConfirmationMessage = "a"

type Message struct {
	separator string
	maxLength int
	filler string
}

func NewMessage() *Message {
	message := &Message{
		separator: "#",
		maxLength: 80, //bytes
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