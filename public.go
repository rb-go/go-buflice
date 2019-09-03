package buflice

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
	defer bs.mu.Unlock()
	bs.stopTicker()
	bs.currentPos = 0
	bs.slice = make([]interface{}, bs.maxSize)
	return nil
}
