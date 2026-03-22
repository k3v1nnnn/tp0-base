package common

import "errors"

var ErrEndOfFile = errors.New("end of file")

type CsvBet struct {
	reader *CsvReader
}

func NewCSVBet(reader *CsvReader) *CsvBet {
	csvBet := &CsvBet{
		reader: reader,
	}
	return csvBet
}

func (cb *CsvBet) NextBet(agencyId string) (Bet, error) {
	row, err := cb.reader.Next()
	if err != nil {
		return Bet{}, err
	}
	if row == nil {
		return Bet{}, ErrEndOfFile
	}
	return Bet{
		AgencyID:  agencyId,
		FirstName: row[0],
		LastName:  row[1],
		Document:  row[2],
		Birthdate: row[3],
		Number:    row[4],
	}, nil
}