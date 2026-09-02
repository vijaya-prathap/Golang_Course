# Level 32: Authentication & JWT - INDEX

Welcome to **Level 32: Authentication & JWT**! This is where your programs start answering a question every real API has to answer: "who is making this request?"

> ⚠️ **This level requires internet access once**, to fetch `golang.org/x/crypto/bcrypt@v0.55.0` and `github.com/golang-jwt/jwt/v5@v5.3.1` - the only third-party packages used here. Both were verified reachable and every example was actually run - see README.md's introduction for the full note. This is only the third level in the course with this requirement, after Level 28 (Gin) and Level 29 (database/sql).

---

## 📖 What You'll Learn

- ✅ Why passwords are hashed (bcrypt), never encrypted
- ✅ The exact three-segment structure of a JWT, and why the payload is readable by anyone
- ✅ Registered and custom claims, issuing a signed HS256 token
- ✅ Validating a JWT - signature, expiration, and the real errors both produce
- ✅ The historical `alg:none` vulnerability and how a correct verifier defeats it
- ✅ Building JWT auth middleware that injects claims into request context
- ✅ Session-based vs token-based auth tradeoffs, and the refresh-token pattern

---

## 🗂️ Level 32 Materials

### 1. **START_HERE.txt** (Begin Here!)
Quick overview, the internet-access note, and welcome message. Start here first!

### 2. **README.md** (Main Theory)
Comprehensive guide to:
- Hashing vs encryption, and why bcrypt for passwords
- JWT structure, claims, issuing, and validating tokens
- The alg:none attack and defensive verification
- JWT auth middleware bridging Level 24 and Level 27
- Session vs JWT tradeoffs and refresh tokens
- Best practices and common mistakes

**Read Time:** 50-65 minutes
**When:** Start your session

### 3. **EXERCISES.md** (Hands-On Practice)
10 detailed, step-by-step exercises:
1. bcrypt hash and verify round trip
2. bcrypt rejects a wrong password (real error captured)
3. Generating a signed JWT with custom claims
4. Decoding a JWT's payload without the secret
5. Validating a legitimate JWT
6. Rejecting a tampered JWT (real signature error)
7. Rejecting an expired JWT (real expiration error)
8. Rejecting the alg:none attack
9. JWT auth middleware, tested end-to-end
10. Comprehensive practice - a login API

Plus 3 bonus challenges (refresh token store, role-based middleware, login rate limiting).

**Time Commitment:** 5-6 hours
**When:** After reading README

### 4. **STUDY_GUIDE.md** (Visual Learning)
- bcrypt hash-and-compare flow diagram
- JWT structure diagram (encoded vs signed)
- Auth middleware request-flow diagram
- Session-vs-JWT comparison table
- Common mistakes guide

**Read Time:** 35-45 minutes
**When:** Use while doing exercises

### 5. **QUICK_REFERENCE.md** (Cheat Sheet)
- bcrypt and JWT API at a glance
- Middleware skeleton
- Common mistakes table

**Read Time:** 10-15 minutes
**When:** Quick lookup during practice

---

## 🎯 Recommended Learning Path

### Day 1: Passwords (2 hours)
1. Read **START_HERE.txt** (10 min)
2. Read **README.md** Sections 1-2 (25 min)
3. Complete Exercises 1-2 (1 hour 25 min)

### Day 2: JWT Fundamentals (2.5 hours)
1. Read **README.md** Sections 3-5 (35 min)
2. Complete Exercises 3-5 (1 hour 55 min)

### Day 3: Security & Rejections (2 hours)
1. Read **README.md** Section 6 (15 min)
2. Complete Exercises 6-8 (1 hour 45 min)

### Day 4: Middleware & Comprehensive Practice (2.5 hours)
1. Read **README.md** Section 7 (20 min)
2. Complete Exercises 9-10 (2 hours 10 min)

### Day 5: Consolidation (1.5 hours)
1. Read **README.md** Sections 8-11 (30 min)
2. Try bonus challenges
3. Review with **QUICK_REFERENCE.md**

---

## 💡 Key Concepts At A Glance

### Password Storage
```go
hash, _ := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
err := bcrypt.CompareHashAndPassword(hash, candidate) // nil == match
```

### JWT Structure
```
header.payload.signature   -- header and payload are ENCODED, not encrypted
```

### Issuing and Validating
```go
signed, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret)
parsed, err := jwt.ParseWithClaims(signed, &AppClaims{}, keyfunc)
```

### Middleware Context Injection
```go
ctx := context.WithValue(r.Context(), claimsContextKey{}, claims)
```

---

## ✅ Prerequisites

Make sure you've completed **Level 31: Dependency Injection**

You need:
- ✅ Comfort with constructor injection and interfaces as seams
- ✅ Level 27's `func(http.Handler) http.Handler` middleware pattern
- ✅ Level 24's `context.WithValue` and unexported context-key types
- ✅ Comfort fetching a third-party package with `go get` (from Level 28/29)

---

## 🎓 Learning Objectives

By the end of Level 32, you'll be able to:

