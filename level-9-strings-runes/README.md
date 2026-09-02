# Level 9: Strings & Runes - Complete Guide

## Introduction

Welcome to Level 9! You've mastered slices (Level 8). Now it's time to take a much deeper look at something you've already been using since Level 1: **strings**.

Level 6 (Loops) briefly introduced range-over-strings and showed that a string like `"héllo"` has `len() == 6` bytes but only 5 runes, with `range` giving byte-offset indices that jump when a multi-byte character appears. This level builds on that foundation and goes all the way down: how strings are represented in memory, the full `byte`/`rune` distinction, every conversion between `string`, `[]byte`, and `[]rune`, the `strings` package in depth, efficient string building, and Unicode-aware text processing.

Strings look simple on the surface - but Go's strings are UTF-8 encoded byte sequences, and treating them like arrays of characters is one of the most common sources of subtle bugs in Go code. By the end of this level, you'll know exactly why, and how to avoid it.

---

## Table of Contents

1. [Strings Are Immutable UTF-8 Byte Sequences](#strings-are-immutable-utf-8-byte-sequences)
2. [byte vs rune: The Core Distinction](#byte-vs-rune-the-core-distinction)
3. [Converting Between string, []byte, and []rune](#converting-between-string-byte-and-rune)
4. [The strings Package Tour](#the-strings-package-tour)
5. [strings.Builder for Efficient Concatenation](#stringsbuilder-for-efficient-concatenation)
6. [strconv in a Strings Context (Recap)](#strconv-in-a-strings-context-recap)
7. [Unicode-Aware Operations](#unicode-aware-operations)
8. [Common String Manipulation Tasks](#common-string-manipulation-tasks)
9. [Best Practices](#best-practices)
10. [Common Mistakes](#common-mistakes)

---

## Strings Are Immutable UTF-8 Byte Sequences

A Go `string` is a **read-only slice of bytes**. Under the hood it's just a pointer to a byte array plus a length - there's no separate "character array" the way some languages have. Whatever text you put in a string literal is stored as its **UTF-8 encoding**.

```go
s := "hello"
fmt.Println(len(s))  // 5 - length in BYTES
```

For ASCII text, "bytes" and "characters" happen to be the same thing, which is why this distinction is easy to miss until you meet non-ASCII text.

### Strings Cannot Be Mutated In Place

Because a string is read-only, you cannot assign to an index of it:

```go
s := "hello"
s[0] = 'H'
```

Trying to compile this gives a real compiler error:

```
./main.go:5:2: cannot assign to s[0] (neither addressable nor a map index expression)
```

If you need to change characters, you must convert to `[]byte` or `[]rune`, modify the copy, and convert back to a new string (covered in the next two sections). The original string is never touched - every "modification" produces a brand-new string.

### Why Immutability Matters

- **Safety:** strings can be freely shared and passed around without anyone accidentally corrupting them.
- **Cheap substrings:** slicing a string (`s[2:5]`) doesn't copy bytes - it just creates a new string header pointing into the same underlying data.
- **The cost:** every "modification" (concatenation, replacement, casing) must allocate a new string. Section 5 covers how to avoid doing this repeatedly and inefficiently.

### UTF-8 Encoding in a Nutshell

Go source files are UTF-8, and so are Go strings. UTF-8 encodes each Unicode character ("code point") as **1 to 4 bytes**:

```
'A'  → 1 byte  (0x41)
'é'  → 2 bytes (0xC3 0xA9)
'€'  → 3 bytes (0xE2 0x82 0xAC)
'😀' → 4 bytes (0xF0 0x9F 0x98 0x80)
```

ASCII characters (0-127) always encode to exactly 1 byte, which is why ASCII-only strings never expose the byte/character difference. Anything outside ASCII - accented letters, currency symbols, CJK characters, emoji - takes multiple bytes, and that's where naive byte-based string code breaks.

---

## byte vs rune: The Core Distinction

Level 6 mentioned this briefly; here's the full picture.

| Type | Underlying type | Represents | Size |
|------|------------------|------------|------|
| `byte` | alias for `uint8` | one raw UTF-8 code **unit** | always 1 byte |
| `rune` | alias for `int32` | one decoded Unicode code **point** (a character) | logically 1 "character", physically 1-4 bytes when encoded |

```go
var b byte = 'A'   // 65, type uint8
var r rune = 'A'   // 65, type int32
var e rune = '€'   // 8364, type int32 - still ONE rune, even though it takes 3 bytes as UTF-8
```

### Indexing a String Gives You a Byte

```go
s := "héllo"
fmt.Println(s[0])  // 104  ('h' - a single-byte character, looks fine)
fmt.Println(s[1])  // 195  (HALF of 'é' - not a character at all!)
fmt.Println(s[2])  // 169  (the other half of 'é')
```

`s[i]` is a **byte index**, not a character index. For a string containing only ASCII, byte index and character index happen to line up. The moment a multi-byte character appears, they diverge - `s[1]` above isn't "the second character," it's an arbitrary half of one.

### len() Counts Bytes, Not Characters

```go
s := "héllo"
fmt.Println(len(s))                          // 6  (bytes)
fmt.Println(utf8.RuneCountInString(s))       // 5  (characters)
```

### Iterating: Bytes vs Runes

```go
for i := 0; i < len(s); i++ {
    fmt.Println(s[i])       // raw bytes - 6 of them
}

for i, r := range s {
    fmt.Println(i, r)       // decoded runes - 5 of them, byte-offset index
}
```

`range` over a string always decodes UTF-8 and hands you runes (as Level 6 showed). A manual `for i := 0; i < len(s); i++` loop instead walks bytes. Reach for whichever one matches what you actually need - raw bytes (e.g., for network I/O) or characters (almost everything text-related).

---

## Converting Between string, []byte, and []rune

You'll frequently need to move between these three representations. Here's the full cheat sheet:

```go
s := "héllo"

b := []byte(s)   // string -> []byte : decodes to UTF-8 bytes, len(b) == 6
r := []rune(s)   // string -> []rune : decodes to code points, len(r) == 5

s2 := string(b)  // []byte -> string : re-encodes bytes as a string
s3 := string(r)  // []rune -> string : re-encodes each rune as UTF-8

var oneRune rune = 'é'
s4 := string(oneRune)  // single rune -> string : "é"
```

Every one of these conversions **copies memory** - a new byte array or rune array is allocated. This is the price you pay for being able to safely mutate the copy without touching the original (immutable) string.

### When to Use []byte

Convert to `[]byte` when you need to:
- Mutate raw bytes in place (then convert back to `string`)
- Pass data to APIs that work with bytes (`io.Writer`, hashing, encoding)
- Avoid caring about character boundaries at all

```go
greeting := "hello"
gb := []byte(greeting)
gb[0] = 'H'
fmt.Println(string(gb))     // "Hello"
fmt.Println(greeting)       // "hello" - original string is untouched
```

### When to Use []rune

Convert to `[]rune` when you need to:
- Index or edit **characters**, not bytes (this is critical once non-ASCII text is possible)
- Count characters accurately
- Reverse, rotate, or otherwise reorder text by character

```go
word := "café"
rw := []rune(word)
rw[3] = 'e'                 // replaces the 4th CHARACTER (é), not the 4th byte
fmt.Println(string(rw))     // "cafe"
```

Compare correct rune indexing with incorrect byte indexing on the same string:

```go
r := []rune("héllo")
fmt.Println(r[1])   // 'é' - correct: the 2nd CHARACTER
fmt.Println("héllo"[1])   // 195 - wrong for characters: just the 2nd BYTE
```

### Conversion Cost Is a Real Cost

`[]byte(s)` and `[]rune(s)` both allocate and copy. If you only need to inspect a string's characters without modifying them, prefer `range` (no allocation for iteration) over converting to `[]rune` just to loop over it. Convert to `[]rune` when you specifically need **indexed, random access** or **in-place editing** by character.

---

## The strings Package Tour

The standard library's `strings` package is where almost all everyday string work happens. Here are the functions you'll use constantly:

### Searching

```go
strings.Contains(s, "fox")        // true if "fox" appears anywhere in s
strings.HasPrefix(s, "The")       // true if s starts with "The"
strings.HasSuffix(s, "dog")       // true if s ends with "dog"
strings.Index(s, "brown")         // byte index of first match, or -1
```

### Splitting and Joining

```go
strings.Split(s, " ")                       // splits on every occurrence of the separator
strings.Split("a,b,,c", ",")                // ["a" "b" "" "c"] - keeps empty fields
strings.Fields("  lots   of spaces  ")      // splits on whitespace, collapses runs, drops empties
strings.Join([]string{"a", "b", "c"}, "-")  // "a-b-c"
```

`Fields` is usually what you want for tokenizing free-form text; `Split` is for structured data (like CSV) where you need exact separators, even if that means empty fields.

### Trimming

```go
strings.TrimSpace("  hi  ")            // "hi" - removes leading/trailing whitespace
strings.Trim("***hi***", "*")          // "hi" - removes any of the given cutset from both ends
strings.TrimLeft("***hi***", "*")      // "hi***"
strings.TrimRight("***hi***", "*")     // "***hi"
strings.TrimPrefix("filename.go", "file")  // "name.go"
strings.TrimSuffix("filename.go", ".go")   // "filename"
```

### Replacing and Casing

```go
strings.Replace(s, "o", "0", 1)   // replace only the first "n" matches (1 here)
strings.ReplaceAll(s, "o", "0")   // replace every match
strings.ToUpper(s)                // "THE QUICK..."
strings.ToLower(s)                // "the quick..."
```

### Repeating and Counting

```go
strings.Repeat("ab", 3)     // "ababab"
strings.Count(s, "o")       // how many times "o" appears
```

All of these were verified in Exercise 1 - run it yourself to see real output for every function above.

---

## strings.Builder for Efficient Concatenation

### The Problem: Naive Concatenation Is O(n²)

Because strings are immutable, every `+=` on a string allocates a brand-new string and copies everything seen so far into it:

```go
result := ""
for i := 0; i < n; i++ {
    result += "x"   // allocates a new string EVERY iteration, copying all prior bytes
}
```

Building a string of length `n` this way does roughly `1 + 2 + 3 + ... + n` byte copies in total - that's **O(n²)** work, not O(n). For small `n` you'll never notice. For large `n` (parsing, log building, code generation) it gets dramatically slower as input grows.

### The Fix: strings.Builder

`strings.Builder` maintains an internal, growable `[]byte` buffer. Writes append to that buffer (doubling its capacity as needed, just like `append` does for slices), and the string is only actually built once, when you call `.String()`.

```go
var b strings.Builder
for i := 0; i < n; i++ {
    b.WriteString("x")
}
result := b.String()
```

This does **O(n)** total work.

### Measured Difference

Building a 50,000-character string both ways (Exercise 4) produced:

```
Naive +=:          109.648458ms
strings.Builder:    134.292µs
strings.Builder was ~816.5x faster
```

Exact numbers vary by machine and Go version, but the pattern - naive concatenation getting dramatically worse as the loop grows, `Builder` staying fast - is completely reliable. Always reach for `strings.Builder` (or `bytes.Buffer`) when building a string across many iterations.

### The Builder API

```go
var b strings.Builder
b.WriteString("Go")     // append a string
b.WriteByte(' ')        // append a single byte
b.WriteRune('世')        // append a single rune (handles multi-byte encoding correctly)
fmt.Println(b.String()) // "Go 世"
```

---

## strconv in a Strings Context (Recap)

Level 3 introduced `strconv.Atoi` (string → int) and `strconv.Itoa` (int → string). They come up constantly once you start combining them with the `strings` package - splitting delimited text into numbers, and rejoining numbers into text.

```go
csvLine := "10,25,7,42,3"
parts := strings.Split(csvLine, ",")

sum := 0
for _, part := range parts {
    val, err := strconv.Atoi(part)
    if err != nil {
        continue // skip anything that isn't a valid number
    }
    sum += val
}
```

Going the other direction - building a delimited string back up from numbers:

```go
var asStrings []string
for _, val := range numbers {
    asStrings = append(asStrings, strconv.Itoa(val))
}
rejoined := strings.Join(asStrings, "-")
```

Always check the `error` that `strconv.Atoi` returns - a malformed number (like `"oops"`) doesn't panic, it just returns a non-nil error, which you should handle rather than ignore.

---

## Unicode-Aware Operations

ASCII-only assumptions (`'a' <= c && c <= 'z'`) break down once your text isn't guaranteed to be English. Go's `unicode` and `unicode/utf8` packages give you correct, Unicode-aware building blocks.

### The unicode Package: Classifying Runes

```go
unicode.IsLetter(r)   // is r a letter, in ANY language/script?
unicode.IsDigit(r)    // is r a digit?
unicode.IsSpace(r)    // is r whitespace (space, tab, newline, ...)?
unicode.IsUpper(r)    // is r an uppercase letter?
unicode.IsLower(r)    // is r a lowercase letter?
unicode.ToUpper(r)    // uppercase version of r (rune -> rune)
unicode.ToLower(r)    // lowercase version of r
```

```go
for _, r := range "Go 101: Héllo, 世界! 42" {
    switch {
    case unicode.IsLetter(r):
        fmt.Printf("%q is a letter\n", r)
    case unicode.IsDigit(r):
        fmt.Printf("%q is a digit\n", r)
    case unicode.IsSpace(r):
        fmt.Printf("%q is whitespace\n", r)
    }
}
```

This correctly classifies `é`, `世`, and `界` as letters right alongside `G` and `o` - something a hardcoded ASCII range check (`'a'-'z'`, `'A'-'Z'`) would get wrong.

### The unicode/utf8 Package: Working With Encoding Directly

```go
utf8.RuneCountInString(s)   // number of runes (characters) in s - like len([]rune(s)) but without allocating
```

```go
s := "Go 101: Héllo, 世界! 42"
fmt.Println(len(s))                      // 26 (bytes)
fmt.Println(utf8.RuneCountInString(s))   // 21 (runes)
```

`utf8.DecodeRuneInString` decodes a single rune starting at a byte offset and tells you how many bytes it consumed - this is what `range` uses internally to walk a string:

```go
str := "日本語"
for i := 0; i < len(str); {
    r, size := utf8.DecodeRuneInString(str[i:])
    fmt.Printf("byte offset %d: rune=%c, encoded in %d bytes\n", i, r, size)
    i += size
}
// byte offset 0: rune=日, encoded in 3 bytes
// byte offset 3: rune=本, encoded in 3 bytes
// byte offset 6: rune=語, encoded in 3 bytes
```

In everyday code you'll almost always just use `range` instead of calling `DecodeRuneInString` manually - but knowing what `range` is doing under the hood explains exactly why the index jumps by more than 1 for multi-byte characters (the gotcha Level 6 introduced).

---

## Common String Manipulation Tasks

### Reversing a String - The Wrong Way and the Right Way

A byte-based reversal looks reasonable and works for ASCII... and silently corrupts anything else:

```go
// ❌ WRONG for non-ASCII text
func reverseBytesWRONG(s string) string {
    b := []byte(s)
    for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
        b[i], b[j] = b[j], b[i]
    }
    return string(b)
}
```

```go
reverseBytesWRONG("Hello")   // "olleH"          - looks fine, it's ASCII
reverseBytesWRONG("héllo")   // "oll\xa9\xc3h"    - BROKEN: é's 2 bytes got swapped apart
reverseBytesWRONG("日本語")   // "\x9e\xaa謜楗\xe6" - BROKEN: not even valid UTF-8 anymore
```

The fix: decode to runes first, reverse the runes, then re-encode.

```go
// ✅ RIGHT - rune-aware
func reverseRunesCorrect(s string) string {
    r := []rune(s)
    for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
        r[i], r[j] = r[j], r[i]
    }
    return string(r)
}
```

```go
reverseRunesCorrect("Hello")   // "olleH"
reverseRunesCorrect("héllo")   // "olléh"   - correct!
reverseRunesCorrect("日本語")   // "語本日"   - correct!
```

This was verified programmatically in Exercise 7: `reverseRunesCorrect("日本語") == "語本日"` is `true`. Never trust a reversal example that only shows ASCII input - it's the classic way this bug hides.

### Palindrome Checking

A rune-aware palindrome check that also ignores case, spaces, and punctuation:

```go
func isPalindrome(s string) bool {
    var cleaned []rune
    for _, r := range s {
        if unicode.IsLetter(r) || unicode.IsDigit(r) {
            cleaned = append(cleaned, unicode.ToLower(r))
        }
    }
    for i, j := 0, len(cleaned)-1; i < j; i, j = i+1, j-1 {
        if cleaned[i] != cleaned[j] {
            return false
        }
    }
    return true
}
```

```go
isPalindrome("racecar")                        // true
isPalindrome("A man a plan a canal Panama")     // true (case/spaces ignored)
isPalindrome("été")                             // true (rune-aware: 'é' compares correctly)
```

### Word and Character Counting

```go
words := strings.Fields(text)          // word count via whitespace-aware splitting
byteCount := len(text)                 // bytes
runeCount := utf8.RuneCountInString(text)  // characters

counts := make(map[string]int)         // frequency map, same pattern as Level 6
for _, w := range words {
    counts[strings.ToLower(w)]++
}
```

---

## Best Practices

### 1. Treat Strings as Character Data via range or []rune, Not Manual Byte Indexing

```go
// ✅ Good - range decodes runes correctly
for _, r := range s {
    process(r)
}

// ❌ Risky - only correct if s is guaranteed ASCII
for i := 0; i < len(s); i++ {
    process(s[i])
}
```

### 2. Use strings.Builder for Concatenation in Loops

```go
// ✅ Good
var b strings.Builder
for _, part := range parts {
    b.WriteString(part)
}
result := b.String()

// ❌ O(n²) for large n
result := ""
for _, part := range parts {
    result += part
}
```

### 3. Prefer strings Package Functions Over Hand-Rolled Loops

```go
// ✅ Good - clear and correct
if strings.Contains(s, "error") { }

// ❌ Reinventing what the standard library already does correctly
found := false
for i := 0; i < len(s)-len("error"); i++ {
    if s[i:i+5] == "error" {
        found = true
    }
}
```

### 4. Use the unicode Package Instead of Hardcoded ASCII Ranges

```go
// ✅ Good - works for any language
if unicode.IsLetter(r) { }

// ❌ WRONG - misses accented letters, CJK, etc.
if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') { }
```

### 5. Convert to []rune Only When You Need Indexed Access

```go
// ✅ Good - no allocation, just iterating
for _, r := range s { }

// Only convert when you need indexing/mutation/length-by-character:
r := []rune(s)
r[2] = 'x'
```

---

## Common Mistakes

### Mistake 1: Assuming s[i] Gives a Character

```go
// ❌ WRONG - s[1] is a raw byte, and a fragment of 'é' at that
s := "héllo"
fmt.Println(s[1])  // 195, not 'é'

// ✅ RIGHT - decode to runes first
r := []rune(s)
fmt.Println(r[1])  // 'é'
```

### Mistake 2: Byte-Based Reversal (or Any Byte-Based Reordering) on Non-ASCII Text

```go
// ❌ WRONG - corrupts multi-byte characters
b := []byte(s)
for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
    b[i], b[j] = b[j], b[i]
}

// ✅ RIGHT - reverse runes, not bytes
r := []rune(s)
for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
    r[i], r[j] = r[j], r[i]
}
```

### Mistake 3: Inefficient Concatenation With += in a Loop

```go
// ❌ WRONG - O(n²) for large n
result := ""
for _, s := range many {
    result += s
}

// ✅ RIGHT - O(n)
var b strings.Builder
for _, s := range many {
    b.WriteString(s)
}
result := b.String()
```

### Mistake 4: Using len() to Mean "Number of Characters"

```go
// ❌ WRONG ASSUMPTION - len() counts bytes
s := "日本語"
fmt.Println(len(s))  // 9, not 3!

// ✅ RIGHT
fmt.Println(utf8.RuneCountInString(s))  // 3
```

### Mistake 5: Trying to Mutate a String Directly

```go
// ❌ WRONG - compile error
s := "hello"
s[0] = 'H'
// ./main.go:5:2: cannot assign to s[0] (neither addressable nor a map index expression)

// ✅ RIGHT - convert, mutate the copy, convert back
b := []byte(s)
b[0] = 'H'
s = string(b)
```

### Mistake 6: Hardcoding ASCII Ranges for Letter/Digit Checks

```go
// ❌ WRONG - misses accented letters and non-Latin scripts entirely
isLetter := (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')

// ✅ RIGHT - Unicode-aware
isLetter := unicode.IsLetter(r)
```

---

## Summary

**Strings in Go:**
- A `string` is an immutable, UTF-8 encoded sequence of bytes - not an array of characters
- `len(s)` counts bytes; `utf8.RuneCountInString(s)` counts characters
- `s[i]` returns a `byte`, never a character

**byte vs rune:**
- `byte` = `uint8`, one raw UTF-8 code unit
- `rune` = `int32`, one decoded Unicode code point (a real character, 1-4 bytes when encoded)

**Conversions:**
- `[]byte(s)` / `string(b)` - for raw byte-level work and mutation
- `[]rune(s)` / `string(r)` - for character-level indexing, editing, and reordering
- Both allocate; use `range` instead when you only need to iterate

**The strings Package:**
- Searching, splitting/joining, trimming, replacing, casing, repeating, counting - covered in Section 4

**Efficient Building:**
- `strings.Builder` avoids the O(n²) cost of repeated `+=` concatenation

**Unicode Awareness:**
- `unicode.IsLetter/IsDigit/IsSpace/IsUpper/IsLower` and `unicode/utf8` work correctly for any language, not just ASCII

---

## Next Steps

You now understand:
- ✅ Why strings are immutable UTF-8 byte sequences
- ✅ The precise difference between byte and rune
- ✅ Every conversion between string, []byte, and []rune, and when to use each
- ✅ The core strings package functions
- ✅ Why naive concatenation is O(n²) and how strings.Builder fixes it
- ✅ Unicode-aware classification with the unicode and unicode/utf8 packages
- ✅ How to reverse and check strings correctly - rune-safe, not byte-based

**Next level:** Level 10 - Maps
- Map declaration, zero values, and the comma-ok idiom
- Adding, updating, and deleting keys
- Maps as sets and frequency counters
- Nested maps and maps of structs

You now understand text at the byte and character level - the same care you learned here will make Level 10's maps (often keyed by strings!) much easier to use correctly. Keep going! 🚀
