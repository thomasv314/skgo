package sidekiq

import (
	"encoding/json"
)

type Job struct {
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

func (sk SidekiqClient) JobsForProcess(procName string) (jobs []Job) {
	jobsHash, _ := sk.Redis.HGetAll(sk.key(procName + ":workers")).Result()

	jobs = make([]Job, len(jobsHash))

	i := 0

	for _, jobJson := range jobsHash {
		job, _ := jobFromJSON(jobJson)
		jobs[i] = job
		i++
	}

	return
}

func jobFromJSON(jsonStr string) (job Job, err error) {
	jsonBytes := []byte(jsonStr)
	err = json.Unmarshal(jsonBytes, &job)
	return
}
