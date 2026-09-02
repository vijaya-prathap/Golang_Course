# Level 32: Authentication & JWT - Exercises

> ⚠️ **This level requires internet access.** Every exercise below installs one or both of two real third-party packages - `golang.org/x/crypto/bcrypt` (pinned at `v0.55.0`) and `github.com/golang-jwt/jwt/v5` (pinned at `v5.3.1`) - the exact versions this course was verified against. Every command, every program, and every "Expected Output" block on this page was actually run with `go run` against those pinned versions - nothing here is hand-computed. If `go get` fails, check your network connection before assuming your code is wrong. Once downloaded, Go caches both packages locally, so later exercises won't need the network again. This is only the third level in the course with this requirement, after Level 28 (Gin) and Level 29 (database/sql).

Complete all exercises in order. Each exercise builds on previous knowledge.

---

## Exercise 1: bcrypt Hash and Verify Round Trip

**Objective:** Hash a password with bcrypt and verify the correct password against it

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level32-exercise1
cd ~/projects/level32-exercise1
go mod init level32.example/exercise1
go get golang.org/x/crypto/bcrypt@v0.55.0
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "fmt"

    "golang.org/x/crypto/bcrypt"
)

func main() {
    password := []byte("correct-horse-battery-staple")

    hash, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
    if err != nil {
        panic(err)
    }
    fmt.Println("=== Password Hash ===")
    fmt.Printf("hash length: %d bytes\n", len(hash))
    fmt.Printf("hash prefix (cost identifier): %s\n", hash[:7])

    fmt.Println("\n=== Verifying the Correct Password ===")
    err = bcrypt.CompareHashAndPassword(hash, password)
    if err != nil {
        fmt.Println("MISMATCH:", err)
    } else {
        fmt.Println("MATCH: password is correct")
    }

    fmt.Println("\n=== Hashing the Same Password Twice ===")
    hash2, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
    if err != nil {
        panic(err)
    }
    fmt.Println("hash1 == hash2:", string(hash) == string(hash2), "(different salts every time)")
    fmt.Println("but both still verify the same password:", bcrypt.CompareHashAndPassword(hash2, password) == nil)
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

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

**Learning Objectives:**
- ✅ Hash a password with `bcrypt.GenerateFromPassword`
- ✅ Verify a password with `bcrypt.CompareHashAndPassword`
- ✅ See that bcrypt salts automatically - the same password never produces the same hash twice

---

## Exercise 2: bcrypt Rejects a Wrong Password

**Objective:** Capture the real error bcrypt produces for a mismatched password

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level32-exercise2
cd ~/projects/level32-exercise2
go mod init level32.example/exercise2
go get golang.org/x/crypto/bcrypt@v0.55.0
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "fmt"

    "golang.org/x/crypto/bcrypt"
)

