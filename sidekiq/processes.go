package sidekiq

import (
	"encoding/json"
)

type Process struct {
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

func (sk SidekiqClient) ProcessNames() (processNames []string) {
	processNames, err := sk.Redis.SMembers(sk.processesKey()).Result()

	if err != nil {
		return
	}

	return processNames
}

func (sk SidekiqClient) Processes() (processes []Process, err error) {
	procNames := sk.ProcessNames()

	processes = make([]Process, len(procNames))

	for i := range procNames {
		procHash, _ := sk.Redis.HGetAll(sk.key(procNames[i])).Result()
		processInfo, _ := processInfoFromJSON(procHash["info"])
		processes[i] = Process{
			Name:      procNames[i],
			Heartbeat: strToInt64(procHash["beat"]),
			Busy:      strToInt(procHash["busy"]),
			Info:      processInfo,
		}
	}

	return
}

func processInfoFromJSON(jsonStr string) (process ProcessInfo, err error) {
	jsonBytes := []byte(jsonStr)
	err = json.Unmarshal(jsonBytes, &process)
	return
}
