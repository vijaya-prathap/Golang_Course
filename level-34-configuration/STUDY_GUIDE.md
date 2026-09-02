# Level 34: Study Guide & Visual Reference

## 📚 Learning Path

### Week 1: Configuration Fundamentals
```
Day 1:  Why configuration matters; os.Getenv vs os.LookupEnv
Day 2:  The flag package: defaults, flag.Parse(), verifying multiple runs
Day 3:  Designing a Config struct and a LoadConfig() function
Day 4:  Config precedence: flags > env > defaults
Day 5:  Fail-fast validation at startup
Day 6:  Type conversion for numeric/boolean env vars
Day 7:  Hand-rolled .env file parsing
```

### Week 2: Practice & Application
```
Day 1:  Exercises 1-3 (Getenv/LookupEnv, flags, Config struct)
Day 2:  Exercises 4-5 (precedence resolver, fail-fast validation)
Day 3:  Exercises 6-7 (type conversion, .env parser)
Day 4:  Exercises 8-9 (typed .env loading, don't-clobber convention)
Day 5:  Exercise 10 (comprehensive service config loader)
Day 6:  Bonus challenges
Day 7:  Review & consolidation
```

---

## 🎯 os.Getenv vs os.LookupEnv

| | `os.Getenv(key)` | `os.LookupEnv(key)` |
|---|---|---|
| **Return type** | `string` | `(string, bool)` |
| **Unset variable** | `""` | `"", false` |
| **Set to `""`** | `""` | `"", true` |
| **Set to `"x"`** | `"x"` | `"x", true` |
| **Can tell unset from empty?** | ❌ No | ✅ Yes (via `ok`) |
| **Typical use** | "give me a value or I'll default it myself" | "I need to know if this was ever configured at all" |

```
os.Getenv("KEY")               os.LookupEnv("KEY")
        │                              │
        ▼                              ▼
   always a string             (string, bool)
   "" means EITHER             ok=false -> definitely unset
   unset OR empty              ok=true  -> was set (maybe to "")
```

---

## 🚀 The flag Package Lifecycle

```
1. flag.String("host", "localhost", "usage text")   -> returns *string
2. flag.Int("port", 8080, "usage text")              -> returns *int
3. flag.Bool("debug", false, "usage text")           -> returns *bool
4. flag.Parse()                                       -> MUST run before reading values
5. *host, *port, *debug now hold either the passed value or the default
```

```
go run main.go                          -> every flag = its default
go run main.go -port=9090               -> port=9090, everything else = default
go run main.go -host=x -port=9 -debug   -> all three explicitly set
```

**Rule:** nothing is read from a flag pointer until after `flag.Parse()` runs.

---

## 🗂️ Config Precedence Flow

```
                    ┌─────────────────────┐
                    │  Hardcoded Default   │  lowest priority
                    │  (always present)     │
                    └──────────┬───────────┘
                               │ overridden by...
                    ┌──────────▼───────────┐
                    │ Environment Variable  │  middle priority
                    │ (os.LookupEnv finds it)│
                    └──────────┬───────────┘
                               │ overridden by...
                    ┌──────────▼───────────┐
                    │  Command-Line Flag    │  highest priority
                    │  (explicitly passed)  │
                    └──────────┬───────────┘
                               │
                               ▼
                     final resolved value
```

```go
func resolve(flagVal string, flagSet bool, envKey, def string) (string, string) {
    if flagSet {
        return flagVal, "flag"
    }
    if envVal, ok := os.LookupEnv(envKey); ok {
        return envVal, "env"
    }
    return def, "default"
}
```

### Verified Truth Table

| Flag passed? | Env var set? | Winning value | Source |
|---|---|---|---|
| No | No | default | `default` |
| No | Yes | env value | `env` |
| Yes | No | flag value | `flag` |
| Yes | Yes | flag value | `flag` |

---

## 🚨 Fail-Fast Validation Flow

