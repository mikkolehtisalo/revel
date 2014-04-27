Cache based sessions for Revel
==============================

This module implements a simple cache based session storage for Revel. The cache is Revel's implementation, which means it can be configured shared for several servers if necessary.

Usage
-----

init.go:

```go
import (
    "cachesession"
)
    revel.Filters = []revel.Filter{
        cachesession.CacheSessionFilter, // Use cache based session implementation.
    }
```
app.conf:

```go
# cachesession
# Allow the session be used only if the requests come from the same IP address
session.iplock=true 
```
