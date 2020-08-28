/**
 * SPDX-License-Identifier: Apache-2.0
 * Copyright 2020 vorteil.io Pty Ltd
 */

package strace

import (
	"os"
)

var osExit = os.Exit

func LogFatalWithExitCode(err error, exitCode int) {
	os.Stdout.WriteString(err.Error())
	osExit(exitCode)
}
