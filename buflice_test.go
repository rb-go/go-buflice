package buflice_test

import (
	"log"
	"sync"
	"testing"
	"time"

	"github.com/rb-pkg/buflice"
)

type Book struct {
	Author string
}

func flushProcessor(chFlush chan []interface{}, chDone chan struct{}, wait *sync.WaitGroup) {
	for {
		select {
		case <-chFlush:
			wait.Add(1)
			// log.Printf("%+v", data)
			wait.Done()
		case <-chDone:
			log.Println("Finished flushProcessor")
			return
		}
	}
}

func TestBuflice_NewBuflice(t *testing.T) {
	chFlush := make(chan []interface{})
	bfl := buflice.NewBuflice(1, 1*time.Second, chFlush)
	bfl.Start()
	bfl.Close()
}

func TestBuflice_GetCap(t *testing.T) {
	chFlush := make(chan []interface{})
	bfl := buflice.NewBuflice(10, 10*time.Second, chFlush)
	bfl.Start()
	cc := bfl.GetCap()
	if cc != 10 {
		t.Errorf("cap %d is not equal to 10", cc)
		return
	}
	bfl.Close()
}

func TestBuflice_GetCurrentLen(t *testing.T) {
	chFlush := make(chan []interface{}, 1)

	bfl := buflice.NewBuflice(20, 60*time.Second, chFlush)

	bfl.Start()

	bfl.Add(Book{Author: "Author #1"})
	bfl.Add(Book{Author: "Author #2"})
	bfl.Add(Book{Author: "Author #3"})
	bfl.Add(Book{Author: "Author #4"})
	bfl.Add(Book{Author: "Author #5"})
	bfl.Add(Book{Author: "Author #6"})
	bfl.Add(Book{Author: "Author #7"})
	bfl.Add(Book{Author: "Author #8"})
	bfl.Add(Book{Author: "Author #9"})
	bfl.Add(Book{Author: "Author #10"})

	ll := bfl.GetCurrentLen()
	if ll != 10 {
		t.Errorf("len %d is not equal to 10", ll)
		return
	}

	bfl.Close()
}

func TestBuflice_AddDur10msec(t *testing.T) {
	chFlush := make(chan []interface{})
	chDone := make(chan struct{})
	wait := sync.WaitGroup{}

	bfl := buflice.NewBuflice(10, 10*time.Millisecond, chFlush)

	go flushProcessor(chFlush, chDone, &wait)
	bfl.Start()

	bfl.Add(Book{Author: "Author #1"})
	bfl.Add(Book{Author: "Author #2"})
	bfl.Add(Book{Author: "Author #3"})
	bfl.Add(Book{Author: "Author #4"})
	bfl.Add(Book{Author: "Author #5"})
	time.Sleep(11 * time.Millisecond)
	bfl.Add(Book{Author: "Author #6"})
	bfl.Add(Book{Author: "Author #7"})
	bfl.Add(Book{Author: "Author #8"})
	bfl.Add(Book{Author: "Author #9"})
	bfl.Add(Book{Author: "Author #10"})
	bfl.Close()

	wait.Wait()
	chDone <- struct{}{}
}

func TestBuflice_Add5(t *testing.T) {
	chFlush := make(chan []interface{})
	chDone := make(chan struct{})
	wait := sync.WaitGroup{}

	bfl := buflice.NewBuflice(5, 5*time.Second, chFlush)

	go flushProcessor(chFlush, chDone, &wait)
	bfl.Start()

	bfl.Add(Book{Author: "Author #1"})
	bfl.Add(Book{Author: "Author #2"})
	bfl.Add(Book{Author: "Author #3"})
	bfl.Add(Book{Author: "Author #4"})
	bfl.Add(Book{Author: "Author #5"})
	bfl.Add(Book{Author: "Author #6"})
	bfl.Add(Book{Author: "Author #7"})
	bfl.Add(Book{Author: "Author #8"})
	bfl.Add(Book{Author: "Author #9"})
	bfl.Add(Book{Author: "Author #10"})

	bfl.Close()
	wait.Wait()
	chDone <- struct{}{}
}

func TestBuflice_Flush(t *testing.T) {
	chFlush := make(chan []interface{})
	chDone := make(chan struct{})
	wait := sync.WaitGroup{}

	bfl := buflice.NewBuflice(100, 10*time.Second, chFlush)

	go flushProcessor(chFlush, chDone, &wait)
	bfl.Start()

	bfl.Add(Book{Author: "Author #1"})
	bfl.Add(Book{Author: "Author #2"})
	bfl.Add(Book{Author: "Author #3"})
	bfl.Add(Book{Author: "Author #4"})
	bfl.Add(Book{Author: "Author #5"})
	bfl.Flush()
	bfl.Add(Book{Author: "Author #6"})
	bfl.Add(Book{Author: "Author #7"})
	bfl.Add(Book{Author: "Author #8"})
	bfl.Add(Book{Author: "Author #9"})
	bfl.Add(Book{Author: "Author #10"})
	bfl.Flush()

	bfl.Close()
	wait.Wait()
	chDone <- struct{}{}
}
