# Level 34: Configuration - Exercises

Complete all exercises in order. Each exercise builds on previous knowledge.

**Important:** Every exercise in this level only ever sets or reads environment variables scoped to its own process - either with `os.Setenv` inside the program itself, or by prefixing `go run` with `VAR=value` in the shell. No exercise ever touches your real shell environment or files outside its own project directory.

---

## Exercise 1: os.Getenv vs os.LookupEnv

**Objective:** Prove that `Getenv` cannot distinguish "unset" from "set to empty," and that `LookupEnv` can

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level34-exercise1
cd ~/projects/level34-exercise1
go mod init level34.example/exercise1
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "os"
)

func main() {
    fmt.Println("=== os.Getenv (returns \"\" if unset) ===")
    val := os.Getenv("APP_NAME")
    fmt.Printf("APP_NAME via Getenv: %q\n", val)

    fmt.Println("\n=== os.LookupEnv (distinguishes unset from empty) ===")
    val2, ok := os.LookupEnv("APP_NAME")
    fmt.Printf("APP_NAME via LookupEnv: value=%q ok=%v\n", val2, ok)

    fmt.Println("\n=== Now set APP_NAME to an EMPTY string (still \"set\") ===")
    os.Setenv("APP_NAME", "")
    val3 := os.Getenv("APP_NAME")
    fmt.Printf("APP_NAME via Getenv: %q\n", val3)
    val4, ok4 := os.LookupEnv("APP_NAME")
    fmt.Printf("APP_NAME via LookupEnv: value=%q ok=%v\n", val4, ok4)

    fmt.Println("\n=== Now set APP_NAME to a real value ===")
    os.Setenv("APP_NAME", "MyService")
    val5 := os.Getenv("APP_NAME")
    fmt.Printf("APP_NAME via Getenv: %q\n", val5)
    val6, ok6 := os.LookupEnv("APP_NAME")
    fmt.Printf("APP_NAME via LookupEnv: value=%q ok=%v\n", val6, ok6)

    fmt.Println("\n=== The key lesson ===")
    fmt.Println("Getenv cannot tell you 'unset' apart from 'set to empty string' - both return \"\".")
    fmt.Println("LookupEnv's second return value (ok) tells you the difference.")
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

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

=== The key lesson ===
Getenv cannot tell you 'unset' apart from 'set to empty string' - both return "".
LookupEnv's second return value (ok) tells you the difference.
```

**Learning Objectives:**
- ✅ Use `os.Getenv` to read an environment variable with an implicit `""` fallback
- ✅ Use `os.LookupEnv` to distinguish "never set" from "set to empty"
- ✅ Recognize why this distinction matters for certain settings

---

## Exercise 2: The flag Package With Defaults

**Objective:** Define flags with defaults and verify the same program with different arguments

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level34-exercise2
cd ~/projects/level34-exercise2
go mod init level34.example/exercise2
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
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

    fmt.Println("=== Resolved Flags ===")
    fmt.Printf("host:  %s\n", *host)
    fmt.Printf("port:  %d\n", *port)
    fmt.Printf("debug: %v\n", *debug)
}
EOF
```

3. Run the program three different ways:

```bash
go run main.go
go run main.go -host=example.com -port=9090 -debug=true
go run main.go -port=3000
```

**Expected Output:**

Run 1 (no arguments - all defaults):
```
=== Resolved Flags ===
host:  localhost
port:  8080
debug: false
```

Run 2 (every flag set):
```
=== Resolved Flags ===
host:  example.com
port:  9090
debug: true
```

Run 3 (only `-port` set):
```
=== Resolved Flags ===
host:  localhost
port:  3000
debug: false
```

**Learning Objectives:**
- ✅ Define string, int, and bool flags with `flag.String`/`Int`/`Bool`
- ✅ Understand that `flag.Parse()` must run before reading any flag value
- ✅ See that an unpassed flag keeps its declared default

---

## Exercise 3: A Config Struct Populated From Environment Variables

**Objective:** Replace scattered `os.Getenv` calls with one typed `Config` struct and one `LoadConfig` function

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level34-exercise3
cd ~/projects/level34-exercise3
go mod init level34.example/exercise3
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "os"
)

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