```
        Start
          │
          ▼
   cfg := LoadConfig()
          │
          ▼
   err := cfg.Validate()
          │
     ┌────┴────┐
     │         │
   err!=nil  err==nil
     │         │
     ▼         ▼
 log.Fatal   start server,
 (exit now,  confident cfg
  never      is complete
  start)
```

The whole point: **every required field is checked once, in one place, before anything else runs.** No handler ever has to defensively check "is this config value actually set?" - by the time any handler runs, `Validate()` has already guaranteed it.

```go
func (c Config) Validate() error {
    if c.DatabaseURL == "" {
        return errors.New("DATABASE_URL is required but not set")
    }
    if c.Port == "" {
        return errors.New("PORT is required but not set")
    }
    return nil
}
```

---

## 🔢 Type Conversion Cheat Sheet

| Env var is really a... | Parse with | Fallback behavior |
|---|---|---|
| `int` | `strconv.Atoi(raw)` | Only call it if the key exists; keep the typed default otherwise |
| `bool` | `strconv.ParseBool(raw)` | Accepts `1,t,T,TRUE,true,True,0,f,F,FALSE,false,False` only |
| `time.Duration` | `time.ParseDuration(raw)` | Accepts `"5s"`, `"2m30s"`, `"1h"` etc. |
| Anything else | `raw` itself | env vars are always strings - no parsing needed |

```go
n, err := strconv.Atoi(raw)
if err != nil {
    return 0, fmt.Errorf("invalid int for %s: %q: %w", key, raw, err)
}
```

**Never ignore the error.** A malformed `MAX_CONNECTIONS=not-a-number` should be a startup failure, not a silently-zeroed setting.

---

## 📄 .env File Parsing Diagram

```
.env file:
┌────────────────────────────────────────────┐
│ # This is a comment - skip it               │
│ APP_NAME=DemoService                        │  ← KEY=VALUE
│                                              │  ← blank line - skip it
│ APP_PORT=4000                               │
│                                              │
│ CONNECTION_STRING=user=admin;password=x=y   │  ← SplitN(line, "=", 2) keeps
└────────────────────────────────────────────┘    everything after the FIRST "="
                    │
                    ▼  bufio.Scanner, line by line
        ┌───────────────────────┐
        │ line == ""?  -> skip   │
        │ starts with #? -> skip │
        │ SplitN(line,"=",2)     │
        │ trim key & value       │
        └───────────┬───────────┘
                    ▼
        map[string]string{
            "APP_NAME": "DemoService",
            "APP_PORT": "4000",
            "CONNECTION_STRING": "user=admin;password=x=y",
        }
```

### The "Don't Clobber Real Env Vars" Rule

```
Before loading .env:              After loading .env:
┌─────────────────────┐          ┌─────────────────────┐
│ LOG_LEVEL=error      │  ─────►  │ LOG_LEVEL=error      │  UNCHANGED - was
│ (set by the shell)   │          │                      │  already real
├─────────────────────┤          ├─────────────────────┤
│ APP_ENV: unset       │  ─────►  │ APP_ENV=development  │  FILLED IN from
│                      │          │                      │  .env - was unset
└─────────────────────┘          └─────────────────────┘

for key, value := range dotenvValues {
    if _, alreadySet := os.LookupEnv(key); alreadySet {
        continue // real environment wins
    }
    os.Setenv(key, value)
}
```

---

## 🧩 Full Precedence Chain (Exercise 10)

```
defaults  <  .env file  <  real env vars  <  command-line flags
  (base)      (dev convenience,        (deployment-level      (per-run
              fills gaps only)          configuration)         override)
```

```
Field        Value                  Source
-----        -----                  ------
Port         8080                   default      <- nothing else set it
LogLevel     warn                   env          <- shell env overrode .env's "debug"
APIKey       dotenv-key-abc123      dotenv       <- only .env supplied this
Environment  development            default      <- nothing else set it
```

---

## 🚨 Common Mistakes

### Mistake 1: Silent Zero-Value Fallback for a Required Setting

```go
// ❌ WRONG - empty DatabaseURL becomes "", program "starts," fails later
cfg.DatabaseURL = os.Getenv("DATABASE_URL")

// ✅ RIGHT - treat absence of a required setting as a startup error
if cfg.DatabaseURL == "" {
    return errors.New("DATABASE_URL is required")
}
```

