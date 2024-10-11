package internal

import (
	"errors"
	"fmt"
	"net"
	"sync"
)

type Client struct {
	nextId uint16
	conn   net.Conn

	mu sync.Mutex
}

type ClientConfig struct {
	Address string
}

func NewClient(cfg ClientConfig) (c Client, err error) {
	if cfg.Address == "" {
		err = errors.New("Address must be specified")
		return
	}

	c.conn, err = net.Dial("udp", cfg.Address)
	if err != nil {
		err = fmt.Errorf(
			"failed to create connection to address %s: %w",
			cfg.Address, err)
		return
	}

	return
}

func (c *Client) LookupAddr(addr string, qType QType) (ips []string, err error) {
	var (
		id      uint16
		payload []byte
	)

	func() {
		c.mu.Lock()
		defer c.mu.Unlock()

		id = c.nextId
		c.nextId += 1
	}()

	queryMsg := &Message{
		Header: Header{
			ID:      id,
			QR:      0,
			Opcode:  OpcodeQuery,
			QDCOUNT: 1,
			RD:      1,
		},
		Questions: []*Question{
			{
				QNAME:  addr,
				QTYPE:  qType,
				QCLASS: QClassIN,
			},
		},
	}

	responseMsg, err := c.GetIps(payload, queryMsg)
	printIps(responseMsg)

	return
}

func printIps(responseMsg *Message) {
	for _, answer := range responseMsg.Answers {
		ip := net.IP(answer.RDATA)
		if len(answer.RDATA) > 4 {
			fmt.Printf("ipv6: %+v", ip.To16().String())
		} else {
			fmt.Printf("ipv4: %+v", ip.To4().String())
		}
		println()
	}
}

func (c *Client) GetIps(payload []byte, queryMsg *Message) (*Message, error) {
	payload, _ = queryMsg.Marshal()
	_, err := c.conn.Write(payload)
	if err != nil {
		err = fmt.Errorf(
			"failed to write query payload %+v: %w",
			queryMsg, err)
		return nil, err
	}

	buf := make([]byte, 1024)
	_, err = c.conn.Read(buf)
	if err != nil {
		err = fmt.Errorf(
			"failed to read from conn")
		return nil, err
	}

	responseMsg := &Message{}
	err = UnmarshallMessage(buf, responseMsg)
	if err != nil {
		err = fmt.Errorf("failed to read message: %w", err)
		return nil, err
	}
	return responseMsg, err
}

func (c *Client) Close() {
	if c.conn != nil {
		c.conn.Close()
	}

	return
}