func main() {
    fmt.Println("=== LoadConfig with NO environment variables set (all defaults) ===")
    cfg := LoadConfig()
    fmt.Printf("%+v\n", cfg)

    fmt.Println("\n=== LoadConfig after setting APP_PORT and APP_ENV ===")
    os.Setenv("APP_PORT", "9000")
    os.Setenv("APP_ENV", "production")
    cfg2 := LoadConfig()
    fmt.Printf("%+v\n", cfg2)
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== LoadConfig with NO environment variables set (all defaults) ===
{Host:0.0.0.0 Port:8080 Env:development}

=== LoadConfig after setting APP_PORT and APP_ENV ===
{Host:0.0.0.0 Port:9000 Env:production}
```

**Learning Objectives:**
- ✅ Design a typed `Config` struct for all of a program's settings
- ✅ Write a single `LoadConfig()` entry point instead of scattered `os.Getenv` calls
- ✅ Confirm defaults apply per-field, independently of each other

---

## Exercise 4: Config Precedence (Flags > Env > Defaults)

**Objective:** Build and test a resolver implementing the real-world precedence order

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level34-exercise4
cd ~/projects/level34-exercise4
go mod init level34.example/exercise4
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "flag"
    "fmt"
    "os"
)

// resolve implements: flag value (if explicitly set) > env var (if set) > default.
func resolve(flagVal string, flagSet bool, envKey, def string) (string, string) {
    if flagSet {
        return flagVal, "flag"
    }
    if envVal, ok := os.LookupEnv(envKey); ok {
        return envVal, "env"
    }
    return def, "default"
}

func main() {
    fs := flag.NewFlagSet("app", flag.ExitOnError)
    portFlag := fs.String("port", "8080", "server port")
    fs.Parse(os.Args[1:])

    // Track which flags were explicitly passed on the command line.
    explicitlySet := map[string]bool{}
    fs.Visit(func(f *flag.Flag) {
        explicitlySet[f.Name] = true
    })

    port, source := resolve(*portFlag, explicitlySet["port"], "APP_PORT", "8080")
    fmt.Printf("port=%s source=%s\n", port, source)
}
EOF
```

3. Run the program against all four combinations of env/flag present or absent:

```bash
go run main.go
APP_PORT=7000 go run main.go
go run main.go -port=9999
APP_PORT=7000 go run main.go -port=9999
```

**Expected Output:**

```
port=8080 source=default
port=7000 source=env
port=9999 source=flag
port=9999 source=flag
```

**Learning Objectives:**
- ✅ Implement the flags > env > defaults precedence order in a small resolver function
- ✅ Use `fs.Visit` to detect whether a flag was actually passed, not just whether it has a value
- ✅ Verify all four presence/absence combinations produce the correct source

---

## Exercise 5: Fail-Fast Validation at Startup

**Objective:** Validate required configuration once, at startup, instead of discovering problems deep inside request handling

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level34-exercise5
cd ~/projects/level34-exercise5
go mod init level34.example/exercise5
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "errors"
    "fmt"
    "os"
)

type Config struct {
    DatabaseURL string
    Port        string
}

func LoadConfig() Config {
    return Config{
        DatabaseURL: os.Getenv("DATABASE_URL"),
        Port:        os.Getenv("PORT"),
    }
}

// Validate fails fast: it checks every required setting up front and returns
// a single descriptive error instead of letting a missing value blow up
// deep inside a request handler later.
func (c Config) Validate() error {
    if c.DatabaseURL == "" {
        return errors.New("DATABASE_URL is required but not set")
    }
    if c.Port == "" {
        return errors.New("PORT is required but not set")
    }
    return nil
}

func run(label string) {
    fmt.Printf("=== %s ===\n", label)
    cfg := LoadConfig()
    if err := cfg.Validate(); err != nil {
        fmt.Println("config error:", err)
        fmt.Println("(in a real program: log.Fatal(err) here, before starting the server)")
        return
    }
    fmt.Printf("config OK: %+v\n", cfg)
}

