// Package emf provides a simple wrapper for AWS CloudWatch Embedded Metric Format (EMF).
// The library allows creating EMF-formatted logs that can be sent to CloudWatch.
package emf

import (
	"encoding/json"
	"time"
)

// MetricLog represents an EMF log that can contain metrics and dimensions.
// It's a simplified interface over the raw EMF format.
type MetricLog struct {
	emf     EmfFormatJson
	metrics map[string]interface{}
}

// NewMetricLog creates a new EMF metric log with the given namespace.
func NewMetricLog(namespace string) *MetricLog {
	timestamp := time.Now().UnixMilli()
	
	// Initialize with default values
	emf := EmfFormatJson{
		Aws: EmfFormatJsonAws{
			Timestamp: int(timestamp),
			CloudWatchMetrics: []EmfFormatJsonAwsCloudWatchMetricsElem{
				{
					Namespace:  namespace,
					Dimensions: [][]string{},
					Metrics:    []EmfFormatJsonAwsCloudWatchMetricsElemMetricsElem{},
				},
			},
		},
	}
	
	return &MetricLog{
		emf:     emf,
		metrics: make(map[string]interface{}),
	}
}

// Builder returns a new MetricLogBuilder for this MetricLog.
func (ml *MetricLog) Builder() *MetricLogBuilder {
	return NewMetricLogBuilder(ml)
}

// WithDimensionSet adds a dimension set to the metric log.
func (ml *MetricLog) WithDimensionSet(dimensions []string) *MetricLog {
	metricDirective := &ml.emf.Aws.CloudWatchMetrics[0]
	metricDirective.Dimensions = append(metricDirective.Dimensions, dimensions)
	return ml
}

// PutDimension adds a dimension key-value pair to the log.
func (ml *MetricLog) PutDimension(key, value string) *MetricLog {
	ml.metrics[key] = value
	return ml
}

// PutMetric adds a metric with the given name and value to the log.
func (ml *MetricLog) PutMetric(name string, value interface{}, unit string) *MetricLog {
	ml.metrics[name] = value
	
	unitPtr := &unit
	
	directive := &ml.emf.Aws.CloudWatchMetrics[0]
	metricDef := EmfFormatJsonAwsCloudWatchMetricsElemMetricsElem{
		Name: name,
		Unit: unitPtr,
	}
	
	directive.Metrics = append(directive.Metrics, metricDef)
	return ml
}

// PutMetricWithResolution adds a metric with the given name, value, unit and storage resolution to the log.
func (ml *MetricLog) PutMetricWithResolution(name string, value interface{}, unit string, resolution int) *MetricLog {
	ml.metrics[name] = value
	
	unitPtr := &unit
	resolutionPtr := &resolution
	
	directive := &ml.emf.Aws.CloudWatchMetrics[0]
	metricDef := EmfFormatJsonAwsCloudWatchMetricsElemMetricsElem{
		Name:              name,
		Unit:              unitPtr,
		StorageResolution: resolutionPtr,
	}
	
	directive.Metrics = append(directive.Metrics, metricDef)
	return ml
}

// MarshalJSON implements the json.Marshaler interface.
func (ml *MetricLog) MarshalJSON() ([]byte, error) {
	// First, validate the metric log
	if err := ml.Validate(); err != nil {
		return nil, err
	}
	
	// Create a map that combines both the EMF format and metrics
	combinedMap := make(map[string]interface{})
	
	// Add the _aws field
	combinedMap["_aws"] = ml.emf.Aws
	
	// Add all metrics and dimensions
	for k, v := range ml.metrics {
		combinedMap[k] = v
	}
	
	return json.Marshal(combinedMap)
}

// String returns the JSON string representation of the metric log.
func (ml *MetricLog) String() string {
	bytes, err := ml.MarshalJSON()
	if err != nil {
		return ""
	}
	return string(bytes)
}