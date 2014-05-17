package deXSS

import "testing"

var (
    allowed map[string]string
)

func init() {
    allowed := make(map[string]string)
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

func TestEmpty(t *testing.T) {
    out := FilterHTML("", allowed)
    if out != "" {
        t.Error("Expected empty string, got ", out)
    }
}

func TestText(t *testing.T) {
    out := FilterHTML("The grass is always greener on the other side of the force", allowed)
    if out != "The grass is always greener on the other side of the force" {
        t.Error("Expected sample string, got ", out)
    }
}

func TestLegalTags(t *testing.T) {
    out := FilterHTML("<p>Hello <a href=\"/snake\">badger</a>!</p><p>Got it?</p>", allowed)
    if out != "<p>Hello <a href=\"/snake\">badger</a>!</p><p>Got it?</p>" {
        t.Error("Expected back the same string, got ", out)
    }
}

func TestIllegalTag(t *testing.T) {
    out := FilterHTML("And how <script>alert(\"Surprise!\");</script> it works!", allowed)
    if out != "And how it works!" {
        t.Error("Expected tag cleanly gone, got ", out)
    }
}

func TestIllegalAttribute(t *testing.T) {
    out := FilterHTML("<p>Hello <a mushroom=\"big\" href=\"/snake\">badger</a>!</p><p>Got it?</p>", allowed)
    if out != "<p>Hello <a href=\"/snake\">badger</a>!</p><p>Got it?</p>" {
        t.Error("Expected back the string without mushrooms, got ", out)
    }
}

func TestIllegalNestedTag(t *testing.T) {
    out := FilterHTML("<p>And how <script>alert(\"Surprise!\");</script> it works!</p>", allowed)
        if out != "<p>And how it works!</p>" {
            t.Error("Expected tag cleanly gone, got ", out)
        }
    }
