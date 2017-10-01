package timedstore

import (
	"math"
	"math/big"

	"github.com/satori/go.uuid"
)

// Value is the struct used to internally represent Value objects, which are
// only active during a given time range. The start time and end time are both
// provided in seconds since January 1, 1970, also known as unix time or epoch
// time.
type Value struct {
	uniqueID string

	StartSeconds int64
	EndSeconds   int64

	Data interface{}
}

// NewValue creates a new Value given a start time and end time. Times are provided in
// seconds since January 1, 1970 UTC (unix/epoch time)
func NewValue(startSeconds int64, endSeconds int64, data interface{}) Value {
	uniqueID := uuid.NewV4().String()
	return Value{uniqueID, startSeconds, endSeconds, data}
}

// NewEternalValue creates a new Value with the minimum start time and maximum
// end time. Except for the case where time t = math.MaxInt64, an eternal value
// will never expire and will always be active.
func NewEternalValue(data interface{}) Value {
	return NewValue(math.MinInt64, math.MaxInt64, data)
}

// NewValueFromDuration creates a new Value given a start time and a duration.
// Start time is provided in seconds since January 1, 1970 UTC (unix/epoch time).
// Duration is also in seconds.
func NewValueFromDuration(startSeconds int64, duration int64, data interface{}) Value {
	endSeconds := startSeconds + duration
	return NewValue(startSeconds, endSeconds, data)
}

// Duration returns the duration of this value in seconds. Calculated by
// subtracting the start time from the end time.
func (v Value) Duration() uint64 {
	start := big.NewInt(v.StartSeconds)
	end := big.NewInt(v.EndSeconds)

	duration := (&big.Int{}).Sub(end, start)

	return duration.Uint64()
}

// IsActiveForTime determines if this value is active during the provided time.
// It is inclusive of the start time and exclusive of the end time.
func (v Value) IsActiveForTime(time int64) bool {
	return time >= v.StartSeconds && time < v.EndSeconds
}

// IsExpiredForTime determines if this value's active period is later than the
// provided time. Can be used to determine if a value's useful time period has
// expired.
func (v Value) IsExpiredForTime(time int64) bool {
	return time >= v.EndSeconds
}

// GetUniqueID returns the unique id supplied at the object's creation time.
func (v Value) GetUniqueID() string {
	return v.uniqueID
}
