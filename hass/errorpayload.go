package hass

import (
	"encoding/json"
	"fmt"
)

type StateOnlyPayload struct {
	State string `json:"state"`
}

func JsonResponseForState(state string) string {
	json, err := json.Marshal(StateOnlyPayload{
		State: state,
	})

	if err != nil {
		panic(fmt.Sprintf("Failed to marshal payload: %s", err))
	}

	return string(json)
}
