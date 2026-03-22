package common

import (
    "encoding/csv"
    "os"
	"io"
)

type CsvReader struct {
	path string
	file *os.File
	csv *csv.Reader
}

func NewCSVReader(path string) *CsvReader {
	csvReader := &CsvReader{
		path: path,
	}
	return csvReader
}

func (cr *CsvReader) Open() error{
	var err error
	cr.file, err = os.Open(cr.path)
	if err == nil {
		cr.csv = csv.NewReader(cr.file)
	}
	return err
}

func (cr *CsvReader) Close() error{
	err := cr.file.Close()
	cr.file = nil
	cr.csv = nil
	return err
}


func (cr *CsvReader) Next() ([]string, error) {
    if cr.csv == nil {
        return nil, nil
    }
    row, err := cr.csv.Read()
    if err != nil {
		if err == io.EOF {
			return nil, nil
		}
        return nil, err
    }
    return row, nil
}