- ✅ Explain why passwords are hashed, not encrypted, and hash/verify one with bcrypt
- ✅ Explain and demonstrate that a JWT's payload is encoded, not encrypted
- ✅ Issue a signed JWT with registered and custom claims
- ✅ Validate a JWT and correctly reject a tampered or expired one
- ✅ Explain the alg:none attack and defend against it in a keyfunc
- ✅ Build and test JWT auth middleware end-to-end
- ✅ Compare session-based and token-based authentication tradeoffs

---

## 📊 Statistics

- **Main Theory:** README.md covering bcrypt, JWT structure, claims, validation, security, and middleware
- **Exercises:** 10 detailed exercises + 3 bonus challenges
- **Study Guide:** flow diagrams, JWT structure diagram, comparison table
- **Quick Reference:** One-page cheat sheet
- **Practice Time:** 5-6 hours
- **Difficulty:** ⭐⭐⭐⭐ (Advanced)

---

## 🚀 How to Use These Materials

### For Reading
1. Open README.md for comprehensive theory
2. Use STUDY_GUIDE.md alongside for the flow diagrams
3. Reference QUICK_REFERENCE.md for fast lookup

### For Practice
1. Read exercise description
2. Copy provided bash commands (or type them)
3. Follow step-by-step instructions
4. Verify output matches expected results
5. Understand what each step does

### For Review
1. Use QUICK_REFERENCE.md for quick recap
2. Refer to the common mistakes section
3. Practice explaining the alg:none attack and defense out loud

---

## 🆘 Common Questions

**Q: Is a JWT encrypted?**
A: No. A JWT's header and payload are base64url-*encoded*, which anyone can decode with no key at all. Only the signature requires the secret - and it proves the claims weren't altered, not that they're hidden.

**Q: Why not just encrypt passwords instead of hashing them?**
A: Encryption is reversible - anyone with the key (including an attacker who steals both the database and the key) gets every plaintext password back. Hashing is one-way; there is no key to steal that reverses it.

**Q: Why bcrypt instead of SHA-256?**
A: SHA-256 is fast by design, which is exactly wrong for passwords - it lets an attacker try billions of guesses per second against a stolen hash. bcrypt is deliberately slow and salts automatically.

**Q: Can I revoke a JWT before it expires?**
A: Not natively - that's JWT's biggest tradeoff versus sessions. Mitigate with short expirations, refresh tokens, and (if truly needed) a server-side denylist for emergencies.

**Q: What was the alg:none vulnerability?**
A: Some early verifiers trusted the token's own header to decide which algorithm to check with. An attacker could set `alg` to `"none"`, strip the signature, and edit the payload freely. The fix: the verifier - never the token - decides which algorithm is acceptable.

---

## 🎯 Before Moving to Level 33

Make sure you can answer these questions:

- [ ] Why is hashing, not encryption, correct for passwords?
- [ ] Why is a JWT's payload readable without the signing secret?
- [ ] What's the difference between the errors for a tampered token and an expired token?
- [ ] How does the alg:none attack work, and how do you defend against it?
- [ ] How does JWT auth middleware get claims from a header into a handler?

---

## 📚 Recommended Reading Order

For Maximum Learning:

1. **START_HERE.txt** (10 min)
   - Gets you oriented, flags the internet-access requirement

2. **README.md** Sections 1-4 (45 min)
   - Passwords, bcrypt, JWT structure, claims

3. **EXERCISES.md** Exercises 1-5 (2.5 hours)
   - Hash/verify passwords, issue and decode tokens

4. **README.md** Sections 5-7 (35 min)
   - Validation, the alg:none attack, middleware

5. **STUDY_GUIDE.md** (40 min)
   - Study the flow diagrams and comparison table

6. **EXERCISES.md** Exercises 6-10 (3+ hours)
   - Real rejections, middleware, the comprehensive login API

7. **QUICK_REFERENCE.md** (10 min)
   - Create your mental cheat sheet

---

## 🏆 Success Indicators

You've mastered Level 32 when:

- ✅ You reach for bcrypt automatically for any password, without hesitation
- ✅ You can explain "signed, not sealed" to someone else in one sentence
- ✅ You never trust a token's own header to decide how it's verified
- ✅ You can write JWT auth middleware from memory
- ✅ You can weigh session-based vs JWT-based auth for a given project
- ✅ You've completed 8+ exercises
- ✅ You can explain refresh tokens' purpose without needing to build the full system

---

## 🚀 What's Next?

After Level 32, you're ready for:

**Level 33: Logging**
- Structured logging and log levels
- Correlating logs with request/user identity carried through context

---

## 💬 Key Takeaway

> **A JWT is signed, not sealed - anyone can read it, only the issuer can vouch for it. Hash passwords, verify signatures, check expiration, and never let the token tell you how to check itself.**

---

## 📖 Quick Links

- [START_HERE.txt](./START_HERE.txt) - Welcome & Overview
- [README.md](./README.md) - Full Theory
- [EXERCISES.md](./EXERCISES.md) - Hands-On Exercises
- [STUDY_GUIDE.md](./STUDY_GUIDE.md) - Visual Patterns
- [QUICK_REFERENCE.md](./QUICK_REFERENCE.md) - Cheat Sheet

---

Happy learning! Level 32 gives your programs the power to know who's asking - safely! 🎉

*Estimated time to complete Level 32: 5-6 hours*
*Difficulty: ⭐⭐⭐⭐ (Advanced)*
*Next Level: Level 33 - Logging*
