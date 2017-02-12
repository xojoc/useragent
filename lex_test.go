// Written by https://xojoc.pw. GPLv3 or later.

package useragent

import (
	"testing"
)

// verify output and state
func (l lex) checkLex(t *testing.T, success, expectSuccess bool, output, expectOutput, expectState string) {
	if success != expectSuccess {
		t.Errorf("success value '%t' does not match expected '%t'", success, expectSuccess)
	}
	if output != expectOutput {
		t.Errorf("output\n'%s' does not match expected\n'%s'", output, expectOutput)
	}
	if l.s[l.p:] != expectState {
		t.Errorf("lexer state\n'%s' does not match expected\n'%s'", l.s[l.p:], expectState)
	}
}

func (l lex) assertLex(t *testing.T, success, expectSuccess bool, output, expectOutput, expectState string) {
	l.checkLex(t, success, expectSuccess, output, expectOutput, expectState)
	if t.Failed() {
		t.Fatalf("assertLex failed for (%t,%s,%s)", expectSuccess, expectOutput, expectState)
	}
}

func TestSpans(t *testing.T) {
	testLex := newLex("☕Mozilla/5.0 (X11; Linux i686; rv:38.0) Gecko/20100101 這是什麼 Firefox/38.0")

	// quick test of match functions
	if testLex.match("Mozilla") {
		t.Fatalf("Matched non-starting string!")
	}
	if !testLex.matchNoConsume("☕M") {
		t.Fatalf("failed to match unicode start")
	}
	if !testLex.match("☕M") {
		t.Fatalf("failed to match unicode start")
	}
	// (The first span test verifies that matchFirst consumed "oz", not "ozilla")
	if !testLex.matchFirst("ill", "oz", "ozilla") {
		t.Fatalf("failed matchFirst")
	}

	// test that all spans consume nothing on failure
	m, ok := testLex.span("M")
	testLex.assertLex(t, ok, false, m, "", "illa/5.0 (X11; Linux i686; rv:38.0) Gecko/20100101 這是什麼 Firefox/38.0")
	m, ok = testLex.spanAny("MZYA")
	testLex.assertLex(t, ok, false, m, "", "illa/5.0 (X11; Linux i686; rv:38.0) Gecko/20100101 這是什麼 Firefox/38.0")
	m, ok = testLex.spanBefore(";", "(")
	testLex.assertLex(t, ok, false, m, "", "illa/5.0 (X11; Linux i686; rv:38.0) Gecko/20100101 這是什麼 Firefox/38.0")

	// test normal span
	m, ok = testLex.span(";")
	testLex.assertLex(t, ok, true, m, "illa/5.0 (X11", " Linux i686; rv:38.0) Gecko/20100101 這是什麼 Firefox/38.0")

	// test spanBefore
	m, ok = testLex.spanBefore("8", ")")
	testLex.assertLex(t, ok, true, m, " Linux i6", "6; rv:38.0) Gecko/20100101 這是什麼 Firefox/38.0")
	m, ok = testLex.spanBefore("8", ")")
	testLex.assertLex(t, ok, true, m, "6; rv:3", ".0) Gecko/20100101 這是什麼 Firefox/38.0")
	m, ok = testLex.spanBefore("8", ")")
	testLex.assertLex(t, ok, false, m, "", ".0) Gecko/20100101 這是什麼 Firefox/38.0")

	// test spanAny
	m, ok = testLex.spanAny("☕/這什")
	testLex.assertLex(t, ok, true, m, ".0) Gecko", "20100101 這是什麼 Firefox/38.0")
	m, ok = testLex.spanAny("☕/這什")
	testLex.assertLex(t, ok, true, m, "20100101 ", "是什麼 Firefox/38.0")
	m, ok = testLex.spanAny("☕這Q")
	testLex.assertLex(t, ok, false, m, "", "是什麼 Firefox/38.0")

	// extra test cases
	m, ok = testLex.spanBefore("是", "Q")
	testLex.assertLex(t, ok, true, m, "", "什麼 Firefox/38.0")
	m, ok = testLex.span("麼")
	testLex.assertLex(t, ok, true, m, "什", " Firefox/38.0")
}
