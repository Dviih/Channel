/*
 *     Channels with multiple receivers and multiple senders capacity.
 *     Copyright (C) 2024  Dviih
 *
 *     This program is free software: you can redistribute it and/or modify
 *     it under the terms of the GNU Affero General Public License as published
 *     by the Free Software Foundation, either version 3 of the License, or
 *     (at your option) any later version.
 *
 *     This program is distributed in the hope that it will be useful,
 *     but WITHOUT ANY WARRANTY; without even the implied warranty of
 *     MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *     GNU Affero General Public License for more details.
 *
 *     You should have received a copy of the GNU Affero General Public License
 *     along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

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

func counter(id int, c <-chan uint64) {
	for {
		select {
		case <-c:
			v, ok := counters.Load(id)
			if !ok {
				counters.Store(id, uint64(1))
				continue
			}

			counters.Store(id, v.(uint64)+1)
		}
	}
}

func BenchmarkChannel(b *testing.B) {
	channel := New[uint64](0)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go counter(1, channel.Receiver())
	go counter(2, channel.Receiver())

	go func(c chan<- uint64) {
		defer wg.Done()

		for i := 0; i < b.N; i++ {
			c <- rand.Uint64()
		}
	}(channel.Sender())

	wg.Wait()

	counters.Range(func(key, value any) bool {
		b.Logf("Receiver %d received: %d", key, value)
		return true
	})
}
