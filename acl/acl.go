package acl

import (
    "github.com/revel/revel/cache"
    "github.com/revel/revel"
    "time"
    "reflect"
    "github.com/mikkolehtisalo/revel/ldapuserdetails"
)

const (
    ACL_ENTRY_ID = "ACL_ENTRY"
    defaultExpiration = 8 * time.Hour
)

// Default is always DENY
// This is not interface{} or anything similar to prevent deep cascades

type ACLEntry struct {
    // Parseable reference to the object ACL entry belongs to, eg. "wiki:example1"
    ObjReference string
    // All the defined ACLs for object
    ACLs []ACL
    // Should be use inheritation with ACLs?
    Inheritation bool
    // Parent for calculating inheritation, eg. "page:level2page" or ""
    Parent string
}

type ACL struct {
    // The permission, eg. "read", "write", "admin"
    Permission string
    // Who has the permission, eg. "u:michael", "g:administrators"
    Principal string
}

func BuildPermissionACLs(permission string, principals []string) []ACL {
    a := []ACL{}
    for _, principal := range principals {
        i := ACL{}
        i.Permission = permission
        i.Principal = principal
        a = append(a, i)
    }
    return a
}

func SetEntry(a ACLEntry) {
     go cache.Set(ACL_ENTRY_ID + a.ObjReference, a, defaultExpiration)
}

func GetEntry(reference string) ACLEntry {
    a := ACLEntry{}
    if err := cache.Get(ACL_ENTRY_ID + reference, &a); err != nil {
        revel.ERROR.Println("Unable to get ACL entry %s", reference)
    }
    return a
}

// TODO: inheritation!
func GetPermissions(principals []string, acl ACLEntry) map[string]bool {
    permissions := make(map[string]bool)

    for _, entry := range acl.ACLs {
        for _, principal := range principals {
            if entry.Principal == principal {
                permissions[entry.Permission] = true
            }
        }
    }

    return permissions
}

type Filterable interface {
    BuildACLReference() string
    GetACLEntry(reference string) ACLEntry
}

func takeSliceArg(arg interface{}) (out []interface{}, ok bool) {
    slice, success := takeArg(arg, reflect.Slice)
    if !success {
        ok = false
          return
    }
    c := slice.Len()
    out = make([]interface{}, c)
    for i := 0; i < c; i++ {
        out[i] = slice.Index(i).Interface()
    }
    return out, true
}

func takeArg(arg interface{}, kind reflect.Kind) (val reflect.Value, ok bool) {
    val = reflect.ValueOf(arg)
    if val.Kind() == kind {
        ok = true
    }
    return
}

// Takes any interface{} and attempt to convert it to []Filterable
func get_filterable (items interface{}) []Filterable {
    slice := reflect.ValueOf(items)
    if slice.Kind() != reflect.Slice {
        // Panic?
    }

    co := slice.Len()
    filterableslice := make([]Filterable, co)
    for i := 0; i < co; i++ {
           filterableslice[i] = slice.Index(i).Interface().(Filterable)
    }

    return filterableslice
}

func Filter(c map[string]interface {}, permission string, i interface{}) {
    // Get the items
    items := get_filterable(i)
    revel.INFO.Printf("Items: %+v", items)

    // Get roles for the user
    dets := c["user_details"].(ldapuserdetails.User_details)
    roles := dets.Roles
    revel.INFO.Printf("Roles: %+v", roles)

    // Get the ACL for item
    for _, item := range items {
        ref := item.BuildACLReference()
        revel.INFO.Printf("Reference: %+v", ref)
        acl := item.GetACLEntry(ref)
        revel.INFO.Printf("ACL entry: %+v", acl)
    }
}

/*


c.Args["user_details"]
() {
    
}
*/