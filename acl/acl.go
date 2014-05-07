package acl

import (
    "github.com/revel/revel/cache"
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
    // Parent for calculating inheritation, eg. "page:level2page" or nil
    Parent string
}

type ACL struct {
    // The permission, eg. "read", "write", "admin"
    Permission string
    // Who has the permission, eg. "u:michael", "g:administrators"
    Principal string
}

func SetEntry(a ACLEntry) {
     go cache.Set(ACL_ENTRY_ID + a.ObjReference, a, defaultExpiration)
}

func GetEntry(reference string) ACLEntry {
    a := ACLEntry{}
    if err := cache.Get(ACL_ENTRY_ID + reference, &a); err != nil {
        // ERROR
    }
    return a
}