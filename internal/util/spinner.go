package utils

import (
	"fmt"
	"time"
)

var spinnerRunning = false

func StartSpinner() chan bool {
	stop := make(chan bool)
	spinnerRunning = true

	go func() {
		frames := []string{"🤖 thinking.", "🤖 thinking..", "🤖 thinking..."}
		i := 0
		for spinnerRunning {
			fmt.Printf("\r%s", frames[i%len(frames)])
			time.Sleep(300 * time.Millisecond)
			i++
			select {
			case <-stop:
				return
			default:
			}
		}
	}()

	return stop
}

func StopSpinner(stop chan bool) {
	spinnerRunning = false
	stop <- true
	fmt.Print("\r") // clear spinner line
}
