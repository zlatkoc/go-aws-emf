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
		MetricWithResolution("ApiLatency", 12.3, emf.UnitMilliseconds, emf.StorageResolutionHigh).
		Property("RequestId", "12345").
		Property("Timestamp", time.Now().String()).
		Build()
	
	// Get the JSON representation
	jsonStr := metricLog.String()
	fmt.Println("Generated EMF with Builder:")
	fmt.Println(jsonStr)
}
