package gssserver

import (
    "testing"
)

// These tests will not free memory (no CGO available)
// To be able to test further would probably require that gss_init_sec_context be implemented

func TestImportName(t *testing.T) {
    gss_name := gss_import_name("HTTP")
    if gss_display_name(gss_name) != "HTTP" {
        t.Error("Expected HTTP, got ", gss_display_name(gss_name))
    }
}

func TestAcquireCred(t *testing.T) {
    gss_name := gss_import_name("HTTP")
    gss_cred := gss_acquire_cred(gss_name)
    gss_cred_name := gss_inquire_cred(gss_cred)

    if gss_display_name(gss_cred_name) != "HTTP/dev.localdomain@LOCALDOMAIN" {
        t.Error("Expected HTTP/dev.localdomain@LOCALDOMAIN, got ", gss_display_name(gss_cred_name))
    }
}

func TestAcceptNoToken(t *testing.T) {
    gss_name := gss_import_name("HTTP")
    gss_cred := gss_acquire_cred(gss_name)

    // Try with empty token
    token, name := gss_accept_sec_context(gss_cred, "")

    // We should get token data
    if token.length == 0 {
        t.Error("Did not get token, length 0!")
    }

    // We should NOT get name yet
    if name != nil {
        t.Error("Expected nil name, got ", name)
    }
}
