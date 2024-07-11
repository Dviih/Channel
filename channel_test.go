package Channel

import (
	"math/rand/v2"
	"sync"
	"testing"
	"time"
)

const (
	times    = 5
	duration = time.Second
)

var (
	expected uint64
	counters sync.Map
)

func TestChannel(t *testing.T) {
	channel := New[interface{}](16)

	go func(c <-chan interface{}) {
		for {
			select {
			case data := <-c:
				if data != expected {
					t.Errorf("Receiver 1 expected %v but got %v", expected, data)
					continue
				}

				t.Logf("Receiver 1: %v", data)
			}
		}
	}(channel.Receiver())

	go func(c <-chan interface{}) {
		for {
			select {
			case data := <-c:
				if data != expected {
					t.Errorf("Receiver 2 expected %v but got %v", expected, data)
					continue
				}

				t.Logf("Receiver 2: %v", data)
			}
		}
	}(channel.Receiver())

	go func(c chan<- interface{}) {
		for i := 0; i < times; i++ {
			expected = rand.Uint64()
			t.Logf("Expected: %d", expected)

			c <- expected

			c <- s

			wg.Done()
			time.Sleep(duration)
		}
	}(channel.Sender())

	wg.Add(times)
	wg.Wait()
}

func BenchmarkChannel(b *testing.B) {
	channel := New[interface{}](16)

	go func(c <-chan interface{}) {
		for {
			select {
			case <-c:
				count1++
			}
		}
	}(channel.Receiver())

	go func(c <-chan interface{}) {
		for {
			select {
			case <-c:
				count2++
			}
		}
	}(channel.Receiver())

	go func(c chan<- interface{}) {
		for i := 0; i < b.N; i++ {
			c <- rand.Uint64()
		}

		close(c)
	}(channel.Sender())

	wg.Add(b.N)
	wg.Wait()

	b.Logf("Receiver 1 received: %v", count1)
	b.Logf("Receiver 2 received: %v", count2)
}
