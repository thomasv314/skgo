package sidekiq

import (
	"encoding/json"
)

type SidekiqProcess struct {
	Id        string             `json:"id"`
	Heartbeat int64              `json:"beat"`
	Busy      int                `json:"busy"`
	Info      SidekiqProcessInfo `json:"info"`
}

type SidekiqProcessInfo struct {
	Identity    string   `json:"identity"`
	Hostname    string   `json:"hostname"`
	Started     float32  `json:"started_at"`
	ProcessID   int      `json:"pid"`
	Tag         string   `json:"tag"`
	Concurrency int      `json:"concurrency"`
	Queues      []string `json:"queues"`
	Labels      []string `json:"labels"`
}

func (sk SidekiqClient) Processes() (processes []SidekiqProcess, err error) {
	procNames, err := sk.Redis.SMembers(sk.processesKey()).Result()

	if err != nil {
		return make([]SidekiqProcess, 0), err
	}

	processes = make([]SidekiqProcess, len(procNames))

	for i := range procNames {
		procHash, _ := sk.Redis.HGetAll(sk.key(procNames[i])).Result()
		processInfo, _ := processInfoFromJSON(procHash["info"])
		processes[i] = SidekiqProcess{
			Id:        procNames[i],
			Heartbeat: strToInt64(procHash["beat"]),
			Busy:      strToInt(procHash["busy"]),
			Info:      processInfo,
		}
	}

	return
}

func processInfoFromJSON(jsonStr string) (process SidekiqProcessInfo, err error) {
	jsonBytes := []byte(jsonStr)
	if err = json.Unmarshal(jsonBytes, &process); err != nil {
		return
	}
	return
}
