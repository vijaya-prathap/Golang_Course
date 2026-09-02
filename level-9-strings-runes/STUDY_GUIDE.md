# Level 9: Study Guide & Visual Reference

## 📚 Learning Path

### Week 1: Strings & Runes Concepts
```
Day 1:  Strings as immutable UTF-8 byte sequences
Day 2:  byte vs rune deep dive
Day 3:  Converting between string, []byte, []rune
Day 4:  The strings package tour
Day 5:  strings.Builder and efficient concatenation
Day 6:  strconv recap + Unicode-aware operations
Day 7:  Reversal, palindromes, and counting tasks
```

### Week 2: Practice & Application
```
Day 1:  Exercises 1-3 (strings tour, byte/rune, conversions)
Day 2:  Exercises 4-6 (Builder, strconv bridge, Unicode)
Day 3:  Exercises 7-8 (reversal, palindrome)
Day 4:  Exercises 9-10 (word frequency, text analyzer)
Day 5:  Bonus challenges
Day 6-7: Practice & review
```

---

## 🧵 String, Byte, Rune: The Conceptual Diagram

```
s := "héllo"

As a STRING (what you write and print):
  h  é  l  l  o                     ← 5 characters, looks simple

As BYTES (what's actually stored - UTF-8 encoded):
  [104] [195 169] [108] [108] [111]
    h      é         l     l     o
   1 byte  2 bytes  1 byte 1 byte 1 byte
                                       → len(s) == 6

As RUNES (what range decodes it into):
  104   233   108   108   111
   h     é     l     l     o
                                       → utf8.RuneCountInString(s) == 5

     byte index:  0   1   2   3   4   5
     s[i]:        h  195 169  l   l   o    ← s[1] and s[2] are HALVES of é!

     range index: 0       1       3   4   5   ← jumps past é's second byte
     range rune:  h       é       l   l   o
```

**The takeaway:** a `string` is bytes on disk/in memory; a `rune` is what those bytes decode to. `len()` measures the first; `utf8.RuneCountInString()` (or counting via `range`) measures the second.

---

## 🎯 byte vs rune At a Glance

| | `byte` | `rune` |
|---|--------|--------|
| Underlying type | `uint8` | `int32` |
| Represents | one raw UTF-8 code unit | one decoded Unicode code point (a character) |
| Range | 0-255 | 0-0x10FFFF |
| Size when encoded as UTF-8 | always 1 byte (it IS a byte) | 1-4 bytes |
| `s[i]` gives you | this | never this directly |
| `for _, x := range s` gives you | never this | this |
| ASCII text | identical to rune | identical to byte |
| Non-ASCII text | a meaningless fragment | the real character |

---

## 🔄 Conversion Cheat Sheet

```
string ──[]byte(s)──▶ []byte     decodes to UTF-8 bytes, len grows to byte count
[]byte ──string(b)──▶ string     re-encodes bytes as a string

string ──[]rune(s)──▶ []rune     decodes to code points, len shrinks/matches char count
[]rune ──string(r)──▶ string     re-encodes each rune as UTF-8

rune ──string(r)──▶ string       single character -> its UTF-8 string form ('é' -> "é")
```

```
Need to...                          Use...
─────────────────────────────────────────────────
mutate raw bytes / binary work      []byte(s), mutate, string(b)
edit/index/reorder by character     []rune(s), mutate, string(r)
just iterate characters (no edit)   for _, r := range s   (no conversion needed!)
just iterate raw bytes              for i := 0; i < len(s); i++
count characters without looping    utf8.RuneCountInString(s)
```

**Remember:** every conversion to `[]byte` or `[]rune` allocates and copies. Don't convert to `[]rune` just to loop over a string read-only - `range` already gives you runes with no extra allocation. Convert only when you need indexed access or mutation.

---

## ❌✅ Reversal: Wrong vs Right

```
INPUT: "héllo"

WRONG (byte-based):
  []byte(s)            = [104 195 169 108 108 111]
  reverse the bytes    = [111 108 108 169 195 104]
  string(...)          = "oll\xa9\xc3h"   ← BROKEN, invalid UTF-8!
                          (é's two bytes 195,169 got separated
                           by the swap and end up out of order)

RIGHT (rune-based):
  []rune(s)            = [104 233 108 108 111]     (h é l l o)
  reverse the runes    = [111 108 108 233 104]     (o l l é h)
  string(...)          = "olléh"                    ← CORRECT
```

