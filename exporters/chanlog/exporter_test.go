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
	"time"

	"go.opentelemetry.io/otel/log"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/log/logtest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/thediveo/success"
)

var _ = Describe("OTel log record channel exporter", func() {

	rf := logtest.RecordFactory{
		EventName: "foo",
	}

	It("returns a channel exporter with an implicitly created log record channel", func() {
		e := Successful(New())
		ch := e.Ch()
		Expect(ch).NotTo(BeNil())
		Expect(ch).To(HaveCap(1))
	})

	It("correctly configures the channel size", func() {
		const size = 42
		e := Successful(New(WithCap(size)))
		ch := e.Ch()
		Expect(ch).NotTo(BeNil())
		Expect(ch).To(HaveCap(size))
	})

	It("configures the buffer with at least room for one", func() {
		e := Successful(New(WithCap(0)))
		ch := e.Ch()
		Expect(ch).NotTo(BeNil())
		Expect(ch).To(HaveCap(1))
	})

	It("closes the channel upon shutdown", func(ctx context.Context) {
		e := Successful(New())
		ch := e.Ch()
		Expect(e.Shutdown(ctx)).To(Succeed())
		Expect(e.Shutdown(ctx)).To(Succeed(), "must be idempotent")
		Expect(ch).To(BeClosed())
		Expect(e.Ch()).To(BeNil())
	})

	It("flushes (not really)", func(ctx context.Context) {
		Expect(Successful(New()).ForceFlush(ctx)).To(Succeed())
	})

	It("exports log records to the buffered channel", func(ctx context.Context) {
		ch := make(RecordsChannel, 2)
		e := Successful(New(WithChannel(ch)))
		Expect(e.Ch()).To(HaveCap(2))

		done := make(chan struct{})
		go func() {
			defer GinkgoRecover()
			for i := range 10 {
				r := rf.NewRecord()
				r.SetBody(log.IntValue(i))
				Expect(e.Export(ctx, []sdklog.Record{r})).To(Succeed())
			}
			close(done)
		}()

		for i := range 10 {
			Eventually(ctx, ch).Within(2 * time.Second).Should(Receive(And(
				HaveField("EventName()", "foo"),
				HaveField("Body()", log.IntValue(i)))))
		}
		Eventually(done).Should(BeClosed())
	})

	It("doesn't export anymore after shutdown", func() {
		e := Successful(New())
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		Expect(e.Shutdown(ctx)).To(Succeed())
		Expect(e.Export(ctx, []sdklog.Record{rf.NewRecord()})).To(Succeed())
		cancel()
		Expect(e.Export(ctx, []sdklog.Record{rf.NewRecord()})).To(MatchError("context canceled"))
	})

	It("cancels exports", func(ctx context.Context) {
		e := Successful(New(WithCap(1)))
		ch := e.Ch()

		exportCtx, cancel := context.WithCancel(ctx)
		defer cancel()
		done := make(chan struct{})
		go func() {
			defer GinkgoRecover()
			Expect(e.Export(exportCtx, []sdklog.Record{
				rf.NewRecord(),
				rf.NewRecord(),
			})).To(MatchError("context canceled"))
			close(done)
		}()

		Eventually(ch).Within(2 * time.Second).Should(HaveLen(1))
		cancel()
		Eventually(done).Should(BeClosed())
	})

})