func main() {
    run("Case 1: DATABASE_URL missing (should fail)")

    os.Setenv("DATABASE_URL", "postgres://localhost/mydb")
    run("Case 2: DATABASE_URL set, PORT missing (should fail)")

    os.Setenv("PORT", "8080")
    run("Case 3: both set (should succeed)")
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

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

**Learning Objectives:**
- ✅ Write a `Validate()` method that checks every required field up front
- ✅ See both a failing case and a succeeding case verified in the same run
- ✅ Understand why `log.Fatal` at startup (bridging Level 33) beats a panic deep in a handler

---

## Exercise 6: Type-Converting Numeric and Boolean Env Vars

**Objective:** Parse env vars into `int`, `bool`, and `time.Duration`, handling parse errors correctly

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level34-exercise6
cd ~/projects/level34-exercise6
go mod init level34.example/exercise6
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "os"
    "strconv"
    "time"
)

func getEnvAsInt(key string, fallback int) (int, error) {
    raw, ok := os.LookupEnv(key)
    if !ok {
        return fallback, nil
    }
    n, err := strconv.Atoi(raw)
    if err != nil {
        return 0, fmt.Errorf("invalid int for %s: %q: %w", key, raw, err)
    }
    return n, nil
}

func getEnvAsBool(key string, fallback bool) (bool, error) {
    raw, ok := os.LookupEnv(key)
    if !ok {
        return fallback, nil
    }
    b, err := strconv.ParseBool(raw)
    if err != nil {
        return false, fmt.Errorf("invalid bool for %s: %q: %w", key, raw, err)
    }
    return b, nil
}

func getEnvAsDuration(key string, fallback time.Duration) (time.Duration, error) {
    raw, ok := os.LookupEnv(key)
    if !ok {
        return fallback, nil
    }
    d, err := time.ParseDuration(raw)
    if err != nil {
        return 0, fmt.Errorf("invalid duration for %s: %q: %w", key, raw, err)
    }
    return d, nil
}

