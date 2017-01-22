package sidekiq

import (
	"encoding/json"
)

type Process struct {
	client *Client

	Name      string      `json:"id"`
	Heartbeat int64       `json:"beat"`
	Busy      int         `json:"busy"`
	Info      ProcessInfo `json:"info"`
}

type ProcessInfo struct {
	Identity    string   `json:"identity"`
	Hostname    string   `json:"hostname"`
	Started     float32  `json:"started_at"`
	ProcessID   int      `json:"pid"`
	Tag         string   `json:"tag"`
	Concurrency int      `json:"concurrency"`
	Queues      []string `json:"queues"`
	Labels      []string `json:"labels"`
}

func (p Process) Jobs() (jobs []Job) {
	return p.client.JobsForProcess(p.Name)
}

func processInfoFromJSON(jsonStr string) (process ProcessInfo, err error) {
	jsonBytes := []byte(jsonStr)
	err = json.Unmarshal(jsonBytes, &process)
	return
}
