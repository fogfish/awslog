//
// Copyright (C) 2021 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/awslog
//

package cmd

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/fogfish/awslog/internal/awslog"
	"github.com/spf13/cobra"
)

var (
	queryFile string
	interval  string
)

func init() {
	rootCmd.AddCommand(latestCmd)
	latestCmd.Flags().StringVarP(&queryFile, "query", "q", "", "path to AWS CloudWatch Log Insight Query")
	latestCmd.Flags().StringVarP(&interval, "time", "t", "10m", "time interval either in seconds (s), minutes (m), hours (h) or days (d)")
}

var latestCmd = &cobra.Command{
	Use:   "latest",
	Short: "execute AWS CloudWatch Log Insight Query",
	Long: `
The command executes AWS CloudWatch Log Insight Query against specified Log Group.
See CloudWatch Logs Insights query syntax:
  https://docs.aws.amazon.com/AmazonCloudWatch/latest/logs/CWL_QuerySyntax.html

Example query is following:

fields @timestamp, @message
| filter @message like /debug/ and ...
| sort @timestamp desc
| limit 20
`,
	Example: `
awslog latest -g "/aws/lambda/myfun" -q query.insight -t 3h
	`,
	SilenceUsage: true,
	PreRunE:      requiredLatestCmdKey,
	RunE:         latest,
}

func requiredLatestCmdKey(cmd *cobra.Command, args []string) error {
	if logGroup == "" {
		return errors.New("undefined AWS CloudWatch Log Group identified, use --log-group flag to specify one")
	}

	if queryFile == "" {
		return errors.New("undefined AWS CloudWatch Log Insight Query")
	}

	return nil
}

func latest(cmd *cobra.Command, args []string) error {
	q, err := os.ReadFile(queryFile)
	if err != nil {
		return fmt.Errorf("unable to read query file %s: %w", queryFile, err)
	}

	sec, err := intervalInSeconds(interval)
	if err != nil {
		return fmt.Errorf("unable to parse time interval: %w", err)
	}

	service := awslog.New(logGroup)
	events, err := service.Query(string(q), time.Now().Add(-sec), time.Now())
	if err != nil {
		return err
	}

	for _, evt := range events {
		stdout(evt.Timestamp, evt.Message)
	}

	return nil
}

func intervalInSeconds(t string) (time.Duration, error) {
	v, err := strconv.Atoi(t[0 : len(t)-1])
	if err != nil {
		return 0, err
	}

	switch t[len(t)-1] {
	case 's':
		return time.Duration(v) * time.Second, nil
	case 'm':
		return time.Duration(v) * time.Minute, nil
	case 'h':
		return time.Duration(v) * time.Hour, nil
	case 'd':
		return time.Duration(24*v) * time.Hour, nil
	case 'w':
		return time.Duration(7*24*v) * time.Hour, nil
	default:
		return 0, fmt.Errorf("time scale %s is not supported", t)
	}
}
