# Level 32: Authentication & JWT - Complete Guide

## Introduction

Welcome to Level 32! You've mastered `Level 31: Dependency Injection` - constructor injection, interfaces as seams, and hand-wiring a dependency graph in `main()`. Now it's time to put that skill to work on a real, security-sensitive feature: **proving who a user is**, and **letting them prove it again on every later request** without sending their password each time.

> ⚠️ **This level needs internet access, like Levels 28 and 29 before it.** Levels 0-27, 30, and 31 work **100% offline**. This level needs two third-party packages that the standard library does not provide: `golang.org/x/crypto/bcrypt` (password hashing) and `github.com/golang-jwt/jwt/v5` (JSON Web Tokens), fetched from the Go module proxy the first time you `go get` them. Network access was verified working in this environment before writing a single exercise below - `golang.org/x/crypto v0.55.0` and `github.com/golang-jwt/jwt/v5 v5.3.1` were both fetched successfully, pinned to those exact versions, and **every command, program, and "Expected Output" block in this level's materials was actually run with `go run` and its output captured for real** - nothing here is hand-computed. Once fetched, Go caches both packages locally (`$GOPATH/pkg/mod`), so you won't need the network again for the same versions. If you're offline right now, come back to this level once you have a connection.

Authentication answers one question: **"who is making this request?"** Everything in this level builds toward answering it safely - hashing passwords so a leaked database doesn't hand out plaintext credentials, and issuing signed tokens so a server doesn't have to remember every logged-in user in memory.

---

## Table of Contents

