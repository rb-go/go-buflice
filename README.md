# buflice

This package need to create buffered slice that can be flushed when reach size or duration limit  

### When it can be needed?
Example: You have a worker for rabbitmq that receives jobs from queue. You receive them one by one and process it. But sometimes 
you need to accumulate data from jobs for batch processing in database.

[Website](https://riftbit.com) | [Blog](https://ergoz.ru)

[![license](https://img.shields.io/github/license/rb-pkg/buflice.svg)](LICENSE)
[![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](https://pkg.go.dev/github.com/rb-pkg/buflice)
[![Coverage Status](https://coveralls.io/repos/github/rb-pkg/buflice/badge.svg?branch=master)](https://coveralls.io/github/rb-pkg/buflice?branch=master)
[![Build Status](https://travis-ci.org/rb-pkg/buflice.svg?branch=master)](https://travis-ci.org/rb-pkg/buflice)
[![Go Report Card](https://goreportcard.com/badge/github.com/rb-pkg/buflice)](https://goreportcard.com/report/github.com/rb-pkg/buflice)

## Installation

```bash
go get -u github.com/riftbit/buflice
```

## Example (dirty example)

```go
package main

import (
	"log"
	"sync"
	"time"

	"github.com/rb-pkg/buflice"
)

type Book struct {
	Author string
}

func flushProcessor(chFlush chan []interface{}, chDone chan struct{}, wait *sync.WaitGroup) {
	for {
		select {
		case data := <-chFlush:
			wait.Add(1)
			log.Printf("%+v", data)
			wait.Done()
		case <-chDone:
			log.Println("Finished flushProcessor")
			return
		}
	}
}

func main() {
	chFlush := make(chan []interface{})
	chDone := make(chan struct{})
	wait := sync.WaitGroup{}

	bfl := buflice.NewBuflice(10, 1000*time.Millisecond, chFlush)

	go flushProcessor(chFlush, chDone, &wait)
	bfl.Start()

	bfl.Add(Book{Author: "Author #1"})
	bfl.Add(Book{Author: "Author #2"})
	bfl.Add(Book{Author: "Author #3"})
	bfl.Add(Book{Author: "Author #4"})
	bfl.Add(Book{Author: "Author #5"})
	time.Sleep(1111 * time.Millisecond)
	bfl.Add(Book{Author: "Author #6"})
	bfl.Add(Book{Author: "Author #7"})
	bfl.Add(Book{Author: "Author #8"})
	bfl.Add(Book{Author: "Author #9"})
	bfl.Add(Book{Author: "Author #10"})
	
	err := bfl.Close()
	if err != nil {
		log.Fatalln(err)
    }

	wait.Wait()
	chDone <- struct{}{}

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
