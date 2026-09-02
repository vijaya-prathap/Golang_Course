# Level 32: Study Guide & Visual Reference

> ⚠️ **Reminder:** this level needs internet access once, to `go get golang.org/x/crypto/bcrypt@v0.55.0` and `go get github.com/golang-jwt/jwt/v5@v5.3.1`. Both were verified reachable and every diagram/output below reflects real, captured program output - see README.md's introduction for details.

## 📚 Learning Path

### Week 1: Passwords and Tokens
```
Day 1:  Hashing vs encryption, why bcrypt for passwords
Day 2:  bcrypt round trip and rejection (Exercises 1-2)
Day 3:  JWT structure - three segments, encoded not encrypted
Day 4:  Registered + custom claims, issuing a signed token
Day 5:  Validating a token - signature and expiration checks
Day 6:  The alg:none attack and defensive keyfuncs
Day 7:  Review Sections 1-6 of README.md
```

### Week 2: Middleware & Practice
```
Day 1:  JWT auth middleware design (Section 7)
Day 2:  Exercises 3-5 (issuing, decoding, validating)
Day 3:  Exercises 6-8 (tampered, expired, alg:none rejections)
Day 4:  Exercise 9 (middleware, tested end-to-end)
Day 5:  Exercise 10 (comprehensive login API)
Day 6:  Session-vs-JWT comparison, refresh tokens (Sections 8-9)
Day 7:  Bonus challenges, review
```

---

## 🔐 bcrypt Hash-and-Compare Flow

```
REGISTRATION (storing a new password)

  plaintext password
        │
        ▼
  bcrypt.GenerateFromPassword(password, cost)
        │
        │   - generates a random salt internally
        │   - runs 2^cost rounds of deliberately slow work
        │
        ▼
  hash  (60 bytes, e.g. $2a$10$N9qo8uLOickgx2ZMRZoMy...)
        │
        ▼
  store ONLY the hash in the database - never the password


LOGIN (verifying a later attempt)

  candidate password  +  stored hash
        │                     │
        └─────────┬───────────┘
                   ▼
  bcrypt.CompareHashAndPassword(hash, candidate)
                   │
        ┌──────────┴──────────┐
        ▼                     ▼
   err == nil             err != nil
   PASSWORD MATCHES       "hashedPassword is not the
                           hash of the given password"
```

**Key fact:** the salt is extracted from inside the hash string itself - `CompareHashAndPassword` never needs it passed separately. That's why the API is only two functions.

---

## 🧩 JWT Structure Diagram

```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9 . eyJzdWIiOiJ1c2VyLTQyIiwuLi59 . dNflZCu_JHBgYMdIukQaxFy52v8cltF9d_Sr77Hz2xw
└──────────── HEADER ─────────────┘   └──────────── PAYLOAD ─────┘   └──────────── SIGNATURE ────────────────┘
    base64url({"alg":"HS256",             base64url({"sub":"user-42",         HMAC-SHA256(
               "typ":"JWT"})                          "exp":1788283767,                base64url(header) + "." +
                                                        "username":"alice"})             base64url(payload),
                                                                                          secret
      ENCODED ONLY                         ENCODED ONLY                                )
      (readable by anyone,                 (readable by anyone,             SIGNED
       no secret needed)                    no secret needed -                (proves the header+payload
                                             NEVER put a secret here!)          weren't altered since signing)
```

```
What's ENCODED (readable by anyone, no key needed):
  ✅ header    - which algorithm was used
  ✅ payload   - all the claims: sub, exp, iat, username, role, anything you put there

What's SIGNED (verifiable only with the secret/key, but NOT hidden):
  ✅ the signature itself proves header+payload weren't tampered with
  ❌ signing does NOT hide the header or payload from view
```

**One-line takeaway:** *Signed, not sealed.* A JWT is a tamper-evident envelope with a clear window, not a locked box.

---

## 🛡️ Auth Middleware Request Flow

```
Incoming Request
      │
      ▼
Authorization: Bearer <token> present and well-formed?
      │
      ├─ NO ──────────────────────────► 401 {"error":"missing or malformed
      │                                        Authorization header"}
      ▼
jwt.ParseWithClaims(token, claims, keyfunc)
      │
      ├─ keyfunc rejects the algorithm (defends against alg:none) ─► 401
      │
      ├─ signature invalid (tampered) ───────────────────────────► 401
      │
      ├─ token expired (exp in the past) ────────────────────────► 401
      │
      ▼ all checks pass
context.WithValue(r.Context(), claimsContextKey{}, claims)
      │
      ▼
next(w, r.WithContext(ctx))   →   downstream handler reads claims
                                    from r.Context(), already authenticated
      │
      ▼
200 (or whatever the handler decides, now that identity is known)
```

