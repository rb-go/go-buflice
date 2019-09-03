package buflice

import (
	"sync"
	"time"
)

// NewBuflice method to initiate Buflice and get it
func NewBuflice(size int, flushDuration time.Duration, notifyChannel chan []interface{}) *Buflice {
	bfl := &Buflice{
		maxSize:        size,
		maxDuration:    flushDuration,
		mu:             sync.Mutex{},
		currentPos:     0,
		slice:          make([]interface{}, 0, size),
		ticker:         time.NewTicker(flushDuration),
		tickerChanDone: make(chan bool),
		flushChan:      notifyChannel,
	}
	go bfl.startTicker()
	return bfl
}
