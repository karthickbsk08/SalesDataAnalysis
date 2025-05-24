package errorcode

import "sync"

var (
	Registry = map[string]*ErrorEntry{}
	Counters = map[string]int{}
	Mu       sync.Mutex
	ErrChan  = make(chan map[string]*ErrorEntry)
)

type ErrorEntry struct {
	Code        string
	Package     string
	Method      string
	Description string
	LineFile    string
}
