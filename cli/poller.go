package cli

import (
	"github.com/thomasv314/skgo/sidekiq"
	"log"
)

type SidekiqPollEvent struct {
	Info      sidekiq.Info
	Processes map[*sidekiq.Process][]sidekiq.Job
}

func pollSidekiq() SidekiqPollEvent {
	info, _ := client.Info()

	processes := buildProcessesJobsMap()

	return SidekiqPollEvent{
		Info:      info,
		Processes: processes,
	}
}

func buildProcessesJobsMap() (processes map[*sidekiq.Process][]sidekiq.Job) {
	procsArr, err := client.Processes()

	if err != nil {
		log.Fatal(err)
	}

	processes = make(map[*sidekiq.Process][]sidekiq.Job, len(procsArr))

	for i := range procsArr {
		proc := procsArr[i]
		processes[&proc] = proc.Jobs()
	}

	return
}
