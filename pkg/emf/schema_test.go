package emf

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/xeipuuv/gojsonschema"
)

// loadJSONSchema loads the EMF JSON schema from the given file path
func loadJSONSchema(t *testing.T) *gojsonschema.Schema {
	t.Helper()
	
	schemaPath := filepath.Join(".", "emf-format.json")
	schemaBytes, err := os.ReadFile(schemaPath)
	if err != nil {
		t.Fatalf("Failed to read JSON schema file: %v", err)
	}
	
	// Load the JSON schema
	schemaLoader := gojsonschema.NewStringLoader(string(schemaBytes))
	schema, err := gojsonschema.NewSchema(schemaLoader)
	if err != nil {
		t.Fatalf("Failed to load JSON schema: %v", err)
	}
	
	return schema
}

// validateAgainstSchema validates the given JSON against the EMF JSON schema
func validateAgainstSchema(t *testing.T, jsonBytes []byte, schema *gojsonschema.Schema) {
	t.Helper()
	
	// Load the JSON data
	documentLoader := gojsonschema.NewStringLoader(string(jsonBytes))
	
	// Validate against the schema
	result, err := schema.Validate(documentLoader)
	if err != nil {
		t.Fatalf("Error validating JSON: %v", err)
	}
	
	// Check validation result
	if !result.Valid() {
		var errMsg string
		for i, desc := range result.Errors() {
			if i > 0 {
				errMsg += ", "
			}
			errMsg += desc.String()
		}
		t.Fatalf("Schema validation failed: %s", errMsg)
	}
}

