//golangcitest:args -Egomegacheck
package gomegacheck

import (
	"time"

	"github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
)

// Copyright (c) 2022 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// This is a test file for executing unit tests for gomegacheck using golang.org/x/tools/go/analysis/analysistest.

func expect() {
	gomega.Expect("foo").To(Match())
	gomega.Expect("foo").Should(Match())
	gomega.Expect("foo").NotTo(Match())
	gomega.Expect("foo").ToNot(Match())
	gomega.Expect("foo").ShouldNot(Match())
	gomega.Expect("foo").WithOffset(1).To(Match())
	gomega.Expect("foo").WithOffset(1).Should(Match())
	gomega.Expect("foo").WithOffset(1).NotTo(Match())
	gomega.Expect("foo").WithOffset(1).ToNot(Match())
	gomega.Expect("foo").WithOffset(1).ShouldNot(Match())

	gomega.Expect("foo")               // want `gomega.Assertion is missing a call to one of Should, ShouldNot, To, ToNot, NotTo`
	gomega.Expect("foo").WithOffset(1) // want `gomega.Assertion is missing a call to one of Should, ShouldNot, To, ToNot, NotTo`

	gomega.Expect("foo").To(Match(), Match())               // want `GomegaMatcher passed to optionalDescription param, you can only pass one matcher to each assertion`
	gomega.Expect("foo").Should(Match(), Match())           // want `GomegaMatcher passed to optionalDescription param, you can only pass one matcher to each assertion`
	gomega.Expect("foo").NotTo(Match(), Match())            // want `GomegaMatcher passed to optionalDescription param, you can only pass one matcher to each assertion`
	gomega.Expect("foo").ToNot(Match(), Match())            // want `GomegaMatcher passed to optionalDescription param, you can only pass one matcher to each assertion`
	gomega.Expect("foo").ShouldNot(Match(), Match())        // want `GomegaMatcher passed to optionalDescription param, you can only pass one matcher to each assertion`
	gomega.Expect("foo").WithOffset(1).To(Match(), Match()) // want `GomegaMatcher passed to optionalDescription param, you can only pass one matcher to each assertion`
}

func expectWithOffset() {
	gomega.ExpectWithOffset(1, "foo").To(Match())
	gomega.ExpectWithOffset(1, "foo").Should(Match())
	gomega.ExpectWithOffset(1, "foo").NotTo(Match())
	gomega.ExpectWithOffset(1, "foo").ToNot(Match())
	gomega.ExpectWithOffset(1, "foo").ShouldNot(Match())
	gomega.ExpectWithOffset(1, "foo").WithOffset(1).To(Match())
	gomega.ExpectWithOffset(1, "foo").WithOffset(1).Should(Match())
	gomega.ExpectWithOffset(1, "foo").WithOffset(1).NotTo(Match())
	gomega.ExpectWithOffset(1, "foo").WithOffset(1).ToNot(Match())
	gomega.ExpectWithOffset(1, "foo").WithOffset(1).ShouldNot(Match())

	gomega.ExpectWithOffset(1, "foo")               // want `gomega.Assertion is missing a call to one of Should, ShouldNot, To, ToNot, NotTo`
	gomega.ExpectWithOffset(1, "foo").WithOffset(1) // want `gomega.Assertion is missing a call to one of Should, ShouldNot, To, ToNot, NotTo`
}

func omega() {
	gomega.Ω("foo").To(Match())
	gomega.Ω("foo").Should(Match())
	gomega.Ω("foo").NotTo(Match())
	gomega.Ω("foo").ToNot(Match())
	gomega.Ω("foo").ShouldNot(Match())
	gomega.Ω("foo").WithOffset(1).To(Match())
	gomega.Ω("foo").WithOffset(1).Should(Match())
	gomega.Ω("foo").WithOffset(1).NotTo(Match())
	gomega.Ω("foo").WithOffset(1).ToNot(Match())
	gomega.Ω("foo").WithOffset(1).ShouldNot(Match())

	gomega.Ω("foo")               // want `gomega.Assertion is missing a call to one of Should, ShouldNot, To, ToNot, NotTo`
	gomega.Ω("foo").WithOffset(1) // want `gomega.Assertion is missing a call to one of Should, ShouldNot, To, ToNot, NotTo`
}

func eventually() {
	c := make(chan bool)

	gomega.Eventually(c).Should(Match())
	gomega.Eventually(c).ShouldNot(Match())
	gomega.Eventually(c).WithOffset(1).Should(Match())
	gomega.Eventually(c).WithOffset(1).ShouldNot(Match())
	gomega.Eventually(c).WithTimeout(time.Second).Should(Match())
	gomega.Eventually(c).WithTimeout(time.Second).ShouldNot(Match())
	gomega.Eventually(c).WithPolling(time.Second).Should(Match())
	gomega.Eventually(c).WithPolling(time.Second).ShouldNot(Match())

	gomega.Eventually(c)                          // want `gomega.AsyncAssertion is missing a call to one of Should, ShouldNot`
	gomega.Eventually(c).WithOffset(1)            // want `gomega.AsyncAssertion is missing a call to one of Should, ShouldNot`
	gomega.Eventually(c).WithTimeout(time.Second) // want `gomega.AsyncAssertion is missing a call to one of Should, ShouldNot`
	gomega.Eventually(c).WithPolling(time.Second) // want `gomega.AsyncAssertion is missing a call to one of Should, ShouldNot`
}

