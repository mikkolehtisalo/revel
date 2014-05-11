package acl

import (
    "github.com/revel/revel/cache"
    "github.com/revel/revel"
    "time"
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
/*

c.Args["user_details"]
() {
    
}
*/