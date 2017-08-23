package contracts

// NOTE: This file was automatically generated.

// Metric data single measurement.
type DataPoint struct {

	// Name of the metric.
	Name string `json:"name"`

	// Metric type. Single measurement or the aggregated value.
	Kind DataPointType `json:"kind"`

	// Single value for measurement. Sum of individual measurements for the
	// aggregation.
	Value float64 `json:"value"`

	// Metric weight of the aggregated metric. Should not be set for a
	// measurement.
	Count int `json:"count"`

	// Minimum value of the aggregated metric. Should not be set for a
	// measurement.
	Min float64 `json:"min"`

	// Maximum value of the aggregated metric. Should not be set for a
	// measurement.
	Max float64 `json:"max"`

	// Standard deviation of the aggregated metric. Should not be set for a
	// measurement.
	StdDev float64 `json:"stdDev"`
}

// Creates a new DataPoint instance with default values set by the schema.
func NewDataPoint() *DataPoint {
	return &DataPoint{
		Kind: Measurement,
	}
}
