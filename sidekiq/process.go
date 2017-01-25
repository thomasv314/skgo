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
	Jobs      []Job
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

func NewProcess(name string, hash map[string]string, info ProcessInfo, sk *Client) (process Process) {

	process = Process{
		client:    sk,
		Name:      name,
		Heartbeat: strToInt64(hash["beat"]),
		Busy:      strToInt(hash["busy"]),
		Info:      info,
	}

	process.FetchJobs()

	return
}

func (p *Process) FetchJobs() {
	jobsKey := p.client.key(p.Name + ":workers")
	jobsHash, _ := p.client.Redis.HGetAll(jobsKey).Result()

	jobs := make([]Job, len(jobsHash))

	i := 0

	for _, jobJson := range jobsHash {
		job, _ := jobFromJSON(jobJson)
		jobs[i] = job
		i++
	}

	p.Jobs = jobs
}

func processInfoFromJSON(jsonStr string) (process ProcessInfo, err error) {
	jsonBytes := []byte(jsonStr)
	err = json.Unmarshal(jsonBytes, &process)
	return
}
