//
// Copyright (C) 2021 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/awslog
//

package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// Execute is entry point for cobra cli application
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		e := err.Error()
		fmt.Println(strings.ToUpper(e[:1]) + e[1:])
		os.Exit(1)
	}
}

var (
	logGroup string
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&logGroup, "log-group", "g", "", "AWS CloudWatch LogGroup (e.g. /aws/lambda/myfun)")
}

var rootCmd = &cobra.Command{
	Use:     "awslog",
	Short:   "command line interface to AWS CloudWatch Logs",
	Long:    `command line interface to AWS CloudWatch Logs`,
	Run:     root,
	Version: "v0",
}

func root(cmd *cobra.Command, args []string) {
	cmd.Help()
}
