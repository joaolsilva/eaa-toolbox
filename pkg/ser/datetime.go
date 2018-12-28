package ser

import "time"

type DateTime uint64

const (
	secondsPerDay       = 24 * 3600
	unixOffset    int64 = (1969*365 + 1969/4 - 1969/100 + 1969/400) * secondsPerDay
)

func DateTimeNow() DateTime {
	return DateTime((time.Now().UnixNano() / 100) + unixOffset*1E7)
}
