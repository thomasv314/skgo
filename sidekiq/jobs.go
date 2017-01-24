package sidekiq

import (
	"encoding/json"
)

type Job struct {
	Process *Process   `json:omit`
	Queue   string     `json:"queue"`
	Payload JobPayload `json:"payload"`
	RunAt   int64      `json:"run_at"`
}

type JobPayload struct {
	Class      string   `json:"class"`
	Arguments  []string `json:"args"`
	Retry      bool     `json:"retry"`
	Queue      string   `json:"queue"`
	Id         string   `json:"jid"`
	CreatedAt  float32  `json:"created_at"`
	Locale     string   `json:"en"`
	EnqueuedAt string   `json:"enqueued_at"`
}

func jobFromJSON(jsonStr string) (job Job, err error) {
	jsonBytes := []byte(jsonStr)
	err = json.Unmarshal(jsonBytes, &job)
	return
}
