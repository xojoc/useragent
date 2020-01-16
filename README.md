# User agent parsing
*useragent* is a library written in [golang](http://golang.org) to parse [user agent strings](http://useragentstring.com/).

# Usage
First install the library with:
```
go get xojoc.pw/useragent
```
*useragent* is simple to use. First parse a string with [useragent.Parse](http://godoc.org/xojoc.pw/useragent#Parse) and then access the fields of [useragent.UserAgent](http://godoc.org/xojoc.pw/useragent#UserAgent) for the required information. Example:
 * [Access fields](http://godoc.org/xojoc.pw/useragent#example-Parse--Access)

see [godoc](http://godoc.org/xojoc.pw/useragent) for the complete documentation.
# How it works?
          Lasciate ogne speranza, voi ch'intrate. -Dante
Parsing user agent strings is a hell. There is no standard for user agent strings, so *useragent* must use some heuristics. The site [http://www.useragentstring.com/](http://www.useragentstring.com/pages/useragentstring.php) has been invaluable during development. Some relevant links are also:

  * https://developer.mozilla.org/en-US/docs/Web/HTTP/Gecko_user_agent_string_reference
  * https://developer.chrome.com/multidevice/user-agent
  * https://developer.apple.com/library/safari/documentation/AppleApplications/Reference/SafariWebContent/OptimizingforSafarioniPhone/OptimizingforSafarioniPhone.html#//apple_ref/doc/uid/TP40006517-SW3
  * http://blogs.msdn.com/b/ieinternals/archive/2013/09/21/internet-explorer-11-user-agent-string-ua-string-sniffing-compatibility-with-gecko-webkit.aspx
  * http://docs.aws.amazon.com/silk/latest/developerguide/user-agent.html

for the supported user agents see:
  * browsers: [browser.go](https://github.com/xojoc/useragent/blob/master/browser.go)
  * crawlers: [crawler.go](https://github.com/xojoc/useragent/blob/master/crawler.go)

If you think *useragent* doesn't parse correctly a particular user agent string, just open an issue :).

# Why this library?
*useragent* doesn't just split the user agent string and look for specific strings like other parsers, but it has specific parser for the most common browsers/crawlers and falls back to a generic parser for everything else. Its main features are:

 * Simple and stable API.
 * High precision in detection of the most common browsers/crawlers.
 * Detects mobile/tablet devices.
 * OS detection.
 * URL with more information about the user agent (usually it's the home page).
 * [Security](http://godoc.org/xojoc.pw/useragent#Security) level detection when reported by browsers.


# Who?
*useragent* was written by [Alexandru Cojocaru](https://xojoc.pw) and uses [blang/semver](https://github.com/blang/semver) to parse versions.

Thanks a lot to [@brendanwalters](https://github.com/brendanwalters) (from http://pendo.io) for the contributions.

# [Donate!](https://xojoc.pw/donate)


# License
*useragent* is released under the GPLv3 or later, see [COPYING](https://github.com/xojoc/useragent/blob/master/COPYING).
