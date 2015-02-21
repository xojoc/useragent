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
)

type lex struct {
	s string
	p int
}

func newLex(s string) *lex {
	return &lex{s, 0}
}

func (l *lex) match(m string) bool {
	if !strings.HasPrefix(l.s[l.p:], m) {
		return false
	}

	l.p += len(m)
	return true
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
	i := strings.IndexAny(l.s[l.p:], chars)
	if i < 0 {
		return "", false
	}
	s := l.s[l.p : l.p+i]
	l.p += i + len(chars)
	return s, true
}
