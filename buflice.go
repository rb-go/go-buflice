package buflice

import (
	"time"
)

// NewBuflice method to initiate Buflice and get it
func NewBuflice(size int, flushDuration time.Duration, notifyChannel chan []interface{}) *Buflice {
	bs := &Buflice{
		flushDuration: flushDuration,
		slice:         make([]interface{}, 0, size),
		chDone:        make(chan struct{}, 1),
		flushChan:     notifyChannel,
	}
	return bs
}

// Start starts ticker and serving data
func (bs *Buflice) Start() {
	bs.ticker = time.NewTicker(bs.flushDuration)
	go func() {
		for {
			select {
			case <-bs.chDone:
				return
			case <-bs.ticker.C:
				bs.Flush()
			}
		}
	}()
}

func (bs *Buflice) flushReset() {
	if len(bs.slice) == 0 {
		return
	}
	bs.wgProc.Add(1)
	sendSlice := make([]interface{}, len(bs.slice))
	copy(sendSlice, bs.slice)
	bs.flushChan <- sendSlice
	bs.slice = bs.slice[:0]
	bs.wgProc.Done()
}

// Add is for adding elements
func (bs *Buflice) Add(element interface{}) {
	bs.mu.Lock()
	defer bs.mu.Unlock()
	if len(bs.slice) < cap(bs.slice) {
		bs.slice = append(bs.slice, element)
	}
	if len(bs.slice) == cap(bs.slice) {
		bs.flushReset()
	}
}

// Flush is for manual flush data to channel
func (bs *Buflice) Flush() {
	bs.mu.Lock()
	defer bs.mu.Unlock()
	bs.flushReset()
}

// Close is for close time ticker, clean slice data and slice position
func (bs *Buflice) Close() error {
	bs.mu.Lock()
	defer bs.mu.Unlock()
	bs.ticker.Stop()
	bs.chDone <- struct{}{}
	bs.flushReset()
	bs.wgProc.Wait()
	return nil
}

func (bs *Buflice) GetCurrentLen() int {
	bs.mu.Lock()
	defer bs.mu.Unlock()
	return len(bs.slice)
}

func (bs *Buflice) GetCap() int {
	bs.mu.Lock()
	defer bs.mu.Unlock()
	return cap(bs.slice)
}