```
INPUT: "日本語" (3 runes, 9 bytes - 3 bytes each)

WRONG (byte-based): "\x9e\xaa謜楗\xe6"   ← not even valid UTF-8 anymore
RIGHT (rune-based):  "語本日"              ← correct, verified: == "語本日"
```

**Rule of thumb:** any algorithm that reorders string data (reverse, shuffle, rotate) must operate on `[]rune`, never on `[]byte` or raw string indices - unless you've deliberately confirmed the input is ASCII-only forever.

---

## 📋 strings Package Function Reference

| Function | Purpose | Example |
|----------|---------|---------|
| `Contains(s, sub)` | does `s` contain `sub`? | `Contains("golang", "lang")` → `true` |
| `HasPrefix(s, p)` | does `s` start with `p`? | `HasPrefix("golang", "go")` → `true` |
| `HasSuffix(s, sfx)` | does `s` end with `sfx`? | `HasSuffix("golang", "lang")` → `true` |
| `Index(s, sub)` | byte index of first match, or -1 | `Index("golang", "lang")` → `2` |
| `Split(s, sep)` | split on every `sep`, keeps empties | `Split("a,,b", ",")` → `["a" "" "b"]` |
| `Fields(s)` | split on whitespace, collapses/drops empties | `Fields(" a  b ")` → `["a" "b"]` |
| `Join(strs, sep)` | join a `[]string` with `sep` between | `Join([]string{"a","b"}, "-")` → `"a-b"` |
| `TrimSpace(s)` | remove leading/trailing whitespace | `TrimSpace(" hi ")` → `"hi"` |
| `Trim(s, cutset)` | remove cutset chars from both ends | `Trim("**hi**", "*")` → `"hi"` |
| `TrimLeft/TrimRight` | remove cutset from one end only | `TrimLeft("**hi", "*")` → `"hi"` |
| `TrimPrefix/TrimSuffix` | remove an exact prefix/suffix if present | `TrimSuffix("f.go", ".go")` → `"f"` |
| `Replace(s, old, new, n)` | replace first `n` matches (-1 = all) | `Replace("aaa","a","b",1)` → `"baa"` |
| `ReplaceAll(s, old, new)` | replace every match | `ReplaceAll("aaa","a","b")` → `"bbb"` |
| `ToUpper(s)` / `ToLower(s)` | change case of every letter | `ToUpper("go")` → `"GO"` |
| `Repeat(s, n)` | repeat `s`, `n` times | `Repeat("ab", 3)` → `"ababab"` |
| `Count(s, sub)` | count non-overlapping matches | `Count("cheese", "e")` → `3` |

---

## 🧮 unicode / unicode/utf8 Reference

| Function | Package | Purpose |
|----------|---------|---------|
| `IsLetter(r)` | `unicode` | is `r` a letter, in any language? |
| `IsDigit(r)` | `unicode` | is `r` a digit? |
| `IsSpace(r)` | `unicode` | is `r` whitespace? |
| `IsUpper(r)` / `IsLower(r)` | `unicode` | letter case check |
| `ToUpper(r)` / `ToLower(r)` | `unicode` | case-convert a single rune |
| `RuneCountInString(s)` | `unicode/utf8` | count characters without allocating a `[]rune` |
| `DecodeRuneInString(s)` | `unicode/utf8` | decode the first rune and its byte width - what `range` does internally |

---

## 🐢 Why Naive Concatenation Is Slow

```
result := ""
for i := 0; i < n; i++ {
    result += "x"      // allocates a NEW string, copies ALL bytes so far
}

Iteration 1: copy 0 bytes, write 1  →  1 byte  total
Iteration 2: copy 1 byte,  write 1  →  2 bytes total
Iteration 3: copy 2 bytes, write 1  →  3 bytes total
...
Iteration n: copy n-1 bytes, write 1

Total copying work: 0+1+2+...+(n-1) ≈ n²/2   →  O(n²)
```

```
var b strings.Builder
for i := 0; i < n; i++ {
    b.WriteString("x")   // appends to a growable buffer (like append on a slice)
}
result := b.String()     // converts to string ONCE, at the end

Total work: O(n)
```

**Measured (Exercise 4, n=50,000):** naive `≈109.6ms` vs Builder `≈134µs` - about 800x faster in that run. Exact numbers vary by machine, but Builder always wins, and the gap widens as `n` grows.

---

## 🚨 Common Mistakes

