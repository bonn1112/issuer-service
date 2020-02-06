package utils_test

import (
	"testing"

	"github.com/lastrust/issuing-service/utils"
	"github.com/stretchr/testify/assert"
)

func TestFileExists(t *testing.T) {
	assert.True(t, utils.FileExists("../test/file"))
	assert.False(t, utils.FileExists("../test/undefined_file"))
}
