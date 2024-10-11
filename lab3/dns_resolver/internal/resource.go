package internal

import "fmt"

/*
https://www.rfc-editor.org/rfc/rfc1035
28 страница
*/

/*
	1  1  1  1  1  1
	0  1  2  3  4  5  6  7  8  9  0  1  2  3  4  5
	+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
	|                                               |
	/                                               /
	/                      NAME                     /
	|                                               |
	+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
	|                      TYPE                     |
	+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
	|                     CLASS                     |
	+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
	|                      TTL                      |
	|                                               |
	+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
	|                   RDLENGTH                    |
	+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--|
	/                     RDATA                     /
	/                                               /
	+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
*/

type Resource struct {
	NAME     string
	TYPE     QType
	CLASS    QClass
	TTL      uint32
	RDLENGTH uint16
	RDATA    []byte
}

func UnmarshallResource(msg []byte, r *Resource) (int, error) {
	if r == nil {
		return 0, fmt.Errorf("resource must be non-nil")
	}

	if len(msg) < 11 {
		return 0, fmt.Errorf("resource msg must be at least 11 bytes long")
	}

	var (
		compressedName = new(CompressedLabel)
	)

	_, err := UnmarshallCompressedLabel(msg[0:2], compressedName)
	if err != nil {
		return 0, fmt.Errorf("failed to parse compressed name: %w", err)
	}

	r.TYPE = QType(uint16(msg[2]) | uint16(msg[3]))
	r.CLASS = QClass(uint16(msg[4]) | uint16(msg[5]))
	r.TTL = uint32(msg[6]) | uint32(msg[7]) | uint32(msg[8]) | uint32(9)
	r.RDLENGTH = uint16(msg[10]) | uint16(msg[11])
	r.RDATA = msg[12 : 12+r.RDLENGTH]

	return int(12 + r.RDLENGTH), nil
}

func (r *Resource) Marshal() (res []byte, err error) {
	// TODO implement
	return
}
