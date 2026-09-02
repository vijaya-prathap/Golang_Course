# Level 34: Quick Reference Card

## 🚀 Quick Start (5 Minutes)

```bash
# Setup
mkdir myapp && cd myapp
go mod init myapp

# Create main.go
cat > main.go << 'EOF'
package main

import (
    "flag"
    "fmt"
    "os"
)

type Config struct {
    Port string
    Env  string
}

func LoadConfig() Config {
    port := flag.String("port", "", "server port")
    flag.Parse()

    resolved := *port
    if resolved == "" {
        resolved = os.Getenv("APP_PORT")
    }
    if resolved == "" {
        resolved = "8080"
    }

    return Config{
        Port: resolved,
        Env:  os.Getenv("APP_ENV"),
    }
}

func main() {
    cfg := LoadConfig()
    fmt.Printf("%+v\n", cfg)
}
EOF

# Run
go run main.go
APP_PORT=9000 go run main.go
go run main.go -port=3000
```

---

## 📋 os.Getenv vs os.LookupEnv

```go
value := os.Getenv("KEY")             // "" if unset OR set to ""
value, ok := os.LookupEnv("KEY")      // ok=false only if truly unset
```

| Case | Getenv | LookupEnv |
|------|--------|-----------|
| Unset | `""` | `"", false` |
| Set to `""` | `""` | `"", true` |
| Set to `"x"` | `"x"` | `"x", true` |

---

## 🚀 The flag Package

```go
host := flag.String("host", "localhost", "usage text")
port := flag.Int("port", 8080, "usage text")
debug := flag.Bool("debug", false, "usage text")

flag.Parse() // call BEFORE reading *host, *port, *debug
```

```bash
go run main.go                  # all defaults
go run main.go -port=9090       # port overridden
go run main.go -debug=true      # debug overridden
```

---

## 🗂️ Config Struct Pattern

```go
type Config struct {
    Host string
    Port string
}

func LoadConfig() Config {
    return Config{
        Host: getEnvOrDefault("APP_HOST", "0.0.0.0"),
        Port: getEnvOrDefault("APP_PORT", "8080"),
    }
}

func getEnvOrDefault(key, fallback string) string {
    if v, ok := os.LookupEnv(key); ok {
        return v
    }
    return fallback
}
```

---

## 🎯 Precedence: Flags > Env > Defaults

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

| Flag set? | Env set? | Winner |
|-----------|----------|--------|
| No | No | default |
| No | Yes | env |
| Yes | No | flag |
| Yes | Yes | flag |

---

## 🚨 Fail-Fast Validation

```go
func (c Config) Validate() error {
    if c.DatabaseURL == "" {
        return errors.New("DATABASE_URL is required")
    }
    return nil
}

cfg := LoadConfig()
if err := cfg.Validate(); err != nil {
    log.Fatal(err) // log, then exit - before starting the server
}
```

---

## 🔢 Type Conversion

```go
n, err := strconv.Atoi(raw)               // string -> int
b, err := strconv.ParseBool(raw)          // string -> bool ("true","1","t"... only)
d, err := time.ParseDuration(raw)         // string -> time.Duration ("5s","2m30s")

if err != nil {
    return fmt.Errorf("invalid value for %s: %q: %w", key, raw, err)
}
```

---

## 📄 Hand-Rolled .env Parser

```go
scanner := bufio.NewScanner(file)
for scanner.Scan() {
    line := strings.TrimSpace(scanner.Text())
    if line == "" || strings.HasPrefix(line, "#") {
        continue // skip blanks and comments
    }
    parts := strings.SplitN(line, "=", 2) // keep any "=" inside the value
    result[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
}
```

```go
// Never override a variable that's already really set
for key, value := range dotenvValues {
    if _, alreadySet := os.LookupEnv(key); !alreadySet {
        os.Setenv(key, value)
    }
}
```

---

## ⚠️ Common Mistakes

| Mistake | Wrong | Right |
|---------|-------|-------|
| Silent zero-value fallback | `cfg.URL = os.Getenv("URL")` and moving on | Check `== ""` and error/fail fast |
| Config mixed into business logic | `os.Getenv` calls buried in handlers | Load once into a `Config`, pass it down |
| Unset vs empty confusion | Using `Getenv` when "explicitly empty" matters | Use `LookupEnv` |
| Undocumented required vars | No list anywhere of what's needed | Keep a `.env.example` |
| Clobbering real env vars | `.env` loader always calls `os.Setenv` | Only set if `os.LookupEnv` says unset |
| Ignoring a bool env var's spelling | `ENABLE_X=yes` (not accepted) | Use `true`/`false`/`1`/`0` |

---

## 🎓 Before Next Level

Can you:
- [ ] Explain the difference between `os.Getenv` and `os.LookupEnv`?
- [ ] Define and parse flags with defaults?
- [ ] Design a `Config` struct with a single `LoadConfig()` function?
- [ ] Implement flags > env > defaults precedence and test all four combinations?
- [ ] Write a `Validate()` that fails fast on a missing required setting?
- [ ] Convert an env var to `int`/`bool`/`time.Duration` with error handling?
- [ ] Hand-roll a `.env` parser that skips comments/blanks and never overrides real env vars?

If YES → You're ready for Level 35!

---

## 📚 Next Level

Level 35: Docker
- Packaging your configured Go binary into a container image
- Passing environment variables into a container at `docker run` time

You've got configuration down! 💪
