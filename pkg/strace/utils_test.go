package strace

import (
	"fmt"
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestLogExit(t *testing.T) {
	// Save current function and restore at the end:
	oldOsExit := osExit
	defer func() { osExit = oldOsExit }()
	expectedCode := 1
	var returnCode int
	myExit := func(code int) {
		returnCode = code
	}

	osExit = myExit
	LogFatalWithExitCode(fmt.Errorf("empty error with exit code %v", expectedCode), expectedCode)
	assert.Equal(t, returnCode, expectedCode, "got unexpected return code")
}
