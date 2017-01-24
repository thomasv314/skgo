package cli

import (
	"github.com/thomasv314/skgo/sidekiq"
	"sort"
)

type SidekiqPollEvent struct {
	Info      sidekiq.Info
	Processes []sidekiq.Process
}

type byName []sidekiq.Job

func (s byName) Len() int           { return len(s) }
func (s byName) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s byName) Less(i, j int) bool { return s[i].Process.Name < s[j].Process.Name }

func (e SidekiqPollEvent) Jobs() (jobs []sidekiq.Job) {
	jobs = make([]sidekiq.Job, 0)

	for _, process := range e.Processes {
		for _, job := range process.Jobs {
			jobs = append(jobs, job)
		}
	}

	sort.Sort(byName(jobs))

	return jobs
}

func pollSidekiq() SidekiqPollEvent {
	info, _ := client.Info()
	processes, _ := client.Processes()

	return SidekiqPollEvent{
		Info:      info,
		Processes: processes,
	}
}
