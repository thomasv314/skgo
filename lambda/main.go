package main

import (
	"encoding/json"
	"github.com/eawsy/aws-lambda-go/service/lambda/runtime"
	"github.com/thomasv314/skgo/sidekiq"
	"os"
)

type LambdaPayload struct {
	Info      sidekiq.Info      `json:"info"`
	Processes []sidekiq.Process `json:"processes"`
}

func handle(evt json.RawMessage, ctx *runtime.Context) (interface{}, error) {
	redisAddr := os.Getenv("SKGO_REDIS_ADDR")
	redisNs := os.Getenv("SKGO_REDIS_NS")

	client := sidekiq.NewClient(redisAddr, redisNs)
	defer client.Close()

	info, err := client.Info()
	if err != nil {
		return "", err
	}

	processes, err := client.Processes()
	if err != nil {
		return "", err
	}

	payload := LambdaPayload{
		Info:      info,
		Processes: processes,
	}

	jsonPayload, err := json.Marshal(&payload)

	return string(jsonPayload), err
}

func init() {
	runtime.HandleFunc(handle)
}

func main() {}
