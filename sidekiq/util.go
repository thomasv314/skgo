package sidekiq

import (
	"strconv"
)

func strToInt(str string) (convertedInt int) {
	i, err := strconv.Atoi(str)

	if err != nil {
		return 0
	}

	return i
}

func strToInt64(str string) (convertedInt int64) {
	return int64(strToInt(str))
}
