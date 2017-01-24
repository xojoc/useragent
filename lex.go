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
	"strings"
	"unicode/utf8"
)

type lex struct {
	s string
	p int
}

func newLex(s string) *lex {
	return &lex{s, 0}
}

// Returns true iff current position matches input string
func (l *lex) matchNoConsume(m string) bool {
	return strings.HasPrefix(l.s[l.p:], m)
}

// If current position matches input string, consumes it and returns true (else false)
func (l *lex) match(m string) bool {
	if !l.matchNoConsume(m) {
		return false
	}

	l.p += len(m)
	return true
}

// Consumes first provided string that matches the current position;
// returns true on finding any match, false otherwise
func (l *lex) matchFirst(args ...string) bool {
	for _, m := range args {
		if l.match(m) {
			return true
		}
	}
	return false
}

func (l *lex) span(m string) (string, bool) {
	i := strings.Index(l.s[l.p:], m)
	if i < 0 {
		return "", false
	}
	s := l.s[l.p : l.p+i]
	l.p += i + len(m)
	return s, true
}

func (l *lex) spanAny(chars string) (string, bool) {
	// should this whole function loop-consume until char doesn't match?
	i := strings.IndexAny(l.s[l.p:], chars)
	if i < 0 {
		return "", false
	}
	s := l.s[l.p : l.p+i]
	_, matchWidth := utf8.DecodeRuneInString(l.s[l.p+i:])
	l.p += i + matchWidth
	return s, true
}

func (l *lex) spanBefore(m, stopAt string) (string, bool) {
	i := strings.Index(l.s[l.p:], m)
	if i < 0 {
		return "", false
	}
	j := strings.Index(l.s[l.p:], stopAt)
	if j >= 0 && j < i {
		return "", false
	}
	s := l.s[l.p : l.p+i]
	l.p += i + len(m)
	return s, true
}
