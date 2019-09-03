package buflice

import (
	"log"
	"sync"
	"testing"
	"time"
)

var bfl *Buflice
var chDone chan bool
var chFlush chan []interface{}
var wgWait sync.WaitGroup

func flushProcessor() {
	for {
		select {
		case data := <-chFlush:
			log.Println(data)
		case <-chDone:
			wgWait.Done()
			log.Println("Finished flushProcessor")
			return
		}
	}
}

func TestNewBuflice(t *testing.T) {
	chDone = make(chan bool)
	chFlush = make(chan []interface{})

	wgWait.Add(1)

	bfl = NewBuflice(6, 5*time.Second, chFlush)

	go flushProcessor()
}

func TestBuflice_Add(t *testing.T) {
	bfl.Add("Record #1")
	bfl.Add("Record #2")
	bfl.Add("Record #3")
	bfl.Add("Record #4")
	bfl.Add("Record #5")
	bfl.Add("Record #6")
	bfl.Add("Record #7")
	bfl.Add("Record #8")
	bfl.Add("Record #9")
	bfl.Add("Record #10")
}

func TestBuflice_Flush(t *testing.T) {
	bfl.Flush()
}

func TestBuflice_Close(t *testing.T) {
	chDone <- true
	wgWait.Wait()
	err := bfl.Close()
	if err != nil {
		t.Error(err)
	}
}
