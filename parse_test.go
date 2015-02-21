/* Copyright (C) 2015 by Alexandru Cojocaru */

/* This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.

   You should have received a copy of the GNU General Public License
   along with this program.  If not, see <http://www.gnu.org/licenses/>. */

package useragent

import (
	"fmt"
	"github.com/blang/semver"
	"testing"
)

func eqUA(a *UserAgent, b *UserAgent) bool {
	if a == nil || b == nil {
		return false
	}

	if a.Type != b.Type ||
		a.OS != b.OS ||
		a.Name != b.Name ||
		!a.Version.EQ(b.Version) ||
		a.Security != b.Security {
		return false
	}
	return true
}

func mustParse(s string) semver.Version {
	v, err := semver.Parse(s)
	if err != nil {
		panic(`semver: Parse(` + s + `): ` + err.Error())
	}
	return v
}

func TestFirefoxLike(t *testing.T) {
	var got *UserAgent
	want := &UserAgent{}

	got = Parse(`Mozilla/5.0 (X11; Linux i686; rv:38.0) Gecko/20100101 Firefox/38.0`)
	want.Type = TypeBrowser
	want.OS = "gnu/linux"
	want.Name = "firefox"
	want.Version = mustParse("38.0.0")
	want.Security = SecurityUnknown
	if !eqUA(want, got) {
		t.Errorf("expected %+v, got %+v\n", want, got)
	}

	got = Parse(`Mozilla/5.0 (X11; U; Linux x86_64; sv-SE; rv:1.9.1.16) Gecko/20120714 Iceweasel/3.5.16 (like Firefox/3.5.16)`)
	want.Type = TypeBrowser
	want.OS = "gnu/linux"
	want.Name = "iceweasel"
	want.Version = mustParse("3.5.16")
	want.Security = SecurityStrong
	if !eqUA(want, got) {
		t.Errorf("expected %+v, got %+v\n", want, got)
	}

	got = Parse(`Mozilla/5.0 (Windows x86; rv:19.0) Gecko/20100101 Firefox/19.0`)
	want.Type = TypeBrowser
	want.OS = "windows"
	want.Name = "firefox"
	want.Version = mustParse("19.0.0")
	want.Security = SecurityUnknown
	if !eqUA(want, got) {
		t.Errorf("expected %+v, got %+v\n", want, got)
	}
}

func TestChrome(t *testing.T) {
	var got *UserAgent
	want := &UserAgent{}

	got = Parse(`Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36`)
	want.Type = TypeBrowser
	want.OS = "windows"
	want.Name = "chrome"
	want.Version = mustParse("41.0.2228")
	want.Security = SecurityUnknown
	if !eqUA(want, got) {
		t.Errorf("expected %+v, got %+v\n", want, got)
	}
	got = Parse(`Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2227.0 Safari/537.36`)
	want.Type = TypeBrowser
	want.OS = "gnu/linux"
	want.Name = "chrome"
	want.Version = mustParse("41.0.2227")
	want.Security = SecurityUnknown
	if !eqUA(want, got) {
		t.Errorf("expected %+v, got %+v\n", want, got)
	}
}

func TestDillo(t *testing.T) {
	var got *UserAgent
	want := &UserAgent{}

	got = Parse(`Dillo/0.8.6-i18n-misc`)
	want.Type = TypeBrowser
	want.OS = "unknown"
	want.Name = "dillo"
	want.Version = mustParse("0.8.6-i18n-misc")
	want.Security = SecurityUnknown
	if !eqUA(want, got) {
		t.Errorf("expected %+v, got %+v\n", want, got)
	}
}

func TestGoogleBot(t *testing.T) {
	var got *UserAgent
	want := &UserAgent{}

	got = Parse(`Googlebot/2.1 (+http://www.googlebot.com/bot.html)`)
	want.Type = TypeCrawler
	want.OS = "unknown"
	want.Name = "googlebot"
	want.Version = mustParse("2.1.0")
	want.Security = SecurityUnknown
	if !eqUA(want, got) {
		t.Errorf("expected %+v, got %+v\n", want, got)
	}

}

func ExampleParse() {
	ua := Parse("Mozilla/5.0 (X11; Linux i686; rv:38.0) Gecko/20100101 Firefox/38.0")
	fmt.Print(ua)
	// Output: Type: Browser
	//Name: firefox
	//Version: 38.0.0
	//OS: gnu/linux
	//Security: Unknown security
}
