package runner

import "time"

type Attributes struct {
	Output    string    `json:"output"`
	Running   bool      `json:"running"`
	ExitCode  int       `json:"exit_code"`
	StartedAt time.Time `json:"started_at"`
	UpdatedAt time.Time `json:"updated_at"`
	EndedAt   time.Time `json:"ended_at"`
	Duration  int       `json:"duration"`
}

type Payload struct {
	State      string     `json:"state"`
	Attributes Attributes `json:"attributes"`
}
