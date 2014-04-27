Revel GSSAPI authentication Filter
==================================

This module implements kerberos/GSS-API (see [RFC417](http://tools.ietf.org/html/rfc4178)) Filter for Revel.

Kerberos is the preferred SSO method for intranet setups because it is secure, robust, and it performs well. It also offers improved usability (most of the logins are transparent to users) and robust interoperatibility between different platforms (Windows, MacOS X, Linux, BSDs, etc) and different browsers (Chrome, Internet Explorer, Firefox).

Prerequisities
--------------

* MIT kerberos (or compatible)
* A working kerberos setup
  * Basic pointers for getting started:
    * https://access.redhat.com/knowledge/docs/en-US/Red_Hat_Enterprise_Linux/6/html/Managing_Smart_Cards/Configuring_a_Kerberos5_Server.html
    * http://www.centos.org/docs/5/html/Deployment_Guide-en-US/ch-kerberos.html
    * http://www.freebsd.org/doc/handbook/kerberos5.html
    * http://tldp.org/HOWTO/Kerberos-Infrastructure-HOWTO/
* keytab & krb5.conf setup correctly on the server
* SPN is HTTP/fqdn@DOMAIN
* Browser setup correctly to provide GSSAPI authentication
  * *Use integrated windows authentication* for IE
  * *network.negotiate-auth.trusted-uris* for Firefox
  * *--auth-server-whitelist* for Chrome

Usage
-----

init.go:

```go
import (
    "gssserver"
)
    revel.Filters = []revel.Filter{
        gssserver.GSSServerFilter,     // GSSAPI authentication
    }
```
The authenticated user can be found from c.Session["username"]. 

Notes
-----

* The session system is used as cache for performance reasons - forcing full authentication for every request would be extremely heavy
* The security of the session system is critical, using a server side session storage is strongly recommended
* Major authentication errors cause panic - this module does not by default allow for unauthenticated use

Debugging
=========

The following environment variables will make Firefox print out extensive debug log:

* export NSPR_LOG_MODULES=negotiateauth:5
* export NSPR_LOG_FILE=/tmp/moz.log

The following environment variable will make krb5-libs print out trace log:

* export KRB5_TRACE=/tmp/krb.log

Unfortunately KRB5_TRACE gets lost within the Go http server's process model. This is an unfortunate POSIX feature, which makes server side tracing of GSSAPI hard. MIT kerberos has however an API for enabling trace on runtime, example can be found from gss_gssserver/accept_sec_context.