// TestJsonSchemaValidation tests that the generated JSON conforms to the EMF schema
func TestJsonSchemaValidation(t *testing.T) {
	schema := loadJSONSchema(t)
	
	tests := []struct {
		name     string
		setup    func() *MetricLog
		expected func(t *testing.T, data map[string]interface{})
	}{
		{
			name: "basic metric with one dimension",
			setup: func() *MetricLog {
				ml := NewMetricLog("TestNamespace")
				ml.PutDimension("Service", "API")
				ml.WithDimensionSet([]string{"Service"})
				ml.PutMetric("Latency", 42.0, UnitMilliseconds)
				return ml
			},
			expected: func(t *testing.T, data map[string]interface{}) {
				if data["Service"] != "API" {
					t.Errorf("Expected Service to be API, got %v", data["Service"])
				}
				if data["Latency"] != 42.0 {
					t.Errorf("Expected Latency to be 42.0, got %v", data["Latency"])
				}
			},
		},
		{
			name: "multiple metrics and dimensions",
			setup: func() *MetricLog {
				ml := NewMetricLog("MultiMetricTest")
				ml.PutDimension("Service", "Payment")
				ml.PutDimension("Region", "us-west-2")
				ml.WithDimensionSet([]string{"Service"})
				ml.WithDimensionSet([]string{"Service", "Region"})
				ml.PutMetric("ProcessingTime", 123.45, UnitMilliseconds)
				ml.PutMetric("SuccessCount", 1, UnitCount)
				ml.PutMetric("FailureCount", 0, UnitCount)
				return ml
			},
			expected: func(t *testing.T, data map[string]interface{}) {
				if data["Service"] != "Payment" {
					t.Errorf("Expected Service to be Payment, got %v", data["Service"])
				}
				if data["Region"] != "us-west-2" {
					t.Errorf("Expected Region to be us-west-2, got %v", data["Region"])
				}
				if data["ProcessingTime"] != 123.45 {
					t.Errorf("Expected ProcessingTime to be 123.45, got %v", data["ProcessingTime"])
				}
				if data["SuccessCount"] != float64(1) {
					t.Errorf("Expected SuccessCount to be 1, got %v", data["SuccessCount"])
				}
				if data["FailureCount"] != float64(0) {
					t.Errorf("Expected FailureCount to be 0, got %v", data["FailureCount"])
				}
			},
		},
		{
			name: "high resolution metrics",
			setup: func() *MetricLog {
				ml := NewMetricLog("HighResTest")
				ml.PutDimension("Service", "API")
				ml.WithDimensionSet([]string{"Service"})
				ml.PutMetricWithResolution("ApiLatency", 12.3, UnitMilliseconds, StorageResolutionHigh)
				return ml
			},
			expected: func(t *testing.T, data map[string]interface{}) {
				if data["Service"] != "API" {
					t.Errorf("Expected Service to be API, got %v", data["Service"])
				}
				if data["ApiLatency"] != 12.3 {
					t.Errorf("Expected ApiLatency to be 12.3, got %v", data["ApiLatency"])
				}
				
				// Verify the storage resolution in the CloudWatchMetrics section
				aws := data["_aws"].(map[string]interface{})
				metrics := aws["CloudWatchMetrics"].([]interface{})[0].(map[string]interface{})
				metricDefs := metrics["Metrics"].([]interface{})
				
				found := false
				for _, m := range metricDefs {
					metric := m.(map[string]interface{})
					if metric["Name"] == "ApiLatency" {
						found = true
						if sr, ok := metric["StorageResolution"]; !ok || sr != float64(StorageResolutionHigh) {
							t.Errorf("Expected StorageResolution to be %d, got %v", StorageResolutionHigh, sr)
						}
					}
				}
				
				if !found {
					t.Error("Expected to find ApiLatency metric definition")
				}
			},
		},
		{
			name: "all metric unit types",
			setup: func() *MetricLog {
				ml := NewMetricLog("AllUnitsTest")
				ml.PutDimension("Service", "API")
				ml.WithDimensionSet([]string{"Service"})
				
				units := []string{
					UnitSeconds, UnitMicroseconds, UnitMilliseconds,
					UnitBytes, UnitKilobytes, UnitMegabytes, UnitGigabytes, UnitTerabytes,
					UnitBits, UnitKilobits, UnitMegabits, UnitGigabits, UnitTerabits,
					UnitPercent, UnitCount,
					UnitBytesPerSecond, UnitKBPerSecond, UnitMBPerSecond, UnitGBPerSecond, UnitTBPerSecond,
					UnitBitsPerSecond, UnitKbitsPerSecond, UnitMbitsPerSecond, UnitGbitsPerSecond, UnitTbitsPerSecond,
					UnitCountPerSecond, UnitNone,
				}
				
				for i, unit := range units {
					ml.PutMetric(fmt.Sprintf("Metric%d", i), float64(i), unit)
				}
				
				return ml
			},
			expected: func(t *testing.T, data map[string]interface{}) {
				aws := data["_aws"].(map[string]interface{})
				metrics := aws["CloudWatchMetrics"].([]interface{})[0].(map[string]interface{})
				metricDefs := metrics["Metrics"].([]interface{})
				
				if len(metricDefs) != 27 {
					t.Errorf("Expected 27 metric definitions, got %d", len(metricDefs))
				}
			},
		},
		{
			name: "builder pattern test",
			setup: func() *MetricLog {
				return NewMetricLog("BuilderTest").Builder().
					Dimension("Service", "Database").
					Dimension("Environment", "Production").
					DimensionSet([]string{"Service"}).
					DimensionSet([]string{"Service", "Environment"}).
					Metric("QueryTime", 15.5, UnitMilliseconds).
					Metric("Throughput", 1000, UnitCountPerSecond).
					MetricWithResolution("DetailedQueryTime", 15.5, UnitMilliseconds, StorageResolutionHigh).
					Property("QueryId", "select-123").
					Build()
			},
			expected: func(t *testing.T, data map[string]interface{}) {
				if data["Service"] != "Database" {
					t.Errorf("Expected Service to be Database, got %v", data["Service"])
				}
				if data["Environment"] != "Production" {
					t.Errorf("Expected Environment to be Production, got %v", data["Environment"])
				}
				if data["QueryTime"] != 15.5 {
					t.Errorf("Expected QueryTime to be 15.5, got %v", data["QueryTime"])
				}
				if data["Throughput"] != float64(1000) {
					t.Errorf("Expected Throughput to be 1000, got %v", data["Throughput"])
				}
				if data["DetailedQueryTime"] != 15.5 {
					t.Errorf("Expected DetailedQueryTime to be 15.5, got %v", data["DetailedQueryTime"])
				}
				if data["QueryId"] != "select-123" {
					t.Errorf("Expected QueryId to be select-123, got %v", data["QueryId"])
				}
			},
		},
	}
	
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ml := test.setup()
			
			// Marshal to JSON
			jsonBytes, err := ml.MarshalJSON()
			if err != nil {
				t.Fatalf("Error marshaling to JSON: %v", err)
			}
			
			// Validate against schema
			validateAgainstSchema(t, jsonBytes, schema)
			
			// Also validate the data structure
			var data map[string]interface{}
			if err := json.Unmarshal(jsonBytes, &data); err != nil {
				t.Fatalf("Error parsing JSON: %v", err)
			}
			
			// Run expected validations
			test.expected(t, data)
		})
	}
}

