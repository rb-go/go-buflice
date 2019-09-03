package buflice

import (
	"sync"
	"time"
)

// NewBuflice ...
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

func (bs *Buflice) startTicker() {
	for {
		select {
		case <-bs.tickerChanDone:
			return
		case <-bs.ticker.C:
			bs.flushResetPosAndCleanSlice()
		}
	}
}

func (bs *Buflice) stopTicker() {
	bs.ticker.Stop()
	bs.tickerChanDone <- true
}

func (bs *Buflice) flushResetPosAndCleanSlice() {
	bs.flushChan <- bs.slice
	bs.currentPos = 0
	bs.slice = make([]interface{}, 0, bs.maxSize)
}

// Add is for adding elements
func (bs *Buflice) Add(element interface{}) {
	bs.mu.Lock()
	defer bs.mu.Unlock()

	if bs.currentPos != bs.maxSize {
		bs.slice = append(bs.slice, element)
		bs.currentPos = bs.currentPos + 1
	}

	if bs.currentPos == bs.maxSize {
		bs.flushResetPosAndCleanSlice()
	}
}

// Flush is for flush data to channel
func (bs *Buflice) Flush() {
	bs.mu.Lock()
	defer bs.mu.Unlock()
	bs.flushResetPosAndCleanSlice()
}

// Close is for close time ticker, clean slice data and slice position
func (bs *Buflice) Close() error {
	bs.mu.Lock()
	bs.mu.Unlock()
	bs.stopTicker()
	bs.currentPos = 0
	bs.slice = make([]interface{}, bs.maxSize)
	return nil
}
