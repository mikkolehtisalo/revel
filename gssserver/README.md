Revel GSSAPI authentication Filter
==================================

This module implements kerberos/GSS-API (see `RFC4178 <http://tools.ietf.org/html/rfc4178>`) Filter for Revel.

Prerequisities
--------------

* A working kerberos setup
* keytab & krb5.conf setup correctly on the server
* SPN is HTTP/fqdn@DOMAIN

Usage
-----

init.go:

```go
import (
    // ...
    "cachesession"
    "gssserver"
)
    revel.Filters = []revel.Filter{
        cachesession.CacheSessionFilter, // Use cache based session implementation.
        gssserver.GSSServerFilter,     // GSSAPI authentication
        // ...
    }
```
app.conf:

```go
# cachesession
session.iplock=true
```

The authenticated user can be found from Session["username"]. 
