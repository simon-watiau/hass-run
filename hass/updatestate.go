package hass

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (h *Hass) UpdateState(json string) error {
	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf(
			"%s/api/states/%s",
			h.endpoint,
			h.entity,
		),
		bytes.NewBuffer([]byte(json)),
	)

	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+h.bearer)

	client := &http.Client{}

	response, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}

	defer response.Body.Close()

	if response.StatusCode != 200 && response.StatusCode != 201 {
		body, _ := ioutil.ReadAll(response.Body)
		return fmt.Errorf("invalid status code: %d != [200, 201]: %s", response.StatusCode, body)
	}

	return nil
}
