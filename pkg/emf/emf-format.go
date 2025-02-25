// Code generated by github.com/atombender/go-jsonschema, DO NOT EDIT.

package emf
import "encoding/json"
import "fmt"
import "regexp"

type EmfFormatJson struct {
	// Aws corresponds to the JSON schema field "_aws".
	Aws EmfFormatJsonAws `json:"_aws" yaml:"_aws" mapstructure:"_aws"`
}

type EmfFormatJsonAws struct {
	// CloudWatchMetrics corresponds to the JSON schema field "CloudWatchMetrics".
	CloudWatchMetrics []EmfFormatJsonAwsCloudWatchMetricsElem `json:"CloudWatchMetrics" yaml:"CloudWatchMetrics" mapstructure:"CloudWatchMetrics"`

	// Timestamp corresponds to the JSON schema field "Timestamp".
	Timestamp int `json:"Timestamp" yaml:"Timestamp" mapstructure:"Timestamp"`
}

type EmfFormatJsonAwsCloudWatchMetricsElem struct {
	// Dimensions corresponds to the JSON schema field "Dimensions".
	Dimensions [][]string `json:"Dimensions" yaml:"Dimensions" mapstructure:"Dimensions"`

	// Metrics corresponds to the JSON schema field "Metrics".
	Metrics []EmfFormatJsonAwsCloudWatchMetricsElemMetricsElem `json:"Metrics" yaml:"Metrics" mapstructure:"Metrics"`

	// Namespace corresponds to the JSON schema field "Namespace".
	Namespace string `json:"Namespace" yaml:"Namespace" mapstructure:"Namespace"`
}

type EmfFormatJsonAwsCloudWatchMetricsElemMetricsElem struct {
	// Name corresponds to the JSON schema field "Name".
	Name string `json:"Name" yaml:"Name" mapstructure:"Name"`

	// StorageResolution corresponds to the JSON schema field "StorageResolution".
	StorageResolution *int `json:"StorageResolution,omitempty" yaml:"StorageResolution,omitempty" mapstructure:"StorageResolution,omitempty"`

	// Unit corresponds to the JSON schema field "Unit".
	Unit *string `json:"Unit,omitempty" yaml:"Unit,omitempty" mapstructure:"Unit,omitempty"`
}


// UnmarshalJSON implements json.Unmarshaler.
func (j *EmfFormatJsonAwsCloudWatchMetricsElemMetricsElem) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil { return err }
	if _, ok := raw["Name"]; raw != nil && !ok {
		return fmt.Errorf("field Name in EmfFormatJsonAwsCloudWatchMetricsElemMetricsElem: required")
	}
	type Plain EmfFormatJsonAwsCloudWatchMetricsElemMetricsElem
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil { return err }
	if matched, _ := regexp.MatchString("^(.*)$", string(plain.Name)); !matched {
		return fmt.Errorf("field %s pattern match: must match %s", "^(.*)$", "Name")
	}
	if len(plain.Name) < 1 {
		return fmt.Errorf("field %s length: must be >= %d", "Name", 1)
	}
	if len(plain.Name) > 1024 {
		return fmt.Errorf("field %s length: must be <= %d", "Name", 1024)
	}
	if plain.Unit != nil {
		if matched, _ := regexp.MatchString("^(Seconds|Microseconds|Milliseconds|Bytes|Kilobytes|Megabytes|Gigabytes|Terabytes|Bits|Kilobits|Megabits|Gigabits|Terabits|Percent|Count|Bytes\\/Second|Kilobytes\\/Second|Megabytes\\/Second|Gigabytes\\/Second|Terabytes\\/Second|Bits\\/Second|Kilobits\\/Second|Megabits\\/Second|Gigabits\\/Second|Terabits\\/Second|Count\\/Second|None)$", string(*plain.Unit)); !matched {
			return fmt.Errorf("field %s pattern match: must match %s", "^(Seconds|Microseconds|Milliseconds|Bytes|Kilobytes|Megabytes|Gigabytes|Terabytes|Bits|Kilobits|Megabits|Gigabits|Terabits|Percent|Count|Bytes\\/Second|Kilobytes\\/Second|Megabytes\\/Second|Gigabytes\\/Second|Terabytes\\/Second|Bits\\/Second|Kilobits\\/Second|Megabits\\/Second|Gigabits\\/Second|Terabits\\/Second|Count\\/Second|None)$", "Unit")
		}
	}
	*j = EmfFormatJsonAwsCloudWatchMetricsElemMetricsElem(plain)
	return nil
}



// UnmarshalJSON implements json.Unmarshaler.
func (j *EmfFormatJsonAwsCloudWatchMetricsElem) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil { return err }
	if _, ok := raw["Dimensions"]; raw != nil && !ok {
		return fmt.Errorf("field Dimensions in EmfFormatJsonAwsCloudWatchMetricsElem: required")
	}
	if _, ok := raw["Metrics"]; raw != nil && !ok {
		return fmt.Errorf("field Metrics in EmfFormatJsonAwsCloudWatchMetricsElem: required")
	}
	if _, ok := raw["Namespace"]; raw != nil && !ok {
		return fmt.Errorf("field Namespace in EmfFormatJsonAwsCloudWatchMetricsElem: required")
	}
	type Plain EmfFormatJsonAwsCloudWatchMetricsElem
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil { return err }
	if plain.Dimensions != nil && len(plain.Dimensions) < 1 {
		return fmt.Errorf("field %s length: must be >= %d", "Dimensions", 1)
	}
	for i1 := range plain.Dimensions {
		if plain.Dimensions[i1] != nil && len(plain.Dimensions[i1]) < 1 {
			return fmt.Errorf("field %s length: must be >= %d", fmt.Sprintf("Dimensions[%d]", i1), 1)
		}
	}
	if matched, _ := regexp.MatchString("^(.*)$", string(plain.Namespace)); !matched {
		return fmt.Errorf("field %s pattern match: must match %s", "^(.*)$", "Namespace")
	}
	if len(plain.Namespace) < 1 {
		return fmt.Errorf("field %s length: must be >= %d", "Namespace", 1)
	}
	if len(plain.Namespace) > 1024 {
		return fmt.Errorf("field %s length: must be <= %d", "Namespace", 1024)
	}
	*j = EmfFormatJsonAwsCloudWatchMetricsElem(plain)
	return nil
}



// UnmarshalJSON implements json.Unmarshaler.
func (j *EmfFormatJsonAws) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil { return err }
	if _, ok := raw["CloudWatchMetrics"]; raw != nil && !ok {
		return fmt.Errorf("field CloudWatchMetrics in EmfFormatJsonAws: required")
	}
	if _, ok := raw["Timestamp"]; raw != nil && !ok {
		return fmt.Errorf("field Timestamp in EmfFormatJsonAws: required")
	}
	type Plain EmfFormatJsonAws
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil { return err }
	*j = EmfFormatJsonAws(plain)
	return nil
}



// UnmarshalJSON implements json.Unmarshaler.
func (j *EmfFormatJson) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil { return err }
	if _, ok := raw["_aws"]; raw != nil && !ok {
		return fmt.Errorf("field _aws in EmfFormatJson: required")
	}
	type Plain EmfFormatJson
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil { return err }
	*j = EmfFormatJson(plain)
	return nil
}

