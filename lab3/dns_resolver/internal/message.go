package internal

import (
	"fmt"
)

// Message represents a DNS message
type Message struct {
	Header
	Questions []*Question
	Answers   []*Resource
}

func (m Message) Marshal() (res []byte, err error) {
	var (
		questionPayload []byte
	)

	res, err = m.Header.Marshall()
	if err != nil {
		err = fmt.Errorf("failed to create header payload %+v: %w", m.Header, err)
		return
	}

	for _, question := range m.Questions {
		questionPayload, err = question.Marshall()
		if err != nil {
			err = fmt.Errorf("failed to marshal question %+v: %w", question, err)
			return
		}

		res = append(res, questionPayload...)
	}

	return
}

func UnmarshallMessage(msg []byte, m *Message) (err error) {
	var (
		header    = &Header{}
		questions []*Question
		resources []*Resource
		i         int = 0
		bytesRead int = 0
		n         int = 0
	)

	n, err = UnmarshallHeader(msg[0:12], header)
	if err != nil {
		err = fmt.Errorf("failed to read header: %w", err)
		return
	}

	bytesRead += n

	//fmt.Printf("%+v\n", header)

	questions = make([]*Question, header.QDCOUNT)
	for i, _ = range questions {
		questions[i] = new(Question)

		n, err = UnmarshallQuestion(msg[bytesRead:], questions[i])
		if err != nil {
			err = fmt.Errorf("failed to read question %d: %w", i, err)
			return
		}

		bytesRead += n
	}

	resources = make([]*Resource, header.ANCOUNT)
	for i, _ = range resources {
		resources[i] = new(Resource)

		n, err = UnmarshallResource(msg[bytesRead:], resources[i])
		if err != nil {
			err = fmt.Errorf("failed to read answer %d: %w", i, err)
			return
		}

		bytesRead += n
	}

	m.Header = *header
	m.Questions = questions
	m.Answers = resources

	return
}
