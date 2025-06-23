package utils

import (
	"encoding/json"
	"fmt"
)

func ConverterStructToMap(data interface{}) (map[string]interface{}, error) {
	var result map[string]interface{}
	j, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal struct: %w", err)
	}
	if err := json.Unmarshal(j, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal into map: %w", err)
	}
	return result, nil
}
