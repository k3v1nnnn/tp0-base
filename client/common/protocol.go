package common

import (
	"encoding/binary"
	"net"
)

type Protocol struct {
	conn   net.Conn
}

func NewProtocol(conn net.Conn) *Protocol {
	return &Protocol{conn: conn}
}

func (p *Protocol) send(data []byte) error {
	sent := 0
	for sent < len(data) {
		n, err := p.conn.Write(data[sent:])
		if err != nil {
			return err
		}
		sent += n
	}
	return nil
}

func (p *Protocol) recv(n int) ([]byte, error) {
	buf := make([]byte, n)
	received := 0
	for received < n {
		read, err := p.conn.Read(buf[received:])
		if err != nil {
			return nil, err
		}
		received += read
	}
	return buf, nil
}

func (p *Protocol) Send(payload string) error {
	data := []byte(payload)
	header := make([]byte, 4)
	binary.BigEndian.PutUint32(header, uint32(len(data)))
	return p.send(append(header, data...))
}

func (p *Protocol) RecvResponse() (string, error) {
	header, err := p.recv(4)
	if err != nil {
		return "", err
	}
	length := binary.BigEndian.Uint32(header)
	data, err := p.recv(int(length))
	if err != nil {
		return "", err
	}
	return string(data), nil
}