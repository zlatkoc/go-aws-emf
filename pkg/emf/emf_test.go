package emf

import (
	"encoding/json"
	"testing"
	"time"
)

func TestNewMetricLog(t *testing.T) {
	ml := NewMetricLog("TestNamespace")
	if ml == nil {
		t.Fatal("Expected non-nil MetricLog")
	}
	if ml.emf.Aws.CloudWatchMetrics[0].Namespace != "TestNamespace" {
		t.Errorf("Expected namespace to be TestNamespace, got %s", ml.emf.Aws.CloudWatchMetrics[0].Namespace)
	}
}

func TestPutDimension(t *testing.T) {
	ml := NewMetricLog("TestNamespace")
	ml.PutDimension("Service", "API")

	if val, ok := ml.metrics["Service"]; !ok || val != "API" {
		t.Errorf("Expected dimension Service:API, got %v", ml.metrics["Service"])
	}
}

func TestWithDimensionSet(t *testing.T) {
	ml := NewMetricLog("TestNamespace")
	dims := []string{"Service", "Region"}
	ml.WithDimensionSet(dims)

	if len(ml.emf.Aws.CloudWatchMetrics[0].Dimensions) != 1 {
		t.Fatalf("Expected 1 dimension set, got %d", len(ml.emf.Aws.CloudWatchMetrics[0].Dimensions))
	}

	if len(ml.emf.Aws.CloudWatchMetrics[0].Dimensions[0]) != 2 {
		t.Errorf("Expected dimension set to have 2 dimensions, got %d", len(ml.emf.Aws.CloudWatchMetrics[0].Dimensions[0]))
	}
}

func TestPutMetric(t *testing.T) {
	ml := NewMetricLog("TestNamespace")
	unit := UnitMilliseconds
	ml.PutMetric("Latency", 42.0, unit)

	if val, ok := ml.metrics["Latency"]; !ok || val != 42.0 {
		t.Errorf("Expected metric Latency:42.0, got %v", ml.metrics["Latency"])
	}

	if len(ml.emf.Aws.CloudWatchMetrics[0].Metrics) != 1 {
		t.Fatalf("Expected 1 metric definition, got %d", len(ml.emf.Aws.CloudWatchMetrics[0].Metrics))
	}

	metricDef := ml.emf.Aws.CloudWatchMetrics[0].Metrics[0]
	if metricDef.Name != "Latency" {
		t.Errorf("Expected metric name Latency, got %s", metricDef.Name)
	}

	if metricDef.Unit == nil || *metricDef.Unit != unit {
		t.Errorf("Expected metric unit %s, got %v", unit, metricDef.Unit)
	}
}

func TestBuilder(t *testing.T) {
	ml := NewMetricLog("TestNamespace")
	builder := ml.Builder()

	builder.Dimension("Service", "API").
		Dimension("Region", "us-west-2").
		DimensionSet([]string{"Service"}).
		DimensionSet([]string{"Service", "Region"}).
		Metric("Latency", 42.0, UnitMilliseconds).
		Metric("Count", 1, UnitCount).
		Build()

	if val, ok := ml.metrics["Service"]; !ok || val != "API" {
		t.Errorf("Expected dimension Service:API, got %v", ml.metrics["Service"])
	}

	if val, ok := ml.metrics["Region"]; !ok || val != "us-west-2" {
		t.Errorf("Expected dimension Region:us-west-2, got %v", ml.metrics["Region"])
	}

	if len(ml.emf.Aws.CloudWatchMetrics[0].Dimensions) != 2 {
		t.Fatalf("Expected 2 dimension sets, got %d", len(ml.emf.Aws.CloudWatchMetrics[0].Dimensions))
	}

	if len(ml.emf.Aws.CloudWatchMetrics[0].Metrics) != 2 {
		t.Fatalf("Expected 2 metric definitions, got %d", len(ml.emf.Aws.CloudWatchMetrics[0].Metrics))
	}
}

