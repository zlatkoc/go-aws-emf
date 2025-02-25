package emf

import (
	"encoding/json"
	"testing"
)

// TestEdgeCaseValidation tests that the validator correctly identifies edge cases and invalid inputs
func TestEdgeCaseValidation(t *testing.T) {
	tests := []struct {
		name          string
		setup         func() *MetricLog
		expectedError bool
		errorContains string
	}{
		{
			name: "namespace too long",
			setup: func() *MetricLog {
				// Create a namespace longer than the maximum allowed (1024 chars)
				longNamespace := string(make([]byte, MaxNamespaceLength+1))
				for i := range longNamespace {
					longNamespace = longNamespace[:i] + "a" + longNamespace[i+1:]
				}

				ml := NewMetricLog(longNamespace)
				ml.PutDimension("Service", "API")
				ml.WithDimensionSet([]string{"Service"})
				ml.PutMetric("Latency", 42.0, UnitMilliseconds)
				return ml
			},
			expectedError: true,
			errorContains: "namespace length",
		},
		{
			name: "metric name too long",
			setup: func() *MetricLog {
				// Create a metric name longer than the maximum allowed (1024 chars)
				longMetricName := string(make([]byte, MaxMetricNameLength+1))
				for i := range longMetricName {
					longMetricName = longMetricName[:i] + "a" + longMetricName[i+1:]
				}

				ml := NewMetricLog("TestNamespace")
				ml.PutDimension("Service", "API")
				ml.WithDimensionSet([]string{"Service"})
				ml.PutMetric(longMetricName, 42.0, UnitMilliseconds)
				return ml
			},
			expectedError: true,
			errorContains: "metric name",
		},
		{
			name: "dimension name too long",
			setup: func() *MetricLog {
				// Create a dimension name longer than the maximum allowed (250 chars)
				longDimName := string(make([]byte, MaxDimensionNameLength+1))
				for i := range longDimName {
					longDimName = longDimName[:i] + "a" + longDimName[i+1:]
				}

				ml := NewMetricLog("TestNamespace")
				ml.PutDimension(longDimName, "Value")
				ml.WithDimensionSet([]string{longDimName})
				ml.PutMetric("Latency", 42.0, UnitMilliseconds)
				return ml
			},
			expectedError: true,
			errorContains: "dimension name",
		},
		{
			name: "too many dimensions in set",
			setup: func() *MetricLog {
				ml := NewMetricLog("TestNamespace")
				
				// Create more than the maximum allowed dimensions in a set (30)
				dimensionNames := make([]string, MaxDimensionSetSize+1)
				for i := 0; i <= MaxDimensionSetSize; i++ {
					dimName := "Dimension" + string(rune('A'+i))
					dimensionNames[i] = dimName
					ml.PutDimension(dimName, "Value")
				}
				
				ml.WithDimensionSet(dimensionNames)
				ml.PutMetric("Latency", 42.0, UnitMilliseconds)
				return ml
			},
			expectedError: true,
			errorContains: "dimension set",
		},
		{
			name: "referenced dimension not provided",
			setup: func() *MetricLog {
				ml := NewMetricLog("TestNamespace")
				ml.WithDimensionSet([]string{"Service"}) // referenced but not provided
				ml.PutMetric("Latency", 42.0, UnitMilliseconds)
				return ml
			},
			expectedError: true,
			errorContains: "dimension 'Service' is referenced but no value is provided",
		},
		{
			name: "invalid storage resolution",
			setup: func() *MetricLog {
				ml := NewMetricLog("TestNamespace")
				ml.PutDimension("Service", "API")
				ml.WithDimensionSet([]string{"Service"})
				// Use an invalid storage resolution - only 1 or 60 are allowed
				ml.PutMetricWithResolution("ApiLatency", 12.3, UnitMilliseconds, 30)
				return ml
			},
			expectedError: true,
			errorContains: "storage resolution",
		},
		{
			name: "empty dimension set",
			setup: func() *MetricLog {
				ml := NewMetricLog("TestNamespace")
				ml.WithDimensionSet([]string{}) // empty dimension set
				ml.PutMetric("Latency", 42.0, UnitMilliseconds)
				return ml
			},
			expectedError: true,
			errorContains: "dimension",
		},
		{
			name: "no metrics provided",
			setup: func() *MetricLog {
				ml := NewMetricLog("TestNamespace")
				ml.PutDimension("Service", "API")
				ml.WithDimensionSet([]string{"Service"})
				// No metrics added
				return ml
			},
			expectedError: true,
			errorContains: "at least one metric must be defined",
		},
		{
			name: "metric defined but no value",
			setup: func() *MetricLog {
				ml := NewMetricLog("TestNamespace")
				ml.PutDimension("Service", "API")
				ml.WithDimensionSet([]string{"Service"})
				
				// Add a metric definition but not the value
				directive := &ml.emf.Aws.CloudWatchMetrics[0]
				metricDef := EmfFormatJsonAwsCloudWatchMetricsElemMetricsElem{
					Name: "Latency",
					Unit: stringPtr(UnitMilliseconds),
				}
				directive.Metrics = append(directive.Metrics, metricDef)
				
				return ml
			},
			expectedError: true,
			errorContains: "metric 'Latency' is defined but no value is provided",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ml := test.setup()
			err := ml.Validate()
			
			if test.expectedError && err == nil {
				t.Fatal("Expected validation error, but got nil")
			}
			
			if !test.expectedError && err != nil {
				t.Fatalf("Expected no validation error, but got: %v", err)
			}
			
			if test.expectedError && err != nil {
				if test.errorContains != "" && !contains(err.Error(), test.errorContains) {
					t.Errorf("Expected error to contain '%s', but got: %v", test.errorContains, err)
				}
			}
			
			// If a validation error is expected, try marshaling which should also fail
			if test.expectedError {
				_, err = ml.MarshalJSON()
				if err == nil {
					t.Error("Expected MarshalJSON to fail with validation error, but it succeeded")
				}
			}
		})
	}
}

