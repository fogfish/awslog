//
// Copyright (C) 2021 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/awslog
//

package awslog

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
)

type Event struct {
	Timestamp time.Time
	Message   string
}

type LogService struct {
	logGroup string
	cli      *cloudwatchlogs.Client
}

func New(logGroup string) *LogService {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		panic(err)
	}

	return &LogService{
		logGroup: logGroup,
		cli:      cloudwatchlogs.NewFromConfig(cfg),
	}
}

func (service *LogService) Events(q string, from time.Time) ([]Event, error) {
	seq := []Event{}
	req := &cloudwatchlogs.FilterLogEventsInput{
		LogGroupName:  aws.String(service.logGroup),
		StartTime:     aws.Int64(from.UnixMilli()),
		FilterPattern: aws.String(q),
	}

	for {
		evt, err := service.cli.FilterLogEvents(context.Background(), req)
		if err != nil {
			return nil, err
		}

		if len(evt.Events) != 0 {
			for _, event := range evt.Events {
				seq = append(seq, Event{
					Timestamp: time.UnixMilli(aws.ToInt64(event.Timestamp)),
					Message:   aws.ToString(event.Message),
				})
			}
		}

		if evt.NextToken == nil {
			return seq, nil
		}

		req.NextToken = evt.NextToken
	}
}

func (service *LogService) Stream(ctx context.Context, q string) <-chan Event {
	ch := make(chan Event)

	go func() {
		ts := time.Now()
		for {
			select {
			case <-ctx.Done():
				close(ch)
				return
			case <-time.After(1 * time.Second):
				seq, err := service.Events(q, ts)
				if err != nil {
					close(ch)
					return
				}

				for _, evt := range seq {
					ch <- evt
				}

				if len(seq) != 0 {
					ts = seq[len(seq)-1].Timestamp.Add(1 * time.Millisecond)
				}
			}
		}
	}()

	return ch
}

func (service *LogService) Query(q string, from, to time.Time) ([]Event, error) {
	insight, err := service.cli.StartQuery(
		context.Background(),
		&cloudwatchlogs.StartQueryInput{
			LogGroupName: aws.String(service.logGroup),
			QueryString:  aws.String(string(q)),
			StartTime:    aws.Int64(from.Unix()),
			EndTime:      aws.Int64(to.Unix()),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query Logs Insight: %w", err)
	}

	var result *cloudwatchlogs.GetQueryResultsOutput
	for running := true; running; {
		result, err = service.cli.GetQueryResults(
			context.Background(),
			&cloudwatchlogs.GetQueryResultsInput{
				QueryId: insight.QueryId,
			},
		)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch Log Insight Query results: %w", err)
		}

		switch result.Status {
		case types.QueryStatusRunning:
			time.Sleep(150 * time.Millisecond)
		case types.QueryStatusComplete:
			running = false
		default:
			return nil, fmt.Errorf("failed to executed Log Insight Query (%s) %+v", result.Status, *result.Statistics)
		}
	}

	seq := []Event{}
	for _, event := range result.Results {
		var evt Event
		for _, field := range event {
			switch aws.ToString(field.Field) {
			case "@message":
				evt.Message = aws.ToString(field.Value)
			case "@timestamp":
				evt.Timestamp, err = time.Parse("2006-01-02 15:04:05.000", aws.ToString(field.Value))
				if err != nil {
					return nil, err
				}
			}
		}
		seq = append(seq, evt)
	}

	return seq, nil
}
