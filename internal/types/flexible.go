package types

import (
	"encoding/json"
	"strconv"
)

// FlexibleInt represents a field that can be either a string or an integer in JSON
type FlexibleInt struct {
	Value int64
	IsSet bool
}

// UnmarshalJSON implements json.Unmarshaler for FlexibleInt
func (f *FlexibleInt) UnmarshalJSON(data []byte) error {
	f.IsSet = true

	// Try as integer first
	var intVal int64
	if err := json.Unmarshal(data, &intVal); err == nil {
		f.Value = intVal
		return nil
	}

	// Try as string
	var strVal string
	if err := json.Unmarshal(data, &strVal); err == nil {
		if strVal == "" {
			f.IsSet = false
			f.Value = 0
			return nil
		}

		intVal, err := strconv.ParseInt(strVal, 10, 64)
		if err != nil {
			f.IsSet = false
			f.Value = 0
			return nil
		}

		f.Value = intVal
		return nil
	}

	f.IsSet = false
	f.Value = 0
	return nil
}

// MarshalJSON implements json.Marshaler for FlexibleInt
func (f FlexibleInt) MarshalJSON() ([]byte, error) {
	if !f.IsSet {
		return json.Marshal("")
	}
	return json.Marshal(f.Value)
}

// Int64 returns the integer value, 0 if not set
func (f FlexibleInt) Int64() int64 {
	if !f.IsSet {
		return 0
	}
	return f.Value
}

// String returns the string representation
func (f FlexibleInt) String() string {
	if !f.IsSet {
		return ""
	}
	return strconv.FormatInt(f.Value, 10)
}