// TestJSONOutput tests specific output formats and edge cases in the JSON output
func TestJSONOutput(t *testing.T) {
	ml := NewMetricLog("OutputTest")
	ml.emf.Aws.Timestamp = 1600000000000 // Set a fixed timestamp for testing
	
	ml.PutDimension("Service", "API")
	ml.WithDimensionSet([]string{"Service"})
	ml.PutMetric("Latency", 42.0, UnitMilliseconds)
	
	// Add some property values of different types
	ml.Builder().
		Property("IntValue", 123).
		Property("FloatValue", 123.456).
		Property("BoolValue", true).
		Property("StringValue", "test").
		Property("NullValue", nil).
		Build()
	
	jsonData, err := ml.MarshalJSON()
	if err != nil {
		t.Fatalf("Error marshaling to JSON: %v", err)
	}
	
	var data map[string]interface{}
	if err := json.Unmarshal(jsonData, &data); err != nil {
		t.Fatalf("Error parsing JSON: %v", err)
	}
	
	// Test property values
	if data["IntValue"] != float64(123) {
		t.Errorf("Expected IntValue to be 123, got %v (%T)", data["IntValue"], data["IntValue"])
	}
	
	if data["FloatValue"] != 123.456 {
		t.Errorf("Expected FloatValue to be 123.456, got %v", data["FloatValue"])
	}
	
	if data["BoolValue"] != true {
		t.Errorf("Expected BoolValue to be true, got %v", data["BoolValue"])
	}
	
	if data["StringValue"] != "test" {
		t.Errorf("Expected StringValue to be 'test', got %v", data["StringValue"])
	}
	
	if _, exists := data["NullValue"]; !exists {
		t.Error("Expected NullValue to exist in JSON")
	}
	
	// Check the _aws structure
	aws, ok := data["_aws"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected _aws to be a map")
	}
	
	if aws["Timestamp"] != float64(1600000000000) {
		t.Errorf("Expected Timestamp to be 1600000000000, got %v", aws["Timestamp"])
	}
	
	metrics, ok := aws["CloudWatchMetrics"].([]interface{})
	if !ok || len(metrics) == 0 {
		t.Fatal("Expected CloudWatchMetrics to be a non-empty array")
	}
	
	metricData, ok := metrics[0].(map[string]interface{})
	if !ok {
		t.Fatal("Expected CloudWatchMetrics[0] to be a map")
	}
	
	if metricData["Namespace"] != "OutputTest" {
		t.Errorf("Expected Namespace to be OutputTest, got %v", metricData["Namespace"])
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// stringPtr returns a pointer to the given string
func stringPtr(s string) *string {
	return &s
}