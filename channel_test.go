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
	"context"
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

	channel     = New[uint64]()
	ctx, cancel = context.WithCancel(context.Background())
)

func receiver(t *testing.T, id int, c <-chan uint64) {
	t.Parallel()

	for {
		select {
		case data := <-c:
			if data != expected {
				t.Fail()
				t.Errorf("Receiver %d expected %d but got %d", id, expected, data)

				return
			}

			t.Logf("Receiver %d: %d", id, data)
		case <-ctx.Done():
			return
		}
	}
}

func TestChannelReceiver1(t *testing.T) {
	receiver(t, 1, channel.Receiver())
}

func TestChannelReceiver2(t *testing.T) {
	receiver(t, 2, channel.Receiver())
}

func TestChannel(t *testing.T) {
	t.Parallel()
	wg.Wait()

	sender := channel.Sender()

	for i = 0; i < times; i++ {
		expected = rand.Uint64()
		t.Logf("Expected: %d", expected)

		sender <- expected

		time.Sleep(duration)
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
