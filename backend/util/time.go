package util

import "time"

func UnixMillis() int64 {
	return time.Now().UnixNano() / (1000 * 1000)
}