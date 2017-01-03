package sidekiq

type Info struct {
	Queues    []string
	Processes []string
	Retries   int64
	Failed    int64
	Processed int64
}

func (sk SidekiqClient) Info() (info Info, err error) {
	pipe := sk.Redis.Pipeline()

	skProcessesCmd := pipe.SMembers(sk.processesKey())
	queuesCmd := pipe.SMembers(sk.queuesKey())
	failedCmd := pipe.Get(sk.failedKey())
	retriesCmd := pipe.Get(sk.retryKey())
	processedCmd := pipe.Get(sk.processedKey())

	_, err = pipe.Exec()

	retriesInt := strToInt64(retriesCmd.Val())
	processedInt := strToInt64(processedCmd.Val())
	failedInt := strToInt64(failedCmd.Val())
	queues := queuesCmd.Val()
	processes := skProcessesCmd.Val()

	info = Info{
		Retries:   retriesInt,
		Processed: processedInt,
		Failed:    failedInt,
		Queues:    queues,
		Processes: processes,
	}

	return
}
