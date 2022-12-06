# awslog

Command line application to query AWS CloudWatch Logs


## Getting Started

```
go install github.com/fogfish/awslog@latest
```

**Note**: go get installs the application to `$GOPATH/bin`. This folder shall be accessible to your user and be part of the `PATH` environment variable. Please see [Golang instructions](https://golang.org/doc/gopath_code.html#GOPATH).


## Stream Log events

```bash
awslog stream -g "/aws/lambda/myfun" -q "some pattern"
```

optionally use [filter and pattern syntax](https://docs.aws.amazon.com/AmazonCloudWatch/latest/logs/FilterAndPatternSyntax.html) to control visibility of events.


## Query Log events

```bash
awslog latest -g "/aws/lambda/myfun" -q query.insight -t 3h
```

it uses Logs Insight Queries to collect events from the logs. See [query syntax](https://docs.aws.amazon.com/AmazonCloudWatch/latest/logs/CWL_QuerySyntax.html) for details

```
fields @timestamp, @message
| filter @message like /debug/ and ...
| sort @timestamp desc
| limit 20
```

## License

[![See LICENSE](https://img.shields.io/github/license/fogfish/awslog.svg?style=for-the-badge)](LICENSE)
