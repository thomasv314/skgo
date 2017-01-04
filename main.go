package main

import (
	"fmt"
	"github.com/thomasv314/skgo/sidekiq"
	"os"
)

func main() {
	redis_addr := os.Getenv("SK_REDIS_ADDR")
	redis_namespace := os.Getenv("SK_REDIS_NAMESPACE")

	client := sidekiq.NewClient(redis_addr, redis_namespace)
	info, _ := client.Info()
	processes, _ := client.Processes()
	defer client.Close()

	fmt.Println("Listening on", redis_addr, "for namespace", redis_namespace)

	fmt.Println(
		len(processes), "processes running",
		len(info.Queues), "queues",
		info.Retries, "retries",
		info.Failed, "failed",
		info.Processed, "processed",
	)

	for i := range processes {
		process := processes[i]
		jobs := client.JobsForProcess(process.Name)

		fmt.Println("Jobs for", process.Name)
		for j := range jobs {
			job := jobs[j]
			fmt.Println("Job", job.Payload.Id, "running in", job.Queue)
			fmt.Println("Payload", jobs[j].Payload)
		}
	}
}
