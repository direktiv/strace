/**
 * SPDX-License-Identifier: Apache-2.0
 * Copyright 2020 vorteil.io Pty Ltd
 */

package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/vorteil/strace/pkg/strace"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "strace [BINARY] [ARGS]...",
	Short: "An strace clone",
	Long: `An strace clone, given a pid this app will try to trace and print the executed it executed system calls.
Current Limitations:
	* Cannot print the names of flags or special values, will only print the raw integer in its place.
	* Cannot display the Errno Names, but can display their associated error message.`,
	DisableFlagParsing: true,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 || args[0] == "--help" {
			cmd.Help()
			os.Exit(1)
		}

		if newTracker, err := strace.NewTracker(args); err != nil {
			log.Fatal(err)
		} else if fin, err := newTracker.Start(); err != nil {
			if fin {
				strace.LogFatalWithExitCode(fmt.Errorf("strace finished on pid=%v with error, %v\n", newTracker.Pid(), err), 0)
			} else {
				strace.LogFatalWithExitCode(fmt.Errorf("strace panicked on pid=%v with error, %v\n", newTracker.Pid(), err), 3)
			}
		}

	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	//rootCmd.Flags().IntVarP(&targetPID, "pid", "p", 0, "pid to target for tracking")
}
