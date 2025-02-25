# go-aws-emf

A simple Go library for creating AWS CloudWatch Embedded Metric Format (EMF) formatted logs.

[![Go](https://github.com/zlatkoc/go-aws-emf/actions/workflows/go.yml/badge.svg)](https://github.com/zlatkoc/go-aws-emf/actions/workflows/go.yml)
[![codecov](https://codecov.io/gh/zlatkoc/go-aws-emf/branch/main/graph/badge.svg)](https://codecov.io/gh/zlatkoc/go-aws-emf)
[![Go Reference](https://pkg.go.dev/badge/github.com/zlatkoc/go-aws-emf.svg)](https://pkg.go.dev/github.com/zlatkoc/go-aws-emf)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Overview

The AWS CloudWatch Embedded Metric Format (EMF) is a JSON specification that enables you to embed custom metrics alongside detailed log event data. This library provides a simple way to create EMF-formatted logs in Go with minimal dependencies.

For full specification details, see:
- [Amazon CloudWatch Embedded Metric Format Specification](https://docs.aws.amazon.com/AmazonCloudWatch/latest/monitoring/CloudWatch_Embedded_Metric_Format_Specification.html)

## Installation

```bash
go get github.com/zlatkoc/go-aws-emf/pkg/emf
```

## Project Structure

The library is organized into the following directory structure:

```
go-aws-emf/
├── pkg/emf/       # Core EMF implementation
├── examples/      # Example applications
│   ├── basic/     # Basic usage example
│   └── builder/   # Builder pattern example
└── internal/      # Internal implementation details
```

## Usage

### Basic Usage

```go
package main

import (
    "fmt"
    "github.com/zlatkoc/go-aws-emf/pkg/emf"
)

func main() {
    // Create a new metric log with a namespace
    metricLog := emf.NewMetricLog("MyApplicationMetrics")

    // Add dimensions
    metricLog.PutDimension("ServiceName", "UserService")
    metricLog.PutDimension("Environment", "Production")

    // Define which dimensions should be used for rollup
    metricLog.WithDimensionSet([]string{"ServiceName"})
    metricLog.WithDimensionSet([]string{"ServiceName", "Environment"})

    // Add metrics
    metricLog.PutMetric("Latency", 42.0, emf.UnitMilliseconds)
    metricLog.PutMetric("RequestCount", 1, emf.UnitCount)
    
    // Get the JSON representation
    jsonStr := metricLog.String()
    fmt.Println(jsonStr)
    
    // You would typically log this JSON string to CloudWatch Logs
    // logger.Info(jsonStr)
}
```

### Using the Builder Pattern

```go
package main

import (
    "fmt"
    "time"
    "github.com/zlatkoc/go-aws-emf/pkg/emf"
)

func main() {
    // Create a new metric log with a namespace
    metricLog := emf.NewMetricLog("MyApplicationMetrics")
    
    // Use the builder pattern
    metricLog.Builder().
        Dimension("ServiceName", "UserService").
        Dimension("Environment", "Production").
        DimensionSet([]string{"ServiceName"}).
        DimensionSet([]string{"ServiceName", "Environment"}).
        Metric("Latency", 42.0, emf.UnitMilliseconds).
        Metric("RequestCount", 1, emf.UnitCount).
        Property("RequestId", "12345").
        Property("Timestamp", time.Now().String()).
        Build()
    
    // Get the JSON representation
    jsonStr := metricLog.String()
    fmt.Println(jsonStr)
}
```

### High-Resolution Metrics

You can use high-resolution metrics (1-second resolution) by specifying the storage resolution:

```go
metricLog.PutMetricWithResolution("ApiLatency", 12.3, emf.UnitMilliseconds, emf.StorageResolutionHigh)

// Or with the builder:
metricLog.Builder().
    MetricWithResolution("ApiLatency", 12.3, emf.UnitMilliseconds, emf.StorageResolutionHigh).
    Build()
```

## Available Units

The library provides constants for all supported CloudWatch metric units:

```go
UnitSeconds
UnitMicroseconds
UnitMilliseconds
UnitBytes
UnitKilobytes
UnitMegabytes
UnitGigabytes
UnitTerabytes
UnitBits
UnitKilobits
UnitMegabits
UnitGigabits
UnitTerabits
UnitPercent
UnitCount
UnitBytesPerSecond
UnitKBPerSecond
UnitMBPerSecond
UnitGBPerSecond
UnitTBPerSecond
UnitBitsPerSecond
UnitKbitsPerSecond
UnitMbitsPerSecond
UnitGbitsPerSecond
UnitTbitsPerSecond
UnitCountPerSecond
UnitNone
```

## Features

- Simple API for generating EMF-formatted JSON logs
- Support for multiple dimensions and metrics in a single log
- Fluent builder pattern for creating logs
- Full validation against the AWS EMF schema
- Comprehensive test suite
- Minimal dependencies (uses only standard library)

## Versioning

This project follows [Semantic Versioning](https://semver.org/). For the versions available, see the [tags on this repository](https://github.com/zlatkoc/go-aws-emf/tags).

See [VERSIONING.md](VERSIONING.md) for more details on our versioning policy.

## Credits

This library was primarily developed by Anthropic Claude Sonnet 3.7 v1, which designed the architecture, implemented the code, created the tests, and structured the project. The human author provided requirements, feedback, and guidance.

## License

MIT