func main() {
    fmt.Println("=== MAX_CONNECTIONS unset -> fallback ===")
    n, err := getEnvAsInt("MAX_CONNECTIONS", 10)
    fmt.Printf("value=%d err=%v\n", n, err)

    fmt.Println("\n=== MAX_CONNECTIONS=\"50\" -> parsed ===")
    os.Setenv("MAX_CONNECTIONS", "50")
    n, err = getEnvAsInt("MAX_CONNECTIONS", 10)
    fmt.Printf("value=%d err=%v\n", n, err)

    fmt.Println("\n=== MAX_CONNECTIONS=\"not-a-number\" -> parse error ===")
    os.Setenv("MAX_CONNECTIONS", "not-a-number")
    n, err = getEnvAsInt("MAX_CONNECTIONS", 10)
    fmt.Printf("value=%d err=%v\n", n, err)

    fmt.Println("\n=== ENABLE_CACHE=\"true\" -> parsed bool ===")
    os.Setenv("ENABLE_CACHE", "true")
    b, err := getEnvAsBool("ENABLE_CACHE", false)
    fmt.Printf("value=%v err=%v\n", b, err)

    fmt.Println("\n=== ENABLE_CACHE=\"yes\" -> invalid bool ===")
    os.Setenv("ENABLE_CACHE", "yes")
    b, err = getEnvAsBool("ENABLE_CACHE", false)
    fmt.Printf("value=%v err=%v\n", b, err)

    fmt.Println("\n=== REQUEST_TIMEOUT=\"5s\" -> parsed duration ===")
    os.Setenv("REQUEST_TIMEOUT", "5s")
    d, err := getEnvAsDuration("REQUEST_TIMEOUT", 30*time.Second)
    fmt.Printf("value=%v err=%v\n", d, err)
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

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

**Learning Objectives:**
- ✅ Convert env var strings to `int`, `bool`, and `time.Duration` with `strconv`/`time` (bridging Level 3)
- ✅ Handle and report a conversion error instead of ignoring it
- ✅ See exactly which boolean spellings `strconv.ParseBool` accepts and rejects

---

## Exercise 7: Hand-Rolled .env File Parser

**Objective:** Parse a real `KEY=VALUE` file, correctly skipping comments and blank lines

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level34-exercise7
cd ~/projects/level34-exercise7
go mod init level34.example/exercise7
```

2. Create a sample `.env` file:

```bash
cat > .env << 'ENVEOF'
# This is a sample .env file for local development
APP_NAME=DemoService

# Blank lines above and below should be skipped
APP_PORT=4000

# A value can contain an equals sign after the first one
CONNECTION_STRING=user=admin;password=secret=123

DEBUG=true
ENVEOF
```

3. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "bufio"
    "fmt"
    "os"
    "strings"
)

// LoadEnvFile reads a simple KEY=VALUE .env file into a map, skipping
// blank lines and lines starting with # (comments).
func LoadEnvFile(path string) (map[string]string, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    result := make(map[string]string)
    scanner := bufio.NewScanner(file)
    lineNum := 0
    for scanner.Scan() {
        lineNum++
        line := strings.TrimSpace(scanner.Text())

        if line == "" {
            continue // skip blank lines
        }
        if strings.HasPrefix(line, "#") {
            continue // skip comments
        }

        parts := strings.SplitN(line, "=", 2)
        if len(parts) != 2 {
            return nil, fmt.Errorf("line %d: invalid format (expected KEY=VALUE): %q", lineNum, line)
        }
        key := strings.TrimSpace(parts[0])
        value := strings.TrimSpace(parts[1])
        result[key] = value
    }
    if err := scanner.Err(); err != nil {
        return nil, err
    }
    return result, nil
}

func main() {
    fmt.Println("=== Parsing .env file ===")
    values, err := LoadEnvFile(".env")
    if err != nil {
        fmt.Println("error:", err)
        return
    }

    for _, key := range []string{"APP_NAME", "APP_PORT", "CONNECTION_STRING", "DEBUG"} {
        fmt.Printf("%s=%q\n", key, values[key])
    }

    fmt.Println("\n=== Loading parsed values into the process environment ===")
    for k, v := range values {
        os.Setenv(k, v)
    }
    fmt.Println("APP_PORT from os.Getenv after loading:", os.Getenv("APP_PORT"))

    fmt.Println("\n=== Missing file error ===")
    _, err = LoadEnvFile("does-not-exist.env")
    fmt.Println("error:", err)
}
EOF
```

4. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Parsing .env file ===
APP_NAME="DemoService"
APP_PORT="4000"
CONNECTION_STRING="user=admin;password=secret=123"
DEBUG="true"

=== Loading parsed values into the process environment ===
APP_PORT from os.Getenv after loading: 4000

=== Missing file error ===
error: open does-not-exist.env: no such file or directory
```

**Learning Objectives:**
- ✅ Parse a real multi-line `.env` file with `bufio.Scanner` (bridging Level 18)
- ✅ Correctly skip comments (`#`) and blank lines
- ✅ Handle a value that itself contains `=` using `strings.SplitN(line, "=", 2)`
- ✅ Load parsed values into the process environment with `os.Setenv`

---

## Exercise 8: Loading a Typed Config Struct Directly From a .env File

**Objective:** Combine the `.env` parser with type conversion to populate a typed `Config`, including error handling

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level34-exercise8
cd ~/projects/level34-exercise8
go mod init level34.example/exercise8
```

2. Create `.env`:

```bash
cat > .env << 'ENVEOF'
# Local dev settings
APP_PORT=5050
MAX_RETRIES=5
ENABLE_DEBUG=true
ENVEOF
```

3. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "bufio"
    "fmt"
    "os"
    "strconv"
    "strings"
)

type Config struct {
    Port        int
    MaxRetries  int
    EnableDebug bool
}

// parseEnvFile reads KEY=VALUE lines into a map, skipping comments and
// blank lines (same shape as the standalone .env parser).
func parseEnvFile(path string) (map[string]string, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    values := make(map[string]string)
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        if line == "" || strings.HasPrefix(line, "#") {
            continue
        }
        parts := strings.SplitN(line, "=", 2)
        if len(parts) != 2 {
            continue
        }
        values[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
    }
    return values, scanner.Err()
}

// LoadConfigFromFile parses the .env file directly into a typed Config,
// converting strings to int/bool and reporting any conversion errors.
func LoadConfigFromFile(path string) (Config, error) {
    values, err := parseEnvFile(path)
    if err != nil {
        return Config{}, fmt.Errorf("reading env file: %w", err)
    }

    var cfg Config

    if raw, ok := values["APP_PORT"]; ok {
        port, err := strconv.Atoi(raw)
        if err != nil {
            return Config{}, fmt.Errorf("APP_PORT: invalid int %q: %w", raw, err)
        }
        cfg.Port = port
    } else {
        cfg.Port = 8080 // default
    }

    if raw, ok := values["MAX_RETRIES"]; ok {
        n, err := strconv.Atoi(raw)
        if err != nil {
            return Config{}, fmt.Errorf("MAX_RETRIES: invalid int %q: %w", raw, err)
        }
        cfg.MaxRetries = n
    } else {
        cfg.MaxRetries = 3 // default
    }

    if raw, ok := values["ENABLE_DEBUG"]; ok {
        b, err := strconv.ParseBool(raw)
        if err != nil {
            return Config{}, fmt.Errorf("ENABLE_DEBUG: invalid bool %q: %w", raw, err)
        }
        cfg.EnableDebug = b
    } else {
        cfg.EnableDebug = false // default
    }

    return cfg, nil
}

func main() {
    fmt.Println("=== Loading typed Config directly from .env ===")
    cfg, err := LoadConfigFromFile(".env")
    if err != nil {
        fmt.Println("error:", err)
        return
    }
    fmt.Printf("%+v\n", cfg)

    fmt.Println("\n=== Same loader, but MAX_RETRIES is malformed ===")
    os.WriteFile("bad.env", []byte("APP_PORT=5050\nMAX_RETRIES=lots\n"), 0644)
    _, err = LoadConfigFromFile("bad.env")
    fmt.Println("error:", err)
    os.Remove("bad.env")
}
EOF
```

4. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Loading typed Config directly from .env ===
{Port:5050 MaxRetries:5 EnableDebug:true}

=== Same loader, but MAX_RETRIES is malformed ===
error: MAX_RETRIES: invalid int "lots": strconv.Atoi: parsing "lots": invalid syntax
```

**Learning Objectives:**
- ✅ Combine `.env` parsing (Exercise 7) with type conversion (Exercise 6) into one typed loader
- ✅ Apply defaults per-field when a key is absent from the file
- ✅ Surface a clear, field-named error when a value fails to parse

---

## Exercise 9: .env Files Must Never Override Real Environment Variables

**Objective:** Verify the standard dotenv convention - `.env` fills gaps, it never overrides an already-set variable

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level34-exercise9
cd ~/projects/level34-exercise9
go mod init level34.example/exercise9
```

2. Create `.env`:

```bash
cat > .env << 'ENVEOF'
# Local dev defaults - real environment variables should win over these
APP_ENV=development
LOG_LEVEL=debug
API_KEY=dev-placeholder-key
ENVEOF
```

3. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "bufio"
    "fmt"
    "os"
    "strings"
)

