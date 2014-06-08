package common

import (
    "strings"
    "fmt"
    "regexp"
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

// Checks whether inputs string looks like UUID
func IsUUID(input string) bool {
    re := regexp.MustCompile("[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}")
    return re.MatchString(input)
}

