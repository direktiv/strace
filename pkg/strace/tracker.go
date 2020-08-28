/**
 * SPDX-License-Identifier: Apache-2.0
 * Copyright 2020 vorteil.io Pty Ltd
 */

package strace

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	sec "github.com/seccomp/libseccomp-golang"
)

const (
	MaxArgLen = 50
	MaxArgs   = 6
)

type tracker struct {
	pid         int
	currReg     syscall.PtraceRegs
	currRegArgs [MaxArgs]uint64
}

func (sT *tracker) Pid() int {
	return sT.pid
}

// Strace loggers
var (
	TraceLogInfo func(string, ...interface{}) = func(s string, i ...interface{}) {
		os.Stdout.WriteString(fmt.Sprintf(s, i...))
	}

	TraceLogError func(string, ...interface{})
)

func NewTracker(args []string) (*tracker, error) {
	// Execute Command
	pCmd := exec.Command(args[0], args[1:]...)

	pCmd.SysProcAttr = &syscall.SysProcAttr{
		Ptrace: true,
	}

	err := pCmd.Start()
	if err != nil {
		return nil, fmt.Errorf("could not start binary \"%s\", err=%v\n", args[0], err)
	}

	err = pCmd.Wait()
	if err != nil {
		fmt.Errorf("Wait returned: %v\n", err)
	}

	return &tracker{
		pid:         pCmd.Process.Pid,
		currReg:     syscall.PtraceRegs{},
		currRegArgs: [6]uint64{},
	}, nil
}

// Start: Starts the strace tracking on a given pid
// Limitations:
//	* Cannot print the names of flags or special values, will only print the raw integer in its place.
//	* Cannot display the Errno Names, but can display their associated error message.
func (sT *tracker) Start() (finished bool, err error) {
	exit := true

	// Tracking Loop
	for {
		if exit {
			err := sT.updateRegister()
			if err != nil {
				// STRACE FINISHES HERE
				return true, err
			}

			sysCallName, err := sec.ScmpSyscall(sT.currReg.Orig_rax).GetName()
			if err != nil {
				return false, err
			}

			sInfo := syscallTable[sysCallName]

			var argPrintString string = ""
			for i := range sInfo.Args {
				switch sInfo.Args[i] {
				case ArgEmpty:
					goto AddSysCallArgs
				case ArgInteger:
					argPrintString += fmt.Sprintf("%d, ", sT.currRegArgs[i])
					break
				case ArgFlags:
					fallthrough
				case ArgPointer:
					argPrintString += fmt.Sprintf("0x%x, ", sT.currRegArgs[i])
					break
				case ArgString:
					count, val := readRegString(sT.pid, sT.currRegArgs[i], MaxArgLen)
					argPrintString += fmt.Sprintf("\"%s\"", strings.ReplaceAll(string(val), "\n", `\n`))
					if count == MaxArgLen {
						argPrintString += "..."
					}
					argPrintString += ", "
					break
				}
			}

		AddSysCallArgs:

			// Trim Suffix ", "
			argPrintString = strings.TrimSuffix(argPrintString, ", ")

			// Handle Return Value
			returnVal := int(sT.currReg.Rax)
			var returnValStr string

			if returnVal < 0 {
				errno := syscall.Errno(-returnVal) // Convert negative return to positive errno index
				returnValStr = "-1 " + errno.Error()
			} else {
				returnValStr = fmt.Sprintf("%d", returnVal)
			}

			TraceLogInfo("%s(%s) = %s\n", sysCallName, argPrintString, returnValStr) // Print Syscall
		}

		err := syscall.PtraceSyscall(sT.pid, 0)
		if err != nil {
			return false, err
		}

		_, err = syscall.Wait4(sT.pid, nil, 0, nil)
		if err != nil {
			return false, err
		}

		exit = !exit
	}
}

// readRegString: Read string at a registry, will return read string up to the length set with max param.
// Will return entire string if max is set to 0
func readRegString(pid int, addr uint64, max uint64) (count uint64, val []byte) {
	buf := make([]byte, 4)

	for {
		_, err := syscall.PtracePeekData(pid, uintptr(addr+count), buf)
		if err != nil {
			return count, []byte{}
		}

		for bI := range buf {
			if buf[bI] == 0 || (count == max && max != 0) {
				// DONE - Have read max data or zero value byte was found
				return count, val
			} else {
				val = append(val, buf[bI])
			}

			if max != 0 {
				count++
			}
		}
	}
}

func (sT *tracker) updateRegister() error {
	// Get new Register
	err := syscall.PtraceGetRegs(sT.pid, &sT.currReg)
	if err != nil {
		return fmt.Errorf("could not get register, err=%v", err) //FIXME:
	}

	// Update Register Args
	sT.currRegArgs[0] = sT.currReg.Rdi
	sT.currRegArgs[1] = sT.currReg.Rsi
	sT.currRegArgs[2] = sT.currReg.Rdx
	sT.currRegArgs[3] = sT.currReg.R10
	sT.currRegArgs[4] = sT.currReg.R8
	sT.currRegArgs[5] = sT.currReg.R9

	return nil
}
