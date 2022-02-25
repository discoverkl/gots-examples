package main

import (
	"context"
	"embed"
	"flag"
	"io/fs"
	"log"
	"sync"

	"github.com/discoverkl/gots/ui"
)

//go:embed fe/dist
var root embed.FS

var port int

var cancelActiveJob context.CancelFunc
var jobLock sync.Mutex

func init() {
	flag.IntVar(&port, "p", 80, "binding port")
	flag.Parse()
}

func main() {
	www, _ := fs.Sub(root, "fe/dist")
	app := ui.New(
		ui.OnlinePort(port),
		ui.Root(www),
	)

	app.BindFunc("js2go", js2go)

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

func js2go(count int, fn *ui.Function) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	jobLock.Lock()
	if cancelActiveJob != nil {
		cancelActiveJob()
	}
	cancelActiveJob = cancel
	jobLock.Unlock()

	ch := Prime(count)
	buffer := []int{}
	i := 0
loop:
	for {
		select {
		case prime, ok := <-ch:
			if !ok {
				break loop
			}
			i++
			buffer = append(buffer, prime)
			if len(buffer) >= 10 {
				fn.Call(buffer)
				buffer = buffer[0:0]
			}
		case <-ctx.Done():
			return
		}
	}
	if len(buffer) > 0 {
		fn.Call(buffer)
	}
}

// A concurrent prime sieve

// Send the sequence 2, 3, 4, ... to channel 'ch'.
func Generate(ch chan<- int) {
	for i := 2; ; i++ {
		ch <- i // Send 'i' to channel 'ch'.
	}
}

// Copy the values from channel 'in' to channel 'out',
// removing those divisible by 'prime'.
func Filter(in <-chan int, out chan<- int, prime int) {
	for {
		i := <-in // Receive value from 'in'.
		if i%prime != 0 {
			out <- i // Send 'i' to 'out'.
		}
	}
}

// The prime sieve: Daisy-chain Filter processes.
func Prime(count int) chan int {
	ret := make(chan int)
	ch := make(chan int) // Create a new channel.
	go Generate(ch)      // Launch Generate goroutine.
	go func() {
		defer close(ret)
		for i := 0; i < count; i++ {
			prime := <-ch
			ret <- prime
			ch1 := make(chan int)
			go Filter(ch, ch1, prime)
			ch = ch1
		}
	}()
	return ret
}
