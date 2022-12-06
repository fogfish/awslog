//
// Copyright (C) 2021 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/awslog
//

package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/TylerBrock/colorjson"
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
	colored  bool
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&logGroup, "log-group", "g", "", "AWS CloudWatch LogGroup (e.g. /aws/lambda/myfun)")
	rootCmd.PersistentFlags().BoolVarP(&colored, "color", "c", false, "enable colored output")
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

func stdout(t time.Time, data string) (err error) {
	os.Stdout.Write([]byte(t.String()))
	os.Stdout.Write([]byte(": "))

	if colored {
		var obj map[string]interface{}
		json.Unmarshal([]byte(data), &obj)

		f := colorjson.NewFormatter()
		f.Indent = 2

		encoded, err := f.Marshal(obj)
		if err != nil {
			return err
		}
		os.Stdout.Write(encoded)
	}

	if !colored {
		os.Stdout.Write([]byte(data))
	}

	os.Stdout.Write([]byte("\n"))

	return nil
}
