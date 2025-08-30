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

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/thediveo/success"
)

var _ = Describe("OTel log record channel exporter", func() {

	It("returns a channel exporter with an implicitly created log record channel", func() {
		e := Successful(New())
		ch := e.Ch()
		Expect(ch).NotTo(BeNil())
		Expect(ch).To(HaveCap(1))
	})

	It("correctly configures the channel size", func() {
		const size = 42
		e := Successful(New(WithSize(size)))
		ch := e.Ch()
		Expect(ch).NotTo(BeNil())
		Expect(ch).To(HaveCap(size))
	})

	It("configures the buffer with at least room for one", func() {
		e := Successful(New(WithSize(0)))
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

	It("exports log records", func() {

	})

})