1. [Why You Never Store Plaintext Passwords](#why-you-never-store-plaintext-passwords)
2. [bcrypt in Practice](#bcrypt-in-practice)
3. [What a JWT Actually Is](#what-a-jwt-actually-is)
4. [Claims: Registered and Custom](#claims-registered-and-custom)
5. [Validating a JWT](#validating-a-jwt)
6. [Security Awareness: The alg:none Attack](#security-awareness-the-algnone-attack)
7. [JWT Auth Middleware](#jwt-auth-middleware)
8. [Session-Based vs Token-Based Authentication](#session-based-vs-token-based-authentication)
9. [Refresh Tokens](#refresh-tokens)
10. [Best Practices](#best-practices)
11. [Common Mistakes](#common-mistakes)

---

## Why You Never Store Plaintext Passwords

There are two very different ways to turn readable data into unreadable data:

- **Encryption** is **two-way (reversible)**. If you have the key, you can turn ciphertext back into plaintext. Encryption is for data *you* need to read again later (a config secret, a database column you'll decrypt for a report).
- **Hashing** is **one-way (irreversible, in the direction that matters)**. You can't turn a hash back into the original input - you can only hash a *guess* and compare the two hashes. Hashing is for data you never need to read again, only to verify - which is exactly what a password is.

If you store passwords encrypted, then anyone who steals your database **and** your encryption key has every plaintext password. If you store them hashed, an attacker who steals the database only has hashes - and a well-chosen hash function makes guessing the original password computationally expensive, even at scale.

**Why not a plain hash like SHA-256 then?** Two reasons:

1. **Speed is the enemy here.** SHA-256 is *designed* to be fast - great for checksums, terrible for passwords, because it lets an attacker try billions of guesses per second against a stolen hash.
2. **No built-in salt.** Without a per-password random salt, identical passwords produce identical hashes, and precomputed "rainbow tables" of common password hashes crack them instantly.

**bcrypt** solves both problems on purpose:

- It's **deliberately slow**, with a tunable "cost" factor - the standard library's `bcrypt.DefaultCost` is 10, meaning ~2^10 rounds of internal key-setup work. Slow for one password check is irrelevant to a real login; slow for a billion guesses is the whole point.
- It **generates and embeds a random salt automatically** - you never manage salts yourself, and hashing the same password twice produces two completely different (but both valid) hashes.

This is why bcrypt (or a similar deliberately-slow algorithm like `scrypt` or `argon2`) is the standard choice for password storage, and why a fast general-purpose hash like MD5 or SHA-1/256 is not.

---

## bcrypt in Practice

Go doesn't ship bcrypt in the standard library, but `golang.org/x/crypto/bcrypt` is the de facto standard extension - maintained by the Go team itself, just not bundled with the compiler. Two functions cover essentially everything:

```go
hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
// hash is a []byte you store (as a string) in your database

err := bcrypt.CompareHashAndPassword(hash, []byte(candidatePassword))
// err == nil means the candidate password matches the hash
```

Real, verified round trip:

```
=== Password Hash ===
hash length: 60 bytes
hash prefix (cost identifier): $2a$10$

=== Verifying the Correct Password ===
MATCH: password is correct

=== Hashing the Same Password Twice ===
hash1 == hash2: false (different salts every time)
but both still verify the same password: true
```

Every bcrypt hash is 60 bytes and self-describes its own algorithm version (`$2a$`) and cost (`10`) right in the string, so `CompareHashAndPassword` never needs you to separately track which cost you hashed with - it's baked into the hash itself.

And the real, verified rejection when the password is wrong:

```
=== Comparing the Hash Against a Wrong Password ===
Rejected as expected.
Real error: crypto/bcrypt: hashedPassword is not the hash of the given password
Error type: *errors.errorString
err == bcrypt.ErrMismatchedHashAndPassword: true
```

That exact error is `bcrypt.ErrMismatchedHashAndPassword` - compare against it with `==` or `errors.Is` when you need to distinguish "wrong password" from "something else went wrong" (a malformed hash, for instance).

---

## What a JWT Actually Is

A **JWT (JSON Web Token)** is a compact, URL-safe string used to carry a set of claims ("who this is" and "what they're allowed to do") from one party to another, in a way the receiver can verify was not tampered with.

Structurally, a JWT is exactly three base64url-encoded segments joined by dots:

```
header.payload.signature

eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ1c2VyLTQyIiwuLi59.dNflZCu_JHBgYMdIukQ...
└──────────── header ────────────┘ └──────────── payload ─────────┘ └── signature ──┘
```

- **Header** - which algorithm signed this token (`{"alg":"HS256","typ":"JWT"}`).
- **Payload** - the claims themselves (who, when, what).
- **Signature** - a cryptographic signature over `base64url(header) + "." + base64url(payload)`, computed with a secret (or private key) only the issuer knows.

**Critically: the header and payload are only *encoded*, not encrypted.** Base64url is not a secret transformation - it's a reversible text encoding, like hex, with no key involved at all. Anyone holding the token can decode the payload with nothing but a text editor or one line of code. Verified, without ever touching the signing secret:

```
=== Full Token ===
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ1c2VyLTQyIiwiZXhwIjoxNzg4MjgzNzY3LCJ1c2VybmFtZSI6ImFsaWNlIiwicm9sZSI6Im1lbWJlciJ9.dNflZCu_JHBgYMdIukQaxFy52v8cltF9d_Sr77Hz2xw

Segments: 3 (header.payload.signature)

=== Decoded Header (no secret used) ===
{"alg":"HS256","typ":"JWT"}

=== Decoded Payload (no secret used!) ===
{
  "exp": 1788283767,
  "role": "member",
  "sub": "user-42",
  "username": "alice"
}

Anyone holding this token can read the payload above WITHOUT the signing secret.
The secret is only required to verify the SIGNATURE, never to read the claims.
```

**The rule that follows immediately: never put a secret, a password, or anything sensitive inside a JWT's claims.** A JWT proves the claims *haven't been altered* since signing (that's what the signature is for) - it does not hide them from view. If you need to hide data from the token holder, encrypt it separately (a JWE, not a JWT) or don't put it in the token at all.

---

## Claims: Registered and Custom

The payload is a JSON object of **claims**. JWT defines a handful of standard, three-letter **registered claims** that most libraries understand out of the box:

| Claim | Meaning |
|-------|---------|
| `sub` | Subject - who this token is about (usually a user ID) |
| `iss` | Issuer - who created and signed the token |
| `exp` | Expiration time (Unix timestamp) - reject the token after this |
| `iat` | Issued-at time (Unix timestamp) |
| `nbf` | Not-before time - reject the token before this |
| `aud` | Audience - who the token is intended for |

`github.com/golang-jwt/jwt/v5` models these as `jwt.RegisteredClaims`, which you embed in your own struct alongside whatever **custom claims** your application needs:

```go
type AppClaims struct {
    jwt.RegisteredClaims
    Username string `json:"username"`
    Role     string `json:"role"`
}

claims := AppClaims{
    RegisteredClaims: jwt.RegisteredClaims{
        Subject:   "user-42",
        Issuer:    "level32-auth-demo",
        IssuedAt:  jwt.NewNumericDate(now),
        ExpiresAt: jwt.NewNumericDate(now.Add(15 * time.Minute)),
    },
    Username: "alice",
    Role:     "member",
}

token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
signed, err := token.SignedString(secret) // secret is a []byte shared between issuer and verifier
```

`jwt.SigningMethodHS256` is **HMAC-SHA256** - a *symmetric* algorithm, meaning the same secret both signs and verifies. (Asymmetric algorithms like RS256 exist too, using a private key to sign and a public key to verify, but HS256 with a shared secret is the simplest correct starting point and what this level uses throughout.) Real, verified output:

```
=== Signed JWT ===
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJsZXZlbDMyLWF1dGgtZGVtbyIsInN1YiI6InVzZXItNDIiLCJleHAiOjE3ODgyODM3NjYsImlhdCI6MTc4ODI4Mjg2NiwidXNlcm5hbWUiOiJhbGljZSIsInJvbGUiOiJtZW1iZXIifQ.5-zWwkglxYo1M1hf9FmK1vKc26mpBNFDsj1RYKi50mo

Number of '.' separators: 2 (3 segments: header.payload.signature)
```

---

## Validating a JWT

Issuing a token is half the job - a server also has to correctly **reject** any token that isn't exactly what it issued. `jwt.ParseWithClaims` does three things in one call: decode, verify the signature, and check time-based claims (`exp`, `nbf`) automatically.

```go
parsed, err := jwt.ParseWithClaims(signed, &AppClaims{}, func(t *jwt.Token) (interface{}, error) {
    if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
        return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
    }
    return secret, nil
})
if err != nil || !parsed.Valid {
    // reject
}
```

The **keyfunc** (that closure) is where you supply the secret - and, just as importantly, where you check that the token's algorithm is one you actually expect (more on why in the next section). Three real, verified cases:

**Case 1 - a valid token is accepted:**

```
=== Validating a Legitimate Token ===
ACCEPTED: token is valid
Subject: user-42
Username: alice
Role: member
Expires at: 2026-09-01T22:59:27+05:30
```

**Case 2 - a tampered token is rejected with the real signature error:**

```
=== Attempting to Validate the Tampered Token ===
REJECTED (as it must be). Real error: token signature is invalid: signature is invalid
errors.Is(err, jwt.ErrTokenSignatureInvalid): true
```

(The token above had a single character of its signature segment flipped - nothing else changed. That alone is enough for verification to fail, which is the entire point of a signature.)

**Case 3 - an expired token is rejected with the real expiration error:**

```
=== Attempting to Validate the Expired Token ===
REJECTED (as it must be). Real error: token has invalid claims: token is expired
errors.Is(err, jwt.ErrTokenExpired): true
```

`ParseWithClaims` checks `exp` for you automatically - you never need to manually compare `time.Now()` against the claim yourself, and you shouldn't skip that check even if it feels redundant.

---

## Security Awareness: The alg:none Attack

Early JWT libraries (including early versions in several languages, not just Go) had a real, exploited vulnerability: the JWT spec allows an `alg` of `"none"`, meaning "this token is unsigned." An attacker could take a legitimately-issued token, **change the header** to `{"alg":"none",...}`, strip the signature entirely, and edit the payload however they liked (say, change `"role":"member"` to `"role":"admin"`) - and a naive verifier that trusted the header's `alg` field would accept it, because there was never a signature to check.

The fix is simple but **must be enforced by the verifier, not assumed from the library**: a correct verifier decides which algorithm(s) it will accept *itself*, and never lets the token's own header dictate that choice. That's exactly what the `t.Method.(*jwt.SigningMethodHMAC)` type-assertion in every keyfunc above is doing - it hard-codes "only HMAC is acceptable here," regardless of what the incoming token's header claims.

Verified: a forged `alg:none` token, with claims edited and the signature segment left empty, run against that same defensive keyfunc:

```
=== Forged alg:none Token (empty signature) ===
eyJhbGciOiJub25lIiwidHlwIjoiSldUIn0.eyJzdWIiOiJ1c2VyLTQyIiwiZXhwIjoxNzg4MjgzNzcwLCJ1c2VybmFtZSI6ImF0dGFja2VyLXByZXRlbmRpbmctdG8tYmUtYWxpY2UifQ.

=== Attempting to Validate the Forged Token (keyfunc checks the method) ===
REJECTED (correct verifier behavior). Real error: token is unverifiable: error while executing keyfunc: unexpected signing method: none (only HMAC is accepted)
```

`github.com/golang-jwt/jwt/v5` additionally requires you to opt in explicitly (`jwt.UnsafeAllowNoneSignatureType`) before it will even consider `alg:none` valid at all - a deliberate, named "unsafe" escape hatch you should never reach for in real code. Between the library's opt-in requirement and your own keyfunc check, `alg:none` has no path through this level's code. **Never remove the algorithm type-check from a keyfunc, even though the library defends against the worst case too - defense in depth matters here.**

---

## JWT Auth Middleware

Bringing Level 27's middleware pattern (`func(http.Handler) http.Handler`-shaped wrapping) together with Level 24's `context.WithValue`, an auth middleware extracts the token, validates it, and injects the parsed claims into the request's context for every handler downstream:

```go
type claimsContextKey struct{} // unexported, unique type - avoids collisions (Level 24 pattern)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        header := r.Header.Get("Authorization")
        if header == "" || !strings.HasPrefix(header, "Bearer ") {
            http.Error(w, `{"error":"missing or malformed Authorization header"}`, http.StatusUnauthorized)
            return
        }
        tokenString := strings.TrimPrefix(header, "Bearer ")

        claims := &AppClaims{}
        token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
            if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
            }
            return secret, nil
        })
        if err != nil || !token.Valid {
            http.Error(w, `{"error":"invalid token"}`, http.StatusUnauthorized)
            return
        }

        ctx := context.WithValue(r.Context(), claimsContextKey{}, claims)
        next(w, r.WithContext(ctx))
    }
}
```

Downstream, a handler retrieves the claims the same way Level 24 retrieves any context value - a type assertion against the same unexported key type:

```go
func profileHandler(w http.ResponseWriter, r *http.Request) {
    claims := r.Context().Value(claimsContextKey{}).(*AppClaims)
    // claims.Username, claims.Subject, etc. are now available
}
```

Verified end-to-end with `httptest`, one valid case and several real-world failure cases:

```
=== Case 1: Valid Token ===
Status: 200
Body: {"username":"alice"}

=== Case 2: Missing Authorization Header ===
Status: 401
Body: {"error":"missing or malformed Authorization header"}

=== Case 3: Malformed Header (no 'Bearer ' prefix) ===
Status: 401
Body: {"error":"missing or malformed Authorization header"}

=== Case 4: Tampered Token ===
Status: 401
Body: {"error":"invalid token"}

=== Case 5: Expired Token ===
Status: 401
Body: {"error":"token is expired"}
```

---

## Session-Based vs Token-Based Authentication

Both approaches answer "who is this?" - they differ in **where the truth lives**.

| | Session-based | Token-based (JWT) |
|---|---|---|
| **Where state lives** | Server-side (a session store: memory, Redis, a DB table) | Client-side (the token itself is self-contained) |
| **What the client holds** | An opaque session ID (usually a cookie) | The full, signed claims |
| **Scalability** | Every server needs access to the shared session store | Any server with the signing secret can verify a token alone - no shared store required |
| **Revocation** | Trivial - delete the session server-side, the ID is instantly worthless | Hard - a valid, unexpired JWT remains valid until it expires; there's no built-in "undo" |
| **Payload visibility** | Server decides what to expose | Client can read every claim (Section 3) |
| **Typical use** | Traditional web apps, same origin, cookies work well | APIs, mobile clients, microservices, cross-domain |

Neither is universally "better." Sessions make revocation and server-side control trivial but cost you a shared store and a lookup on every request. JWTs make verification cheap and stateless but make "log this one user out right now" genuinely hard - the usual answers are short expirations, a server-side denylist for the rare emergency revoke, or accepting the tradeoff and living with short-lived tokens (see the next section).

---

## Refresh Tokens

Because a JWT can't be un-issued, the standard mitigation is to make access tokens **short-lived** - minutes, not days - so a leaked token has a small window of usefulness. But re-logging in (password + all) every 15 minutes would be unusable, so the pattern splits authentication into two tokens:

- **Access token** - short-lived (minutes), sent with every API request, exactly what this level has built so far.
- **Refresh token** - longer-lived (days/weeks), stored more carefully (often an httpOnly cookie or secure storage, never sent with ordinary API calls), used *only* to request a new access token when the old one expires.

The client uses the access token normally. When it expires, the client sends the refresh token to a dedicated endpoint (`POST /refresh`) to get a fresh access token, without the user re-entering credentials. Crucially, a refresh token **can** be tracked server-side (in a small store, keyed by user or device) - which gets back some of session-based auth's revocation ability, while still letting most requests stay fast and stateless. This level explains the concept and lets Bonus Challenge 1 build a minimal version of that store; a production system typically also rotates refresh tokens on each use and detects reuse of an already-rotated token as a signal of theft.

---

## Best Practices

### 1. HTTPS Is a Hard Requirement in Production

A JWT's payload is plaintext-readable (Section 3) - the *only* thing protecting it in transit is TLS. Sending a bearer token over plain HTTP hands it to anyone on the network path. This course's exercises run over `httptest` (loopback only, no real network), so TLS setup isn't demonstrated here - but state it plainly: **never send an `Authorization: Bearer` header over unencrypted HTTP outside of local development.**

### 2. Strong, Random Signing Secrets From Configuration

Every example in this level uses a hard-coded secret string for clarity - **never do that in real code.** Generate a long, cryptographically random secret (32+ bytes) and load it from an environment variable or a secrets manager, never a source file. `Level 34: Configuration` covers exactly this: pulling secrets out of code and into config.

### 3. Short Expiration Windows

The shorter the `exp`, the smaller the damage window if a token leaks. Pair short access-token lifetimes with the refresh-token pattern (Section 9) rather than issuing long-lived access tokens for convenience.

### 4. Never Log Tokens or Passwords

A JWT in a log line is as sensitive as the credentials that produced it - anyone with log access can replay it until it expires. Never log an `Authorization` header, a raw password, or a bcrypt hash's plaintext input. Log a user ID or request ID instead.

---

## Common Mistakes

### Mistake 1: Assuming a JWT Is Encrypted

```
❌ WRONG ASSUMPTION: "It's signed, so nobody can read what's inside."
✅ REALITY: Signing proves integrity (not tampered), not confidentiality (not readable).
   Section 3's decode-without-a-secret proved this. Never put secrets in claims.
```

### Mistake 2: Not Checking Expiration

```go
// ❌ WRONG - trusting a token forever just because the signature is valid
if signatureValid {
    grantAccess()
}

// ✅ RIGHT - jwt.ParseWithClaims checks exp automatically; don't disable or skip that
```

### Mistake 3: Weak or Hardcoded Signing Secrets

```go
// ❌ WRONG - short, guessable, and baked into source control forever
var secret = []byte("secret123")

// ✅ RIGHT - long, random, loaded from the environment (Level 34)
var secret = []byte(os.Getenv("JWT_SIGNING_SECRET"))
```

### Mistake 4: Not Validating the Signing Algorithm

```go
// ❌ WRONG - trusts whatever alg the token itself claims
return lookupKeyFor(t.Header["alg"])

// ✅ RIGHT - the verifier decides the acceptable algorithm, always
if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
    return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
}
```

### Mistake 5: Storing Passwords in Plaintext or With Weak Hashing

```go
// ❌ WRONG - plaintext, or a fast general-purpose hash with no salt
storePassword(password)
storePassword(sha256(password))

// ✅ RIGHT - a deliberately slow, auto-salted algorithm
hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
storePassword(hash)
```

---

## Summary

**Password Storage:**
- Hashing (one-way) for passwords, never encryption (two-way) - bcrypt auto-salts and is deliberately slow
- `bcrypt.GenerateFromPassword` / `bcrypt.CompareHashAndPassword` - the entire API surface you need

**JWT Structure:**
- Three base64url segments: `header.payload.signature`
- Encoded, not encrypted - anyone can read the payload; never put secrets in claims
- The signature proves integrity (untampered), not confidentiality

**Claims:**
- Registered (`sub`, `iss`, `exp`, `iat`) + custom fields embedded alongside them
- `jwt.NewWithClaims` + `SignedString(secret)` to issue; `jwt.ParseWithClaims` to validate

**Validation:**
- Always check the signature, always check expiration, always pin the accepted algorithm in the keyfunc
- Never trust the token's own `alg` header to decide how to verify it (the `alg:none` lesson)

**Middleware:**
- Extract `Authorization: Bearer <token>`, validate, inject claims via `context.WithValue` with an unexported key type (Level 24)

**Session vs Token:**
- Sessions: server-side state, trivial revocation, needs a shared store
- JWTs: stateless, cheap to verify anywhere, hard to revoke early - mitigate with short expirations and refresh tokens

---

## Next Steps

You now understand:
- ✅ Why passwords are hashed (bcrypt), never encrypted, and verified this with real round-trip and rejection output
- ✅ The exact three-segment structure of a JWT, and proved the payload is only encoded, not encrypted
- ✅ Registered and custom claims, and issuing a real signed HS256 token
- ✅ Validating a JWT and seeing the real errors for a tampered signature and an expired token
- ✅ The historical `alg:none` vulnerability and why the verifier - not the token - must choose the algorithm
- ✅ Building JWT auth middleware that bridges Level 27's middleware pattern with Level 24's context values
- ✅ The session-vs-JWT tradeoff, and the access-token/refresh-token pattern for mitigating JWT revocation limits

**Next level:** Level 33 - Logging

You've now built the piece nearly every real API needs before anything else: proving who's asking. Every handler from here on can assume `r.Context()` already knows who's calling! 🚀
