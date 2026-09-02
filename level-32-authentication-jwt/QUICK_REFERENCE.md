# Level 32: Quick Reference Card

> ⚠️ **Needs internet access once:** `go get golang.org/x/crypto/bcrypt@v0.55.0` and `go get github.com/golang-jwt/jwt/v5@v5.3.1`. Verified working; both packages cache locally after the first fetch.

## 🚀 Quick Start (5 Minutes)

```bash
# Setup
mkdir myapp && cd myapp
go mod init myapp
go get golang.org/x/crypto/bcrypt@v0.55.0
go get github.com/golang-jwt/jwt/v5@v5.3.1

# Create main.go
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "time"

    "github.com/golang-jwt/jwt/v5"
    "golang.org/x/crypto/bcrypt"
)

func main() {
    hash, _ := bcrypt.GenerateFromPassword([]byte("hunter2"), bcrypt.DefaultCost)
    fmt.Println(bcrypt.CompareHashAndPassword(hash, []byte("hunter2")) == nil)

    claims := jwt.RegisteredClaims{
        Subject:   "user-1",
        ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
    }
    tok, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte("secret"))
    fmt.Println(tok)
}
EOF

# Run
go run main.go
```

---

## 🔑 bcrypt: The Whole API

```go
hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
// hash is []byte, 60 bytes, self-describing (e.g. $2a$10$...) - store it as a string

err := bcrypt.CompareHashAndPassword(hash, []byte(candidate))
// err == nil            -> match
// err == bcrypt.ErrMismatchedHashAndPassword -> real, exact mismatch error
```

---

## 🪙 JWT: Issue a Token

```go
type AppClaims struct {
    jwt.RegisteredClaims
    Username string `json:"username"`
}

claims := AppClaims{
    RegisteredClaims: jwt.RegisteredClaims{
        Subject:   "user-42",
        IssuedAt:  jwt.NewNumericDate(time.Now()),
        ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
    },
    Username: "alice",
}
token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
signed, err := token.SignedString(secret) // secret is []byte
```

---

## 🔍 JWT: Validate a Token

```go
parsed, err := jwt.ParseWithClaims(signed, &AppClaims{}, func(t *jwt.Token) (interface{}, error) {
    if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
        return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
    }
    return secret, nil
})
if err != nil || !parsed.Valid {
    // reject - err distinguishes WHY:
    //   errors.Is(err, jwt.ErrTokenSignatureInvalid) -> tampered
    //   errors.Is(err, jwt.ErrTokenExpired)           -> expired
}
claims := parsed.Claims.(*AppClaims)
```

---

## 🧱 JWT Structure

```
header.payload.signature      (3 base64url segments)

header + payload   -> ENCODED only, readable by anyone, no secret needed
signature          -> proves header+payload weren't tampered with

NEVER put secrets/passwords inside claims - they are NOT hidden.
```

---

## 🛡️ Auth Middleware Skeleton

```go
type claimsContextKey struct{} // unexported type - Level 24 pattern

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        header := r.Header.Get("Authorization")
        if !strings.HasPrefix(header, "Bearer ") {
            http.Error(w, `{"error":"missing or malformed Authorization header"}`, http.StatusUnauthorized)
            return
        }
        tokenString := strings.TrimPrefix(header, "Bearer ")
        claims := &AppClaims{}
        token, err := jwt.ParseWithClaims(tokenString, claims, keyfunc)
        if err != nil || !token.Valid {
            http.Error(w, `{"error":"invalid token"}`, http.StatusUnauthorized)
            return
        }
        ctx := context.WithValue(r.Context(), claimsContextKey{}, claims)
        next(w, r.WithContext(ctx))
    }
}

// downstream handler:
claims := r.Context().Value(claimsContextKey{}).(*AppClaims)
```

---

## ⚖️ Session vs JWT

| | Session | JWT |
|---|---|---|
| State | Server-side store | Inside the token |
| Revoke | Trivial | Hard (mitigate with short `exp` + refresh tokens) |
| Scale | Needs shared store | Verify anywhere with the secret |

---

## ⚠️ Common Mistakes

| Mistake | Wrong | Right |
|---------|-------|-------|
| Thinking JWTs are encrypted | assuming payload is hidden | it's base64url-encoded only - readable by anyone |
| Skipping expiration | trusting signature alone | `ParseWithClaims` checks `exp` automatically - never bypass it |
| Weak/hardcoded secret | `[]byte("secret123")` | long random secret from env/secrets manager (Level 34) |
| Trusting the token's `alg` | `lookupKeyFor(t.Header["alg"])` | verifier hard-codes the accepted method type |
| Weak password storage | plaintext or `sha256(password)` | `bcrypt.GenerateFromPassword` |

---

## 🎓 Before Next Level

Can you:
- [ ] Hash and verify a password with bcrypt from memory?
- [ ] Explain why a JWT payload is readable without the secret?
- [ ] Issue and validate a JWT with custom claims?
- [ ] Explain the `alg:none` attack and how a correct keyfunc defeats it?
- [ ] Write and test JWT auth middleware end-to-end?

If YES → You're ready for Level 33!

---

## 📚 Next Level

Level 33: Logging
- Structured logging and log levels
- Correlating logs with the identity this level now carries in context

You've built real authentication! 💪
