package emf

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"

	"github.com/xeipuuv/gojsonschema"
)

// TestComplianceWithEmfFormat verifies that all generated EMF logs comply with the official EMF schema
func TestComplianceWithEmfFormat(t *testing.T) {
	// Load the schema directly from the file in the repository
	schemaBytes, err := os.ReadFile("emf-format.json")
	if err != nil {
		t.Fatalf("Failed to read the EMF JSON schema: %v", err)
	}

	// Create a schema loader
	schemaLoader := gojsonschema.NewStringLoader(string(schemaBytes))
	schema, err := gojsonschema.NewSchema(schemaLoader)
	if err != nil {
		t.Fatalf("Failed to parse the EMF JSON schema: %v", err)
	}

	// Define test cases that generate various EMF logs
	testCases := []struct {
		name        string
		generateLog func() *MetricLog
	}{
		{
			name: "minimal valid log",
			generateLog: func() *MetricLog {
				ml := NewMetricLog("MinimalTest")
				ml.PutDimension("Service", "Test")
				ml.WithDimensionSet([]string{"Service"})
				ml.PutMetric("Count", 1, UnitCount)
				return ml
			},
		},
		{
			name: "multiple metrics with dimensions",
			generateLog: func() *MetricLog {
				ml := NewMetricLog("ComplexTest")
				ml.PutDimension("Service", "Orders")
				ml.PutDimension("Region", "us-west-2")
				ml.PutDimension("Environment", "Production")
				
				ml.WithDimensionSet([]string{"Service"})
				ml.WithDimensionSet([]string{"Service", "Region"})
				ml.WithDimensionSet([]string{"Service", "Environment"})
				
				ml.PutMetric("ProcessingTime", 123.45, UnitMilliseconds)
				ml.PutMetric("SuccessCount", 100, UnitCount)
				ml.PutMetric("FailureCount", 5, UnitCount)
				ml.PutMetricWithResolution("HighResTime", 10.5, UnitMilliseconds, StorageResolutionHigh)
				
				ml.Builder().
					Property("OrderId", "ord-12345").
					Property("CustomerId", "cust-6789").
					Property("Status", "Completed").
					Build()
				
				return ml
			},
		},
		{
			name: "maximum dimensions",
			generateLog: func() *MetricLog {
				ml := NewMetricLog("MaxDimTest")
				
				// Add the maximum allowed dimensions (30)
				dimensionNames := make([]string, MaxDimensionSetSize)
				for i := 0; i < MaxDimensionSetSize; i++ {
					dimName := "Dim" + string(rune('A'+i))
					dimensionNames[i] = dimName
					ml.PutDimension(dimName, "Value"+string(rune('A'+i)))
				}
				
				ml.WithDimensionSet(dimensionNames)
				ml.PutMetric("TestMetric", 1.0, UnitCount)
				
				return ml
			},
		},
		{
			name: "all unit types",
			generateLog: func() *MetricLog {
				ml := NewMetricLog("AllUnitsTest")
				ml.PutDimension("Service", "API")
				ml.WithDimensionSet([]string{"Service"})
				
				// Add all supported unit types
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
					ml.PutMetric("Metric"+string(rune('A'+i)), float64(i), unit)
				}
				
				return ml
			},
		},
	}

	// Execute each test case
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Generate the EMF log
			ml := tc.generateLog()
			
			// Convert to JSON
			jsonBytes, err := ml.MarshalJSON()
			if err != nil {
				t.Fatalf("Failed to marshal EMF log to JSON: %v", err)
			}
			
			// Print the JSON for debugging
			if testing.Verbose() {
				var pretty bytes.Buffer
				if err := json.Indent(&pretty, jsonBytes, "", "  "); err == nil {
					t.Logf("Generated EMF JSON:\n%s", pretty.String())
				}
			}
			
			// Create a document loader for the generated JSON
			documentLoader := gojsonschema.NewStringLoader(string(jsonBytes))
			
			// Validate against the schema
			result, err := schema.Validate(documentLoader)
			if err != nil {
				t.Fatalf("Schema validation error: %v", err)
			}
			
			// Check if validation passed
			if !result.Valid() {
				var errDetails string
				for i, desc := range result.Errors() {
					if i > 0 {
						errDetails += ", "
					}
					errDetails += desc.String()
				}
				t.Errorf("EMF log doesn't comply with schema: %s", errDetails)
			}
		})
	}
}