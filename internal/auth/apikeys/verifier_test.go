package apikeys

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVerifyApiKey(t *testing.T) {
	tc := []struct {
		apiKey        string
		hash          string
		expectedMatch bool
		expectedError bool
	}{
		{
			apiKey:        "65RIGoN3TM-AhqowJAPIoWvI1IArq4WPgDvyKFHPJZ-3rr5ZhLafxUibfXCK774RRhGVO-1VYYG2cWlPmieyVA",
			hash:          "foo",
			expectedMatch: false,
			expectedError: true,
		},
		{
			apiKey:        "65RIGoN3TM-AhqowJAPIoWvI1IArq4WPgDvyKFHPJZ-3rr5ZhLafxUibfXCK774RRhGVO-1VYYG2cWlPmieyVA",
			hash:          "$argon2id$v=19$m=65536,t=1,p=2$V+VI24cKNaEDrXdz0xI3Lg$epL8hNnvWkNiK1BPnqRrLqoZk/KvAM1HHK1HrtxMwyw",
			expectedMatch: false,
			expectedError: false,
		},
		{
			apiKey:        "Cp9MyxL2YQM6EygSOwkDaB8-avi_sL2OpqxrKamvgmhKidPiqESpWVb6FDTXZlpOgii0c9TEMrNk0jqbn0rQyw",
			hash:          "$argon2id$v=19$m=65536,t=1,p=2$V+VI24cKNaEDrXdz0xI3Lg$epL8hNnvWkNiK1BPnqRrLqoZk/KvAM1HHK1HrtxMwyw",
			expectedMatch: true,
			expectedError: false,
		},
	}

	for _, tc := range tc {
		result, err := CompareArgon2Hash(tc.apiKey, tc.hash)
		if tc.expectedError {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
			assert.Equal(t, tc.expectedMatch, result)
		}
	}
}

func TestDecodeArgon2Hash(t *testing.T) {
	tc := []struct {
		encoded          string
		memExpected      uint32
		timeExpected     uint32
		parallelExpected uint8
		saltExpected     string
		hashExpected     string
		errorExpected    bool
	}{
		{
			errorExpected: true,
		},
		{
			encoded:          "$argon2id$v=19$m=65536,t=1,p=2$V+VI24cKNaEDrXdz0xI3Lg$epL8hNnvWkNiK1BPnqRrLqoZk/KvAM1HHK1HrtxMwyw",
			memExpected:      65536,
			timeExpected:     1,
			parallelExpected: 2,
			saltExpected:     "V+VI24cKNaEDrXdz0xI3Lg",
			hashExpected:     "epL8hNnvWkNiK1BPnqRrLqoZk/KvAM1HHK1HrtxMwyw",
			errorExpected:    false,
		},
	}

	for _, tc := range tc {
		mem, time, p, salt, hash, err := DecodeArgon2Hash(tc.encoded)

		if tc.errorExpected {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)

			expectedHash, err := base64.RawStdEncoding.Strict().DecodeString(tc.hashExpected)
			assert.Nil(t, err)

			expectedSalt, err := base64.RawStdEncoding.Strict().DecodeString(tc.saltExpected)
			assert.Nil(t, err)

			assert.Equal(t, tc.memExpected, mem)
			assert.Equal(t, tc.timeExpected, time)
			assert.Equal(t, tc.parallelExpected, p)
			assert.Equal(t, expectedSalt, salt)
			assert.Equal(t, expectedHash, hash)
		}
	}
}
