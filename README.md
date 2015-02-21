# User agent parsing
*useragent* is a library written in [golang](http://golang.org) to parse [user agent strings](http://useragentstring.com/).

# Usage

```
go get github.com/xojoc/useragent       
```

```
package main

import (
	"github.com/xojoc/useragent"
	"log"
)

func main() {
	ua := useragent.Parse("Mozilla/5.0 (X11; Linux i686; rv:38.0) Gecko/20100101 Firefox/38.0")
	if ua == nil {
		log.Fatal("cannot parse user agent string")
	}
	log.Println(ua.OS)
	log.Println(ua.Name)
	log.Println(ua.Version)
}

```

# How it works

There is no standard for user agent strings, so *useragent* must use some heuristics. The site [http://www.useragentstring.com/](http://www.useragentstring.com/pages/useragentstring.php) has been invaluable during development. This parser so far recognizes:
 * Firefox and derivatives (IceCat, IceWeasel, etc.)
 * Dillo
 * Chrome
 * GoogleBot

More is coming...
If you need support for a particular user agent just open an issue :).

# Who?
*useragent* was written by Alexandru Cojocaru ([http://xojoc.pw](http://xojoc.pw)), [blang/semver](https://github.com/blang/semver) is used to parse versions.