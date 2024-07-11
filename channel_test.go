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

func receiver(t *testing.T, wg *sync.WaitGroup, id int, c <-chan uint64) {
	for {
		select {
		case data := <-c:
			if data != expected {
				t.Fail()
				t.Errorf("Receiver %d expected %d but got %d", id, expected, data)

				wg.Done()
				return
			}

			t.Logf("Receiver %d: %d", id, data)
		}
	}
}

func TestChannel(t *testing.T) {
	channel := New[uint64](0)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go receiver(t, wg, 1, channel.Receiver())
	go receiver(t, wg, 2, channel.Receiver())

	go func(c chan<- uint64) {
		defer wg.Done()

		for i := 0; i < times; i++ {
			expected = rand.Uint64()
			t.Logf("Expected: %d", expected)

			c <- expected

			time.Sleep(duration)
		}
	}(channel.Sender())

	wg.Wait()
}

func BenchmarkChannel(b *testing.B) {
	channel := New[uint64](0)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func(c <-chan interface{}) {
		for {
			select {
			case <-c:
				count2++
			}
		}
	}(channel.Receiver())

	go func(c chan<- uint64) {
		defer wg.Done()

		for i := 0; i < b.N; i++ {
			c <- rand.Uint64()
		}

		close(c)
	}(channel.Sender())

	wg.Wait()

	counters.Range(func(key, value any) bool {
		b.Logf("Receiver %d received: %d", key, value)
		return true
	})
}
