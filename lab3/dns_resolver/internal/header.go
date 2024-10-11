package internal

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

/*
https://www.rfc-editor.org/rfc/rfc1035
26 страница
*/

/*
	1  1  1  1  1  1
	0  1  2  3  4  5  6  7  8  9  0  1  2  3  4  5
	+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
	|                      ID                       |
	+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
	|QR|   Opcode  |AA|TC|RD|RA|   Z    |   RCODE   |
	+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
	|                    QDCOUNT                    |
	+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
	|                    ANCOUNT                    |
	+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
	|                    NSCOUNT                    |
	+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
	|                    ARCOUNT                    |
	+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
*/

type Header struct {
	ID      uint16
	QR      byte
	Opcode  Opcode
	AA      byte
	TC      byte
	RD      byte
	RA      byte
	Z       byte
	RCODE   RCODE
	QDCOUNT uint16
	ANCOUNT uint16
	NSCOUNT uint16
	ARCOUNT uint16
}

type RCODE byte

const (
	RCODENoError RCODE = iota
	RCODEFormatError
	RCODEServerFailure
	RCODENameError
	RCODENotImplemented
	RCODERefused
)

type Opcode byte

const (
	OpcodeQuery = iota
	OpcodeIquery
	OpcodeStatus
)

func UnmarshallHeader(raw []byte, h *Header) (int, error) {
	if h == nil {
		err := errors.New("header must be non-nil")
		return 0, err
	}

	if len(raw) != 12 {
		err := fmt.Errorf(
			"raw message does not have the expected size - %d",
			len(raw))
		return 0, err
	}

	var hb byte

	h.ID = uint16(raw[0])<<8 | uint16(raw[1])

	// read QR|OPCODE|AA|TC
	hb = raw[2]

	h.RD = hb & masks[0]
	h.TC = (hb >> 1) & masks[0]
	h.AA = (hb >> 2) & masks[0]
	h.Opcode = Opcode((hb >> 3) & masks[3])
	h.QR = (hb >> 7) & masks[0]

	// read RD|RA|Z|RCODE
	hb &= 0
	hb = raw[3]
	h.RCODE = RCODE(hb & masks[3])
	h.Z = (hb >> 4) & masks[2]
	h.RA = (hb >> 7) & masks[0]

	h.QDCOUNT = uint16(raw[4])<<8 | uint16(raw[5])
	h.ANCOUNT = uint16(raw[6])<<8 | uint16(raw[7])
	h.NSCOUNT = uint16(raw[8])<<8 | uint16(raw[9])
	h.ARCOUNT = uint16(raw[10])<<8 | uint16(raw[11])

	return 12, nil
}

func (h Header) Marshall() ([]byte, error) {
	buf := new(bytes.Buffer)

	err := binary.Write(buf, binary.BigEndian, h.ID)
	if err != nil {
		return nil, errors.New("failed to write ID field of header")
	}

	var hb byte // represent merged QR|OPCODE|AA|TC|RD header bits
	hb |= h.QR << 7
	hb |= byte(h.Opcode) << 3
	hb |= h.AA << 2
	hb |= h.TC << 1
	hb |= h.RD << 0

	err = binary.Write(buf, binary.BigEndian, hb)
	if err != nil {
		return nil, errors.New("failed to write QR|OPCODE|AA|TC|RD fields of header")
	}

	hb &= 0 // represent merged RA|Z|RCODE header bits
	hb |= h.RA << 7
	hb |= h.Z << 4
	hb |= byte(h.RCODE) << 0

	err = binary.Write(buf, binary.BigEndian, hb)
	if err != nil {
		return nil, errors.New("failed to write RA|Z|RCODE fields of header")
	}

	err = binary.Write(buf, binary.BigEndian, h.QDCOUNT)
	if err != nil {
		return nil, errors.New("failed to write QDCOUNT field of header")
	}

	err = binary.Write(buf, binary.BigEndian, h.ANCOUNT)
	if err != nil {
		return nil, errors.New("failed to write ANCOUNT field of header")
	}

	err = binary.Write(buf, binary.BigEndian, h.NSCOUNT)
	if err != nil {
		return nil, errors.New("failed to write NSCOUNT field of header")
	}

	err = binary.Write(buf, binary.BigEndian, h.ARCOUNT)
	if err != nil {
		return nil, errors.New("failed to write ARCOUNT field of header")
	}

	return buf.Bytes(), nil
}
