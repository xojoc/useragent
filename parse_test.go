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
	"testing"
)

func TestFirefoxLike(t *testing.T) {
	var got *UserAgent
	want := &UserAgent{}

	got = Parse(`Mozilla/5.0 (X11; Linux i686; rv:38.0) Gecko/20100101 Firefox/38.0`)
	want.Type = TypeBrowser
	want.OS = "gnu/linux"
	want.Name = "firefox"
	want.Version = Version{38, 0}
	if *want != *got {
		t.Errorf("expected %+v, got %+v\n", want, got)
	}

	got = Parse(`Mozilla/5.0 (X11; U; Linux x86_64; sv-SE; rv:1.9.1.16) Gecko/20120714 Iceweasel/3.5.16 (like Firefox/3.5.16)`)
	want.Type = TypeBrowser
	want.OS = "gnu/linux"
	want.Name = "iceweasel"
	want.Version = Version{3, 5}
	if *want != *got {
		t.Errorf("expected %+v, got %+v\n", want, got)
	}

	got = Parse(`Mozilla/5.0 (Windows x86; rv:19.0) Gecko/20100101 Firefox/19.0`)
	want.Type = TypeBrowser
	want.OS = "windows"
	want.Name = "firefox"
	want.Version = Version{19, 0}
	if *want != *got {
		t.Errorf("expected %+v, got %+v\n", want, got)
	}
}
