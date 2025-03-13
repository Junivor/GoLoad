package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type JSON struct {
	Data any `json:"data"`
}

func (j JSON) Value() (driver.Value, error) {
	bytes, err := json.Marshal(j.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %w", err)
	}
	return string(bytes), nil
}

func (j *JSON) Scan(value interface{}) error {
	if value == nil {
		j.Data = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan JSON: incompatible type")
	}

	err := json.Unmarshal(bytes, &j.Data)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return nil
}
