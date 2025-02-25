package emf

// MetricLogBuilder provides a fluent builder interface for creating EMF metric logs.
type MetricLogBuilder struct {
	metricLog *MetricLog
}

// NewMetricLogBuilder creates a new MetricLogBuilder for the given MetricLog.
func NewMetricLogBuilder(ml *MetricLog) *MetricLogBuilder {
	return &MetricLogBuilder{
		metricLog: ml,
	}
}

// Dimension adds a dimension key-value pair to the log.
func (b *MetricLogBuilder) Dimension(key, value string) *MetricLogBuilder {
	b.metricLog.PutDimension(key, value)
	return b
}

// DimensionSet adds a dimension set to the log.
// Dimension sets define the available dimensions by which metrics can be rolled up.
func (b *MetricLogBuilder) DimensionSet(dimensions []string) *MetricLogBuilder {
	b.metricLog.WithDimensionSet(dimensions)
	return b
}

// Metric adds a metric with the given name, value, and unit to the log.
func (b *MetricLogBuilder) Metric(name string, value interface{}, unit string) *MetricLogBuilder {
	b.metricLog.PutMetric(name, value, unit)
	return b
}

// MetricWithResolution adds a metric with the given name, value, unit, and storage resolution to the log.
// Use StorageResolutionStandard (60) for standard resolution metrics.
// Use StorageResolutionHigh (1) for high-resolution metrics.
func (b *MetricLogBuilder) MetricWithResolution(name string, value interface{}, unit string, resolution int) *MetricLogBuilder {
	b.metricLog.PutMetricWithResolution(name, value, unit, resolution)
	return b
}

// Property adds a custom property to the log.
// Properties are not reported as metrics but appear in the log events.
func (b *MetricLogBuilder) Property(key string, value interface{}) *MetricLogBuilder {
	b.metricLog.metrics[key] = value
	return b
}

// Build returns the built MetricLog.
func (b *MetricLogBuilder) Build() *MetricLog {
	return b.metricLog
}