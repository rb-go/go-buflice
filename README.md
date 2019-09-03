# buflice

This package need to create buffered slice that can be flushed when reach size or duration limit  

### When it can be needed?
Example: You have a worker for rabbitmq that receives jobs from queue. You receive them one by one and process it. But sometimes 
you need to accumulate data from jobs for batch processing in database.

[Website](https://riftbit.com) | [Blog](https://ergoz.ru)

[![license](https://img.shields.io/github/license/riftbit/buflice.svg)](LICENSE)
[![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](https://godoc.org/github.com/riftbit/buflice)
[![Coverage Status](https://coveralls.io/repos/github/riftbit/buflice/badge.svg?branch=master)](https://coveralls.io/github/riftbit/buflice?branch=master)
[![Build Status](https://travis-ci.org/riftbit/buflice.svg?branch=master)](https://travis-ci.org/riftbit/buflice)
[![Go Report Card](https://goreportcard.com/badge/github.com/riftbit/buflice)](https://goreportcard.com/report/github.com/riftbit/buflice)

## Installation

```bash
go get -u github.com/riftbit/buflice
```

## Example

```go

package main

import (
	"log"
	"sync"
	"testing"
	"time"

	"github.com/riftbit/buflice"
)

var bfl *buflice.Buflice
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

func main() {
	chDone = make(chan bool)
	chFlush = make(chan []interface{})

	wgWait.Add(1)

	bfl = buflice.NewBuflice(6, 5*time.Second, chFlush)
	go flushProcessor()

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

	bfl.Flush()

    chDone <- true
	wgWait.Wait()
	err := bfl.Close()
	if err != nil {
		log.Fatalln(err)
	}

}
```

Will print:

```bash
2019/09/03 14:56:28 [Record #1 Record #2 Record #3 Record #4 Record #5 Record #6]
2019/09/03 14:56:28 [Record #7 Record #8 Record #9 Record #10]
2019/09/03 14:56:28 Finished flushProcessor
```


## Credits

Thanks to:

- Everyone that [gave this repo a star](https://github.com/riftbit/buflice/stargazers) :star: - *you keep me motivated* :slightly_smiling_face: 
- [Contributors](https://github.com/riftbit/buflice/graphs/contributors) that submitted useful [pull-requests](https://github.com/riftbit/buflice/pulls?utf8=%E2%9C%93&q=is%3Apr+is%3Aclosed+is%3Amerged) or opened good issues with suggestions or a detailed bug report.
