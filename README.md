# User agent parsing
*useragent* is a library written in [golang](http://golang.org) to parse [user agent strings](http://useragentstring.com/).

# Usage
First install the library with:
```
go get github.com/xojoc/useragent       
```
*useragent* is simple to use. First parse a string with [useragent.Parse](http://godoc.org/github.com/xojoc/useragent#Parse) and then access the fields of [useragent.UserAgent](http://godoc.org/github.com/xojoc/useragent#UserAgent) for the required information. Example:
 * [Access fields](http://godoc.org/github.com/xojoc/useragent#example-Parse--Access)

see [godoc](http://godoc.org/github.com/xojoc/useragent) for the complete documentation.
# How it works

There is no standard for user agent strings, so *useragent* must use some heuristics. The site [http://www.useragentstring.com/](http://www.useragentstring.com/pages/useragentstring.php) has been invaluable during development. This parser so far recognizes:
 * Firefox and derivatives (IceCat, IceWeasel, etc.)
 * Dillo
 * Chrome
 * GoogleBot

**More is coming...**
If you need support for a particular user agent just open an issue :).

# Who?
*useragent* was written by Alexandru Cojocaru (http://xojoc.pw), [blang/semver](https://github.com/blang/semver) is used to parse versions.

# License
*useragent* is released under the GPLv3 or later, see [COPYING](COPYING).