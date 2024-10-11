package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResourceMarshallingAndUnmarshalling(t *testing.T) {
	var testCases = []struct {
		desc       string
		entity     *Resource
		shouldFail bool
	}{
		{
			desc: "0-ed case",
			entity: &Resource{
				NAME:     "test.com",
				RDLENGTH: 4,
				RDATA:    []byte{192, 168, 0, 1},
			},
		},
	}

	var (
		msg          []byte
		err          error
		unmarshalled *Resource
	)

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			msg, err = tc.entity.Marshal()

			require.NoError(t, err)

			unmarshalled = new(Resource)
			_, err = UnmarshallResource(msg, unmarshalled)
			require.NoError(t, err)

			assert.Equal(t, tc.entity.RDATA, unmarshalled.RDATA)
			assert.Equal(t, tc.entity.TTL, unmarshalled.TTL)
			assert.Equal(t, tc.entity.RDLENGTH, unmarshalled.RDLENGTH)
			assert.Equal(t, tc.entity.TYPE, unmarshalled.TYPE)

		})
	}
}
