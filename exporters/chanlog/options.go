// Copyright 2025 Harald Albrecht.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may not
// use this file except in compliance with the License. You may obtain a copy
// of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package chanlog

// Option configures a log record channel [Exporter].
type Option func(*options)

type options struct {
	capacity int
	ch       RecordsChannel
}

// WithCap configures the capacity of the implicit log record channel, unless an
// explicit log record channel is configured using [WithChannel]. The specified
// capacity is clamped to at least 1.
func WithCap(capacity int) func(o *options) {
	return func(o *options) {
		o.capacity = max(capacity, 1)
	}
}

// WithChannel configures an explicit log record channel. Any [WithCap]
// configuration is ignored.
func WithChannel(ch RecordsChannel) func(o *options) {
	return func(o *options) {
		o.ch = ch
	}
}
