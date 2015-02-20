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
	"strconv"
	"strings"
)

type AgentType int

const (
	TypeBrowser AgentType = iota
	TypeCrawler
	TypeLinkChecker
	TypeValidator
	TypeFeedReader
	TypeLibrary
)

type Version struct {
	Major uint64
	Minor uint64
}

type UserAgent struct {
	Type    AgentType
	OS      string
	Name    string
	Version Version
}

type parseFn func(l *lex) *UserAgent

func Parse(uas string) *UserAgent {
	// we try each user agent parser in order until we get one that succeeds
	for _, f := range []parseFn{parseFirefoxLike} {
		if ua := f(newLex(uas)); ua != nil {
			return ua
		}
	}
	return nil
}

func parseFirefoxLike(l *lex) *UserAgent {
	var err error
	var s string
	var ok bool
	ua := &UserAgent{}
	ua.Type = TypeBrowser

	if !l.match("Mozilla/5.0 (") {
		return nil
	}

	s, ok = l.span(";")
	if !ok {
		return nil
	}

	switch {
	case s == "X11":
		l.match(" ")
		l.match("N; ")
		l.match("U; ")
		l.match("I; ")
		switch {
		case l.match("Linux") || l.match("Ubuntu"):
			ua.OS = "gnu/linux"
		case l.match("OpenBSD"):
			ua.OS = "openbsd"
		default:
			return nil
		}
	case strings.HasPrefix(s, "Windows"):
		ua.OS = "windows"
	case s == "Macintosh":
		ua.OS = "mac os"
	default:
		return nil
	}

	if _, ok = l.span(") "); !ok {
		return nil
	}

	// skip 'gecko/version ' field
	if _, ok = l.span(" "); !ok {
		return nil
	}

	s, ok = l.span("/")
	if !ok {
		return nil
	}
	ua.Name = strings.ToLower(s)

	if s, ok = l.span("."); !ok {
		return nil
	}
	ua.Version.Major, err = strconv.ParseUint(s, 10, 64)
	if err != nil {
		return nil
	}

	if s, ok = l.span("."); !ok {
		s = l.s[l.p:]
		if s == "" {
			return nil
		}
	}
	ua.Version.Minor, err = strconv.ParseUint(s, 10, 64)
	if err != nil {
		return nil
	}

	return ua
}
