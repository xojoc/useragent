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
	"net/url"
	"strings"
)

type Type int

const (
	Unknown Type = iota
	Browser
	Crawler
	LinkChecker
	Validator
	FeedReader
	Library
)

func (a Type) String() string {
	switch a {
	case Unknown:
		return "Unkonwn Agent type"
	case Browser:
		return "Browser"
	case Crawler:
		return "Crawler"
	case LinkChecker:
		return "Link Checker"
	case Validator:
		return "Validator"
	case FeedReader:
		return "Feed Reader"
	case Library:
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
	Type     Type
	// The browser/crawler/etc. name. For example:
	//  Firefox, IceCat, Iceweasel
	//  Dillo
	//  Chrome
	//  MSIE
	//  Googlebot
	//   etc.
	// If the name is not known, Name will be `unknown'.
	Name    string
	Version semver.Version
	// The OS name. Can be one of:
	//  GNU/Linux
	//  FreeBSD
	//  OpenBSD
	//  NetBSD
	//  Windows
	//  Mac OS X
	//  Android
	//  Firefox OS
	//  CrOS
	//   etc.
	// If the os is not known, OS will be `unknown'.
	OS       string
	Security Security
	// URL with more information about the user agent (in most cases it's the home page).
	// If unknown is nil.
	URL *url.URL
	// Is it a phone device?
	Mobile bool
	// Is it a tablet device?
	Tablet bool
}

func (ua *UserAgent) String() string {
	return fmt.Sprintf(`Type: %v
Name: %v
Version: %v
OS: %v
Security: %v
Mobile: %v
Tablet: %v`, ua.Type, ua.Name, ua.Version, ua.OS, ua.Security, ua.Mobile, ua.Tablet)
}

func new() *UserAgent {
	ua := &UserAgent{}
	ua.Name = "unknown"
	ua.OS = "unknown"
	return ua
}

func u(rawurl string) *url.URL {
	url, err := url.Parse(rawurl)
	if err != nil {
		panic("useragent: Parse(" + rawurl + "): " + err.Error())
	}
	return url
}

type parseFn func(l *lex) *UserAgent

// Try to extract information about an user agent from uas.
// Since user agent strings don't have a standard, this function uses heuristics.
func Parse(uas string) *UserAgent {
	for _, f := range []parseFn{parseBrowser, parseGeneric} {
		if ua := f(newLex(uas)); ua != nil {
			ua.Original = uas
			return ua
		}
	}
	return nil
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
	ua.Name = s

	return parseVersion(l, ua, " ")
}

func parseGeneric(l *lex) *UserAgent {
	ua := new()
	if !parseNameVersion(l, ua) {
		return nil
	}
	if url, ok := browsers[ua.Name]; ok {
		ua.Type = Browser
		ua.URL = url
		return ua
	}

	if url, ok := crawlers[ua.Name]; ok {
		ua.Type = Crawler
		ua.URL = url
		return ua
	}

	return nil
}