### Mistake 2: Config-Loading Code Scattered Through Business Logic

```go
// ❌ WRONG - os.Getenv calls buried in handlers/services
func ProcessOrder() {
    timeout := os.Getenv("ORDER_TIMEOUT")
}

// ✅ RIGHT - load once, pass a typed Config down
func ProcessOrder(cfg Config) { /* uses cfg.OrderTimeout */ }
```

### Mistake 3: Confusing Unset With Empty String

```go
// ❌ WRONG - Getenv can't tell "never configured" from "explicitly empty"
flags := os.Getenv("FEATURE_FLAGS")

// ✅ RIGHT - use LookupEnv when the distinction matters
if flags, ok := os.LookupEnv("FEATURE_FLAGS"); ok { /* respect it exactly */ }
```

### Mistake 4: Undocumented Required Environment Variables

```
❌ WRONG: a new developer greps the source for every os.Getenv call
✅ RIGHT: a .env.example (or README table) lists every variable and whether it's required
```

---

## 📈 Progression Summary

### Understanding Level 34

Level 34 teaches how one compiled binary adapts to every environment it runs in:

1. **Reading settings** - `os.Getenv`/`os.LookupEnv`, `flag`
2. **Structuring settings** - a typed `Config` struct, one `LoadConfig()` entry point
3. **Layering settings** - flags > env > defaults, tested against every combination
4. **Trusting settings** - fail-fast validation before anything else runs
5. **Typed settings** - `strconv`/`time.ParseDuration` with real error handling
6. **Local convenience** - a hand-rolled `.env` parser that respects the real environment

### Prerequisites for Level 35

Before moving to Level 35 (Docker), you need:

- ✅ Comfortable distinguishing `os.Getenv` from `os.LookupEnv`
- ✅ Comfortable defining and parsing flags with `flag`
- ✅ Can design a `Config` struct with a single `LoadConfig()` function
- ✅ Understand and can implement flags > env > defaults precedence
- ✅ Understand why configuration should fail fast at startup
- ✅ Can parse and convert env var strings into typed values with error handling
- ✅ Can hand-roll a `.env` file parser, including the "don't clobber real env vars" rule

### Ready for Level 35?

Level 35 teaches how to package the binary you've been configuring into a portable container image:
- Writing a `Dockerfile` for a Go binary
- Passing environment variables into a container at `docker run` time - directly building on this level's `os.Getenv`/`LoadConfig` patterns
- Why a "twelve-factor," environment-configured binary is exactly what makes a container image portable across dev/staging/production

---

## ✅ Checklist Before Level 35

- [ ] Can explain the difference between `os.Getenv` and `os.LookupEnv` from memory
- [ ] Can define flags with `flag.String`/`Int`/`Bool` and call `flag.Parse()`
- [ ] Can design a `Config` struct and a `LoadConfig()` function
- [ ] Can implement and test flags > env > defaults precedence
- [ ] Can write a `Validate()` method that fails fast on a missing required setting
- [ ] Can convert env var strings to `int`/`bool`/`time.Duration` with error handling
- [ ] Can hand-roll a `.env` parser that skips comments/blanks and never overrides real env vars
- [ ] Completed 8+ exercises

---

## 💡 Key Takeaways

### The Core Distinction
`os.Getenv` gives you a string with `""` for "unset"; `os.LookupEnv`'s second return value tells you whether it was ever set at all.

### The Precedence
Flags override environment variables, which override hardcoded defaults - in that exact order, every time.

### The Discipline
Validate everything required at startup and fail fast - a broken config should never make it as far as a request handler.

### The Convenience
A `.env` file is for local dev only, is never committed, and never overrides a real environment variable that's already set.

---

## 📚 Next Level

Level 35: Docker
- Packaging your configured Go binary into a container image
- Passing environment variables into a container at `docker run` time
- Why the configuration habits from this level make a Go binary "container-friendly"

You've got configuration down! Keep going! 🚀
