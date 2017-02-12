// Written by https://xojoc.pw. GPLv3 or later.

package useragent

import (
	"net/url"

	"github.com/blang/semver"
)

// Keep them sorted
var browsers = map[string]*url.URL{
	"Chrome":    u("http://www.chromium.org/"),
	"Dillo":     u("http://www.dillo.org/"),
	"Edge":      u("https://www.microsoft.com/en-us/windows/microsoft-edge"),
	"Firefox":   u("https://www.mozilla.org/en-US/firefox"),
	"IceCat":    u("https://www.gnu.org/software/gnuzilla/"),
	"Iceweasel": u("https://wiki.debian.org/Iceweasel"),
	"NetSurf":   u("http://www.netsurf-browser.org/"),
	"Opera":     u("http://www.opera.com/"),
	"PhantomJS": u("http://phantomjs.org/"),
	"Silk":      u("http://aws.amazon.com/documentation/silk/"),
	"WebView":   u("http://developer.android.com/guide/webapps/webview.html"),
}

const (
	OSAndroid = "Android"
	OSMacOS   = "Mac OS X"
	OSiOS     = "iOS"
	OSLinux   = "GNU/Linux"
	OSWindows = "Windows"
)

func parseBrowser(l *lex) *UserAgent {
	for _, f := range []parseFn{parseGecko, parseChromeSafari, parseIE1, parseIE2, parseOperaClassic} {
		if ua := f(newLex(l.s)); ua != nil {
			return ua
		}
	}
	return nil
}

func parseSecurity(l *lex) Security {
	switch {
	case l.match("U"):
		if l.matchFirst("; ", ";") || l.matchNoConsume(")") {
			return SecurityStrong
		}
	case l.match("I"):
		if l.matchFirst("; ", ";") || l.matchNoConsume(")") {
			return SecurityWeak
		}
	case l.match("N"):
		if l.matchFirst("; ", ";") || l.matchNoConsume(")") {
			return SecurityNone
		}
	}
	return SecurityUnknown
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
		ua.OS = OSAndroid
		if l.match("; Mobile") {
			ua.Mobile = true
		} else if l.match("; Tablet") {
			ua.Tablet = true
		}
	case l.match("Linux; "):
		ua.Security = parseSecurity(l)
		if l.match("Android") {
			ua.OS = OSAndroid
		} else {
			return false
		}
	case l.match("Windows"):
		ua.OS = OSWindows
		// Windows has version before security
		_ = parseOSVersion(l, ua)
		l.span("; ")
		ua.Security = parseSecurity(l)
	case l.match("Macintosh"):
		l.span("; ")
		ua.Security = parseSecurity(l)
		ua.OS = OSMacOS
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
		ua.OS = OSiOS
		ua.Tablet = true
	case l.match("iPhone; ") || l.match("iPod; ") || l.match("iPod touch; "):
		ua.Security = parseSecurity(l)
		ua.OS = OSiOS
		ua.Mobile = true
	case l.match("Unknown; "):
		ua.Security = parseSecurity(l)
		parseUnixLike(l, ua)
	default:
		return false
	}

	// OS Version is not required, and may be set above
	if ua.OSVersion.Equals(semver.Version{}) {
		_ = parseOSVersion(l, ua)
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
		ua.OS = OSLinux
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
		// Various distros use "... Distro; Linux x86_64) "
		if _, ok := l.spanBefore("Linux", ") "); ok {
			ua.OS = OSLinux
			return true
		}
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
	if _, ok := l.span("Opera "); ok {
		if !parseVersion(l, ua, " ") {
			return nil
		}
		ua.Name = "Opera"
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

	if l.match("Mobile") && !ua.Tablet {
		ua.Mobile = true
	}
	if ua.OS == OSAndroid && !ua.Mobile {
		ua.Tablet = true
	}

	// Identify non-Chrome browsers with Chromelike UAs:
	if _, ok := l.span("OPR/"); ok {
		if !parseVersion(l, ua, " ") {
			return nil
		}
		ua.Name = "Opera"
	}
	if _, ok := l.span("Edge/"); ok {
		if !parseVersion(l, ua, " ") {
			return nil
		}
		ua.Name = "Edge"
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
	if !parseVersion(l, ua, ";") {
		return nil
	}
	if !l.match(" Windows NT") {
		return nil
	}

	ua.OS = OSWindows
	// swallow the error to preserve backwards compatibility
	_ = parseOSVersion(l, ua)

	if _, ok := l.span("Opera "); ok {
		if !parseVersion(l, ua, " ") {
			return nil
		}
		ua.Name = "Opera"
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
	if _, ok := l.span("(Windows NT"); !ok {
		return nil
	}
	ua.OS = OSWindows

	// swallow the error to preserve backwards compatibility
	_ = parseOSVersion(l, ua)

	if _, ok := l.span("Trident/"); !ok {
		return nil
	}
	if _, ok := l.span("rv:"); !ok {
		return nil
	}
	ua.Name = "MSIE"
	if !parseVersion(l, ua, ")") {
		return nil
	}

	return ua
}

// Non-Mozilla Opera UAs
func parseOperaClassic(l *lex) *UserAgent {
	ua := new()

	if !l.match("Opera/") {
		return nil
	}
	ua.Type = Browser
	ua.Name = "Opera"
	// Start with the Opera version (versions 9.80+ will overwrite later)
	if !parseVersion(l, ua, " ") {
		return nil
	}
	if _, ok := l.span("("); !ok {
		return nil
	}
	switch {
	case l.match("Windows"):
		ua.OS = OSWindows
	case l.match("Macintosh"):
		ua.OS = OSMacOS
		l.spanBefore("OS X", ")")
	default:
		if !parseUnixLike(l, ua) {
			return nil
		}
	}

	// swallow the error to preserve backwards compatibility
	_ = parseOSVersion(l, ua)

	// Get security if present, then get to the end of the parens block
	if _, ok := l.spanBefore("; ", ")"); ok {
		ua.Security = parseSecurity(l)
	}
	if _, ok := l.span(")"); !ok {
		return nil
	} else {
		// Opera occasionally uses nested parens; for simplicity, skip over all instead of matching
		for ok {
			_, ok = l.span(")")
		}
	}
	if l.match(" Presto/") {
		l.span(" ")
	}
	if l.match("Version/") && !parseVersion(l, ua, " ") {
		return nil
	}

	return ua
}
