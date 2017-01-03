# skgo - sidekiq api client in go

API Client for monitoring the Sidekiq API.

Includes an [AWS Lambda function](lambda/README.md).

## Go API usage

```go
package main

import (
	"fmt"
	"github.com/thomasv314/skgo/sidekiq"
)

func main() {
	client := sidekiq.NewClient("localhost:6379", "mynamespace")
	info, _ := client.Info()
	processes, _ := client.Processes()
	defer client.Close()

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
```
