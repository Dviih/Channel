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
	"time"
)

type Channel[T interface{}] struct {
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
		current := channel.receivers
		go func() {
			time.Sleep(channel.options.timeout)

			for i, receiver := range channel.receivers {
				x := false
				for _, c := range current {
					if receiver == c {
						x = true
						break
					}
				}

				if x {
					continue
				}

				for _, data := range t {
					if !Try(receiver, data, channel.options.timeout) {
						channel.receivers = append(channel.receivers[:i], channel.receivers[i+1:]...)
					}
				}
			}
		}()
	}

	for _, data := range t {
		for i, receiver := range channel.receivers {
			if !Try(receiver, data, channel.options.timeout) {
				channel.receivers = append(channel.receivers[:i], channel.receivers[i+1:]...)
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
				for i, receiver := range channel.receivers {
					if !Try(receiver, data, channel.options.timeout) {
						channel.receivers = append(channel.receivers[:i], channel.receivers[i+1:]...)
					}
				}
			}
		}
	}()

	return c
}

func (channel *Channel[T]) resend(t ...T) {
	current := channel.receivers
	time.Sleep(channel.options.timeout)

	for i, receiver := range channel.receivers {
		x := false
		for _, c := range current {
			if receiver == c {
				x = true
				break
			}
		}

		if x {
			continue
		}

		for _, data := range t {
			if !Try(receiver, data, channel.options.timeout) {
				channel.receivers = append(channel.receivers[:i], channel.receivers[i+1:]...)
			}
		}
	}
}

func (channel *Channel[T]) Receiver() <-chan T {
	c := make(chan T, channel.options.size)

	channel.receivers = append(channel.receivers, c)
	return c
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