// TestSchemaValidationWithComplexData tests complex scenarios and edge cases
func TestSchemaValidationWithComplexData(t *testing.T) {
	schema := loadJSONSchema(t)
	
	tests := []struct {
		name     string
		setup    func() *MetricLog
		expected func(t *testing.T, data map[string]interface{})
	}{
		{
			name: "maximum dimensions",
			setup: func() *MetricLog {
				ml := NewMetricLog("MaxDimensionsTest")
				
				// Add the maximum allowed dimensions (30)
				dimensionNames := make([]string, MaxDimensionSetSize)
				for i := 0; i < MaxDimensionSetSize; i++ {
					dimName := fmt.Sprintf("Dimension%d", i)
					dimensionNames[i] = dimName
					ml.PutDimension(dimName, fmt.Sprintf("Value%d", i))
				}
				
				ml.WithDimensionSet(dimensionNames)
				ml.PutMetric("TestMetric", 1.0, UnitCount)
				
				return ml
			},
			expected: func(t *testing.T, data map[string]interface{}) {
				aws := data["_aws"].(map[string]interface{})
				metrics := aws["CloudWatchMetrics"].([]interface{})[0].(map[string]interface{})
				dimensions := metrics["Dimensions"].([]interface{})[0].([]interface{})
				
				if len(dimensions) != MaxDimensionSetSize {
					t.Errorf("Expected %d dimensions, got %d", MaxDimensionSetSize, len(dimensions))
				}
			},
		},
		{
			name: "multiple dimension sets",
			setup: func() *MetricLog {
				ml := NewMetricLog("MultiDimSetsTest")
				
				// Add dimensions
				ml.PutDimension("Service", "API")
				ml.PutDimension("Region", "us-west-2")
				ml.PutDimension("Environment", "Production")
				ml.PutDimension("Host", "server-123")
				
				// Add multiple dimension sets for different rollup views
				ml.WithDimensionSet([]string{"Service"})
				ml.WithDimensionSet([]string{"Service", "Region"})
				ml.WithDimensionSet([]string{"Service", "Environment"})
				ml.WithDimensionSet([]string{"Service", "Host"})
				ml.WithDimensionSet([]string{"Service", "Region", "Environment"})
				
				ml.PutMetric("Latency", 42.0, UnitMilliseconds)
				
				return ml
			},
			expected: func(t *testing.T, data map[string]interface{}) {
				aws := data["_aws"].(map[string]interface{})
				metrics := aws["CloudWatchMetrics"].([]interface{})[0].(map[string]interface{})
				dimensions := metrics["Dimensions"].([]interface{})
				
				if len(dimensions) != 5 {
					t.Errorf("Expected 5 dimension sets, got %d", len(dimensions))
				}
			},
		},
		{
			name: "integer and float metrics",
			setup: func() *MetricLog {
				ml := NewMetricLog("MixedTypesTest")
				ml.PutDimension("Service", "API")
				ml.WithDimensionSet([]string{"Service"})
				
				// Add various metric types
				ml.PutMetric("IntMetric", 42, UnitCount)
				ml.PutMetric("FloatMetric", 42.5, UnitMilliseconds)
				ml.PutMetric("ZeroMetric", 0, UnitCount)
				ml.PutMetric("LargeMetric", 1000000, UnitBytes)
				ml.PutMetric("SmallMetric", 0.0001, UnitSeconds)
				
				return ml
			},
			expected: func(t *testing.T, data map[string]interface{}) {
				if data["IntMetric"] != float64(42) {
					t.Errorf("Expected IntMetric to be 42, got %v", data["IntMetric"])
				}
				if data["FloatMetric"] != 42.5 {
					t.Errorf("Expected FloatMetric to be 42.5, got %v", data["FloatMetric"])
				}
				if data["ZeroMetric"] != float64(0) {
					t.Errorf("Expected ZeroMetric to be 0, got %v", data["ZeroMetric"])
				}
				if data["LargeMetric"] != float64(1000000) {
					t.Errorf("Expected LargeMetric to be 1000000, got %v", data["LargeMetric"])
				}
				if data["SmallMetric"] != 0.0001 {
					t.Errorf("Expected SmallMetric to be 0.0001, got %v", data["SmallMetric"])
				}
			},
		},
		{
			name: "custom properties",
			setup: func() *MetricLog {
				ml := NewMetricLog("PropertiesTest")
				ml.PutDimension("Service", "API")
				ml.WithDimensionSet([]string{"Service"})
				ml.PutMetric("Latency", 42.0, UnitMilliseconds)
				
				// Add custom properties using the builder
				ml.Builder().
					Property("RequestId", "req-123").
					Property("UserId", "user-456").
					Property("CorrelationId", "corr-789").
					Property("BoolProperty", true).
					Property("NullProperty", nil).
					Property("NumberProperty", 12345).
					Build()
				
				return ml
			},
			expected: func(t *testing.T, data map[string]interface{}) {
				if data["RequestId"] != "req-123" {
					t.Errorf("Expected RequestId to be req-123, got %v", data["RequestId"])
				}
				if data["UserId"] != "user-456" {
					t.Errorf("Expected UserId to be user-456, got %v", data["UserId"])
				}
				if data["BoolProperty"] != true {
					t.Errorf("Expected BoolProperty to be true, got %v", data["BoolProperty"])
				}
				if data["NumberProperty"] != float64(12345) {
					t.Errorf("Expected NumberProperty to be 12345, got %v", data["NumberProperty"])
				}
			},
		},
	}
	
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ml := test.setup()
			
			// Marshal to JSON
			jsonBytes, err := ml.MarshalJSON()
			if err != nil {
				t.Fatalf("Error marshaling to JSON: %v", err)
			}
			
			// Validate against schema
			validateAgainstSchema(t, jsonBytes, schema)
			
			// Also validate the data structure
			var data map[string]interface{}
			if err := json.Unmarshal(jsonBytes, &data); err != nil {
				t.Fatalf("Error parsing JSON: %v", err)
			}
			
			// Run expected validations
			test.expected(t, data)
		})
	}
}