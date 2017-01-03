package sidekiq

import (
	"gopkg.in/redis.v5"
)

func NewClient(addr string, redisNs string) (sidekiqClient SidekiqClient) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})

	sidekiqClient = SidekiqClient{
		Redis:     redisClient,
		Namespace: redisNs,
	}

	return SidekiqClient{
		Redis:     redisClient,
		Namespace: redisNs,
	}
}

type SidekiqClient struct {
	Redis     *redis.Client
	Namespace string
}

func (sk SidekiqClient) Close() {
	err := sk.Redis.Close()

	if err != nil {
		panic(err)
	}
}

func (sk SidekiqClient) key(key string) string {
	return sk.Namespace + ":" + key
}

func (sk SidekiqClient) processesKey() string {
	return sk.key("processes")
}

func (sk SidekiqClient) retryKey() string {
	return sk.key("retry")
}

func (sk SidekiqClient) processedKey() string {
	return sk.key("stat:processed")
}

func (sk SidekiqClient) failedKey() string {
	return sk.key("stat:failed")
}

func (sk SidekiqClient) queuesKey() string {
	return sk.key("queues")
}
