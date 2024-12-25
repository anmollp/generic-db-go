package utils

import (
	"fmt"
	"runtime"
)

// GetGoroutineID Utility function to get the current goroutine ID.
func GetGoroutineID() int {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	var goroutineID int
	fmt.Sscanf(string(buf[:n]), "goroutine %d ", &goroutineID)
	return goroutineID
}
