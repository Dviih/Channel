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
	duration = time.Microsecond
)

var (
	expected uint64

	channel     = New[uint64](OptionTimeout(10 * time.Microsecond))
	ctx, cancel = context.WithCancel(context.Background())
)

func receiver(t *testing.T, id int, c <-chan uint64) {
	t.Parallel()

	for {
		select {
		case data := <-c:
			if data != expected {
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

	for i := 0; i < times; i++ {
		expected = rand.Uint64()

		channel.Send(expected)
		time.Sleep(duration)
	}

	cancel()
}

func BenchmarkChannel(b *testing.B) {
	channel := New[uint64](OptionSize(b.N))

	var wg sync.WaitGroup
	wg.Add(b.N)

	go func(receiver <-chan uint64) {
		for {
			select {
			case <-receiver:
				wg.Done()
			}
		}
	}(channel.Receiver())

	go func(sender chan<- uint64) {
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			sender <- rand.Uint64()
		}

		b.StopTimer()
	}(channel.Sender())

	wg.Wait()
}