func parseEnvFile(path string) (map[string]string, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    values := make(map[string]string)
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        if line == "" || strings.HasPrefix(line, "#") {
            continue
        }
        parts := strings.SplitN(line, "=", 2)
        if len(parts) != 2 {
            continue
        }
        values[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
    }
    return values, scanner.Err()
}

// LoadDotEnv loads KEY=VALUE pairs from a .env file into the process
// environment, but never overwrites a variable that is ALREADY set in
// the real environment. .env is for filling gaps during local dev, not
// for overriding a value someone deliberately exported.
func LoadDotEnv(path string) error {
    values, err := parseEnvFile(path)
    if err != nil {
        return err
    }
    for key, value := range values {
        if _, alreadySet := os.LookupEnv(key); alreadySet {
            continue // real environment wins - do not clobber it
        }
        os.Setenv(key, value)
    }
    return nil
}

func main() {
    fmt.Println("=== Before loading .env: LOG_LEVEL was pre-set by the shell ===")
    fmt.Printf("LOG_LEVEL=%q APP_ENV=%q API_KEY=%q\n",
        os.Getenv("LOG_LEVEL"), os.Getenv("APP_ENV"), os.Getenv("API_KEY"))

    if err := LoadDotEnv(".env"); err != nil {
        fmt.Println("error:", err)
        return
    }

    fmt.Println("\n=== After loading .env ===")
    fmt.Printf("LOG_LEVEL=%q  (pre-set value preserved, .env value ignored)\n", os.Getenv("LOG_LEVEL"))
    fmt.Printf("APP_ENV=%q    (filled in from .env, was unset)\n", os.Getenv("APP_ENV"))
    fmt.Printf("API_KEY=%q    (filled in from .env, was unset)\n", os.Getenv("API_KEY"))
}
EOF
```

4. Run the program with `LOG_LEVEL` already set in the shell, scoped only to this one process:

```bash
LOG_LEVEL=error go run main.go
```

**Expected Output:**

```
=== Before loading .env: LOG_LEVEL was pre-set by the shell ===
LOG_LEVEL="error" APP_ENV="" API_KEY=""

