package internal

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
)

/*
https://www.rfc-editor.org/rfc/rfc1035
27 страница
*/

/*
0  1  2  3  4  5  6  7  8  9  0  1  2  3  4  5
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
|                                               |
/                     QNAME                     /
/                                               /
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
|                     QTYPE                     |
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
|                     QCLASS                    |
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
*/
type Question struct {
	QNAME  string
	QTYPE  QType
	QCLASS QClass
}

type QType uint16

const (
	QTypeUnknown QType = iota

	// Host address
	QTypeA

	// Authoritative name server
	QTypeNS

	QTypeMD
	QTypeMF

	// Canonical name for an alias
	QTypeCNAME

	// Marks the start of a zone of authority
	QTypeSOA

	QTypeMB
	QTypeMG
	QTypeMR
	QTypeNULL
	QTypeWKS

	// Domain name pointer
	QTypePTR
	QTypeHINFO
	QTypeMINFO

	// Mail exchange
	QTypeMX
	QTypeTXT
	QTypeAXFR  QType = 252
	QTypeMAILB QType = 253
	QTypeMAILA QType = 254

	// All records
	QTypeWildcard QType = 255
)

const QTypeAAAA QType = 28

type QClass uint16

const (
	QClassUnknown QClass = iota

	// Internet
	QClassIN

	QClassCS
	QClassCH
	QClassHS

	// Any class
	QClassWildcard QClass = 255
)

func UnmarshallQuestion(raw []byte, q *Question) (int, error) {
	if q == nil {
		return 0, fmt.Errorf("question must be non-nil")
	}

	var (
		idx    = 0
		size   = 0
		labels []string
	)

	for {
		size = int(raw[idx])
		if size == 0 {
			idx += 1
			break
		}

		labels = append(labels, string(raw[idx+1:idx+size+1]))
		idx += size + 1
	}

	q.QNAME = strings.Join(labels, ".")
	q.QTYPE = QType(uint16(raw[idx]<<8) | uint16(raw[idx+1]))
	q.QCLASS = QClass(uint16(raw[idx+2]<<8) | uint16(raw[idx+3]))

	return idx + 4, nil
}

func (q Question) Marshall() ([]byte, error) {
	var (
		buf    = new(bytes.Buffer)
		labels []string
	)

	labels = strings.Split(q.QNAME, ".")
	if len(labels) < 2 {
		return nil, fmt.Errorf(
			"malformed qname %s",
			q.QNAME)
	}

	for _, label := range labels {
		if len(label) == 0 {
			return nil, fmt.Errorf("can't have empty label")
		}

		buf.WriteByte(uint8(len(label)))
		buf.Write([]byte(label))
	}

	buf.WriteByte(0)

	binary.Write(buf, binary.BigEndian, q.QTYPE)
	binary.Write(buf, binary.BigEndian, q.QCLASS)

	return buf.Bytes(), nil
}
