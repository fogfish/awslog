//
// Copyright (C) 2021 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/awslog
//

package cmd

import (
	"context"

	"github.com/fogfish/awslog/internal/awslog"
	"github.com/spf13/cobra"
)

var (
	streamQuery string
)

func init() {
	rootCmd.AddCommand(streamCmd)
	streamCmd.Flags().StringVarP(&streamQuery, "query", "q", "", "AWS CloudWatch Logs Filter pattern")
}

var streamCmd = &cobra.Command{
	Use:   "stream",
	Short: "stream AWS CloudWatch Log events",
	Long: `
The command streams AWS CloudWatch Log events to console.
It supports filtering messages using filter pattern:
  https://docs.aws.amazon.com/AmazonCloudWatch/latest/logs/FilterAndPatternSyntax.html

`,
	Example: `
awslog stream -g "/aws/lambda/myfun"
awslog stream -g "/aws/lambda/myfun" -q "some pattern"
	`,
	SilenceUsage: true,
	RunE:         stream,
}

func stream(cmd *cobra.Command, args []string) error {
	service := awslog.New(logGroup)
	for evt := range service.Stream(context.Background(), streamQuery) {
		stdout(evt.Timestamp, evt.Message)
	}

	return nil
}