### Mistake 1: Assuming s[i] Is a Character

```go
s := "héllo"
s[1]              // 195 - HALF of é, not a character

[]rune(s)[1]      // 'é' - correct
```

### Mistake 2: Byte-Based Reordering on Non-ASCII Text

```go
// ❌ corrupts multi-byte runes
b := []byte(s)
// swap b[i], b[j] ...

// ✅ reorder runes instead
r := []rune(s)
// swap r[i], r[j] ...
```

### Mistake 3: Naive += Concatenation in a Loop

```go
// ❌ O(n²)
result := ""
for _, s := range many { result += s }

// ✅ O(n)
var b strings.Builder
for _, s := range many { b.WriteString(s) }
result := b.String()
```

### Mistake 4: Using len() for Character Count

```go
len("日本語")                       // 9 (bytes)
utf8.RuneCountInString("日本語")    // 3 (characters) ← what you probably wanted
```

### Mistake 5: Hardcoded ASCII Ranges Instead of the unicode Package

```go
// ❌ misses accented letters and non-Latin scripts
(r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')

// ✅ correct for any language
unicode.IsLetter(r)
```

### Mistake 6: Trying to Mutate a String In Place

```go
s := "hello"
s[0] = 'H'   // compile error: cannot assign to s[0]
             // (neither addressable nor a map index expression)
```

---

## 📈 Progression Summary

### Understanding Level 9

Level 9 takes the rune/byte hint from Level 6 and goes all the way down:

1. **Representation** - strings are immutable, UTF-8 encoded byte sequences
2. **byte vs rune** - a raw code unit vs a decoded character
3. **Conversions** - string/[]byte/[]rune, and when each is warranted
4. **The strings package** - the everyday toolbox for text
5. **Efficient building** - strings.Builder over naive concatenation
6. **Unicode awareness** - unicode / unicode/utf8 for correct, language-agnostic logic
7. **Applied tasks** - reversal, palindromes, word/character counting

### Prerequisites for Level 10

Before moving to Level 10 (Maps), you need:

- ✅ Comfortable with slices (Level 8) - []byte and []rune are just slices
- ✅ Understand the byte vs rune distinction and why s[i] gives a byte
- ✅ Can convert between string, []byte, and []rune correctly
- ✅ Can explain why naive string concatenation is inefficient
- ✅ Comfortable using the strings, unicode, and unicode/utf8 packages

### Ready for Level 10?

Level 10 teaches Go's key-value collection type, which you've already used informally for frequency counting in this level and Level 6:
- Map declaration, zero values, and the comma-ok idiom
- Adding, updating, and deleting keys safely
- Maps as sets and frequency counters (formalizing the pattern from Exercise 9)
- Nested maps and maps of structs

---

## ✅ Checklist Before Level 10

- [ ] Can explain why a Go string is an immutable, UTF-8 encoded byte sequence
- [ ] Can state the exact difference between byte (uint8) and rune (int32)
- [ ] Know that s[i] returns a byte, and range over a string returns runes
- [ ] Can convert string ↔ []byte ↔ []rune and explain when each conversion is needed
- [ ] Can explain why naive += concatenation is O(n²) and fix it with strings.Builder
- [ ] Can use unicode.IsLetter/IsDigit/IsSpace and utf8.RuneCountInString correctly
- [ ] Can implement a rune-safe string reversal and explain why a byte-based one breaks
- [ ] Completed 8+ exercises

---

## 💡 Key Takeaways

### The Representation
A Go string is immutable bytes, UTF-8 encoded. `len()` counts those bytes, not characters.

### The Distinction
`byte` is a raw code unit (`uint8`); `rune` is a decoded character (`int32`). ASCII blurs the difference; everything else exposes it.

### The Conversions
`[]byte` for raw/binary work, `[]rune` for character-level indexing and editing - both allocate, so don't convert just to iterate.

### The Efficiency
`strings.Builder` turns O(n²) concatenation into O(n) - always prefer it inside loops.

### The Correctness
Unicode-aware functions (`unicode.IsLetter`, `utf8.RuneCountInString`, rune-based reversal) work for every language; ASCII-range checks and byte-based reordering quietly fail the moment real-world text shows up.

---

## 📚 Next Level

Level 10: Maps
- Map declaration, zero values, and the comma-ok idiom
- Adding, updating, and deleting keys
- Maps as sets and frequency counters
- Nested maps and maps of structs

You've mastered strings and runes! Keep going! 🚀
