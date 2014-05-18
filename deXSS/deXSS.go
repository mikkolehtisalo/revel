package deXSS

import (
    "code.google.com/p/go.net/html"
    "github.com/revel/revel"
    "bytes"
    "regexp"
    "strings"
)

var (
    stripped_content *regexp.Regexp = regexp.MustCompile("<html><head></head><body>(.*)</body></html>")
)

// Strips the html/head/body tags
func strip_html(in string) string {
    out := in
    if stripped_content.MatchString(in) {  
        out = stripped_content.FindStringSubmatch(in)[1]

    }
    return out
}

// Render the document or node to string
func get_html(n *html.Node) string {
    tmp := bytes.Buffer{}
    html.Render(&tmp, n)
    return tmp.String()
}

// Check whether node can be found from the allowed map
func is_legal_node(n *html.Node, allowed map[string]string) bool {
    legal := false
    _, present := allowed[n.Data]
    if present {
        legal = true
    }
    return legal
}

// Remove node from document
func remove_node(n *html.Node) {
    revel.TRACE.Printf("remove_node() %+v", n)
    var parent *html.Node
    // Nil if root node, but probably don't need to handle that
    parent = n.Parent
    parent.RemoveChild(n)
}

func filter_attributes(n *html.Node, allowed map[string]string) {
    revel.TRACE.Printf("filter_attributes() %+v, %+v", n, allowed)
    result := []html.Attribute{}
    // Loop all the node's attributes against the allowed list for this node
    for _, att := range n.Attr {
        legal := false
        for _, okay := range strings.Split(allowed[n.Data],",") {
            if att.Key == okay {
                legal = true
            }
        }
        // The attribute seems ok, add it to results
        if legal {
            result = append(result, att)
        }
    }
    revel.TRACE.Printf("filter_attributes() result %+v", result)
    n.Attr = result

}

// Filters HTML, returns filtered version. Please note that go.net/htmls parsing might change many minor things in document.
// Key of allowed is tag, its value is comma separated list of allowed attributes for that tag.
// If stip is set, removes the head/body/html tags that html.Parse always ensures in results.
func FilterHTML(h string, allowed map[string]string, strip bool) string {
    revel.TRACE.Printf("FilterHTML() %+v, %+v", h, allowed)

    // Make sure allowed contains html/head/body, since the go.net/html always adds them to the parsed document tree!
    if _, present := allowed["html"]; !present {
        allowed["html"] = ""
    }
    if _, present := allowed["head"]; !present {
        allowed["head"] = ""
    }
    if _, present := allowed["body"]; !present {
        allowed["body"] = ""
    }

    meh := bytes.NewBufferString(h)
    doc, err := html.Parse(meh)
    if err != nil {
        revel.ERROR.Printf("Unable to parse HTML: %+v", err)
    }

    var f func(*html.Node)
    f = func(n *html.Node) {
        if n.Type == html.ElementNode {
            // Filter only ElementNodes
            if !is_legal_node(n, allowed) {
                remove_node(n)
            } else {
                // Still have to filter Attributes
                filter_attributes(n, allowed)
            }
        }
        for c := n.FirstChild; c != nil; c = c.NextSibling {
            f(c)
        }
    }

    f(doc)

    // Strip html/head/body?
    var result string
    if strip {
        result = strip_html(get_html(doc))
    } else {
        result = get_html(doc)
    }
    revel.TRACE.Printf("FilterHTML() returning %+v", result)

    return result
}