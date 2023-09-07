package common

import "strings"

type Bet struct {
	agency string
	firstName string
	lastName  string
	dni string
	birthdate string
	number string
}

func NewBet(agency string, firstName string, lastName string, dni string, birthdate string, number string) *Bet {
	bet := &Bet{
		agency: agency,
		firstName: firstName,
		lastName: lastName,
		dni: dni,
		birthdate: birthdate,
		number: number,
	}
	return bet
}

func (b *Bet) information() []string {
	return []string{b.agency,
		strings.Replace(b.firstName, " ", "_", -1),
		strings.Replace(b.lastName, " ", "_", -1),
		b.dni, b.birthdate, b.number}
}
