package sqlbuilder

import (
	"database/sql/driver"
	"encoding/json"
)

type MapAny map[string]any

func (m *MapAny) Scan(value interface{}) error {
	bytesValue, _ := value.([]byte)
	return json.Unmarshal(bytesValue, m)
}

func (m MapAny) Value() (driver.Value, error) {
	return json.Marshal(m)
}

type MapString map[string]string

func (m *MapString) Scan(value interface{}) error {
	bytesValue, _ := value.([]byte)
	return json.Unmarshal(bytesValue, m)
}

func (m MapString) Value() (driver.Value, error) {
	return json.Marshal(m)
}

type StringSlice []string

func (s *StringSlice) Scan(value interface{}) error {
	bytesValue, _ := value.([]byte)
	return json.Unmarshal(bytesValue, s)
}

func (s StringSlice) Value() (driver.Value, error) {
	value, err := json.Marshal(s)
	return value, err
}
