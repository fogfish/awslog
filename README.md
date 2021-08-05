# awslog

Command line application to query AWS CloudWatch Logs


## Getting Started

```
go get github.com/fogfish/awslog
```

**Note**: go get installs the application to `$GOPATH/bin`. This folder shall be accessible to your user and be part of the `PATH` environment variable. Please see [Golang instructions](https://golang.org/doc/gopath_code.html#GOPATH).


## AWS CloudWatch Log Insight Queries

The utility supports "streaming" of CloudWatch Logs Groups using Insight Queries stored in the files:

```
fields @timestamp, @message
| filter @message like /debug/ and ...
| sort @timestamp desc
| limit 20
```

```bash
awslog latest -g "/aws/lambda/myfun" -q query.insight -t 3h
```

## License

[![See LICENSE](https://img.shields.io/github/license/fogfish/awslog.svg?style=for-the-badge)](LICENSE)
