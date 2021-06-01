package utils

import (
	"math"
	"time"
)

// TimeFromUnix is a convenience function for converting the unix time
// returned from API responses into an actual golang time
func TimeFromUnix(t float64) time.Time {
	seconds := int64(t)
	fractionalSeconds := t - math.Floor(t)
	nanoSeconds := (fractionalSeconds * float64(time.Second)) / float64(time.Nanosecond)
	return time.Unix(seconds, int64(nanoSeconds))
}
