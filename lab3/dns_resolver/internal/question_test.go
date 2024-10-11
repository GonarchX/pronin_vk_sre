package internal

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestQuestion_UnmarshallAndMarshall(t *testing.T) {
	var testCases = []struct {
		desc       string
		entity     *Question
		shouldFail bool
	}{
		{
			desc:   "0-ed case",
			entity: &Question{},
		},
		{
			desc: "empty label should fail",
			entity: &Question{
				QNAME: "..",
			},
			shouldFail: true,
		},
		{
			desc: "well formed",
			entity: &Question{
				QNAME:  "test.com",
				QTYPE:  QTypeTXT,
				QCLASS: QClassIN,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			msg, err := tc.entity.Marshall()
			if tc.shouldFail {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			unmarshalled := new(Question)
			_, err = UnmarshallQuestion(msg, unmarshalled)
			require.NoError(t, err)

			assert.Equal(t, tc.entity.QNAME, unmarshalled.QNAME)
			assert.Equal(t, tc.entity.QCLASS, unmarshalled.QCLASS)
			assert.Equal(t, tc.entity.QTYPE, unmarshalled.QTYPE)
		})
	}
}

func TestQuestion_MarshallAndUnmarshall(t *testing.T) {
	var testCases = []struct {
		desc   string
		entity []byte
	}{
		{
			desc:   "Random query from wireshark",
			entity: []byte{0x3, 0x77, 0x77, 0x77, 0x7, 0x79, 0x6f, 0x75, 0x74, 0x75, 0x62, 0x65, 0x3, 0x63, 0x6f, 0x6d, 0x0, 0x0, 0x41, 0x0, 0x1},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			unmarshalled := new(Question)
			_, err := UnmarshallQuestion(tc.entity, unmarshalled)
			require.NoError(t, err)

			msg, err := unmarshalled.Marshall()
			require.NoError(t, err)
			assert.Equal(t, len(tc.entity), len(msg))

			assert.Equal(t, tc.entity, msg)
		})
	}
}