=== After loading .env ===
LOG_LEVEL="error"  (pre-set value preserved, .env value ignored)
APP_ENV="development"    (filled in from .env, was unset)
API_KEY="dev-placeholder-key"    (filled in from .env, was unset)
```

**Learning Objectives:**
- ✅ Implement the "don't clobber real env vars" convention used by real dotenv loaders
- ✅ Verify a pre-set variable survives `.env` loading unchanged
- ✅ Verify unset variables are correctly filled in from the file

---

## Exercise 10: Comprehensive Practice — Full Service Config Loading

**Objective:** Combine every technique from this level into one realistic config loader: `.env` for local dev, environment variable overrides, flag overrides, fail-fast validation, and a summary of where each value came from

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level34-exercise10
cd ~/projects/level34-exercise10
go mod init level34.example/exercise10
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "bufio"
    "errors"
    "flag"
    "fmt"
    "os"
    "strconv"
    "strings"
)

type Config struct {
    Port        int
    LogLevel    string
    APIKey      string
    Environment string
}

// sources records, per field name, which layer supplied the final value:
// "default", "dotenv", "env", or "flag".
type sources map[string]string

func parseEnvFile(path string) (map[string]string, error) {
    file, err := os.Open(path)
    if err != nil {
        if os.IsNotExist(err) {
            return map[string]string{}, nil // no .env file is fine
        }
        return nil, err
    }
    defer file.Close()

    values := make(map[string]string)
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        if line == "" || strings.HasPrefix(line, "#") {
            continue
        }
        parts := strings.SplitN(line, "=", 2)
        if len(parts) != 2 {
            continue
        }
        values[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
    }
    return values, scanner.Err()
}

var trackedKeys = []string{"APP_PORT", "LOG_LEVEL", "API_KEY", "APP_ENV"}

// LoadConfig resolves the full config in precedence order:
// defaults < .env file < real env vars < command-line flags.
// It records, per field, exactly which layer supplied the final value.
func LoadConfig(args []string) (Config, sources, error) {
    src := sources{}
    cfg := Config{
        Port:        8080,
        LogLevel:    "info",
        APIKey:      "",
        Environment: "development",
    }
    for _, field := range []string{"Port", "LogLevel", "APIKey", "Environment"} {
        src[field] = "default"
    }

    // Snapshot which tracked keys were ALREADY real environment variables
    // before touching anything - this is what lets us tell "env" apart
    // from "dotenv" later, even though .env values get loaded into the
    // same process environment.
    preExisting := map[string]bool{}
    for _, key := range trackedKeys {
        _, ok := os.LookupEnv(key)
        preExisting[key] = ok
    }

    // Layer 2: .env file fills gaps in the process environment.
    dotenvValues, err := parseEnvFile(".env")
    if err != nil {
        return cfg, src, fmt.Errorf("loading .env: %w", err)
    }
    for key, value := range dotenvValues {
        if !preExisting[key] {
            os.Setenv(key, value)
        }
    }

    // Layer 3: read the (now possibly .env-filled) environment.
    applyEnv := func(field, key string, apply func(string) error) error {
        raw, ok := os.LookupEnv(key)
        if !ok {
            return nil
        }
        if err := apply(raw); err != nil {
            return err
        }
        if preExisting[key] {
            src[field] = "env"
        } else if _, fromDotenv := dotenvValues[key]; fromDotenv {
            src[field] = "dotenv"
        }
        return nil
    }

    if err := applyEnv("Port", "APP_PORT", func(v string) error {
        n, err := strconv.Atoi(v)
        if err != nil {
            return fmt.Errorf("APP_PORT: invalid int %q: %w", v, err)
        }
        cfg.Port = n
        return nil
    }); err != nil {
        return cfg, src, err
    }
    if err := applyEnv("LogLevel", "LOG_LEVEL", func(v string) error {
        cfg.LogLevel = v
        return nil
    }); err != nil {
        return cfg, src, err
    }
    if err := applyEnv("APIKey", "API_KEY", func(v string) error {
        cfg.APIKey = v
        return nil
    }); err != nil {
        return cfg, src, err
    }
    if err := applyEnv("Environment", "APP_ENV", func(v string) error {
        cfg.Environment = v
        return nil
    }); err != nil {
        return cfg, src, err
    }

    // Layer 4: command-line flags override everything.
    fs := flag.NewFlagSet("service", flag.ContinueOnError)
    portFlag := fs.Int("port", cfg.Port, "server port")
    logLevelFlag := fs.String("log-level", cfg.LogLevel, "log level")
    apiKeyFlag := fs.String("api-key", cfg.APIKey, "API key")
    envFlag := fs.String("env", cfg.Environment, "environment")
    if err := fs.Parse(args); err != nil {
        return cfg, src, err
    }
    fs.Visit(func(f *flag.Flag) {
        switch f.Name {
        case "port":
            cfg.Port = *portFlag
            src["Port"] = "flag"
        case "log-level":
            cfg.LogLevel = *logLevelFlag
            src["LogLevel"] = "flag"
        case "api-key":
            cfg.APIKey = *apiKeyFlag
            src["APIKey"] = "flag"
        case "env":
            cfg.Environment = *envFlag
            src["Environment"] = "flag"
        }
    })

    return cfg, src, nil
}

// Validate fails fast if required settings are missing or invalid.
func (c Config) Validate() error {
    if c.APIKey == "" {
        return errors.New("API_KEY is required (set the API_KEY env var, put it in .env, or pass -api-key)")
    }
    if c.Port <= 0 || c.Port > 65535 {
        return fmt.Errorf("invalid port %d: must be between 1 and 65535", c.Port)
    }
    return nil
}

func printSummary(cfg Config, src sources) {
    fmt.Println("Field        Value                  Source")
    fmt.Println("-----        -----                  ------")
    fmt.Printf("%-12s %-22v %s\n", "Port", cfg.Port, src["Port"])
    fmt.Printf("%-12s %-22v %s\n", "LogLevel", cfg.LogLevel, src["LogLevel"])
    fmt.Printf("%-12s %-22v %s\n", "APIKey", cfg.APIKey, src["APIKey"])
    fmt.Printf("%-12s %-22v %s\n", "Environment", cfg.Environment, src["Environment"])
}

func main() {
    cfg, src, err := LoadConfig(os.Args[1:])
    if err != nil {
        fmt.Println("fatal: failed to load config:", err)
        os.Exit(1)
    }
    if err := cfg.Validate(); err != nil {
        fmt.Println("fatal: invalid config:", err)
        os.Exit(1)
    }
    fmt.Println("Config loaded successfully:")
    printSummary(cfg, src)
}
EOF
```

