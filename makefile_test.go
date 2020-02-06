package main_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/lastrust/issuing-service/certissuer"
	"github.com/stretchr/testify/assert"
)

func TestCmdIssue(t *testing.T) {
	_ = os.Setenv("APP_ENV", "test")

	const testFilepath = "/path/to/test/file/"

	cmd := certissuer.CmdIssue(testFilepath)

	assert.Equal(t, cmd.String(), fmt.Sprintf("/usr/bin/make issue CONF_PATH=%s", testFilepath))

	out, err := cmd.Output()

	assert.Nil(t, err)
	assert.Equal(t, string(out), fmt.Sprintf("/usr/bin/cert-issuer -c %s\n", testFilepath))
}
