# Channel

### Channels are useful, so does your application, but in some contexts that need the same channel to send data across many receivers, here is the solution.

---

## Install: `go get -u github.com/Dviih/Channel`

## Usage
- `Send[T]` - sends as a parameter.
- `Sender[T]` - creates a new sender channel.
- `Receiver[T]` - creates a new receiver channel.
- `New[T](Options)` - creates a `*Channel`.

## Options
- `Size` - Creates a buffered channel.
- `Timeout` - Sets timeout for Try.
- `Resend` - Might want to resend after timeout.

## Example
```go
package main

import (
	"fmt"
	"github.com/Dviih/Channel"
)

func main() {
	channel := Channel.New[string](Channel.OptionSize(16))

	sender := channel.Sender()
	receiver := channel.Receiver()

	go func() {
		for {
			select {
			case data := <-receiver:
				fmt.Println("Received", data)
			}
		}
	}()

	for {
		sender <- "Hello, World!"
	}
}
```
The code from the example creates a `*Channel` instance and gets both a sender and a receiver, then creates a coroutine for a receiver, which prints when receives data, the last part is a for loop sending `"Hello, World"` non-stop.<br>
`T` represents a type, if unsure use `interface{}`.

---
#### Made for Gophers by Dviih
