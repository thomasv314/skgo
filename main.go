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

	fmt.Println("procs", procs)
	fmt.Println("\n\nerr", err)
}
