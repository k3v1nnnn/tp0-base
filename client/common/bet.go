package common
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
	return []string{b.agency, b.firstName, b.lastName, b.dni, b.birthdate, b.number}
}