func main() {
    correct := []byte("correct-horse-battery-staple")
    wrong := []byte("Tr0ub4dor&3")

    hash, err := bcrypt.GenerateFromPassword(correct, bcrypt.DefaultCost)
    if err != nil {
        panic(err)
    }

    fmt.Println("=== Comparing the Hash Against a Wrong Password ===")
    err = bcrypt.CompareHashAndPassword(hash, wrong)
    if err != nil {
        fmt.Printf("Rejected as expected.\nReal error: %v\n", err)
        fmt.Printf("Error type: %T\n", err)
        fmt.Println("err == bcrypt.ErrMismatchedHashAndPassword:", err == bcrypt.ErrMismatchedHashAndPassword)
    } else {
        fmt.Println("UNEXPECTED: wrong password matched!")
    }
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Comparing the Hash Against a Wrong Password ===
Rejected as expected.
Real error: crypto/bcrypt: hashedPassword is not the hash of the given password
Error type: *errors.errorString
err == bcrypt.ErrMismatchedHashAndPassword: true
```

**Learning Objectives:**
- ✅ See bcrypt's real, exact error text for a wrong password
- ✅ Compare against the specific sentinel error `bcrypt.ErrMismatchedHashAndPassword`
- ✅ Understand that a mismatch is a normal, expected outcome - not a panic-worthy failure

---

## Exercise 3: Generating a Signed JWT With Custom Claims

**Objective:** Issue an HS256-signed JWT carrying both registered and custom claims

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level32-exercise3
cd ~/projects/level32-exercise3
go mod init level32.example/exercise3
go get github.com/golang-jwt/jwt/v5@v5.3.1
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "time"

    "github.com/golang-jwt/jwt/v5"
)

type AppClaims struct {
    jwt.RegisteredClaims
    Username string `json:"username"`
    Role     string `json:"role"`
}

var secret = []byte("this-is-a-demo-secret-do-not-use-in-production")

func main() {
    now := time.Now()
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
    signed, err := token.SignedString(secret)
    if err != nil {
        panic(err)
    }

    fmt.Println("=== Signed JWT ===")
    fmt.Println(signed)

    dots := 0
    for _, r := range signed {
        if r == '.' {
            dots++
        }
    }
    fmt.Printf("\nNumber of '.' separators: %d (%d segments: header.payload.signature)\n", dots, dots+1)
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Signed JWT ===
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJsZXZlbDMyLWF1dGgtZGVtbyIsInN1YiI6InVzZXItNDIiLCJleHAiOjE3ODgyODM3NjYsImlhdCI6MTc4ODI4Mjg2NiwidXNlcm5hbWUiOiJhbGljZSIsInJvbGUiOiJtZW1iZXIifQ.5-zWwkglxYo1M1hf9FmK1vKc26mpBNFDsj1RYKi50mo

Number of '.' separators: 2 (3 segments: header.payload.signature)
```

Note: your own run will produce a different token string, since `exp`/`iat` are timestamps of when you ran it and the HMAC signature depends on them - the shape (3 segments) is what matters.

**Learning Objectives:**
- ✅ Embed `jwt.RegisteredClaims` alongside custom fields in your own claims struct
- ✅ Sign a token with `jwt.NewWithClaims` + `SignedString` using HS256
- ✅ Confirm a JWT is exactly three dot-separated segments

---

## Exercise 4: Decoding a JWT's Payload Without the Secret

**Objective:** Prove that a JWT's payload is only encoded, never encrypted

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level32-exercise4
cd ~/projects/level32-exercise4
go mod init level32.example/exercise4
go get github.com/golang-jwt/jwt/v5@v5.3.1
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "encoding/base64"
    "encoding/json"
    "fmt"
    "strings"
    "time"

    "github.com/golang-jwt/jwt/v5"
)

type AppClaims struct {
    jwt.RegisteredClaims
    Username string `json:"username"`
    Role     string `json:"role"`
}

var secret = []byte("this-is-a-demo-secret-do-not-use-in-production")

