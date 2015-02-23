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

// Package useragent parses a user agent string.
package useragent

import (
	"fmt"
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

func (a AgentType) String() string {
	switch a {
	case TypeUnknown:
		return "Unkonwn Agent type"
	case TypeBrowser:
		return "Browser"
	case TypeCrawler:
		return "Crawler"
	case TypeLinkChecker:
		return "Link Checker"
	case TypeValidator:
		return "Validator"
	case TypeFeedReader:
		return "Feed Reader"
	case TypeLibrary:
		return "Library"
	default:
		panic("")
	}
}

// Some browsers may put security level information in their user agent string.
type Security int

const (
	SecurityUnknown Security = iota
	SecurityNone
	SecurityWeak
	SecurityStrong
)

func (s Security) String() string {
	switch s {
	case SecurityUnknown:
		return "Unknown security"
	case SecurityNone:
		return "No security"
	case SecurityWeak:
		return "Weak security"
	case SecurityStrong:
		return "Strong Security"
	default:
		panic("cannot happen")
	}
}

type UserAgent struct {
	// The original user agent string.
	Original string
	Type     AgentType
	// The browser/crawler/etc. name in lowercase. For example:
	//  firefox, iceweasel, icecat
	//  dillo
	//  chrome
	//  ie
	//  googlebot
	// If the name is not known, Name will be `unknown'.
	Name    string
	Version semver.Version
	// The OS name in lowercase. Can be one of:
	//  gnu/linux
	//  openbsd
	//  windows
	//  macosx
	//  unknown
	OS       string
	Security Security
}

func (ua *UserAgent) String() string {
	return fmt.Sprintf(`Type: %v
Name: %v
Version: %v
OS: %v
Security: %v`, ua.Type, ua.Name, ua.Version, ua.OS, ua.Security)
}

func new() *UserAgent {
	ua := &UserAgent{}
	ua.Name = "unknown"
	ua.OS = "unknown"
	return ua
}

type parseFn func(l *lex) *UserAgent

// Try to extract information about an user agent from uas.
// Since user agent strings don't have a standard, this function uses heuristics.
func Parse(uas string) *UserAgent {
	// we try each user agent parser in order until we get one that succeeds
	for _, f := range []parseFn{parseFirefoxLike, parseChrome, parseDillo, parseIE, parseGoogleBot} {
		if ua := f(newLex(strings.ToLower(uas))); ua != nil {
			ua.Original = uas
			return ua
		}
	}
	return nil
}

func parseSecurity(l *lex) Security {
	switch {
	case l.match("u; "):
		return SecurityStrong
	case l.match("i; "):
		return SecurityWeak
	case l.match("n; "):
		return SecurityNone
	default:
		return SecurityUnknown
	}
}

func parseMozillaLike(l *lex, ua *UserAgent) bool {
	var ok bool

	ua.Type = TypeBrowser

	if !l.match("mozilla/5.0 (") {
		return false
	}

	l.match("compatible; ")

	switch {
	case l.match("x11; "):
		ua.Security = parseSecurity(l)
		switch {
		case l.match("linux") || l.match("ubuntu"):
			ua.OS = "gnu/linux"
		case l.match("openbsd"):
			ua.OS = "openbsd"
		default:
			return false
		}
	case l.match("windows"):
		l.span("; ")
		ua.Security = parseSecurity(l)
		ua.OS = "windows"
	case l.match("macintosh; "):
		ua.Security = parseSecurity(l)
		ua.OS = "macosx"
	default:
		return false
	}

	if _, ok = l.span(") "); !ok {
		return false
	}

	return true
}

func parseVersion(l *lex, ua *UserAgent, sep string) bool {
	var err error
	var s string
	var ok bool

	if s, ok = l.span(sep); !ok {
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

func parseNameVersion(l *lex, ua *UserAgent) bool {
	var s string
	var ok bool

	s, ok = l.span("/")
	if !ok {
		return false
	}
	ua.Name = strings.ToLower(s)

	return parseVersion(l, ua, " ")
}

func parseFirefoxLike(l *lex) *UserAgent {
	var ok bool
	ua := new()

	if !parseMozillaLike(l, ua) {
		return nil
	}
	if !l.match("gecko") {
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
	ua := new()

	if !parseMozillaLike(l, ua) {
		return nil
	}
	if !l.match("applewebkit") {
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
	ua := new()
	ua.Type = TypeBrowser
	if !parseNameVersion(l, ua) {
		return nil
	}
	if ua.Name != "dillo" || l.s[l.p:] != "" {
		return nil
	}
	return ua
}

func parseIE(l *lex) *UserAgent {
	ua := new()
	ua.Type = TypeBrowser

	if !l.match("mozilla") {
		return nil
	}
	if _, ok := l.span(" ("); !ok {
		return nil
	}
	l.match("compatible; ")
	if !l.match("msie ") {
		return nil
	}
	ua.Name = "ie"
	ua.OS = "windows"

	if !parseVersion(l, ua, ";") {
		return nil
	}

	return ua
}

func parseGoogleBot(l *lex) *UserAgent {
	ua := new()
	ua.Type = TypeCrawler
	if !parseNameVersion(l, ua) {
		return nil
	}
	if ua.Name != "googlebot" {
		return nil
	}
	return ua

}
