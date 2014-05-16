package acl

import (
    "github.com/revel/revel/cache"
    "github.com/revel/revel"
    "time"
    "reflect"
    "github.com/mikkolehtisalo/revel/ldapuserdetails"
    . "github.com/mikkolehtisalo/revel/common"
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
        revel.TRACE.Println("Unable to get ACL entry %s", reference)
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
    SetMatched(permissions []string) interface{}
    BuildACLReference() string
    BuildACLEntry(reference string) ACLEntry
}

// Takes any interface{} and attempts to convert it to []Filterable
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

// Get ACL entry from cache. If not available, build new (and set in cache)
func GetACLEntry(reference string, item Filterable) ACLEntry {
    entry := GetEntry(reference)

    if len(entry.ObjReference)==0 {
        // Wasn't found from cache, build new!
        entry = item.BuildACLEntry(reference)
        // Save to cache
        SetEntry(entry)
    }

    return entry;
}

// Filters items based on user's roles. Returns those that match any of the listed permissions.
func Filter(c map[string]interface {}, permissions []string, i interface{}) []Filterable {
    revel.TRACE.Printf("Filter(): %s : %+v", permissions, i)
    result := []Filterable{}

    // Get the items
    items := get_filterable(i)
    revel.TRACE.Printf("Items: %+v", items)

    // Get roles for the user
    dets := c["user_details"].(ldapuserdetails.User_details)
    roles := dets.Roles
    revel.TRACE.Printf("Roles: %+v", roles)

    // Compare all items against all ACLs, add matches to result
    for _, item := range items {
        ref := item.BuildACLReference()
        //aclentry := item.GetACLEntry(ref)
        aclentry := GetACLEntry(ref, item)
        matched := []string{}
        for _, acl := range aclentry.ACLs {
            if StringInSlice(acl.Permission, permissions) && StringInSlice(acl.Principal, roles) {
                matched = append(matched, acl.Permission)
            }
        }
        if len(matched) > 0 {
            // Set the matched, produces copy so get it back...
            item2 := item.SetMatched(matched).(Filterable)
            result = append(result, item2)
        }
    }

    revel.TRACE.Printf("Filter returning: %+v", result)
    return result
}

/*



*/