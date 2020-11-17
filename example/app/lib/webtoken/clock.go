package webtoken

import "time"

// IClock represents a system clock.
type IClock interface {
	Now() time.Time
}

// Clock is a standard system clock.
type Clock struct{}

// Now returns the current time.
func (c *Clock) Now() time.Time {
	return time.Now()
}