---

## ⚖️ Session-Based vs JWT (Token-Based) Comparison

| Dimension | Session-Based | JWT (Token-Based) |
|---|---|---|
| Truth lives | Server-side store | Inside the token itself |
| Client holds | Opaque session ID (cookie) | Full signed claims |
| Verify a request | Look up the store | Verify the signature locally - no lookup |
| Scale across servers | Needs a shared store (Redis, DB) | Any server with the secret can verify alone |
| Revoke immediately | Trivial - delete the session | Hard - valid until `exp`, no built-in undo |
| Claims visible to client | No (server controls exposure) | Yes - anyone can decode the payload |
| Typical fit | Traditional web apps, same-origin cookies | APIs, mobile, cross-domain, microservices |

---

## 🚨 Common Mistakes

### Mistake 1: Assuming a JWT Is Encrypted

```
❌ "It's signed, so the contents are hidden."
✅ Signing proves integrity, not secrecy - the payload decodes with zero effort.
```

### Mistake 2: Skipping the Expiration Check

```go
// ❌ trusting a valid signature forever
if signatureValid { grantAccess() }

// ✅ jwt.ParseWithClaims checks exp for you - never bypass it
```

### Mistake 3: Weak or Hardcoded Signing Secrets

```go
// ❌
var secret = []byte("secret123")

// ✅ long, random, from environment/secrets management (Level 34)
var secret = []byte(os.Getenv("JWT_SIGNING_SECRET"))
```

### Mistake 4: Not Validating the Signing Algorithm

```go
// ❌ trusts the token's own header
return lookupKeyFor(t.Header["alg"])

// ✅ the verifier decides, always
if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
    return nil, fmt.Errorf("unexpected signing method")
}
```

### Mistake 5: Plaintext or Weakly-Hashed Passwords

```go
// ❌
storePassword(password)          // plaintext
storePassword(sha256(password))  // fast hash, no salt

// ✅
hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
storePassword(hash)
```

---

## 📈 Progression Summary

### Understanding Level 32

Level 32 teaches how to answer "who is this?" safely, and keep answering it on every later request without storing server-side state for each user:

1. **Hashing vs encryption** - one-way for passwords, always
2. **bcrypt** - auto-salted, deliberately slow, the standard choice
3. **JWT structure** - three segments, encoded (not encrypted), signed (not sealed)
4. **Claims** - registered + custom, issued with `jwt.NewWithClaims`
5. **Validation** - signature, expiration, and algorithm, all enforced by the verifier
6. **The alg:none lesson** - never let the token's header decide how it's checked
7. **Middleware** - bridging Level 27's wrapping pattern with Level 24's context values

### Prerequisites for Level 33

Before moving to Level 33 (Logging), you need:

- ✅ Comfortable hashing and verifying passwords with bcrypt
- ✅ Can explain why a JWT's payload is readable without the secret
- ✅ Can issue, validate, and correctly reject (tampered/expired) a JWT
- ✅ Understand why the verifier - not the token - must choose the accepted algorithm
- ✅ Can build and test JWT auth middleware with `httptest`
- ✅ Can explain the session-vs-JWT tradeoff and the refresh-token pattern

### Ready for Level 33?

Level 33 teaches structured logging - the natural next step once requests carry an authenticated identity (and a request ID, from Level 24) worth writing into every log line:
- Log levels and when to use each
- Structured (key/value) logging vs plain text
- Correlating related log lines across a request's lifetime

---

## ✅ Checklist Before Level 33

- [ ] Can hash a password and verify it with bcrypt, from memory
- [ ] Can explain, out loud, why a JWT is encoded but not encrypted
- [ ] Can issue a signed JWT with both registered and custom claims
- [ ] Can explain what makes the `alg:none` attack work and how to defend against it
- [ ] Can write middleware that validates a bearer token and injects claims into context
- [ ] Completed 8+ exercises

---

## 💡 Key Takeaways

### The One Rule for Passwords
Hash, never encrypt - bcrypt handles salting and slowness for you.

### The One Rule for JWTs
Signed proves untampered; it does not mean hidden. Never put secrets in claims.

### The One Rule for Verifiers
The verifier decides the accepted algorithm and checks expiration - never trust the token to tell you either.

---

## 📚 Next Level

Level 33: Logging
- Structured logging and log levels
- Correlating logs with request/user identity carried through context

You've built real authentication from first principles - keep going! 🚀
