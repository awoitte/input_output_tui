# input_output_tui
Minimalistic tui for reading input and showing lines of output simultaneously.

Compatability tested and working on Linux (tty and virtual terminal) and Windows (CMD and Command)

example useage:
```go
package main

import "github.com/awoitte/input_output_tui"

func main() {
	input := make(chan string)
	output := make(chan string)
	quit := make(chan bool)

	go input_output_tui.Start(input, output, quit)

	for {
		user_input := <-input
		if user_input == "quit" {
			quit <- true
			return
		}
		output <- user_input
	}
}
```
