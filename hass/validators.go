package hass

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)

func ValidateEntityName(entity string) error {
	validEntity, err := regexp.MatchString(`^[a-zA-Z0-9_]+\.[a-zA-Z0-9_]+$`, entity)

	if err != nil {
		return fmt.Errorf("invalid entity name: %w", err)
	}

	if !validEntity {
		return errors.New("invalid entity name")
	}

	return nil
}

func ValidateHostAndBearer(host string, bearer string) error {
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf(
			"%s/api/",
			host,
		),
		nil,
	)

	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+bearer)

	client := &http.Client{}

	response, err := client.Do(req)

	if err != nil {
		return fmt.Errorf("validation request failed: %w", err)
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		return fmt.Errorf("invalid API validation status code (%d != 200)", response.StatusCode)
	}

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return fmt.Errorf("failed to parse API validation response: %w", err)
	}

	var jsonResp map[string]string

	err = json.Unmarshal(body, &jsonResp)

	if err != nil {
		return fmt.Errorf("failed to unmarshal API validation response: %w (%s)", err, body)
	}

	if val, ok := jsonResp["message"]; !ok || val != "API running." {
		return fmt.Errorf("HomeAssistant API is not running: %s", body)
	}

	return nil
}