func consistently() {
	c := make(chan bool)

	gomega.Consistently(c).Should(Match())
	gomega.Consistently(c).ShouldNot(Match())
	gomega.Consistently(c).WithOffset(1).Should(Match())
	gomega.Consistently(c).WithOffset(1).ShouldNot(Match())
	gomega.Consistently(c).WithTimeout(time.Second).Should(Match())
	gomega.Consistently(c).WithTimeout(time.Second).ShouldNot(Match())
	gomega.Consistently(c).WithPolling(time.Second).Should(Match())
	gomega.Consistently(c).WithPolling(time.Second).ShouldNot(Match())

	gomega.Consistently(c)                          // want `gomega.AsyncAssertion is missing a call to one of Should, ShouldNot`
	gomega.Consistently(c).WithOffset(1)            // want `gomega.AsyncAssertion is missing a call to one of Should, ShouldNot`
	gomega.Consistently(c).WithTimeout(time.Second) // want `gomega.AsyncAssertion is missing a call to one of Should, ShouldNot`
	gomega.Consistently(c).WithPolling(time.Second) // want `gomega.AsyncAssertion is missing a call to one of Should, ShouldNot`
}

func assertionInAsync() {
	gomega.Eventually(func(g gomega.Gomega) {
		g.Expect("foo").To(Match())
	}).Should(Match())
	gomega.Eventually(func(g gomega.Gomega) {
		g.Expect("foo").Should(Match())
	}).Should(Match())
	gomega.Eventually(func(g gomega.Gomega) {
		g.Expect("foo").NotTo(Match())
	}).Should(Match())
	gomega.Eventually(func(g gomega.Gomega) {
		g.Expect("foo").ToNot(Match())
	}).Should(Match())
	gomega.Eventually(func(g gomega.Gomega) {
		g.Expect("foo").ShouldNot(Match())
	}).Should(Match())
	gomega.Eventually(func(g gomega.Gomega) {
		g.Expect("foo").WithOffset(1).To(Match())
	}).Should(Match())
	gomega.Eventually(func(g gomega.Gomega) {
		g.Expect("foo").WithOffset(1).Should(Match())
	}).Should(Match())
	gomega.Eventually(func(g gomega.Gomega) {
		g.Expect("foo").WithOffset(1).NotTo(Match())
	}).Should(Match())
	gomega.Eventually(func(g gomega.Gomega) {
		g.Expect("foo").WithOffset(1).ToNot(Match())
	}).Should(Match())
	gomega.Eventually(func(g gomega.Gomega) {
		g.Expect("foo").WithOffset(1).ShouldNot(Match())
	}).Should(Match())

	gomega.Eventually(func(g gomega.Gomega) {
		g.Expect("foo") // want `gomega.Assertion is missing a call to one of Should, ShouldNot, To, ToNot, NotTo`
	}).Should(Match())
	gomega.Eventually(func(g gomega.Gomega) {
		g.Expect("foo").WithOffset(1) // want `gomega.Assertion is missing a call to one of Should, ShouldNot, To, ToNot, NotTo`
	}).Should(Match())
}

// returning Assertion/AsyncAssertion should not yield an error.

func helper(actual interface{}) gomega.Assertion {
	return gomega.Expect(actual).WithOffset(1)
}

func helperAsync(actual interface{}) gomega.AsyncAssertion {
	return gomega.Eventually(actual).WithOffset(1)
}

// calling helpers should however yield an error

func callingHelper() {
	helper("foo").To(Match())
	helper("foo").Should(Match())
	helper("foo").NotTo(Match())
	helper("foo").ToNot(Match())
	helper("foo").ShouldNot(Match())
	helper("foo").WithOffset(1).To(Match())
	helper("foo").WithOffset(1).Should(Match())
	helper("foo").WithOffset(1).NotTo(Match())
	helper("foo").WithOffset(1).ToNot(Match())
	helper("foo").WithOffset(1).ShouldNot(Match())

	helper("foo")               // want `gomega.Assertion is missing a call to one of Should, ShouldNot, To, ToNot, NotTo`
	helper("foo").WithOffset(1) // want `gomega.Assertion is missing a call to one of Should, ShouldNot, To, ToNot, NotTo`
}

func callingHelperAsync() {
	helperAsync("foo").Should(Match())
	helperAsync("foo").ShouldNot(Match())
	helperAsync("foo").WithOffset(1).Should(Match())
	helperAsync("foo").WithOffset(1).ShouldNot(Match())

	helperAsync("foo")               // want `gomega.AsyncAssertion is missing a call to one of Should, ShouldNot`
	helperAsync("foo").WithOffset(1) // want `gomega.AsyncAssertion is missing a call to one of Should, ShouldNot`
}

func Match() types.GomegaMatcher {
	return dummyMatcher{}
}

type dummyMatcher struct{}

func (m dummyMatcher) Match(actual interface{}) (success bool, err error)        { return true, nil }
func (m dummyMatcher) FailureMessage(actual interface{}) (message string)        { return "fail" }
func (m dummyMatcher) NegatedFailureMessage(actual interface{}) (message string) { return "whatever" }
