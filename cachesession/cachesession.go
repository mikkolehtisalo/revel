package cachesession

import (
    "github.com/robfig/revel/cache"
    "github.com/robfig/revel"
    "time"
    "net/http"
    "fmt"
    "crypto/sha1"
    "crypto/rand"
    "io"
    "encoding/hex"
    "strings"
)

const (
    // The default expiration time of the session if session.expires is not defined
    // 8 hours is sane default because it covers usually one normal working day
    defaultExpiration = 8 * time.Hour

    // Used for distinguishing different types of data in cache
    SESSION_CACHE_ID = "SESSIONCACHE_"

    // Used for finding the designated session ID from the map
    SESSION_ID_KEY = "SESSION_ID"

    // Used for finding the session IP from the map
    SESSION_IP = "SESSION_IP"
)

var (
    expireAfterDuration time.Duration
    sessionIPLock bool
)

func init() {
    revel.OnAppStart(func() {
        var err error

        if expiresString, ok := revel.Config.String("session.expires"); !ok {
            expireAfterDuration = defaultExpiration
        } else if expiresString == "session" {
            expireAfterDuration = 0
        } else if expireAfterDuration, err = time.ParseDuration(expiresString); err != nil {
            panic(fmt.Errorf("session.expires invalid: %s", err))
        }

        if IPLock, ok := revel.Config.Bool("session.iplock"); !ok {
            sessionIPLock = false
        } else {
            sessionIPLock = IPLock
        }

    })
}


// Filter for revel to load & save the session with requests
func CacheSessionFilter(c *revel.Controller, fc []revel.Filter) {
    c.Session = restoreSession(c.Request.Request)

    // Make session vars available in templates as {{.session.xyz}}
    c.RenderArgs["session"] = c.Session

    // Next filter
    fc[0](c, fc[1:])

    // Save the session
    saveSession(c)
}

// Returns a Session pulled from cache by id in cookie
func getSessionFromCache(cookie *http.Cookie) revel.Session {
    var session revel.Session

    // Get session id from cookie
    var id string = SESSION_CACHE_ID + cookie.Value

    // Restore session from cache.
    if err := cache.Get(id, &session); err != nil {
        // Generate new
        session = make(revel.Session)
    }
    return session
}


// Restores session by the cookie
func restoreSession(req *http.Request) revel.Session {

    // Find session cookie
    cookie, err := req.Cookie(revel.CookiePrefix + "_SESSION")
    if err != nil {
        // There is no session cookie present, so return new session
        return make(revel.Session)
    }
    
    sess := getSessionFromCache(cookie)

    // Check that the user IP is not changed, if enabled
    if sessionIPLock {
        // Will remove : and port number from the remote address
        remote := req.RemoteAddr[:strings.LastIndex(req.RemoteAddr, ":")]
        if ipStr, ok := sess[SESSION_IP]; ok {
            if ipStr != remote {
                revel.WARN.Printf("Session IP mismatch! Present: %s, Restored: %s, Session: %s", remote, ipStr, sess)
                // SESSION_IP is set but doesn't match restored. Will return empty session instead!
                return make(revel.Session)
            }
        }
    }

    return sess
}

// Random SHA1 hash generator
func randomHash() string { 
        buf := make([]byte, sha1.Size, sha1.Size) 
        _, err := io.ReadFull(rand.Reader, buf) 
        if err != nil { 
                panic(fmt.Errorf("random read failed: %v", err)) 
        } 
        h := sha1.New() 
        h.Write(buf) 
        return hex.EncodeToString(h.Sum(buf))
} 

// Get the hash random hash identifying the session.
func HashId(s revel.Session) string { 

    // ID already generated and saved into this session?
    if idStr, ok := s[SESSION_ID_KEY]; ok {
        return idStr
    }

    // Generate new!
    idStr := randomHash()

    s[SESSION_ID_KEY] = idStr
    return s[SESSION_ID_KEY]
}

// Saves the session to the cache
func saveSession(c *revel.Controller) {

    generatedId := HashId(c.Session)

    // Save the current remote user's address
    if sessionIPLock {
        remote := c.Request.RemoteAddr[:strings.LastIndex(c.Request.RemoteAddr, ":")]
        c.Session[SESSION_IP] = remote
    }

    // Save session only if there are items in it. SESSION_ID_KEY should always be present.
    if len(c.Session) > 1 {

            // Store the session in cache.
            go cache.Set(SESSION_CACHE_ID + generatedId, c.Session, expireAfterDuration)

            // Set the cookie. It gets re-set because every access pushes expiry time forward.
            c.SetCookie(&http.Cookie{
                Name: revel.CookiePrefix + "_SESSION",
                Value: generatedId,
                Expires: time.Now().Add(expireAfterDuration),
                Path: "/",
                HttpOnly: revel.CookieHttpOnly,
                Secure: revel.CookieSecure,
            })
    }
}

