package common

import "strings"

const separator = "|"

func SerializeBet(bet Bet) string {
	return strings.Join([]string{
		bet.FirstName,
		bet.LastName,
		bet.Document,
		bet.Birthdate,
		bet.Number,
		bet.AgencyID,
	}, separator)
}

func SerializeWinnerFlag(agency string) string {
	return strings.Join([]string{"WINNERS", agency}, separator)
}

func UnserializeWinners(winners string) []string {
	return strings.Split(winners, separator)
}
