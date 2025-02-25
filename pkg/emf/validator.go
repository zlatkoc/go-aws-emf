package emf

import (
	"fmt"
	"regexp"
)

// Pre-compiled regular expressions for validation
var (
	unitRegex = regexp.MustCompile(`^(Seconds|Microseconds|Milliseconds|Bytes|Kilobytes|Megabytes|Gigabytes|Terabytes|Bits|Kilobits|Megabits|Gigabits|Terabits|Percent|Count|Bytes\/Second|Kilobytes\/Second|Megabytes\/Second|Gigabytes\/Second|Terabytes\/Second|Bits\/Second|Kilobits\/Second|Megabits\/Second|Gigabits\/Second|Terabits\/Second|Count\/Second|None)$`)
)

// Validate performs validation on the metric log to ensure it conforms to the EMF spec.
func (ml *MetricLog) Validate() error {
	// Validate namespace
	directive := ml.emf.Aws.CloudWatchMetrics[0]
	if len(directive.Namespace) < MinNamespaceLength {
		return fmt.Errorf("namespace length must be at least %d characters", MinNamespaceLength)
	}
	if len(directive.Namespace) > MaxNamespaceLength {
		return fmt.Errorf("namespace length must be at most %d characters", MaxNamespaceLength)
	}

	// Validate metrics and dimensions
	if len(directive.Metrics) == 0 {
		return fmt.Errorf("at least one metric must be defined")
	}

	// Check that the metric dimensions are valid
	if len(directive.Dimensions) < MinDimensions {
		return fmt.Errorf("at least one dimension set must be defined")
	}
	
	// Check that no dimension set is empty
	for i, dimSet := range directive.Dimensions {
		if len(dimSet) == 0 {
			return fmt.Errorf("dimension set %d is empty, must contain at least one dimension", i)
		}
	}

	// Validate metric names
	for _, metric := range directive.Metrics {
		if len(metric.Name) < MinMetricNameLength {
			return fmt.Errorf("metric name '%s' length must be at least %d characters", metric.Name, MinMetricNameLength)
		}
		if len(metric.Name) > MaxMetricNameLength {
			return fmt.Errorf("metric name '%s' length must be at most %d characters", metric.Name, MaxMetricNameLength)
		}

		// Validate that we have a metric value
		if _, exists := ml.metrics[metric.Name]; !exists {
			return fmt.Errorf("metric '%s' is defined but no value is provided", metric.Name)
		}

		// Validate unit if provided
		if metric.Unit != nil {
			if !unitRegex.MatchString(*metric.Unit) {
				return fmt.Errorf("invalid unit '%s' for metric '%s'", *metric.Unit, metric.Name)
			}
		}

		// Validate storage resolution if provided
		if metric.StorageResolution != nil {
			if *metric.StorageResolution != StorageResolutionStandard && *metric.StorageResolution != StorageResolutionHigh {
				return fmt.Errorf("invalid storage resolution for metric '%s'. Must be either %d (standard) or %d (high resolution)", 
					metric.Name, StorageResolutionStandard, StorageResolutionHigh)
			}
		}
	}

	// Validate dimension sets
	for i, dimSet := range directive.Dimensions {
		if len(dimSet) > MaxDimensionSetSize {
			return fmt.Errorf("dimension set %d exceeds maximum size of %d", i, MaxDimensionSetSize)
		}

		// Validate each dimension in the set
		for _, dim := range dimSet {
			if len(dim) > MaxDimensionNameLength {
				return fmt.Errorf("dimension name '%s' exceeds maximum length of %d", dim, MaxDimensionNameLength)
			}

			// Ensure the dimension has a value
			if _, exists := ml.metrics[dim]; !exists {
				return fmt.Errorf("dimension '%s' is referenced but no value is provided", dim)
			}
		}
	}

	return nil
}