// Written by https://xojoc.pw. GPLv3 or later.

package useragent

import (
	"net/url"
)

// Keep them sorted
var crawlers = map[string]*url.URL{
	"Google AdsBot":    u("https://support.google.com/webmasters/answer/1061943"),
	"Google AdSense":   u("https://support.google.com/webmasters/answer/1061943"),
	"Googlebot":        u("http://www.google.com/bot.html"),
	"Googlebot Images": u("https://support.google.com/webmasters/answer/1061943"),
	"Googlebot News":   u("https://support.google.com/news/publisher/answer/93977"),
	"Googlebot Video":  u("https://support.google.com/webmasters/answer/1061943"),
}

func parseCrawler(l *lex) *UserAgent {
	for _, f := range []parseFn{parseGooglebot, parseGooglebotSmartphone} {
		if ua := f(newLex(l.s)); ua != nil {
			return ua
		}
	}
	return nil
}

func parseGooglebot(l *lex) *UserAgent {
	ua := new()
	ua.Type = Crawler

	// Alternate Googlebots
	if l.match("Googlebot") {
		if l.match("-News") {
			ua.Name = "Googlebot News"
		} else if parseNameVersion(l, ua) {
			switch ua.Name {
			case "":
				ua.Name = "Googlebot"
			case "-Image":
				ua.Name = "Googlebot Images"
			default:
				ua.Name = "Googlebot " + ua.Name[1:]
			}
		} else {
			return nil
		}
		return ua
	} else if l.match("Mediapartners-Google") {
		ua.Name = "Google AdSense"
		return ua
	} else if l.match("AdsBot-Google") {
		ua.Name = "Google AdsBot"
		return ua
	}

	// Googlebot
	if !l.match("Mozilla/5.0 (compatible; Googlebot/") {
		return nil
	}

	ua.Name = "Googlebot"

	if !parseVersion(l, ua, ";") {
		return nil
	}
	if !l.match(" +http://www.google.com/bot.html)") {
		return nil
	}
	return ua
}

func parseGooglebotSmartphone(l *lex) *UserAgent {
	ua := new()

	if _, ok := l.span("Mozilla"); !ok {
		return nil
	}

	if _, ok := l.span("Linux"); !ok {
		return nil
	}

	if _, ok := l.span("Android"); !ok {
		return nil
	}

	if _, ok := l.span("AppleWebKit"); !ok {
		return nil
	}

	if _, ok := l.span("Chrome"); !ok {
		return nil
	}

	if _, ok := l.span("Mobile Safari"); !ok {
		return nil
	}

	if _, ok := l.span("Googlebot/"); !ok {
		return nil
	}

	if !parseVersion(l, ua, ";") {
		return nil
	}
	ua.Type = Crawler
	ua.Name = "Googlebot"
	ua.Mobile = true
	return ua
}