3. Run through five scenarios that build on each other:

```bash
# Scenario A: no .env, no env vars, no flags -> validation fails (API_KEY missing)
go run main.go

# Scenario B: create a .env file providing API_KEY and LOG_LEVEL
cat > .env << 'ENVEOF'
API_KEY=dotenv-key-abc123
LOG_LEVEL=debug
ENVEOF
go run main.go

# Scenario C: a real env var overrides the .env value
LOG_LEVEL=warn go run main.go

# Scenario D: a flag overrides everything
LOG_LEVEL=warn go run main.go -port=9090 -log-level=trace

# Scenario E: remove .env, provide API_KEY via env var only
rm .env
API_KEY=env-only-key go run main.go
```

**Expected Output:**

Scenario A:
```
fatal: invalid config: API_KEY is required (set the API_KEY env var, put it in .env, or pass -api-key)
```

Scenario B:
```
Config loaded successfully:
Field        Value                  Source
-----        -----                  ------
Port         8080                   default
LogLevel     debug                  dotenv
APIKey       dotenv-key-abc123      dotenv
Environment  development            default
```

Scenario C:
```
Config loaded successfully:
Field        Value                  Source
-----        -----                  ------
Port         8080                   default
LogLevel     warn                   env
APIKey       dotenv-key-abc123      dotenv
Environment  development            default
```

Scenario D:
```
Config loaded successfully:
Field        Value                  Source
-----        -----                  ------
Port         9090                   flag
LogLevel     trace                  flag
APIKey       dotenv-key-abc123      dotenv
Environment  development            default
```

