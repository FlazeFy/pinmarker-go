package e2e

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Positive - Test Case
func TestSuccessPostCreateTrackWithValidInput(t *testing.T) {
	// Test Data
	payload := map[string]interface{}{
		"battery_indicator": 80,
		"track_lat":         "-6.228755",
		"track_long":        "106.820035",
		"track_type":        "live",
		"app_source":        "pinmarker",
		"created_by":        "fcd3f23e-e5aa-11ee-892a-3216422910e9",
	}

	jsonPayload, _ := json.Marshal(payload)

	// Exec
	url := "http://127.0.0.1:9000/api/v1/tracks"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	assert.NoError(t, err)

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Prepare Response Test
	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	assert.NoError(t, err)

	// Template Response
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, "success", result["status"])
	assert.Equal(t, "Track created", result["message"])

	// Data Object
	data, ok := result["data"].(map[string]interface{})
	assert.True(t, ok, "data should be a JSON object")

	// Check Data Fields
	assert.NotEmpty(t, data["id"])
	assert.Equal(t, float64(payload["battery_indicator"].(int)), data["battery_indicator"])
	assert.Equal(t, payload["track_lat"], data["track_lat"])
	assert.Equal(t, payload["track_long"], data["track_long"])
	assert.Equal(t, payload["track_type"], data["track_type"])
	assert.Equal(t, payload["app_source"], data["app_source"])
	assert.Equal(t, payload["created_by"], data["created_by"])

	// Check Data Types
	assert.IsType(t, "", data["id"])
	assert.IsType(t, float64(0), data["battery_indicator"])
	assert.IsType(t, "", data["track_lat"])
	assert.IsType(t, "", data["track_long"])
	assert.IsType(t, "", data["track_type"])
	assert.IsType(t, "", data["app_source"])
	assert.IsType(t, "", data["created_by"])
}

// Negative - Test Case
func TestFailedPostCreateTrackWithInvalidInput(t *testing.T) {
	// Test Data
	payload := map[string]interface{}{
		"battery_indicator": 80,
		"track_lat":         "-6.228755",
		"track_long":        "106.820035",
		"track_type":        "live-loc",
		"app_source":        "pinmarker",
		"created_by":        "fcd3f23e-e5aa-11ee-892a-3216422910e9",
	}

	jsonPayload, _ := json.Marshal(payload)

	// Exec
	url := "http://127.0.0.1:9000/api/v1/tracks"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	assert.NoError(t, err)

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Prepare Response Test
	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	assert.NoError(t, err)

	// Template Response
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "failed", result["status"])
	assert.Equal(t, "track type is not valid", result["message"])
}
