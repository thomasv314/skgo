package sidekiq

import (
	"gopkg.in/redis.v5"
)

func NewClient(addr string, redisNs string) (sidekiqClient Client) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})

	sidekiqClient = Client{
		Redis:     redisClient,
		Namespace: redisNs,
	}

	return Client{
		Redis:     redisClient,
		Namespace: redisNs,
	}
}

type Client struct {
	Redis     *redis.Client
	Namespace string
}

func (sk Client) ProcessNames() (processNames []string) {
	processNames, err := sk.Redis.SMembers(sk.processesKey()).Result()

	if err != nil {
		return
	}

	return processNames
}

func (sk Client) Processes() (processes []Process, err error) {
	procNames := sk.ProcessNames()

	processes = make([]Process, len(procNames))

	for i := range procNames {
		procHash, _ := sk.Redis.HGetAll(sk.key(procNames[i])).Result()
		processInfo, _ := processInfoFromJSON(procHash["info"])
		processes[i] = Process{
			client:    &sk,
			Name:      procNames[i],
			Heartbeat: strToInt64(procHash["beat"]),
			Busy:      strToInt(procHash["busy"]),
			Info:      processInfo,
		}
	}

	return
}

func (sk Client) Close() {
	err := sk.Redis.Close()

	if err != nil {
		panic(err)
	}
}

func (sk Client) key(key string) string {
	return sk.Namespace + ":" + key
}

func (sk Client) processesKey() string {
	return sk.key("processes")
}

func (sk Client) retryKey() string {
	return sk.key("retry")
}

func (sk Client) processedKey() string {
	return sk.key("stat:processed")
}

func (sk Client) failedKey() string {
	return sk.key("stat:failed")
}

func (sk Client) queuesKey() string {
	return sk.key("queues")
}
