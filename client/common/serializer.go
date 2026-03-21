package common

import (
	"strings"
)

const separator = "|"

func SerializeBet(bet Bet) string{
	return strings.Join([]string{
		bet.FirstName,
		bet.LastName,
		bet.Document,
		bet.Birthdate,
		bet.Number,
		bet.AgencyID,
	}, separator)
}