func TestMarshalJSON(t *testing.T) {
	ml := NewMetricLog("TestNamespace")
	// Set a fixed timestamp for testing
	ml.emf.Aws.Timestamp = 1600000000000

	ml.PutDimension("Service", "API")
	ml.WithDimensionSet([]string{"Service"})
	ml.PutMetric("Latency", 42.0, UnitMilliseconds)

	jsonData, err := ml.MarshalJSON()
	if err != nil {
		t.Fatalf("Error marshaling to JSON: %v", err)
	}

	// Parse the JSON to validate it
	var parsed map[string]interface{}
	if err := json.Unmarshal(jsonData, &parsed); err != nil {
		t.Fatalf("Error parsing JSON: %v", err)
	}

	// Check that we have both the _aws field and the metrics
	if _, ok := parsed["_aws"]; !ok {
		t.Error("Expected _aws field in JSON")
	}

	if val, ok := parsed["Service"]; !ok || val != "API" {
		t.Errorf("Expected Service:API, got %v", val)
	}

	if val, ok := parsed["Latency"]; !ok || val != 42.0 {
		t.Errorf("Expected Latency:42.0, got %v", val)
	}
}

func TestValidation(t *testing.T) {
	tests := []struct {
		name          string
		setup         func() *MetricLog
		expectedError bool
	}{
		{
			name: "valid metric log",
			setup: func() *MetricLog {
				ml := NewMetricLog("TestNamespace")
				ml.PutDimension("Service", "API")
				ml.WithDimensionSet([]string{"Service"})
				ml.PutMetric("Latency", 42.0, UnitMilliseconds)
				return ml
			},
			expectedError: false,
		},
		{
			name: "missing dimension set",
			setup: func() *MetricLog {
				ml := NewMetricLog("TestNamespace")
				ml.PutDimension("Service", "API")
				ml.PutMetric("Latency", 42.0, UnitMilliseconds)
				return ml
			},
			expectedError: true,
		},
		{
			name: "missing metric",
			setup: func() *MetricLog {
				ml := NewMetricLog("TestNamespace")
				ml.PutDimension("Service", "API")
				ml.WithDimensionSet([]string{"Service"})
				return ml
			},
			expectedError: true,
		},
		{
			name: "dimension referenced but not provided",
			setup: func() *MetricLog {
				ml := NewMetricLog("TestNamespace")
				ml.WithDimensionSet([]string{"Service"})
				ml.PutMetric("Latency", 42.0, UnitMilliseconds)
				return ml
			},
			expectedError: true,
		},
		{
			name: "invalid unit",
			setup: func() *MetricLog {
				ml := NewMetricLog("TestNamespace")
				ml.PutDimension("Service", "API")
				ml.WithDimensionSet([]string{"Service"})
				ml.PutMetric("Latency", 42.0, "InvalidUnit")
				return ml
			},
			expectedError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ml := test.setup()
			err := ml.Validate()

			if test.expectedError && err == nil {
				t.Error("Expected validation error, but got nil")
			}

			if !test.expectedError && err != nil {
				t.Errorf("Expected no validation error, but got: %v", err)
			}
		})
	}
}

func TestEndToEnd(t *testing.T) {
	// Create a new metric log
	metricLog := NewMetricLog("ApplicationMetrics")

	// Add dimensions and metrics using builder
	metricLog.Builder().
		Dimension("Service", "PaymentService").
		Dimension("Environment", "Production").
		DimensionSet([]string{"Service"}).
		DimensionSet([]string{"Service", "Environment"}).
		Metric("ProcessingTime", 123.45, UnitMilliseconds).
		Metric("SuccessCount", 1, UnitCount).
		Property("RequestId", "req-123").
		Property("Timestamp", time.Now().String()).
		Build()

	// Validate and marshal to JSON
	jsonData, err := metricLog.MarshalJSON()
	if err != nil {
		t.Fatalf("Error marshaling to JSON: %v", err)
	}

	// Verify we can parse it back
	var parsed map[string]interface{}
	if err := json.Unmarshal(jsonData, &parsed); err != nil {
		t.Fatalf("Error parsing JSON: %v", err)
	}

	// Check expected fields
	awsField, ok := parsed["_aws"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected _aws field to be a map")
	}

	metrics, ok := awsField["CloudWatchMetrics"].([]interface{})
	if !ok || len(metrics) == 0 {
		t.Fatal("Expected CloudWatchMetrics field to be a non-empty array")
	}

	// Verify custom properties
	if val, ok := parsed["RequestId"]; !ok || val != "req-123" {
		t.Errorf("Expected RequestId:req-123, got %v", val)
	}
}
