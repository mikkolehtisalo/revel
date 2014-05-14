ACLs for Revel
==============

This will be a simple ACL system for Revel. Do not use, very early phase. 

Introduction
------------

TBD.

Usage
-----

```go
import "github.com/mikkolehtisalo/revel/acl"

// Imaginary Model
type Opinion struct {
    Uuid string
    Message string

}

// Must be implemented, gets ACL entry for arbitrary item from cache
func (o Opinion) GetACLEntry(reference string) acl.ACLEntry {
    entry := acl.GetEntry(reference)

    if len(entry.ObjReference)==0 {
        // Wasn't found from cache, build new!
        entry = o.BuildACLEntry(reference)
        // Save to cache
        acl.SetEntry(entry)
    }

    return entry;
}

// Must be implemented, builds ACL entry if one can not be fetched from cache
func (w Wiki) BuildACLEntry(reference string) acl.ACLEntry {
    entry := acl.ACLEntry{}
    // Should probably be built from some data from database etc
    acls := acl.BuildPermissionACLs("criticize", []string{"u:mikkolehtisalo"})
    entry.ObjReference = reference
    entry.ACLs = acls
    entry.Inheritation = w.BuildACLInheritation()
    entry.Parent = w.BuildACLParent()
    return entry
}

// Generate a cache key that will distinguish different types of items
func (w Wiki) BuildACLReference() string {
    return "opinion:"+w.Wiki_id
}

// Inheritation in use?
func (w Wiki) BuildACLInheritation() bool {
    return false
}

// Reference to parent - used with inheritation
func (w Wiki) BuildACLParent() string {
    return ""
}

// Imaginary Controller
type Opinions struct {
    *revel.Controller
}

// Method in Controller
func (c Opinions) List revel.Result {
    // Gets []Opinion from database
    opins := db.GetManyOpinions() 
    // Logged on user must have permission to criticize in order to see an item!
    filtered := acl.Filter(c.Args, []string{"criticize"}, opins)
    // Return the filtered list
    return c.RenderJson(filtered)
}

```

Gotchas
-------

acl.Filter will return []Filterable. If the slice contents will be used for anything else besides printing, asserting type will be probably needed. For example:

```go
    filtered := acl.Filter(c.Args, []string{"criticize"}, opins)
    // Naive handling of the slice
    op := filtered[0].(Opinion)
    // After previous the item will handle like Opinion
```