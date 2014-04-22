package gssserver
/*
#include <gssapi.h>
#include <gssapi/gssapi_krb5.h>
#include <stdio.h>
#include <string.h>
#cgo LDFLAGS: -lgssapi_krb5 -lkrb5
OM_uint32 iserr(OM_uint32 status) {
 return GSS_ERROR(status);  
}
gss_channel_bindings_t getbindings() {
    return GSS_C_NO_CHANNEL_BINDINGS;
}
gss_cred_id_t getcredential() {
    return GSS_C_NO_CREDENTIAL;
}
gss_ctx_id_t getnocontext() {
    gss_ctx_id_t val = GSS_C_NO_CONTEXT;
    return val;
}
gss_OID getnooid() {
    gss_OID val = GSS_C_NO_OID;
    return val;
}
gss_buffer_t getnobuffer() {
    gss_buffer_t val = GSS_C_NO_BUFFER;
    return val;
}
*/
import "C"

import (
    "github.com/robfig/revel"
    "unsafe"
    "fmt"
    "strings"
    "errors"
    "net/http"
    "encoding/base64"
)

const (
    // Used for cache lookups
    AUTH_USER_ID = "AUTH_USER_ID"
)

func gss_acquire_cred(desired_name C.gss_name_t) C.gss_cred_id_t {
    revel.TRACE.Printf("gss_acquire_cred(): %s", desired_name)
    major_status := C.OM_uint32(0)
    minor_status := C.OM_uint32(0)
    tmp := &[0]byte{}
    
    output_cred_handle := C.gss_cred_id_t(tmp)
    major_status = C.gss_acquire_cred(
        &minor_status,
        desired_name,
        C.GSS_C_INDEFINITE,
        nil,
        C.GSS_C_ACCEPT,
        &output_cred_handle,
        nil,
        nil)

    if int(C.iserr(major_status)) > 0 {
        panic(fmt.Sprintf("gss_acquire_cred() panic! Major: %d Msg: %s - Minor: %d Msg: %s", int(major_status), gss_display_status(major_status, 1), int(minor_status), gss_display_status(minor_status, 2)))
    }

    revel.TRACE.Printf("gss_acquire_cred(): returning %+v", C.gss_cred_id_t(output_cred_handle))

    return output_cred_handle
}

func gss_accept_sec_context(acceptor_cred_handle C.gss_cred_id_t, token string) (C.struct_gss_buffer_desc_struct, C.gss_name_t) {
    revel.TRACE.Printf("gss_accept_sec_context(): %+v, %s", C.gss_cred_id_t(acceptor_cred_handle), token)
    major_status := C.OM_uint32(0)
    minor_status := C.OM_uint32(0)
    context_handle := C.getnocontext()
    input_token := C.struct_gss_buffer_desc_struct{}

    if len(token)>0 {
        data, err := base64.StdEncoding.DecodeString(token)
        if err != nil {
            panic(fmt.Sprintf("Unable to Base64 decode token %+v", token))
        }
        revel.TRACE.Printf("Base64 decoded token: %+v", data)
        input_token.length = C.size_t(len(data))
        input_token.value = unsafe.Pointer(&data[0])
    }

    tmp2 := &[0]byte{}
    src_name := C.gss_name_t(tmp2)
    output_token := C.struct_gss_buffer_desc_struct{}
    tmp3 := &[0]byte{}
    credential := C.gss_cred_id_t(tmp3)
    bindings := C.struct_gss_channel_bindings_struct{}

    /*
    A way to attempt tracing
    da := &[0]byte{}
    ctx := C.krb5_context(da)
    C.krb5_init_context(&ctx);
    C.krb5_set_trace_filename(ctx, C.CString("/tmp/lol.txt"))
    */

    major_status = C.gss_accept_sec_context(
        &minor_status,  
        &context_handle, 
        acceptor_cred_handle, 
        &input_token,  
        &bindings, 
        &src_name, 
        nil, 
        &output_token,
        nil, 
        nil,
        &credential) 

    defer C.gss_delete_sec_context(&major_status, &context_handle, C.getnobuffer())
    defer C.gss_release_cred(&major_status, &credential)

    if int(C.iserr(major_status)) > 0 {
        panic(fmt.Sprintf("gss_accept_sec_context() panic! Major: %d Msg: %s - Minor: %d Msg: %s", int(major_status), gss_display_status(major_status, 1), int(minor_status), gss_display_status(minor_status, 2)))
    }

    return output_token, src_name
}


func gss_display_name(input_name C.gss_name_t) string {
    revel.TRACE.Printf("gss_display_name(): %+v", input_name)
    major_status := C.OM_uint32(0)
    minor_status := C.OM_uint32(0)
    output_name_buffer := C.struct_gss_buffer_desc_struct{}

    major_status = C.gss_display_name(
        &minor_status,
        input_name,
        &output_name_buffer,
        nil)

    defer C.gss_release_buffer(&minor_status, &output_name_buffer)

    if int(C.iserr(major_status)) > 0 {
        panic(fmt.Sprintf("gss_display_name panic! Major: %d Msg: %s - Minor: %d Msg: %s", int(major_status), gss_display_status(major_status, 1), int(minor_status), gss_display_status(minor_status, 2)))
    }

    return C.GoString((*C.char)(output_name_buffer.value))
}

