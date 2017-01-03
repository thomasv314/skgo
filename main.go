package main

import (
	"fmt"
	"github.com/thomasv314/skgo/sidekiq"
)

func main() {
	client := sidekiq.NewClient("localhost:6379", "wl_ops_app_sidekiq")

	info, err := client.Info()

	fmt.Println("info", info, err)
	fmt.Println("\n")

	procs, err := client.Processes()

	for i := range procs {
		proc := procs[i]
		jobs := client.JobsForProcess(proc.Name)
		fmt.Println("Jobs for", proc.Name)
		for j := range jobs {
			job := jobs[j]
			fmt.Println("Job", job.Payload.Id, "running in", job.Queue)
			fmt.Println("Payload", jobs[j].Payload)
			fmt.Println("--")
		}
	}
}
