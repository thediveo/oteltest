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

import (
	"context"
	"sync/atomic"

	sdklog "go.opentelemetry.io/otel/sdk/log"
)

type LogRecordsChannel chan sdklog.Record

// Exporter writes log records to a Go chan. Use [New] to create an Exporter.
type Exporter struct {
	ch atomic.Pointer[LogRecordsChannel]
}

var _ (sdklog.Exporter) = (*Exporter)(nil)

// New returns a new log record exporter, configured with the passed options. If
// no log record channel has been explicitly configured using [WithChannel], a
// suitable channel will be automatically created and can later be retrieved
// using [Exporter.Ch].
func New(opts ...Option) (*Exporter, error) {
	var o options
	for _, opt := range opts {
		opt(&o)
	}

	if o.ch == nil {
		o.ch = make(LogRecordsChannel, max(o.size, 1))
	}

	e := &Exporter{}
	ch := o.ch
	e.ch.Store(&ch)
	return e, nil
}

// Ch returns the log record channel, or nil after [Exporter.Shutdown] has been
// called.
func (e *Exporter) Ch() LogRecordsChannel {
	ch := e.ch.Load()
	if ch == nil {
		return nil
	}
	return *ch
}

// Export log records to the configured channel.
func (e *Exporter) Export(ctx context.Context, records []sdklog.Record) error {
	ch := e.ch.Load()
	if ch == nil {
		return nil
	}

	for _, rec := range records {
		if err := ctx.Err(); err != nil {
			return err
		}
		*ch <- rec
	}
	return nil
}

// Shutdown the Exporter so that any later calls to [Exporter.Export] will
// perform no operation anymore and closes the writing end of the exporter's log
// record channel.
func (e *Exporter) Shutdown(context.Context) error {
	e.ch.Store(nil)
	return nil
}

// ForceFlush is a no-op.
func (*Exporter) ForceFlush(context.Context) error { return nil }
