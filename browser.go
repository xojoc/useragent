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
	"net/url"
)

// keep them sorted
var browsers = map[string]*url.URL{
	"Chrome":    u("http://www.chromium.org/"),
	"Dillo":     u("http://www.dillo.org/"),
	"Firefox":   u("https://www.mozilla.org/en-US/firefox"),
	"IceCat":    u("https://www.gnu.org/software/gnuzilla/"),
	"Iceweasel": u("https://wiki.debian.org/Iceweasel"),
	"NetSurf":   u("http://www.netsurf-browser.org/"),
	"PhantomJS": u("http://phantomjs.org/"),
	"Silk":      u("http://aws.amazon.com/documentation/silk/"),
	"WebView":   u("http://developer.android.com/guide/webapps/webview.html"),
}

func parseBrowser(l *lex) *UserAgent {
	for _, f := range []parseFn{parseGecko, parseChromeSafari, parseIE1, parseIE2} {
		if ua := f(newLex(l.s)); ua != nil {
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
	ua.Type = Browser

	if !l.match("Mozilla/5.0 (") {
		return false
	}

	switch {
	case l.match("X11; "):
		ua.Security = parseSecurity(l)
		parseUnixLike(l, ua)
	case l.match("Android"):
		ua.Security = parseSecurity(l)
		ua.OS = "Android"
		if l.match("; Mobile") {
			ua.Mobile = true
		} else if l.match("; Tablet") {
			ua.Tablet = true
		}
	case l.match("Linux; "):
		ua.Security = parseSecurity(l)
		if l.match("Android") {
			ua.OS = "Android"
		} else {
			return false
		}
	case l.match("Windows"):
		ua.Security = parseSecurity(l)
		ua.OS = "Windows"
	case l.match("Macintosh"):
		ua.Security = parseSecurity(l)
		ua.OS = "Mac OS X"
	case l.match("Mobile; "):
		ua.Security = parseSecurity(l)
		ua.OS = "Firefox OS"
		ua.Mobile = true
	case l.match("Tablet; "):
		ua.Security = parseSecurity(l)
		ua.OS = "Firefox OS"
		ua.Tablet = true
	case l.match("iPad; "):
		ua.Security = parseSecurity(l)
		ua.OS = "iOS"
		ua.Tablet = true
	case l.match("iPhone; ") || l.match("iPod; ") || l.match("iPod touch; "):
		ua.Security = parseSecurity(l)
		ua.OS = "iOS"
		ua.Mobile = true
	case l.match("Unknown; "):
		ua.Security = parseSecurity(l)
		parseUnixLike(l, ua)
	default:
		return false
	}

	if _, ok := l.span(") "); !ok {
		return false
	}

	return true
}

// Parse *nix variants (eg inside of a MozillaLike)
func parseUnixLike(l *lex, ua *UserAgent) bool {
	switch {
	case l.match("Linux") || l.match("Ubuntu"):
		ua.OS = "GNU/Linux"
	case l.match("FreeBSD"):
		ua.OS = "FreeBSD"
	case l.match("OpenBSD"):
		ua.OS = "OpenBSD"
	case l.match("NetBSD"):
		ua.OS = "NetBSD"
	case l.match("Maemo"):
		// FIXME: should it be GNU/Linux?
		ua.OS = "Maemo"
	case l.match("CrOS"):
		ua.OS = "CrOS"
	default:
		return false
	}
	return true
}

// https://developer.mozilla.org/en-US/docs/Web/HTTP/Gecko_user_agent_string_reference
func parseGecko(l *lex) *UserAgent {
	ua := new()

	if !parseMozillaLike(l, ua) {
		return nil
	}
	if !l.match("Gecko/") {
		return nil
	}
	if _, ok := l.span(" "); !ok {
		return nil
	}
	if !parseNameVersion(l, ua) {
		return nil
	}

	return ua
}

// Includes WebKit-based Firefox for iOS
func parseChromeSafari(l *lex) *UserAgent {
	ua := new()

	if !parseMozillaLike(l, ua) {
		return nil
	}
	if !l.match("AppleWebKit/") {
		return nil
	}
	if _, ok := l.span(" "); !ok {
		return nil
	}
	if !l.match("(KHTML, like Gecko) ") {
		return nil
	}
	if !parseNameVersion(l, ua) {
		return nil
	}
	if ua.Name == "CriOS" {
		ua.Name = "Chrome"
	} else if ua.Name == "FxiOS" {
		ua.Name = "Firefox"
	} else if ua.Name == "Version" {
		if l.match("Chrome/") {
			if !parseVersion(l, ua, " ") {
				return nil
			}
			ua.Name = "WebView"
			ua.Type = Library
		} else {
			if l.match("Mobile/") {
				if _, ok := l.span(" "); !ok {
					return nil
				}
			}
			if !l.match("Safari/") {
				return nil
			}
			ua.Name = "Safari"
		}
	} else if ua.Name == "Silk" {
		if l.match("like Chrome/") {
			if _, ok := l.span(" "); !ok {
				return nil
			}
		} else {
			return nil
		}
	}
	if ua.OS == "Android" {
		if l.match("Mobile") {
			ua.Mobile = true
		} else {
			ua.Tablet = true
		}
	}

	return ua
}

// pre IE11 uas
func parseIE1(l *lex) *UserAgent {
	ua := new()

	ua.Type = Browser
	if !l.match("Mozilla") {
		return nil
	}
	if _, ok := l.span(" ("); !ok {
		return nil
	}
	l.match("compatible; ")
	l.match("Compatible; ")
	if !l.match("MSIE ") {
		return nil
	}
	ua.Name = "MSIE"
	ua.OS = "Windows"
	if !parseVersion(l, ua, ";") {
		return nil
	}

	return ua
}

// IE11 changed its uas http://blogs.msdn.com/b/ieinternals/archive/2013/09/21/internet-explorer-11-user-agent-string-ua-string-sniffing-compatibility-with-gecko-webkit.aspx
func parseIE2(l *lex) *UserAgent {
	ua := new()

	ua.Type = Browser
	if !l.match("Mozilla") {
		return nil
	}
	if _, ok := l.span("Trident/"); !ok {
		return nil
	}
	if _, ok := l.span("rv:"); !ok {
		return nil
	}
	ua.Name = "MSIE"
	ua.OS = "Windows"
	if !parseVersion(l, ua, ")") {
		return nil
	}

	return ua
}
