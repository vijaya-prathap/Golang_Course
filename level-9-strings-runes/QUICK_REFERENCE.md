# Level 9: Quick Reference Card

## 🚀 Quick Start (5 Minutes)

```bash
# Setup
mkdir myapp && cd myapp
go mod init myapp

# Create main.go
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "strings"
    "unicode/utf8"
)

func main() {
    s := "héllo"
    fmt.Println("bytes:", len(s))                        // 6
    fmt.Println("runes:", utf8.RuneCountInString(s))      // 5

    fmt.Println(strings.ToUpper(s))
    fmt.Println(strings.Split("a,b,c", ","))

    var b strings.Builder
    b.WriteString("Go")
    b.WriteString(" rocks")
    fmt.Println(b.String())
}
EOF

# Run
go run main.go
```

---

## 📋 byte vs rune

```go
var b byte = 'A'   // uint8  - one raw UTF-8 code unit (0-255)
var r rune = 'A'   // int32  - one decoded Unicode code point (1-4 bytes encoded)

s := "héllo"
s[1]                  // byte:  195 (HALF of é - not a character!)
[]rune(s)[1]          // rune:  'é' (the real character)

len(s)                          // 6  - byte count
utf8.RuneCountInString(s)       // 5  - character count
```

---

## 🔄 Conversions

```go
b := []byte(s)      // string -> []byte
r := []rune(s)      // string -> []rune

s1 := string(b)      // []byte -> string
s2 := string(r)      // []rune -> string

s3 := string('é')    // single rune -> string
```

---

## 📦 strings Package Syntax

```go
strings.Contains(s, "sub")           // true/false
strings.HasPrefix(s, "pre")          // true/false
strings.HasSuffix(s, "suf")          // true/false
strings.Index(s, "sub")              // byte index, or -1

strings.Split(s, ",")                // split on separator, keeps empties
strings.Fields(s)                    // split on whitespace, drops empties
strings.Join(parts, "-")             // join with separator

strings.TrimSpace(s)                 // trim leading/trailing whitespace
strings.Trim(s, "*")                 // trim cutset from both ends
strings.TrimLeft(s, "*")             // trim cutset from left only
strings.TrimRight(s, "*")            // trim cutset from right only
strings.TrimPrefix(s, "pre")         // remove exact prefix if present
strings.TrimSuffix(s, "suf")         // remove exact suffix if present

strings.Replace(s, "old", "new", 1)  // replace first N matches
strings.ReplaceAll(s, "old", "new")  // replace all matches
strings.ToUpper(s)                   // uppercase
strings.ToLower(s)                   // lowercase

strings.Repeat(s, 3)                 // repeat string N times
strings.Count(s, "sub")              // count occurrences
```

---

## 🏗️ strings.Builder

```go
var b strings.Builder
b.WriteString("Go")     // append a string
b.WriteByte(' ')        // append a byte
b.WriteRune('世')        // append a rune
result := b.String()    // get the final string
```

```go
// ❌ O(n^2) - avoid in loops
result := ""
for _, s := range many {
    result += s
}

// ✅ O(n) - use Builder in loops
var b strings.Builder
for _, s := range many {
    b.WriteString(s)
}
result := b.String()
```

---

## 🧮 unicode / unicode/utf8

```go
unicode.IsLetter(r)   // any-language letter check
unicode.IsDigit(r)    // digit check
unicode.IsSpace(r)    // whitespace check
unicode.IsUpper(r)    // uppercase check
unicode.IsLower(r)    // lowercase check
unicode.ToUpper(r)    // uppercase a rune
unicode.ToLower(r)    // lowercase a rune

utf8.RuneCountInString(s)      // character count, no allocation
utf8.DecodeRuneInString(s)     // (rune, byteWidth) of first rune
```

---

## 🔁 strconv Bridge (Recap From Level 3)

```go
n, err := strconv.Atoi("42")   // string -> int
s := strconv.Itoa(42)          // int -> string
```

```go
parts := strings.Split("1,2,3", ",")
for _, p := range parts {
    n, err := strconv.Atoi(p)
    if err != nil {
        continue
    }
    // use n
}
```

---

## 🔃 Correct String Reversal

```go
// ✅ RIGHT - rune-based, safe for all Unicode text
func reverse(s string) string {
    r := []rune(s)
    for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
        r[i], r[j] = r[j], r[i]
    }
    return string(r)
}

reverse("héllo")   // "olléh"  - correct
reverse("日本語")   // "語本日" - correct

// ❌ WRONG - byte-based, corrupts multi-byte characters
// (swapping []byte(s) instead of []rune(s))
```

---

## ⚠️ Common Mistakes

| Mistake | Wrong | Right |
|---------|-------|-------|
| String index as character | `s[i]` for a character | `[]rune(s)[i]` |
| Byte-based reversal | swap `[]byte(s)` | swap `[]rune(s)` |
| Naive concatenation in a loop | `result += s` | `strings.Builder` |
| len() as character count | `len(s)` for characters | `utf8.RuneCountInString(s)` |
| ASCII-only letter check | `'a' <= r && r <= 'z'` | `unicode.IsLetter(r)` |
| Mutating a string directly | `s[0] = 'H'` | `b := []byte(s); b[0] = 'H'; s = string(b)` |

---

## 🎓 Before Next Level

Can you:
- [ ] Explain why a Go string is immutable, UTF-8-encoded bytes?
- [ ] State the exact difference between byte and rune?
- [ ] Convert correctly between string, []byte, and []rune?
- [ ] Explain why naive concatenation is O(n²) and fix it with strings.Builder?
- [ ] Use the unicode package instead of hardcoded ASCII ranges?
- [ ] Write a rune-safe string reversal and explain why a byte-based one breaks?

If YES → You're ready for Level 10!

---

## 📚 Next Level

Level 10: Maps
- Map declaration, zero values, and the comma-ok idiom
- Adding, updating, and deleting keys
- Maps as sets and frequency counters

You've got strings and runes down! 💪
