<h1 align="center">LogFlake Client Go</h1>

> This repository contains the sources for the client-side components of the LogFlake product suite for applications logs and performance collection for Golang applications.

<h3 align="center">🏠 [LogFlake Website](https://logflake.io) |  🔥 [CloudPhoenix Website](https://cloudphoenix.it)</h3>

## Downloads

|                       Package Name                       |                                          Version                                           |
|:--------------------------------------------------------:|:------------------------------------------------------------------------------------------:|
| [logflake-client-go](https://github.com/CloudPhoenix/logflake-client-go) |     ![GitHub Tag](https://img.shields.io/github/v/tag/cloudphoenix/logflake-client-go)     |

## Usage
Retrieve your _application-key_ from Application Settings in LogFlake UI.

```go
import "github.com/CloudPhoenix/logflake-client-go/logflake"
```

```go
i := logflake.New("application-key")

i.SendLog(logflake.Log{
    Content: "Hello World",
    Level:   logflake.LevelInfo,
})
```
