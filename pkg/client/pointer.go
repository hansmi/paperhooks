package client

import "time"

// Bool allocates a new bool value to store v and returns a pointer to it.
func Bool(v bool) *bool {
	return &v
}

// Int allocates a new int value to store v and returns a pointer to it.
func Int(v int) *int {
	return &v
}

// Int64 allocates a new int64 value to store v and returns a pointer to it.
func Int64(v int64) *int64 {
	return &v
}

// String allocates a new string value to store v and returns a pointer to it.
func String(v string) *string {
	return &v
}

// Time allocates a new time.Time value to store v and returns a pointer to it.
func Time(v time.Time) *time.Time {
	return &v
}