func gss_inquire_cred(cred_handle C.gss_cred_id_t) C.gss_name_t {
    revel.TRACE.Printf("gss_inquire_cred(): %+v", cred_handle)
    major_status := C.OM_uint32(0)
    minor_status := C.OM_uint32(0)
    tmp := &[0]byte{}
    name := C.gss_name_t(tmp)

    major_status = C.gss_inquire_cred(
        &minor_status,
        cred_handle,
        &name,  // gss_release_name ->
        nil, // lifetime
        nil, // cred_usage
        nil) // mechanisms

    if int(C.iserr(major_status)) > 0 {
        panic(fmt.Sprintf("gss_inquire_cred panic! Major: %d Msg: %s - Minor: %d Msg: %s", int(major_status), gss_display_status(major_status, 1), int(minor_status), gss_display_status(minor_status, 2)))
    }

    return name
}

func gss_import_name(s string) C.gss_name_t {
    revel.TRACE.Printf("gss_import_name(): %s", s)

    major_status := C.OM_uint32(0)
    minor_status := C.OM_uint32(0)
    input_name_buffer := C.struct_gss_buffer_desc_struct{}
    input_name_type := C.GSS_C_NT_HOSTBASED_SERVICE
    tmp := &[0]byte{}
    cs := C.CString(s)
    defer C.free(unsafe.Pointer(cs))
    input_name_buffer.length = C.strlen(cs)
    input_name_buffer.value = unsafe.Pointer(cs)

    output_name := C.gss_name_t(tmp)
    major_status = C.gss_import_name(
        &minor_status,
        &input_name_buffer,
        input_name_type,
        &output_name)

    if int(C.iserr(major_status)) > 0 {
        panic(fmt.Sprintf("gss_import_name panic! Major: %d Msg: %s - Minor: %d Msg: %s", int(major_status), gss_display_status(major_status, 1), int(minor_status), gss_display_status(minor_status, 2)))
    }

    revel.TRACE.Printf("gss_import_name(): returning %+v", C.gss_name_t(output_name))
    return output_name
}

func gss_display_status(status_value C.OM_uint32, status_type int) string {
    revel.TRACE.Printf("gss_display_status() for: %s, %s", status_value, status_type)
    status := ""
    minor_status := C.OM_uint32(0)
    message_context := C.OM_uint32(1)
    mech_type := C.getnooid()
    status_string := C.struct_gss_buffer_desc_struct{}

    for int(message_context) > 0 { // if != 0 after call, ask for more!
        C.gss_display_status(
            &minor_status,
            status_value,
            C.int(status_type), // either GSS_C_GSS_CODE 1 or GSS_C_MECH_CODE 2
            mech_type,
            &message_context,
            &status_string) // gss_buffer_t

        defer C.gss_release_buffer(&minor_status, &status_string)

        errmsg := C.GoString((*C.char)(status_string.value))
        status = fmt.Sprint(status, errmsg, " - ")
    }

    revel.TRACE.Printf("gss_display_status() returning: %s", status)
    return status
}

func getAuthorization(c *revel.Controller) string {
    revel.TRACE.Printf("getAuthorization()")
    authorization := ""
    authorizations := c.Request.Header["Authorization"]
    if len(authorizations) > 0 {
        authorization = authorizations[0][10:]
        
    }
    revel.TRACE.Printf("getAuthorization() returning token %s", authorization)
    return authorization
}

// Strips the @DOMAIN part out
func cleanUsername(input string) string {
    parts := strings.Split(input, "@")
    return parts[0]
}

func GSSServerFilter(c *revel.Controller, fc []revel.Filter) {
    revel.TRACE.Printf("GSSServerFilter() starts")

    if len(c.Session["username"]) == 0 {
        revel.TRACE.Printf("Username not found from session, authentication steps...")

        status := C.OM_uint32(0)
        continue_filters := true

        gss_name := gss_import_name("HTTP")
        defer C.gss_release_name(&status, &gss_name)

        gss_cred := gss_acquire_cred(gss_name)
        defer C.gss_release_cred(&status, &gss_cred)
        
        authorization := getAuthorization(c)
        token, name := gss_accept_sec_context(gss_cred, authorization)
        defer C.gss_release_name(&status, &name)
        defer C.gss_release_buffer(&status, &token)

        if name==nil {
            revel.TRACE.Printf("Did not get name, building challenge")
            base64_token := base64.StdEncoding.EncodeToString(C.GoBytes(token.value, C.int(token.length)))
            c.Response.Out.Header().Add("WWW-Authenticate",fmt.Sprintf("Negotiate %s", base64_token))
            c.Response.Status = http.StatusUnauthorized
            c.Result = c.RenderError(errors.New("401: Not authorized"))
            continue_filters = false
        } else {
            revel.TRACE.Printf("Got name, retrieving authenticated user")
            base64_token := base64.StdEncoding.EncodeToString(C.GoBytes(token.value, C.int(token.length)))
            c.Response.Out.Header().Add("WWW-Authenticate",fmt.Sprintf("Negotiate %s", base64_token))
            c.Session["username"] = cleanUsername(gss_display_name(name))
        }

        // Next filter
        if continue_filters {
            fc[0](c, fc[1:])
        }
    
    } else {
        revel.TRACE.Printf("Username %s found from session, proceeding", c.Session["username"])
        fc[0](c, fc[1:])
    }
   
}
