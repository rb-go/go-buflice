package buflice

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
