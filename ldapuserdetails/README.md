Revel LDAP user details Filter
==================================

Simple Filter that gets details about user from LDAP.

Usage
-----

init.go:

```go
import (
    "github.com/mikkolehtisalo/ldapuserdetails"
)
    revel.Filters = []revel.Filter{
        ldapuserdetails.UserDetailsLoadFilter,     // Load user details from LDAP
    }
```

app.conf:

```go
ldap.server=freeipa.localdomain
ldap.port=389
ldap.base=cn=users,cn=accounts,dc=localdomain
ldap.user_filter=(&(uid=*)(objectClass=inetUser))
ldap.user_uid_attr=uid
ldap.user_cn_attr=cn
ldap.user_photo_attr=photo;binary
ldap.user_group_attr=memberOf
ldap.group_filter=(&(cn=*)(objectClass=groupOfNames))
ldap.group_attributes=cn,member
ldap.user=uid=admin,cn=users,cn=accounts,dc=localdomain
ldap.passwd=password
```

The details can be accessed from struct

```go
type User_details struct {
    Username string
    Visiblename string
    Photo []byte
    Groups []string
    Roles []string
}
```

that will be saved into *c.Args["user_details"]*. Type assertion .(ldapuserdetails.User_details) is usually required.
