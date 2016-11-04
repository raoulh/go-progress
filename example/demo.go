package main

import (
	"time"

	"github.com/raoulh/go-progress"
)

func testProgressBar(f int) {
	bar := progress.New(100)
	bar.Format = progress.ProgressFormats[f]

	for bar.Inc() {
		time.Sleep(time.Millisecond * 20)
	}
}

func main() {
	for i := 0; i < len(progress.ProgressFormats); i++ {
		testProgressBar(i)
	}
}
