package core

import (
	"encoding/json"

	"go.yaml.in/yaml/v4"
)

// MaybeStringArray is a custom type that can unmarshal from either a single string or an array of strings.
// It also marshals back to the same format, preserving the original structure.
// origin: https://carlosbecker.com/posts/go-custom-marshaling/
type MaybeStringArray []string

var (
	_ yaml.Unmarshaler = &MaybeStringArray{}
	_ yaml.Marshaler   = MaybeStringArray{}
	_ json.Unmarshaler = &MaybeStringArray{}
	_ json.Marshaler   = MaybeStringArray{}
)

func (a *MaybeStringArray) ToSlice() []string {
	if a == nil {
		return nil
	}
	return []string(*a)
}

func (a *MaybeStringArray) UnmarshalYAML(value *yaml.Node) error {
	var slice []string
	if err := value.Decode(&slice); err == nil {
		*a = slice
		return nil
	}

	var single string
	if err := value.Decode(&single); err != nil {
		return err
	}
	*a = []string{single}
	return nil
}

func (a MaybeStringArray) MarshalYAML() (any, error) {
	switch len(a) {
	case 0:
		return nil, nil
	case 1:
		return a[0], nil
	default:
		return []string(a), nil
	}
}

func (a *MaybeStringArray) UnmarshalJSON(data []byte) error {
	var slice []string
	if err := json.Unmarshal(data, &slice); err == nil {
		*a = slice
		return nil
	}

	var single string
	if err := json.Unmarshal(data, &single); err != nil {
		return err
	}
	*a = []string{single}
	return nil
}

func (a MaybeStringArray) MarshalJSON() ([]byte, error) {
	switch len(a) {
	case 0:
		return json.Marshal(nil)
	case 1:
		return json.Marshal(a[0])
	default:
		return json.Marshal([]string(a))
	}
}
