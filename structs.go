package buflice

import (
	"sync"
	"time"
)

// Buflice main struct that contains configs and methods
type Buflice struct {
	mu            sync.Mutex
	wgProc        sync.WaitGroup
	flushDuration time.Duration
	slice         []interface{}
	ticker        *time.Ticker
	chDone        chan struct{}
	flushChan     chan []interface{}
}
