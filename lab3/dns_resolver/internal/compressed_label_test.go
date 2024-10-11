package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCompressedLabelMarshallingAndUnmarshalling(t *testing.T) {
	var testCases = []struct {
		desc   string
		entity *CompressedLabel
	}{
		{
			desc: "pointer",
			entity: &CompressedLabel{
				IsPointer: true,
				Offset:    10,
			},
		},
	}

	var (
		msg          []byte
		err          error
		unmarshalled *CompressedLabel
	)

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			msg, err = tc.entity.Marshal()
			require.NoError(t, err)

			unmarshalled = new(CompressedLabel)
			_, err = UnmarshallCompressedLabel(msg, unmarshalled)
			require.NoError(t, err)

			assert.Equal(t, tc.entity.IsPointer, unmarshalled.IsPointer)
			assert.Equal(t, tc.entity.Offset, unmarshalled.Offset)
		})
	}
}
