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
