package common

import (
    "strings"
    "fmt"
    //"github.com/revel/revel"
)

// Checks whether string can be found from slice
func StringInSlice(a string, list []string) bool {
    for _, b := range list {
        if b == a {
            return true
        }
    }
    return false
}

func comma(r rune) bool {
    return r == ','
}

func AddUserToACLList(user string, acl *string) {
    ustr := fmt.Sprintf("u:%s", user)
    a := strings.FieldsFunc(*acl, comma)

    if !StringInSlice(ustr, a) {
        a = append(a, ustr)
        astr := strings.Join(a, ",")
        *acl = astr
    }
}

