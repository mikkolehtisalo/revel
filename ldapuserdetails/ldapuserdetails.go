package ldapuserdetails

import (
    "github.com/revel/revel/cache"
    "github.com/mikkolehtisalo/ldap"
    "github.com/revel/revel"
    "fmt"
    "strings"
)

var (
    ldap_server string   = "localhost"
    ldap_port   uint16   = 389
    ldap_user       string   = "*"
    ldap_passwd     string   = "*"
    ldap_user_base     string   = "dc=*,dc=*"
    ldap_user_filter     string   = "&(objectClass=*)"
    ldap_user_uid_attr string = "*"
    ldap_user_cn_attr string = "*"
    ldap_user_photo_attr string = "*"
    ldap_user_group_attr string = "*"
    ldap_group_filter string = "&(objectClass=*)"
    ldap_group_base string = "dc=*,dc=*"
    ldap_group_cn_attr string = "*"
    ldap_group_dn_attr string = "*"
)

// Struct for holding details about user
type User_details struct {
    username string
    visiblename string
    photo []byte
    groups []string
    roles []string
}

func get_c_str(name string) string {
    if tmp, ok := revel.Config.String(name); !ok {
        panic(fmt.Errorf("%s invalid", name))
    } else {
        return tmp
    }
}

func get_c_uint16(name string) uint16 {
    if tmp, ok := revel.Config.Int(name); !ok {
        panic(fmt.Errorf("%s invalid", name))
    } else {
        return uint16(tmp)
    }
}

func init() {
    revel.OnAppStart(func() {
        ldap_server = get_c_str("ldap.server")
        ldap_port = get_c_uint16("ldap.port")
        ldap_user = get_c_str("ldap.user")
        ldap_passwd = get_c_str("ldap.passwd")
        ldap_user_base = get_c_str("ldap.user_base")
        ldap_user_filter = get_c_str("ldap.user_filter")
        ldap_user_uid_attr = get_c_str("ldap.user_uid_attr")
        ldap_user_cn_attr = get_c_str("ldap.user_cn_attr")
        ldap_user_photo_attr = get_c_str("ldap.user_photo_attr")
        ldap_user_group_attr = get_c_str("ldap.user_group_attr")
        ldap_group_filter = get_c_str("ldap.group_filter")
        ldap_group_base = get_c_str("ldap.group_base")
        ldap_group_cn_attr = get_c_str("ldap.group_cn_attr")
        ldap_group_dn_attr = get_c_str("ldap.group_dn_attr")
    })
}

// Gets new LDAP connection. Must be Close()d later!
func Get_connection() *ldap.Conn {
    revel.TRACE.Printf("Get_connection(): %s:%s, %s/%s", ldap_server, ldap_port, ldap_user, ldap_passwd)
    l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", ldap_server, ldap_port))
    if err != nil {
        panic(fmt.Errorf("ldap.Dial panic: %s",err.Error()))
    }

    err = l.Bind(ldap_user, ldap_passwd)
    if err != nil {
        panic(fmt.Errorf("ldap.Bind panic: %s",err.Error()))
    }

    revel.TRACE.Printf("Get_connection() ok")
    return l

}

// Performs basic query
func QueryLdap(base string, filter string, attributes []string) *ldap.SearchResult {
    revel.TRACE.Printf("Query_Ldap(): filter: %s attributes: %s", filter, attributes)
    l := Get_connection()
    defer l.Close()

    search := ldap.NewSearchRequest(
            base,
            ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
            filter,
            attributes,
            nil)

    sr, err := l.Search(search)
    if err != nil {
        panic(fmt.Errorf("ldap.Search panic: %s",err.Error()))
    }

    revel.TRACE.Printf("Query_Ldap() returning %+v", sr)
    return sr
}

// Build the struct from LDAP search result
func Build_user_details(entry *ldap.Entry) User_details {
    details := User_details{}
    details.username = entry.GetAttributeValue(ldap_user_uid_attr) 
    details.visiblename = entry.GetAttributeValue(ldap_user_cn_attr) 
    details.photo = entry.GetRawAttributeValue(ldap_user_photo_attr)
    details.groups = entry.GetAttributeValues(ldap_user_group_attr)
    details.roles = append(details.roles, fmt.Sprintf("u:%s", details.username))
    for _, elem := range details.groups {
        details.roles = append(details.roles, fmt.Sprintf("g:%s", elem))
    }
    return details
}

// Retrieve user details from LDAP
func Get_user_details(username string) User_details {
    revel.TRACE.Printf("Get_user_details(): %s", username)

    sr := QueryLdap(ldap_user_base, strings.Replace(ldap_user_filter, "*", username, -1), []string{ldap_user_uid_attr, ldap_user_cn_attr, ldap_user_photo_attr, ldap_user_group_attr})
    // We expect exactly one result. If that's not true, something is probably really really wrong!
    if len(sr.Entries) != 1 {
        panic(fmt.Errorf("ldap query for %s returned %s hits", username, len(sr.Entries)))
    }

    details := Build_user_details(sr.Entries[0])
    revel.TRACE.Printf("Get_user_details() returning: %+v", details)
    return details
}

// Filter that loads the user's LDAP groups
func UserDetailsLoadFilter(c *revel.Controller, fc []revel.Filter) {
    username := c.Session["username"]
    revel.TRACE.Printf("UserDetailsLoadFilter() for %s", username)
    var dets User_details

    err := cache.Get(fmt.Sprintf("user_details_%s",username), &dets)

    if err != nil {
        // Unable to get from cache, recreate
        dets = Get_user_details(username)
        go cache.Set(fmt.Sprintf("user_details_%s",username), dets, cache.DEFAULT)
    }

    // Now we should have the details so pass them forward in Args
    c.Args["user_details"] = dets

    revel.TRACE.Printf("UserDetailsLoadFilter() saved Args[\"user_details\"]")

    // Next filter
    fc[0](c, fc[1:])
}
