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

// RecordsChannel channels OTel SDK log records as emitted by an OTel
// [sdklog.Logger]. The log records coming through the channel are of type
// [sdklog.Record], not to be confused with the
// [go.opentelemetry.io/otel/log.Record] type used when applications create
// their log records.
type RecordsChannel chan sdklog.Record

// Exporter writes log records to a Go channel of type [RecordsChannel] (chan of
// [sdklog.Record]). Use [New] to create an Exporter.
type Exporter struct {
	ch atomic.Pointer[RecordsChannel]
}

// statically ensure that we fulfill the OTel logging SDK's Exporter interface.
var _ (sdklog.Exporter) = (*Exporter)(nil)

// New returns a new log record exporter, configured with the passed options.
//
// If no log record channel has been explicitly configured using [WithChannel],
// a suitable channel will be implicitly created and can later be retrieved
// using [Exporter.Ch]. Please note that the minimum configurable buffer size of
// an implicitly created channel is 1.
func New(opts ...Option) (*Exporter, error) {
	var o options
	for _, opt := range opts {
		opt(&o)
	}

	if o.ch == nil {
		o.ch = make(RecordsChannel, max(o.capacity, 1))
	}

	e := &Exporter{}
	ch := o.ch
	e.ch.Store(&ch)
	return e, nil
}

// Ch returns the log record channel, or nil after [Exporter.Shutdown] has been
// called.
func (e *Exporter) Ch() RecordsChannel {
	ch := e.ch.Load()
	if ch == nil {
		return nil
	}
	return *ch
}

// Export log records to the configured channel. It does nothing after
// [Exporter.Shutdown] has been called.
func (e *Exporter) Export(ctx context.Context, records []sdklog.Record) error {
	ch := e.ch.Load()
	if ch == nil {
		return ctx.Err()
	}

	for _, rec := range records {
		select {
		case *ch <- rec:
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return ctx.Err()
}

// Shutdown the Exporter so that any later calls to [Exporter.Export] will
// perform no operation anymore and additionally closes the writing end of the
// exporter's log record channel.
func (e *Exporter) Shutdown(context.Context) error {
	ch := e.ch.Swap(nil)
	if ch != nil {
		close(*ch)
	}
	return nil
}

// ForceFlush is a no-op.
func (*Exporter) ForceFlush(context.Context) error { return nil }
