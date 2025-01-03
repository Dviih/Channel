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
	"sync"
	"time"
)

type Channel[T interface{}] struct {
	m sync.Mutex

	options   *Options
	receivers []chan T
}

type Options struct {
	size    int
	timeout time.Duration
	resend  bool
}

func (channel *Channel[T]) Send(t ...T) {
	if channel.options.resend {
		go channel.resend(t...)
	}

	for _, data := range t {
		for _, receiver := range channel.receivers {
			if !Try(receiver, data, channel.options.timeout) {
				channel.Close(receiver)
			}
		}
	}
}

func (channel *Channel[T]) Sender() chan<- T {
	c := make(chan T, channel.options.size)

	go func() {
		for {
			select {
			case data := <-c:
				channel.Send(data)
			}
		}
	}()

	return c
}

func (channel *Channel[T]) resend(t ...T) {
	current := channel.receivers
	time.Sleep(channel.options.timeout)

	for _, receiver := range cmp(current, channel.receivers) {
		for _, data := range t {
			if !Try(receiver, data, channel.options.timeout) {
				channel.Close(receiver)
			}
		}
	}
}

func (channel *Channel[T]) Receiver() <-chan T {
	defer channel.m.Unlock()

	c := make(chan T, channel.options.size)

	channel.m.Lock()
	channel.receivers = append(channel.receivers, c)

	return c
}

func (channel *Channel[T]) Close(c <-chan T) {
	defer channel.m.Unlock()
	channel.m.Lock()

	for i, receiver := range channel.receivers {
		if receiver == c {
			channel.receivers = append(channel.receivers[:i], channel.receivers[i+1:]...)
			break
		}
	}
}

func New[T interface{}](v ...Option) *Channel[T] {
	options := &Options{}

	for _, option := range v {
		switch option.Name() {
		case "size":
			options.size = option.Value().(int)
		case "timeout":
			options.timeout = option.Value().(time.Duration)
		case "resend":
			options.resend = true
		}
	}

	return &Channel[T]{
		options: options,
	}
}

func cmp[T interface{}](old []chan T, new []chan T) (cmp []chan T) {
	m := make(map[chan T]bool)

	for _, v := range old {
		m[v] = true
	}

	for _, v := range new {
		if !m[v] {
			cmp = append(cmp, v)
		}
	}

	return cmp
}