Scenario E:
```
Config loaded successfully:
Field        Value                  Source
-----        -----                  ------
Port         8080                   default
LogLevel     info                   default
APIKey       env-only-key           env
Environment  development            default
```

**Learning Objectives:**
- ✅ Combine `.env` loading, environment variables, and flags into one precedence-respecting loader
- ✅ Correctly attribute each field's final value to the layer that actually supplied it
- ✅ Fail fast with a clear error when a required setting (`API_KEY`) is missing from every layer
- ✅ Produce a startup summary an operator could use to debug "why is this config wrong"

---

## Bonus Challenges

### Challenge 1: Config Reload on SIGHUP (Conceptual Sketch)

Many long-running Unix services reload their configuration when they receive a `SIGHUP` signal, instead of requiring a restart. Sketch (simplified, doesn't need to be production-ready) a program that re-runs `LoadConfig()` whenever it receives `SIGHUP`, and safely swaps in the new config for future use.

```bash
mkdir -p ~/projects/level34-bonus1
cd ~/projects/level34-bonus1
go mod init level34.example/bonus1
```

**Hints:**
- `os/signal.Notify(ch, syscall.SIGHUP)` delivers the signal to a channel you select on in a loop
- Store the live `*Config` behind a `sync.RWMutex` (or `atomic.Pointer[Config]`) so readers never see a half-updated config while a reload is in progress
- On reload, call `LoadConfig()` and `Validate()` again - if validation fails, log the error and **keep serving with the old config** rather than crashing a running service
- Test it by sending yourself the signal: `kill -HUP <pid>` from another terminal while the program runs

### Challenge 2: Redacting Secrets in a Config Summary

Bridge Level 33's "never put secrets in logs" rule: write a `redact` function that masks secret-looking values (API keys, passwords, tokens) before they're ever printed in a config summary like Exercise 10's.

```bash
mkdir -p ~/projects/level34-bonus2
cd ~/projects/level34-bonus2
go mod init level34.example/bonus2
```

**Hints:**
- A simple version: keep the first 4 characters, replace the rest with `****` (e.g. `sk-live-abcdef123456` -> `sk-l****`)
- Guard the short-value case - don't reveal a whole value that's 4 characters or shorter; just return `****`
- Apply it only to fields you know are secret-shaped (by field name, e.g. anything containing `KEY`, `SECRET`, `PASSWORD`, `TOKEN`) - never redact something like `Port` or `Environment` that's safe and useful to see in full

### Challenge 3: A Typed Environment Enum With Validation

Instead of storing `Environment` as a plain `string` that could be any typo, define a proper `Environment` type restricted to `"dev"`, `"staging"`, or `"prod"`.

```bash
mkdir -p ~/projects/level34-bonus3
cd ~/projects/level34-bonus3
go mod init level34.example/bonus3
```

**Hints:**
- `type Environment string` plus `const ( Dev Environment = "dev"; Staging Environment = "staging"; Prod Environment = "prod" )`
- A `ParseEnvironment(s string) (Environment, error)` function that only accepts the three valid spellings and returns a descriptive error otherwise (this is exactly the "validate at startup" idea from Exercise 5, applied to one field)
- Call `ParseEnvironment` from inside `LoadConfig`/`Validate` so an invalid `APP_ENV=produciton` (typo) is caught at startup, not three deploys later

---

## What You've Learned

After completing these 10 exercises, you can:

✅ Distinguish `os.Getenv` from `os.LookupEnv` and know when the difference matters
✅ Define command-line flags with defaults and verify them across multiple runs
✅ Design a typed `Config` struct with a single `LoadConfig()` entry point
✅ Implement and test the flags > env > defaults precedence order
✅ Fail fast with startup validation instead of discovering problems in production
✅ Convert env var strings to `int`, `bool`, and `time.Duration` with proper error handling
✅ Hand-roll a `.env` file parser that handles comments, blank lines, and embedded `=`
✅ Respect the "don't clobber real env vars" convention when loading `.env` files
✅ Combine every technique into one realistic, fully-attributed service config loader

---

## Next Level

Level 35: Docker
- Packaging your configured Go binary into a container image
- Passing environment variables into a container at `docker run` time
- Why the configuration habits from this level make a Go binary "container-friendly"

Great work! You can now build one binary that adapts correctly to any environment! 🚀
