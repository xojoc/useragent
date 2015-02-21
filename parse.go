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
	"github.com/blang/semver"
	"strings"
)

type AgentType int

const (
	TypeUnknown AgentType = iota
	TypeBrowser
	TypeCrawler
	TypeLinkChecker
	TypeValidator
	TypeFeedReader
	TypeLibrary
)

type Security int

const (
	SecurityUnknown Security = iota
	SecurityNone
	SecurityWeak
	SecurityStrong
)

type Version struct {
	Major uint64
	Minor uint64
}

type UserAgent struct {
	Type     AgentType
	OS       string
	Name     string
	Version  semver.Version
	Security Security
}

func newUserAgent() *UserAgent {
	ua := &UserAgent{}
	ua.OS = "unknown"
	ua.Name = "unknown"
	return ua
}

type parseFn func(l *lex) *UserAgent

func Parse(uas string) *UserAgent {
	// we try each user agent parser in order until we get one that succeeds
	for _, f := range []parseFn{parseFirefoxLike, parseChrome, parseDillo, parseGoogleBot} {
		if ua := f(newLex(uas)); ua != nil {
			return ua
		}
	}
	return nil
}

func parseSecurity(l *lex) Security {
	switch {
	case l.match("U; "):
		return SecurityStrong
	case l.match("I; "):
		return SecurityWeak
	case l.match("N; "):
		return SecurityNone
	default:
		return SecurityUnknown
	}
}

func parseMozillaLike(l *lex, ua *UserAgent) bool {
	var ok bool

	ua.Type = TypeBrowser

	if !l.match("Mozilla/5.0 (") {
		return false
	}

	switch {
	case l.match("X11"):
		l.match("; ")
		ua.Security = parseSecurity(l)
		switch {
		case l.match("Linux") || l.match("Ubuntu"):
			ua.OS = "gnu/linux"
		case l.match("OpenBSD"):
			ua.OS = "openbsd"
		default:
			return false
		}
	case l.match("Windows"):
		l.match("; ")
		ua.Security = parseSecurity(l)
		ua.OS = "windows"
	case l.match("Macintosh"):
		l.match("; ")
		ua.Security = parseSecurity(l)
		ua.OS = "mac os"
	default:
		return false
	}

	if _, ok = l.span(") "); !ok {
		return false
	}

	return true
}

func parseNameVersion(l *lex, ua *UserAgent) bool {
	var err error
	var s string
	var ok bool

	s, ok = l.span("/")
	if !ok {
		return false
	}
	ua.Name = strings.ToLower(s)

	if s, ok = l.span(" "); !ok {
		s = l.s[l.p:]
		l.p = len(l.s)
		if s == "" {
			return false
		}
	}

	// kludge:
	//  some versions have extra dot fields (instead of only 3)
	//  we try to detect this and remove all the extra stuff
	//   e.g. X.Y.Z.Q.W-beta -> X.Y.Z-beta
	//  others miss the `patch` field in their version
	//  so we add a fictious one
	//   e.g. X.Y -> X.Y.0

	hypen := strings.SplitN(s, "-", 2)
	fs := strings.Split(hypen[0], ".")
	maxfs := 3
	if len(fs) < 3 {
		if len(fs) == 2 {
			fs = append(fs, "0")
		} else {
			maxfs = len(fs)
		}
	}
	s = strings.Join(fs[:maxfs], ".")
	if len(hypen) > 1 {
		s += "-" + hypen[1]
	}

	ua.Version, err = semver.Parse(s)
	if err != nil {
		return false
	}

	return true
}

func parseFirefoxLike(l *lex) *UserAgent {
	var ok bool
	ua := newUserAgent()

	if !parseMozillaLike(l, ua) {
		return nil
	}
	if !l.match("Gecko") {
		return nil
	}
	if _, ok = l.span(" "); !ok {
		return nil
	}
	if !parseNameVersion(l, ua) {
		return nil
	}

	return ua
}

func parseChrome(l *lex) *UserAgent {
	var ok bool
	ua := newUserAgent()

	if !parseMozillaLike(l, ua) {
		return nil
	}
	if !l.match("AppleWebKit") {
		return nil
	}
	if _, ok = l.span(" "); !ok {
		return nil
	}
	if l.match("(") {
		l.span(") ")
	}
	if !parseNameVersion(l, ua) {
		return nil
	}

	return ua
}

func parseDillo(l *lex) *UserAgent {
	ua := newUserAgent()
	ua.Type = TypeBrowser
	if !parseNameVersion(l, ua) {
		return nil
	}
	if ua.Name != "dillo" || l.s[l.p:] != "" {
		return nil
	}
	return ua
}

func parseGoogleBot(l *lex) *UserAgent {
	ua := newUserAgent()
	ua.Type = TypeCrawler
	if !parseNameVersion(l, ua) {
		return nil
	}
	if ua.Name != "googlebot" {
		return nil
	}
	return ua

}
