package strace

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testCommand = []string{"echo", "hello world"}

func TestTracker(t *testing.T) {
	newTracker, err := NewTracker(testCommand)
	assert.NoError(t, err, fmt.Sprintf("error creating new tracker"))

	lgAssert := func(format string, values ...interface{}) {
		assert.NotEmpty(t, fmt.Sprintf(format, values...), "captured Sys call was empty")
		t.Logf(format, values...) // In case user wants to see strace output
	}

	TraceLogInfo = lgAssert

	fin, err := newTracker.Start()
	assert.Error(t, err, fmt.Sprintf("error starting tracker"))
	assert.True(t, fin, "tracker did not finish gracefully")

	assert.NotZero(t, newTracker.Pid(), "pid is somehow zero")
}

func TestTrackerMissingExec(t *testing.T) {
	_, err := NewTracker([]string{"path-/to-/fake-/file"})
	assert.Error(t, err, fmt.Sprintf("somehow created invalid tracker"))
}
