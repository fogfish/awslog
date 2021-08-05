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
	"io/ioutil"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
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
executes AWS CloudWatch Log Insight Query against specified Log Group.
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
	q, err := ioutil.ReadFile(queryFile)
	if err != nil {
		return fmt.Errorf("Unable to read query file %s: %w", queryFile, err)
	}

	io, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		return fmt.Errorf("Unable to access AWS: %w", err)
	}

	sec, err := intervalInSeconds(interval)
	if err != nil {
		return fmt.Errorf("Unable to parse time interval: %w", err)
	}

	cwlog := cloudwatchlogs.New(io)

	insight, err := cwlog.StartQuery(
		&cloudwatchlogs.StartQueryInput{
			LogGroupName: aws.String(logGroup),
			QueryString:  aws.String(string(q)),
			StartTime:    aws.Int64(time.Now().Unix() - sec),
			EndTime:      aws.Int64(time.Now().Unix()),
		},
	)
	if err != nil {
		return fmt.Errorf("Unable to build AWS CloudWatch Log Insight query: %w", err)
	}

	result, err := cwlog.GetQueryResults(
		&cloudwatchlogs.GetQueryResultsInput{
			QueryId: insight.QueryId,
		},
	)
	if err != nil {
		return fmt.Errorf("Unable to fetch results of AWS CloudWatch Log Insight query: %w", err)
	}

	for aws.StringValue(result.Status) == "Running" {
		time.Sleep(1 * time.Second)
		result, err = cwlog.GetQueryResults(
			&cloudwatchlogs.GetQueryResultsInput{
				QueryId: insight.QueryId,
			},
		)
		if err != nil {
			return fmt.Errorf("Unable to fetch results of AWS CloudWatch Log Insight query: %w", err)
		}
	}

	if aws.StringValue(result.Status) != "Complete" {
		return fmt.Errorf("Failed (%s) to execute query %+v", *result.Status, *result.Statistics)
	}

	for _, event := range result.Results {
		for _, field := range event {
			if aws.StringValue(field.Field) == "@message" {
				fmt.Print(aws.StringValue(field.Value))
			}
		}
	}

	return nil
}

func intervalInSeconds(t string) (int64, error) {
	v, err := strconv.Atoi(t[0 : len(t)-1])
	if err != nil {
		return 0, err
	}

	switch t[len(t)-1] {
	case 's':
		return int64(v), nil
	case 'm':
		return int64(v * 60), nil
	case 'h':
		return int64(v * 3600), nil
	case 'd':
		return int64(v * 86400), nil
	default:
		return 0, fmt.Errorf("time scale %s is not supported", t)
	}
}
