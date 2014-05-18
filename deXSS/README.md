deXSS - HTML Stripping for Revel
================================

Options for sanitizing HTML input:
* Escaping everything, e.g. with [html.EscapeString] [1]
* Parse HTML input, and filter the nodes using pre-defined rules

This library does the latter for both tags and attributes.

Usage example
-------------

```go
import (
    "github.com/mikkolehtisalo/revel/deXSS"
    "github.com/revel/revel"
)

var (
    allowed map[string]string
)

func init() {
    allowed = make(map[string]string)
    // This is actually what most basic editing functions of CKEditor require
    allowed["p"] = "class,id"
    allowed["div"] = "class,id"
    allowed["h1"] = "class,id"
    allowed["h2"] = "class,id"
    allowed["h3"] = "class,id"
    allowed["ul"] = "class,id"
    allowed["li"] = "class,id"
    allowed["a"] = "class,id,href,rel"
    allowed["img"] = "class,id,src,alt,hspace,vspace,width,height"
    allowed["span"] = "class,id,style"
}

func blahblah() {
    out := FilterHTML("<p>Hello <a mushroom=\"big\" href=\"/snake\">badger</a>!</p><p>Got it?</p>", allowed, true)
    // The attribute "mushroom" was not in allowed, so it will be gone!
    revel.INFO.Printf("Result of filtering: %+v", out)
}

```


[1]:http://golang.org/pkg/html/#EscapeString