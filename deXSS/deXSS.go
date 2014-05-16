package deXSS

import (
    "code.google.com/p/go.net/html"
    "github.com/revel/revel"
    "bytes"
)

// Key of allowed is tag, its value is comma separated list of allowed attributes for that tag
func FilterHTML(h string, allowed map[string]string) {
    meh := bytes.NewBufferString(h)
    doc, err := html.Parse(meh)

    if err != nil {
        revel.ERROR.Printf("Unable to parse HTML: %+v", err)
    }

    revel.ERROR.Printf("doc: %+v\n", doc)


}