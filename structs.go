package buflice

import (
	"sync"
	"time"
)

// Buflice main struct that contains configs and methods
type Buflice struct {
	maxSize        int
	maxDuration    time.Duration
	mu             sync.Mutex
	currentPos     int
	slice          []interface{}
	ticker         *time.Ticker
	tickerChanDone chan bool
	flushChan      chan []interface{}
}
