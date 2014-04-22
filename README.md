Revel GSSAPI authentication Filter
==================================

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
