# Level 34: Configuration - Complete Guide

## Introduction

Welcome to Level 34! Level 33 (Logging) taught you how to make your program tell you what it's doing. This level teaches you how to make your program *behave differently* in different places - dev, staging, production - without touching a single line of code or maintaining separate binaries for each environment.

Every real service needs to know things that change depending on where it runs: which database to connect to, which port to listen on, whether debug logging is on, what the API key for a third-party service is. Hardcoding any of that means recompiling every time it changes, or worse, `if environment == "production"` branches scattered through your business logic. Configuration is how you pull all of that out into one place, loaded once at startup.

Go's standard library already has everything you need for this: `os` to read environment variables, `flag` to read command-line arguments, and `strconv`/`time` (from Levels 3 and earlier) to turn the strings they give you into real typed values. This level is entirely standard library - no third-party config packages - because understanding the fundamentals this way is what lets you evaluate (or outgrow) a config library later with your eyes open.

---

## Table of Contents

1. [Why Configuration Matters](#why-configuration-matters)
2. [Environment Variables: Getenv vs LookupEnv](#environment-variables-getenv-vs-lookupenv)
3. [The flag Package](#the-flag-package)
4. [Designing a Config Struct](#designing-a-config-struct)
5. [Config Precedence: Flags > Env > Defaults](#config-precedence-flags--env--defaults)
6. [Validating Configuration at Startup](#validating-configuration-at-startup)
7. [Type Conversion and Validation for Env Vars](#type-conversion-and-validation-for-env-vars)
8. [Hand-Rolled .env File Support](#hand-rolled-env-file-support)
9. [Third-Party Options: Further Reading Only](#third-party-options-further-reading-only)
10. [Best Practices](#best-practices)
11. [Common Mistakes](#common-mistakes)

---

## Why Configuration Matters

Imagine your service needs a database connection string. In dev, that's a local Postgres instance. In staging, it's a shared test database. In production, it's a real database with real credentials that must never appear in source control. If that connection string is hardcoded, you need three different binaries (or three branches of `if` statements) just to run in three places.

Configuration solves this by moving anything environment-specific **out of the code and into something read at startup**:

```go
// ❌ Hardcoded - the same binary can never work in both places
const databaseURL = "postgres://localhost/myapp_dev"

// ✅ Configured - the same binary works everywhere; only the environment changes
databaseURL := os.Getenv("DATABASE_URL")
```

The goal for this whole level: **one compiled binary, deployed unchanged everywhere, that behaves correctly because of what's fed to it at startup** - not because of what branch of code it happens to be running.

---

## Environment Variables: Getenv vs LookupEnv

Environment variables are key-value pairs that live in the process's environment - set by the shell, a `Dockerfile`, a deployment platform, or a CI pipeline. Go gives you two ways to read them, and the difference matters more than it looks.

### os.Getenv: Simple, But Can't Tell "Unset" From "Empty"

```go
value := os.Getenv("APP_NAME")
// value is "" whether APP_NAME was never set, OR was set to an empty string
```

`os.Getenv` always returns a `string`. If the variable isn't set at all, you get `""` - exactly the same thing you'd get if someone ran `APP_NAME= go run main.go` and explicitly set it to empty. From `Getenv`'s return value alone, you cannot tell these two situations apart.

### os.LookupEnv: Distinguishes Unset From Empty

```go
value, ok := os.LookupEnv("APP_NAME")
// ok == false -> APP_NAME was never set
// ok == true, value == "" -> APP_NAME was explicitly set to an empty string
// ok == true, value == "MyApp" -> APP_NAME was set to "MyApp"
```

Verified, showing all three cases in the same run:

```
=== os.Getenv (returns "" if unset) ===
APP_NAME via Getenv: ""

=== os.LookupEnv (distinguishes unset from empty) ===
APP_NAME via LookupEnv: value="" ok=false

=== Now set APP_NAME to an EMPTY string (still "set") ===
APP_NAME via Getenv: ""
APP_NAME via LookupEnv: value="" ok=true

=== Now set APP_NAME to a real value ===
APP_NAME via Getenv: "MyService"
APP_NAME via LookupEnv: value="MyService" ok=true
```

### When the Difference Actually Matters

Most of the time you just want "give me a value, or a default" and `Getenv` (or a small helper wrapping it) is fine. But `LookupEnv` matters whenever "the operator deliberately set this to nothing" needs to mean something different from "the operator never thought about this setting" - for example, an empty `FEATURE_FLAGS` meaning "explicitly disable all flags" versus an unset one meaning "use the built-in default flag set."

---

## The flag Package

Environment variables aren't the only way to configure a program. `flag` reads settings from command-line arguments, which are often more convenient for local runs, one-off overrides, and CLI tools.

### Defining and Parsing Flags

```go
host := flag.String("host", "localhost", "server host")
port := flag.Int("port", 8080, "server port")
debug := flag.Bool("debug", false, "enable debug mode")

flag.Parse() // MUST be called before reading *host, *port, *debug
```

`flag.String`, `flag.Int`, and `flag.Bool` each take three arguments: the flag's name, its **default value**, and a usage string (shown by `-h`). Each returns a pointer - the flag's final value is only filled in once `flag.Parse()` runs, so always call `Parse()` before reading any of them.

### Verified: Same Program, Different Arguments

```go
package main

import (
    "flag"
    "fmt"
)

func main() {
    host := flag.String("host", "localhost", "server host")
    port := flag.Int("port", 8080, "server port")
    debug := flag.Bool("debug", false, "enable debug mode")
    flag.Parse()

    fmt.Printf("host:  %s\n", *host)
    fmt.Printf("port:  %d\n", *port)
    fmt.Printf("debug: %v\n", *debug)
}
```

Run with no arguments (uses every default):

```
host:  localhost
port:  8080
debug: false
```

Run with `-host=example.com -port=9090 -debug=true`:

```
host:  example.com
port:  9090
debug: true
```

Run with only `-port=3000` (everything else still defaults):

```
host:  localhost
port:  3000
debug: false
```

Flags you don't pass simply keep their declared default - you never need an `if` to fill in a "missing" flag yourself.

---

## Designing a Config Struct

Reading individual `os.Getenv` calls scattered across a codebase is exactly the kind of thing this level is trying to prevent. Instead, collect every setting into one typed struct, populated by a single `LoadConfig` function.

```go
type Config struct {
    Host string
    Port string
    Env  string
}

func LoadConfig() Config {
    return Config{
        Host: getEnvOrDefault("APP_HOST", "0.0.0.0"),
        Port: getEnvOrDefault("APP_PORT", "8080"),
        Env:  getEnvOrDefault("APP_ENV", "development"),
    }
}

func getEnvOrDefault(key, fallback string) string {
    if value, ok := os.LookupEnv(key); ok {
        return value
    }
    return fallback
}
```

Verified - with no environment variables set, every field falls back to its default; after setting two of them, only those two change:

```
=== LoadConfig with NO environment variables set (all defaults) ===
{Host:0.0.0.0 Port:8080 Env:development}

=== LoadConfig after setting APP_PORT and APP_ENV ===
{Host:0.0.0.0 Port:9000 Env:production}
```

The rest of the program then depends only on the `Config` struct - never on `os.Getenv` directly. That's the whole point: **one function knows where settings come from; everything else just receives a `Config`.**

---

## Config Precedence: Flags > Env > Defaults

Real services almost always support all three sources of configuration at once, layered by priority:

```
1. Hardcoded defaults    (lowest priority  - always present as a fallback)
2. Environment variables (middle priority  - set per-deployment, e.g. by Docker/Kubernetes)
3. Command-line flags    (highest priority - explicit, wins every time)
```

The idea: defaults make the program runnable out of the box, environment variables let an ops team configure a deployment without touching the command line, and flags let a developer override any single value for one specific run without changing anything else.

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

`flagSet` here comes from `flag.Visit`, which only visits flags the user *actually passed* - this is what lets "flag not passed" mean something different from "flag passed with its default value."

Verified against all four combinations of env-set/unset crossed with flag-set/unset:

```
--- Combo 1: nothing set -> default ---
port=8080 source=default

--- Combo 2: only env set -> env ---
port=7000 source=env

--- Combo 3: only flag set -> flag ---
port=9999 source=flag

--- Combo 4: both env and flag set -> flag wins ---
port=9999 source=flag
```

---

## Validating Configuration at Startup

A missing or malformed setting is far cheaper to catch the moment the program starts than three requests into production traffic. **Fail fast**: validate every required setting up front, and refuse to start if anything is wrong.

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

func main() {
    cfg := LoadConfig()
    if err := cfg.Validate(); err != nil {
        log.Fatal(err) // bridges Level 33: log the failure, then exit - don't limp along
    }
    // ... start the server, confident cfg is complete
}
```

Verified across three cases - missing `DATABASE_URL`, missing `PORT`, and both present:

```
=== Case 1: DATABASE_URL missing (should fail) ===
config error: DATABASE_URL is required but not set
(in a real program: log.Fatal(err) here, before starting the server)
=== Case 2: DATABASE_URL set, PORT missing (should fail) ===
config error: PORT is required but not set
(in a real program: log.Fatal(err) here, before starting the server)
=== Case 3: both set (should succeed) ===
config OK: {DatabaseURL:postgres://localhost/mydb Port:8080}
```

`log.Fatal` (Level 33's territory) calls `os.Exit(1)` after logging - exactly the behavior you want for a startup-time configuration failure: log it clearly, then stop, before any request handler ever runs with a broken config.

---

## Type Conversion and Validation for Env Vars

Every environment variable is a `string` - always, no exceptions. If a setting is really a number, a boolean, or a duration, you must parse it yourself and handle the case where someone typed it wrong.

```go
func getEnvAsInt(key string, fallback int) (int, error) {
    raw, ok := os.LookupEnv(key)
    if !ok {
        return fallback, nil
    }
    n, err := strconv.Atoi(raw) // bridges Level 3's string/number conversions
    if err != nil {
        return 0, fmt.Errorf("invalid int for %s: %q: %w", key, raw, err)
    }
    return n, nil
}
```

The same shape works for `strconv.ParseBool` and `time.ParseDuration`. Verified, including both the happy path and a deliberately bad value for each type:

```
=== MAX_CONNECTIONS unset -> fallback ===
value=10 err=<nil>

=== MAX_CONNECTIONS="50" -> parsed ===
value=50 err=<nil>

=== MAX_CONNECTIONS="not-a-number" -> parse error ===
value=0 err=invalid int for MAX_CONNECTIONS: "not-a-number": strconv.Atoi: parsing "not-a-number": invalid syntax

=== ENABLE_CACHE="true" -> parsed bool ===
value=true err=<nil>

=== ENABLE_CACHE="yes" -> invalid bool ===
value=false err=invalid bool for ENABLE_CACHE: "yes": strconv.ParseBool: parsing "yes": invalid syntax

=== REQUEST_TIMEOUT="5s" -> parsed duration ===
value=5s err=<nil>
```

Note that `strconv.ParseBool` only accepts a specific set of spellings (`1`, `t`, `T`, `TRUE`, `true`, `True`, `0`, `f`, `F`, `FALSE`, `false`, `False`) - `"yes"` is not one of them, and that's a real, common source of "why won't my boolean env var work" bugs.

---

## Hand-Rolled .env File Support

Typing a dozen `export` commands (or `VAR=value` prefixes) before every local run gets old fast. A `.env` file - a plain text file of `KEY=VALUE` lines - is the standard convenience for local development. You don't need a library to read one; it's a handful of lines with `bufio` (bridging Level 18's file I/O).

```go
func LoadEnvFile(path string) (map[string]string, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    result := make(map[string]string)
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())

        if line == "" || strings.HasPrefix(line, "#") {
            continue // skip blank lines and comments
        }

        parts := strings.SplitN(line, "=", 2) // SplitN(..., 2) keeps '=' inside the value intact
        if len(parts) != 2 {
            return nil, fmt.Errorf("invalid format: %q", line)
        }
        result[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
    }
    return result, scanner.Err()
}
```

Verified against a real multi-line file with comments, blank lines, and a value containing an extra `=`:

```
.env contents:
    # This is a sample .env file for local development
    APP_NAME=DemoService

    # Blank lines above and below should be skipped
    APP_PORT=4000

    # A value can contain an equals sign after the first one
    CONNECTION_STRING=user=admin;password=secret=123

    DEBUG=true

=== Parsing .env file ===
APP_NAME="DemoService"
APP_PORT="4000"
CONNECTION_STRING="user=admin;password=secret=123"
DEBUG="true"
```

### A Real-World Convention: Don't Clobber Real Env Vars

Popular `.env` loaders (like the JavaScript/Ruby `dotenv` libraries) share one convention worth copying: **a `.env` file should fill in gaps, never override a variable someone already exported for real.** That way a real, deliberately-set production environment variable always wins over a stray `.env` file that happens to be sitting in the working directory.

```go
func LoadDotEnv(path string) error {
    values, err := parseEnvFile(path)
    if err != nil {
        return err
    }
    for key, value := range values {
        if _, alreadySet := os.LookupEnv(key); alreadySet {
            continue // real environment wins - never clobber it
        }
        os.Setenv(key, value)
    }
    return nil
}
```

Verified with `LOG_LEVEL` pre-set by the shell before the `.env` file (which also defines `LOG_LEVEL`) is loaded:

```
=== Before loading .env: LOG_LEVEL was pre-set by the shell ===
LOG_LEVEL="error" APP_ENV="" API_KEY=""

=== After loading .env ===
LOG_LEVEL="error"  (pre-set value preserved, .env value ignored)
APP_ENV="development"    (filled in from .env, was unset)
API_KEY="dev-placeholder-key"    (filled in from .env, was unset)
```

---

## Third-Party Options: Further Reading Only

Everything in this level is standard library on purpose - `os`, `flag`, `strconv`, `time`, `bufio` - so every exercise runs fully offline with zero dependencies. For larger projects, several well-known third-party libraries build on these same ideas with more features (nested config files, live reloading, automatic struct binding via tags):

- **[spf13/viper](https://github.com/spf13/viper)** - reads YAML/JSON/TOML/env/flags into one unified config source, with live-reload support
- **[kelseyhightower/envconfig](https://github.com/kelseyhightower/envconfig)** - populates a struct from environment variables using struct tags
- **[caarlos0/env](https://github.com/caarlos0/env)** - similar struct-tag-based environment variable binding, with built-in type parsing

These are worth knowing exist for when a project's configuration genuinely outgrows hand-rolled loading (many file formats, nested structures, hot-reload requirements). They are **not imported or used anywhere in this level** - understanding `os.Getenv`, `flag`, and a hand-rolled `.env` parser first is what lets you actually evaluate whether a library like this is solving a real problem for you, or just adding a dependency for something a 20-line function already does.

---

## Best Practices

### 1. Never Commit Secrets or .env Files to Version Control

```
# .gitignore
.env
```

A `.env` file holding real API keys or database passwords should never reach a repository. Commit a `.env.example` instead (see below) and let each developer/environment create their own real `.env` locally or via their deployment platform's secret store.

### 2. Keep a Documented .env.example

```
# .env.example - copy to .env and fill in real values for local dev
APP_PORT=8080
DATABASE_URL=
API_KEY=
LOG_LEVEL=info
```

This is the single source of truth for "what environment variables does this program need" - anyone cloning the repo can see exactly what to set without reading the source code.

### 3. Validate Config Early and Fail Fast

```go
// ✅ Good - caught in the first second of the program's life
cfg := LoadConfig()
if err := cfg.Validate(); err != nil {
    log.Fatal(err)
}

// ❌ Bad - the missing setting is discovered three requests deep,
// possibly in production, possibly at 3am
```

### 4. Keep Config Loading Separate From Business Logic

```go
// ✅ Good - one place loads config; everything else just receives it
func main() {
    cfg := LoadConfig()
    server := NewServer(cfg)
    server.Run()
}

// ❌ Bad - os.Getenv calls scattered through handlers, services, and
// helpers, each with its own ad-hoc default and no single source of truth
```

---

## Common Mistakes

### Mistake 1: Silently Falling Back to a Zero-Value Default for a Required Setting

```go
// ❌ WRONG - an empty DatabaseURL silently becomes "", the program
// starts "successfully," and fails confusingly on the first query
cfg.DatabaseURL = os.Getenv("DATABASE_URL")

// ✅ RIGHT - treat a required setting's absence as a startup error
if cfg.DatabaseURL == "" {
    return errors.New("DATABASE_URL is required")
}
```

### Mistake 2: Mixing Config-Loading Code Into Business Logic

```go
// ❌ WRONG - os.Getenv calls scattered through the codebase, with no
// single place to see every setting the program depends on
func ProcessOrder() {
    timeout := os.Getenv("ORDER_TIMEOUT") // buried three files deep
    // ...
}

// ✅ RIGHT - load once, pass the typed Config down through constructors
func ProcessOrder(cfg Config) { /* uses cfg.OrderTimeout */ }
```

### Mistake 3: Not Distinguishing "Unset" From "Empty String"

```go
// ❌ WRONG - can't tell "no feature flags configured, use defaults"
// apart from "explicitly disable every feature flag"
flags := os.Getenv("FEATURE_FLAGS")

// ✅ RIGHT - when the difference matters, use LookupEnv
if flags, ok := os.LookupEnv("FEATURE_FLAGS"); ok {
    // explicitly set (even if empty) - respect it exactly
} else {
    // never set - fall back to the built-in default set
}
```

### Mistake 4: Forgetting to Document Required Environment Variables

```
# ❌ WRONG - a new developer has to grep the source for every os.Getenv
# call just to figure out what the program needs to run

# ✅ RIGHT - a .env.example (or a README table) lists every variable,
# whether it's required, and what it does
```

---

## Summary

**Why Configuration Matters:**
- One binary, many environments - no recompiling, no environment-specific branches

**Reading Config:**
- `os.Getenv` - simple, but `""` means both "unset" and "empty"
- `os.LookupEnv` - the `ok` return value distinguishes the two
- `flag.String`/`Int`/`Bool` + `flag.Parse()` - command-line overrides with built-in defaults

**Structuring Config:**
- One `Config` struct, one `LoadConfig()` function - never scattered `os.Getenv` calls
- Precedence: **flags > environment variables > hardcoded defaults**

**Correctness:**
- `Validate()` at startup, `log.Fatal` on failure - fail fast, never limp along with a broken config
- `strconv`/`time.ParseDuration` for typed values, always checking the error

**Local Dev Convenience:**
- A hand-rolled `.env` parser (comments, blank lines, `KEY=VALUE`) - no library needed
- Never let `.env` values override real, already-set environment variables

**Bigger Projects:**
- viper, envconfig, caarlos0/env exist for nested formats and hot-reload - further reading only, never imported here

---

## Next Steps

You now understand:
- ✅ Why the same binary needs to behave differently across dev/staging/production
- ✅ `os.Getenv` vs `os.LookupEnv`, and when the difference actually matters
- ✅ The `flag` package: definitions, defaults, and `flag.Parse()`
- ✅ Designing a typed `Config` struct with a single `LoadConfig()` entry point
- ✅ The flags > env > defaults precedence pattern used by real services
- ✅ Fail-fast validation at startup, bridging Level 33's `log.Fatal`
- ✅ Parsing and hand-rolling a `.env` file loader with zero dependencies

**Next level:** Level 35 - Docker
- Packaging your configured Go binary into a container image
- Passing environment variables into a container at `docker run` time
- Why the configuration habits from this level are exactly what makes a Go binary "container-friendly"

You can now build one binary that adapts to wherever it's deployed! Keep going! 🚀