func main() {
    now := time.Now()
    claims := AppClaims{
        RegisteredClaims: jwt.RegisteredClaims{
            Subject:   "user-42",
            ExpiresAt: jwt.NewNumericDate(now.Add(15 * time.Minute)),
        },
        Username: "alice",
        Role:     "member",
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    signed, err := token.SignedString(secret)
    if err != nil {
        panic(err)
    }

    fmt.Println("=== Full Token ===")
    fmt.Println(signed)

    segments := strings.Split(signed, ".")
    fmt.Printf("\nSegments: %d (header.payload.signature)\n", len(segments))

    headerJSON, err := base64.RawURLEncoding.DecodeString(segments[0])
    if err != nil {
        panic(err)
    }
    fmt.Println("\n=== Decoded Header (no secret used) ===")
    fmt.Println(string(headerJSON))

    payloadJSON, err := base64.RawURLEncoding.DecodeString(segments[1])
    if err != nil {
        panic(err)
    }
    var pretty map[string]interface{}
    if err := json.Unmarshal(payloadJSON, &pretty); err != nil {
        panic(err)
    }
    prettyBytes, _ := json.MarshalIndent(pretty, "", "  ")
    fmt.Println("\n=== Decoded Payload (no secret used!) ===")
    fmt.Println(string(prettyBytes))

    fmt.Println("\nAnyone holding this token can read the payload above WITHOUT the signing secret.")
    fmt.Println("The secret is only required to verify the SIGNATURE, never to read the claims.")
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output** (yours will have different timestamps/signature, but identical shape):

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

**Learning Objectives:**
- ✅ Split a JWT into its three segments and base64url-decode them directly
- ✅ Prove, by direct observation, that the payload is readable without any secret
- ✅ Internalize why secrets must never be placed inside JWT claims

---

## Exercise 5: Validating a Legitimate JWT

**Objective:** Parse and verify a correctly-signed, unexpired token

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level32-exercise5
cd ~/projects/level32-exercise5
go mod init level32.example/exercise5
go get github.com/golang-jwt/jwt/v5@v5.3.1
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "time"

    "github.com/golang-jwt/jwt/v5"
)

type AppClaims struct {
    jwt.RegisteredClaims
    Username string `json:"username"`
    Role     string `json:"role"`
}

var secret = []byte("this-is-a-demo-secret-do-not-use-in-production")

func main() {
    now := time.Now()
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
    signed, err := token.SignedString(secret)
    if err != nil {
        panic(err)
    }

    fmt.Println("=== Validating a Legitimate Token ===")
    parsed, err := jwt.ParseWithClaims(signed, &AppClaims{}, func(t *jwt.Token) (interface{}, error) {
        if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
        }
        return secret, nil
    })
    if err != nil {
        fmt.Println("REJECTED:", err)
        return
    }
    if !parsed.Valid {
        fmt.Println("Token parsed but reported invalid")
        return
    }

    got := parsed.Claims.(*AppClaims)
    fmt.Println("ACCEPTED: token is valid")
    fmt.Println("Subject:", got.Subject)
    fmt.Println("Username:", got.Username)
    fmt.Println("Role:", got.Role)
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Validating a Legitimate Token ===
ACCEPTED: token is valid
Subject: user-42
Username: alice
Role: member
```

**Learning Objectives:**
- ✅ Use `jwt.ParseWithClaims` with a keyfunc that supplies the shared secret
- ✅ Recover strongly-typed claims back out of `parsed.Claims`
- ✅ Confirm a correctly-signed, unexpired token is accepted

---

## Exercise 6: Rejecting a Tampered JWT

**Objective:** Capture the real signature-verification error for a tampered token

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level32-exercise6
cd ~/projects/level32-exercise6
go mod init level32.example/exercise6
go get github.com/golang-jwt/jwt/v5@v5.3.1
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "errors"
    "fmt"
    "strings"
    "time"

    "github.com/golang-jwt/jwt/v5"
)

type AppClaims struct {
    jwt.RegisteredClaims
    Username string `json:"username"`
}

var secret = []byte("this-is-a-demo-secret-do-not-use-in-production")

func main() {
    now := time.Now()
    claims := AppClaims{
        RegisteredClaims: jwt.RegisteredClaims{
            Subject:   "user-42",
            ExpiresAt: jwt.NewNumericDate(now.Add(15 * time.Minute)),
        },
        Username: "alice",
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    signed, err := token.SignedString(secret)
    if err != nil {
        panic(err)
    }

    parts := strings.Split(signed, ".")
    sig := []byte(parts[2])
    if sig[len(sig)-1] == 'A' {
        sig[len(sig)-1] = 'B'
    } else {
        sig[len(sig)-1] = 'A'
    }
    tampered := parts[0] + "." + parts[1] + "." + string(sig)

    fmt.Println("=== Original Token ===")
    fmt.Println(signed)
    fmt.Println("\n=== Tampered Token (last signature character flipped) ===")
    fmt.Println(tampered)

    fmt.Println("\n=== Attempting to Validate the Tampered Token ===")
    _, err = jwt.ParseWithClaims(tampered, &AppClaims{}, func(t *jwt.Token) (interface{}, error) {
        return secret, nil
    })
    if err != nil {
        fmt.Println("REJECTED (as it must be). Real error:", err)
        fmt.Println("errors.Is(err, jwt.ErrTokenSignatureInvalid):", errors.Is(err, jwt.ErrTokenSignatureInvalid))
    } else {
        fmt.Println("UNEXPECTED: tampered token was accepted!")
    }
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output** (your token text will differ, the error text will not):

```
=== Attempting to Validate the Tampered Token ===
REJECTED (as it must be). Real error: token signature is invalid: signature is invalid
errors.Is(err, jwt.ErrTokenSignatureInvalid): true
```

**Learning Objectives:**
- ✅ See that flipping even one character anywhere in a signed token invalidates it
- ✅ Capture the library's real, exact rejection error text
- ✅ Use `errors.Is` to check against the library's sentinel error, not string matching

---

## Exercise 7: Rejecting an Expired JWT

**Objective:** Capture the real expiration error for a token whose `exp` has passed

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level32-exercise7
cd ~/projects/level32-exercise7
go mod init level32.example/exercise7
go get github.com/golang-jwt/jwt/v5@v5.3.1
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "errors"
    "fmt"
    "time"

    "github.com/golang-jwt/jwt/v5"
)

type AppClaims struct {
    jwt.RegisteredClaims
    Username string `json:"username"`
}

var secret = []byte("this-is-a-demo-secret-do-not-use-in-production")

func main() {
    past := time.Now().Add(-1 * time.Hour)
    claims := AppClaims{
        RegisteredClaims: jwt.RegisteredClaims{
            Subject:   "user-42",
            IssuedAt:  jwt.NewNumericDate(past.Add(-15 * time.Minute)),
            ExpiresAt: jwt.NewNumericDate(past),
        },
        Username: "alice",
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    signed, err := token.SignedString(secret)
    if err != nil {
        panic(err)
    }

    fmt.Println("=== Token That Expired One Hour Ago ===")
    fmt.Println(signed)

    fmt.Println("\n=== Attempting to Validate the Expired Token ===")
    _, err = jwt.ParseWithClaims(signed, &AppClaims{}, func(t *jwt.Token) (interface{}, error) {
        return secret, nil
    })
    if err != nil {
        fmt.Println("REJECTED (as it must be). Real error:", err)
        fmt.Println("errors.Is(err, jwt.ErrTokenExpired):", errors.Is(err, jwt.ErrTokenExpired))
    } else {
        fmt.Println("UNEXPECTED: expired token was accepted!")
    }
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output** (your token text will differ, the error text will not):

```
=== Attempting to Validate the Expired Token ===
REJECTED (as it must be). Real error: token has invalid claims: token is expired
errors.Is(err, jwt.ErrTokenExpired): true
```

**Learning Objectives:**
- ✅ Sign a token whose `exp` claim is already in the past
- ✅ See that `jwt.ParseWithClaims` checks expiration automatically, with no manual comparison needed
- ✅ Recognize the library's real expiration error text and its sentinel `jwt.ErrTokenExpired`

---

## Exercise 8: Rejecting the alg:none Attack

**Objective:** Forge an unsigned `alg:none` token and confirm a correctly-written verifier rejects it

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level32-exercise8
cd ~/projects/level32-exercise8
go mod init level32.example/exercise8
go get github.com/golang-jwt/jwt/v5@v5.3.1
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "encoding/base64"
    "encoding/json"
    "fmt"
    "time"

    "github.com/golang-jwt/jwt/v5"
)

type AppClaims struct {
    jwt.RegisteredClaims
    Username string `json:"username"`
}

var secret = []byte("this-is-a-demo-secret-do-not-use-in-production")

func b64(v interface{}) string {
    b, err := json.Marshal(v)
    if err != nil {
        panic(err)
    }
    return base64.RawURLEncoding.EncodeToString(b)
}

func main() {
    header := map[string]string{"alg": "none", "typ": "JWT"}
    claims := AppClaims{
        RegisteredClaims: jwt.RegisteredClaims{
            Subject:   "user-42",
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
        },
        Username: "attacker-pretending-to-be-alice",
    }

    forged := b64(header) + "." + b64(claims) + "."
    fmt.Println("=== Forged alg:none Token (empty signature) ===")
    fmt.Println(forged)

    fmt.Println("\n=== Attempting to Validate the Forged Token (keyfunc checks the method) ===")
    _, err := jwt.ParseWithClaims(forged, &AppClaims{}, func(t *jwt.Token) (interface{}, error) {
        if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v (only HMAC is accepted)", t.Header["alg"])
        }
        return secret, nil
    })
    if err != nil {
        fmt.Println("REJECTED (correct verifier behavior). Real error:", err)
    } else {
        fmt.Println("VULNERABLE: forged alg:none token was accepted!")
    }
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output** (timestamps will differ, the rejection text will not):

```
=== Attempting to Validate the Forged Token (keyfunc checks the method) ===
REJECTED (correct verifier behavior). Real error: token is unverifiable: error while executing keyfunc: unexpected signing method: none (only HMAC is accepted)
```

**Learning Objectives:**
- ✅ Understand exactly how the historical `alg:none` attack was constructed
- ✅ See why the keyfunc, not the token's own header, must decide the accepted algorithm
- ✅ Confirm the defensive keyfunc pattern actually rejects the forged token, for real

---

## Exercise 9: JWT Auth Middleware, Tested End-to-End

**Objective:** Build middleware that extracts, validates, and injects JWT claims into the request context

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level32-exercise9
cd ~/projects/level32-exercise9
go mod init level32.example/exercise9
go get github.com/golang-jwt/jwt/v5@v5.3.1
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "context"
    "encoding/json"
    "errors"
    "fmt"
    "net/http"
    "net/http/httptest"
    "strings"
    "time"

    "github.com/golang-jwt/jwt/v5"
)

type claimsContextKey struct{} // unexported, unique type - avoids collisions (Level 24 pattern)

type AppClaims struct {
    jwt.RegisteredClaims
    Username string `json:"username"`
}

var secret = []byte("this-is-a-demo-secret-do-not-use-in-production")

func makeToken(username string, ttl time.Duration) string {
    now := time.Now()
    claims := AppClaims{
        RegisteredClaims: jwt.RegisteredClaims{
            Subject:   username,
            IssuedAt:  jwt.NewNumericDate(now),
            ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
        },
        Username: username,
    }
    tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    signed, _ := tok.SignedString(secret)
    return signed
}

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
            msg := "invalid token"
            if errors.Is(err, jwt.ErrTokenExpired) {
                msg = "token is expired"
            }
            http.Error(w, fmt.Sprintf(`{"error":%q}`, msg), http.StatusUnauthorized)
            return
        }

        ctx := context.WithValue(r.Context(), claimsContextKey{}, claims)
        next(w, r.WithContext(ctx))
    }
}

func profileHandler(w http.ResponseWriter, r *http.Request) {
    claims, ok := r.Context().Value(claimsContextKey{}).(*AppClaims)
    if !ok {
        http.Error(w, `{"error":"no claims in context"}`, http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"username": claims.Username})
}

func doRequest(mux http.Handler, authHeader string) *httptest.ResponseRecorder {
    req := httptest.NewRequest(http.MethodGet, "/profile", nil)
    if authHeader != "" {
        req.Header.Set("Authorization", authHeader)
    }
    rec := httptest.NewRecorder()
    mux.ServeHTTP(rec, req)
    return rec
}

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("GET /profile", AuthMiddleware(profileHandler))

    valid := makeToken("alice", 15*time.Minute)

    fmt.Println("=== Case 1: Valid Token ===")
    rec := doRequest(mux, "Bearer "+valid)
    fmt.Println("Status:", rec.Code)
    fmt.Println("Body:", strings.TrimSpace(rec.Body.String()))

    fmt.Println("\n=== Case 2: Missing Authorization Header ===")
    rec = doRequest(mux, "")
    fmt.Println("Status:", rec.Code)
    fmt.Println("Body:", strings.TrimSpace(rec.Body.String()))

    fmt.Println("\n=== Case 3: Malformed Header (no 'Bearer ' prefix) ===")
    rec = doRequest(mux, valid)
    fmt.Println("Status:", rec.Code)
    fmt.Println("Body:", strings.TrimSpace(rec.Body.String()))

    fmt.Println("\n=== Case 4: Tampered Token ===")
    tampered := valid[:len(valid)-2] + "xx"
    rec = doRequest(mux, "Bearer "+tampered)
    fmt.Println("Status:", rec.Code)
    fmt.Println("Body:", strings.TrimSpace(rec.Body.String()))

    fmt.Println("\n=== Case 5: Expired Token ===")
    expired := makeToken("alice", -1*time.Minute)
    rec = doRequest(mux, "Bearer "+expired)
    fmt.Println("Status:", rec.Code)
    fmt.Println("Body:", strings.TrimSpace(rec.Body.String()))
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

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

**Learning Objectives:**
- ✅ Build middleware following Level 27's `func(http.HandlerFunc) http.HandlerFunc` wrapping pattern
- ✅ Inject parsed claims into the request context with an unexported key type (Level 24)
- ✅ Cover one success path and four distinct real-world failure paths with `httptest`

---

## Exercise 10: Comprehensive Practice - A Login API

**Objective:** Build a complete, tested "login" API: bcrypt-checked login issuing a JWT, and a JWT-protected profile route

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level32-exercise10
cd ~/projects/level32-exercise10
go mod init level32.example/exercise10
go get github.com/golang-jwt/jwt/v5@v5.3.1
go get golang.org/x/crypto/bcrypt@v0.55.0
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "net/http/httptest"
    "strings"
    "time"

    "github.com/golang-jwt/jwt/v5"
    "golang.org/x/crypto/bcrypt"
)

type claimsKey struct{} // unexported, unique type - avoids collisions (Level 24 pattern)

var jwtSecret = []byte("final-exercise-demo-secret-change-me")

type User struct {
    Username     string
    PasswordHash []byte
}

type UserStore struct {
    users map[string]User
}

func NewUserStore() *UserStore {
    return &UserStore{users: make(map[string]User)}
}

func (s *UserStore) Register(username, password string) error {
    hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }
    s.users[username] = User{Username: username, PasswordHash: hash}
    return nil
}

func (s *UserStore) Verify(username, password string) bool {
    u, ok := s.users[username]
    if !ok {
        return false
    }
    return bcrypt.CompareHashAndPassword(u.PasswordHash, []byte(password)) == nil
}

type AppClaims struct {
    jwt.RegisteredClaims
    Username string `json:"username"`
}

func issueToken(username string) (string, error) {
    now := time.Now()
    claims := AppClaims{
        RegisteredClaims: jwt.RegisteredClaims{
            Subject:   username,
            IssuedAt:  jwt.NewNumericDate(now),
            ExpiresAt: jwt.NewNumericDate(now.Add(15 * time.Minute)),
        },
        Username: username,
    }
    tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return tok.SignedString(jwtSecret)
}

type loginRequest struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

func loginHandler(store *UserStore) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var req loginRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
            return
        }
        if !store.Verify(req.Username, req.Password) {
            http.Error(w, `{"error":"invalid username or password"}`, http.StatusUnauthorized)
            return
        }
        token, err := issueToken(req.Username)
        if err != nil {
            http.Error(w, `{"error":"could not issue token"}`, http.StatusInternalServerError)
            return
        }
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string]string{"token": token})
    }
}

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        header := r.Header.Get("Authorization")
        if !strings.HasPrefix(header, "Bearer ") {
            http.Error(w, `{"error":"missing or malformed Authorization header"}`, http.StatusUnauthorized)
            return
        }
        tokenString := strings.TrimPrefix(header, "Bearer ")
        claims := &AppClaims{}
        token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
            if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("unexpected signing method")
            }
            return jwtSecret, nil
        })
        if err != nil || !token.Valid {
            http.Error(w, `{"error":"invalid or expired token"}`, http.StatusUnauthorized)
            return
        }
        ctx := context.WithValue(r.Context(), claimsKey{}, claims)
        next(w, r.WithContext(ctx))
    }
}

func profileHandler(w http.ResponseWriter, r *http.Request) {
    claims := r.Context().Value(claimsKey{}).(*AppClaims)
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"logged_in_as": claims.Username})
}

func newMux() http.Handler {
    store := NewUserStore()
    if err := store.Register("alice", "S3cur3P@ssw0rd!"); err != nil {
        panic(err)
    }

    mux := http.NewServeMux()
    mux.HandleFunc("POST /login", loginHandler(store))
    mux.HandleFunc("GET /profile", authMiddleware(profileHandler))
    return mux
}

func main() {
    mux := newMux()
    server := httptest.NewServer(mux)
    defer server.Close()
    client := server.Client()

    fmt.Println("=== Step 1: Login With Correct Credentials ===")
    body, _ := json.Marshal(loginRequest{Username: "alice", Password: "S3cur3P@ssw0rd!"})
    resp, err := client.Post(server.URL+"/login", "application/json", bytes.NewReader(body))
    if err != nil {
        panic(err)
    }
    var loginResp map[string]string
    if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
        panic(err)
    }
    resp.Body.Close()
    fmt.Println("Status:", resp.StatusCode)
    token := loginResp["token"]
    fmt.Println("Received a token:", token != "")

    fmt.Println("\n=== Step 2: Login With the Wrong Password ===")
    badBody, _ := json.Marshal(loginRequest{Username: "alice", Password: "wrong-password"})
    resp2, err := client.Post(server.URL+"/login", "application/json", bytes.NewReader(badBody))
    if err != nil {
        panic(err)
    }
    b2 := new(bytes.Buffer)
    b2.ReadFrom(resp2.Body)
    resp2.Body.Close()
    fmt.Println("Status:", resp2.StatusCode)
    fmt.Println("Body:", strings.TrimSpace(b2.String()))

    fmt.Println("\n=== Step 3: Access /profile Without a Token ===")
    resp3, err := client.Get(server.URL + "/profile")
    if err != nil {
        panic(err)
    }
    b3 := new(bytes.Buffer)
    b3.ReadFrom(resp3.Body)
    resp3.Body.Close()
    fmt.Println("Status:", resp3.StatusCode)
    fmt.Println("Body:", strings.TrimSpace(b3.String()))

    fmt.Println("\n=== Step 4: Access /profile With the Real Token ===")
    req4, _ := http.NewRequest(http.MethodGet, server.URL+"/profile", nil)
    req4.Header.Set("Authorization", "Bearer "+token)
    resp4, err := client.Do(req4)
    if err != nil {
        panic(err)
    }
    b4 := new(bytes.Buffer)
    b4.ReadFrom(resp4.Body)
    resp4.Body.Close()
    fmt.Println("Status:", resp4.StatusCode)
    fmt.Println("Body:", strings.TrimSpace(b4.String()))
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Step 1: Login With Correct Credentials ===
Status: 200
Received a token: true

=== Step 2: Login With the Wrong Password ===
Status: 401
Body: {"error":"invalid username or password"}

=== Step 3: Access /profile Without a Token ===
Status: 401
Body: {"error":"missing or malformed Authorization header"}

=== Step 4: Access /profile With the Real Token ===
Status: 200
Body: {"logged_in_as":"alice"}
```

**Learning Objectives:**
- ✅ Combine bcrypt password verification with JWT issuance in one `POST /login` handler
- ✅ Protect a second route with the auth middleware built in Exercise 9
- ✅ Test a full, realistic login-then-access flow end-to-end with `httptest.NewServer`

---

## Bonus Challenges

### Challenge 1: In-Memory Refresh Token Store

Build a minimal refresh-token system: `POST /login` returns both a short-lived access token and a longer-lived refresh token (a random opaque string, not a JWT); a new `POST /refresh` endpoint accepts a refresh token and, if it's found in an in-memory store and not expired, issues a fresh access token.

```bash
mkdir -p ~/projects/level32-bonus1
cd ~/projects/level32-bonus1
go mod init level32.example/bonus1
```

**Hints:**
- A `map[string]refreshTokenRecord{UserID string; ExpiresAt time.Time}` guarded by a `sync.Mutex` (Level 23) is enough for an in-memory store
- Generate the opaque refresh token string with `crypto/rand`, not `math/rand` - it's a credential
- On `/refresh`, look up the token, check it hasn't expired, then call your existing `issueToken` for a new access token
- Consider what happens to an already-used refresh token: real systems often delete-and-reissue ("rotate") it and treat a reused old one as a signal of theft

### Challenge 2: Role-Based Claim Checked by a Second Middleware Layer

Add a `Role string` custom claim to `AppClaims`. Build a second middleware, `RequireRole(role string, next http.HandlerFunc) http.HandlerFunc`, that runs *after* `AuthMiddleware` (so it can read claims already placed in the context) and returns `403 Forbidden` if the claim's role doesn't match.

```bash
mkdir -p ~/projects/level32-bonus2
cd ~/projects/level32-bonus2
go mod init level32.example/bonus2
```

**Hints:**
- Compose the two middlewares like Level 27's logging example: `AuthMiddleware(RequireRole("admin", adminOnlyHandler))`
- Read the claims the same way `profileHandler` does - a type assertion on the context value
- Test three cases: no token (401 from `AuthMiddleware`), valid token with the wrong role (403 from `RequireRole`), valid token with the right role (200)

### Challenge 3: Rate-Limiting Failed Login Attempts

Extend Exercise 10's login handler to track failed attempts per username (recall the token-bucket/counting patterns from Level 22's `select` exercises and Level 27's middleware) and return `429 Too Many Requests` after too many failures within a short window, resetting on a successful login.

```bash
mkdir -p ~/projects/level32-bonus3
cd ~/projects/level32-bonus3
go mod init level32.example/bonus3
```

**Hints:**
- A `map[string][]time.Time` of recent failure timestamps per username, pruned of anything older than your window, is enough for a simple version
- Guard the map with a `sync.Mutex` - login requests can arrive concurrently
- Reset (delete) a username's failure history the moment a login for that username succeeds

---

## What You've Learned

After completing these 10 exercises, you can:

✅ Hash and verify passwords with bcrypt, and recognize its real mismatch error
✅ Explain and demonstrate that a JWT's payload is encoded, not encrypted
✅ Issue an HS256-signed JWT with registered and custom claims
✅ Validate a JWT and correctly reject a tampered signature and an expired token, with real captured errors
✅ Explain the historical `alg:none` vulnerability and verify a correct keyfunc defeats it
✅ Build and test JWT auth middleware that bridges Level 27's middleware pattern and Level 24's context values
✅ Build and fully test a small login API combining bcrypt and JWT issuance

---

## Next Level

Level 33: Logging
- Structured logging in Go
- Log levels and when to use each
- Correlating logs with request IDs (building on this level's context-carried claims)

Great work! You've built the security foundation almost every real API depends on! 🚀
