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

	// Add a high-resolution metric
	metricLog.PutMetricWithResolution("DetailedLatency", 12.3, emf.UnitMilliseconds, emf.StorageResolutionHigh)

	// Get the JSON representation
	jsonStr := metricLog.String()
	fmt.Println("Generated EMF:")
	fmt.Println(jsonStr)
}
