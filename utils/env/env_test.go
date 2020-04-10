package env_test

import (
	"os"
	"testing"

	"github.com/lastrust/template-service/utils/env"
	"github.com/stretchr/testify/assert"
)

func TestGetDefault(t *testing.T) {
	tests := []struct {
		desc     string
		key      string
		value    string
		defaults string
		expected string
	}{
		{
			"test empty environment variable",
			"TEST_EMPTY_VAR",
			"",
			"default_value",
			"default_value",
		},
		{
			"test non empty environment variable",
			"TEST_NON_EMPTY_VAR",
			"1",
			"default_value",
			"1",
		},
		{
			"test empty defaults value",
			"TEST_EMPTY_DEFAULTS_VAR",
			"",
			"",
			"",
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			_ = os.Setenv(test.key, test.value)
			actual := env.GetDefault(test.key, test.defaults)
			assert.Equal(t, test.expected, actual)
		})
	}
}
