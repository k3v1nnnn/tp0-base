package common

import (
	"bufio"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"strings"
)

type File struct {
	path string
	fd *os.File
	reader *bufio.Reader
}

func NewFile(filepath string) *File {
	file := &File{
		path: filepath,
	}
	return file
}

func (f *File) Open() {
	file, err := os.Open(f.path)
	if err != nil {
		log.Fatalf(
			"action: read_bets | result: fail | error: %v",
			err,
		)
	}
	f.fd = file
	f.reader = bufio.NewReader(f.fd)
}

func (f *File) Close()  {
	err := f.fd.Close()
	if err != nil {
		log.Fatalf(
			"action: read_bets | result: fail | error: %v",
			err,
		)
	}
}

func (f *File) ReadLine() string {
	line, _, err := f.reader.ReadLine()
	if err != nil {
		if err == io.EOF {
			return ""
		}
		log.Fatalf("read file line error: %v", err)
		return ""
	}
	return string(line)
}

func (f *File) getBet(client string, line string) Bet {
	infoBet := strings.Split(line, ",")
	return Bet{client, infoBet[0], infoBet[1], infoBet[2], infoBet[3], infoBet[4]}
}
