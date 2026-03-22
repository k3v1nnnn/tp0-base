package common

type CsvBet struct {
	reader *CsvReader
}

func NewCSVBet(reader *CsvReader) *CsvBet {
	return &CsvBet{reader: reader}
}

func (cb *CsvBet) NextBet(agencyId string) (Bet, bool, error) {
	row, err := cb.reader.Next()
	if row == nil {
		return Bet{}, true, nil
	}
	if err != nil {
		return Bet{}, false, err
	}
	return Bet{
		AgencyID:  agencyId,
		FirstName: row[0],
		LastName:  row[1],
		Document:  row[2],
		Birthdate: row[3],
		Number:    row[4],
	}, false, nil
}