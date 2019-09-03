package buflice

import (
	"sync"
	"time"
)

// Buflice ...
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
