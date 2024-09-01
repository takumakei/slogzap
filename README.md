# slogzap

slogzap is a Go package that provides integration between the standard
library's `log/slog` package and Uber's `zap` logging library. It allows you to
use a zap logger as a backend for slog, combining the simplicity of slog's API
with the performance of zap.

## Features

- Converts between `slog.Level` and `zapcore.Level`
- Implements `slog.Handler` interface using a zap logger
- Supports all slog features including groups and attributes

## Installation

To install slogzap, use `go get`:

```
go get github.com/takumakei/slogzap
```

## Usage

Here's a basic example of how to use slogzap:

```go
package main

import (
    "log/slog"

    "github.com/takumakei/slogzap"
    "go.uber.org/zap"
)

func main() {
    // Create a zap logger
    zapLogger, _ := zap.NewProduction()

    // Create a slog handler using the zap logger
    handler := slogzap.New(zapLogger)

    // Create a slog logger with the custom handler
    logger := slog.New(handler)

    // Use the logger
    logger.Info("Hello, world!", "key", "value")
}
```
