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

import "time"

type Option interface {
	Name() string
	Value() interface{}
}

type optionSize struct {
	size int
}

func (*optionSize) Name() string {
	return "size"
}

func (size *optionSize) Value() interface{} {
	return size.size
}

func OptionSize(size int) Option {
	return &optionSize{size: size}
}

type optionTimeout struct {
	timeout time.Duration
}

func (*optionTimeout) Name() string {
	return "timeout"
}

func (timeout *optionTimeout) Value() interface{} {
	return timeout.timeout
}

func OptionTimeout(timeout time.Duration) Option {
	return &optionTimeout{timeout: timeout}
}

type optionResend struct{}

func (*optionResend) Name() string {
	return "resend"
}

func (*optionResend) Value() interface{} {
	return true
}

func OptionResend() Option {
	return &optionResend{}
}
