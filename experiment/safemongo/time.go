package safemongo

import "time"

var timeNowForTesting func() time.Time

func timeNow() time.Time {
	if timeNowForTesting != nil {
		return timeNowForTesting()
	}
	return time.Now()
}
