# Level 9: Strings & Runes - Exercises

Complete all exercises in order. Each exercise builds on previous knowledge.

---

## Exercise 1: The strings Package Tour

**Objective:** Practice the core functions of the strings package

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level9-exercise1
cd ~/projects/level9-exercise1
go mod init level9.example/exercise1
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "strings"
)

func main() {
    s := "The quick brown fox jumps over the lazy dog"

    fmt.Println("=== Searching ===")
    fmt.Println("Contains \"fox\":", strings.Contains(s, "fox"))
    fmt.Println("HasPrefix \"The\":", strings.HasPrefix(s, "The"))
    fmt.Println("HasSuffix \"dog\":", strings.HasSuffix(s, "dog"))
    fmt.Println("Index of \"brown\":", strings.Index(s, "brown"))
    fmt.Println("Index of \"cat\" (not found):", strings.Index(s, "cat"))

    fmt.Println("\n=== Splitting and Joining ===")
    words := strings.Split(s, " ")
    fmt.Printf("Split by space: %q\n", words)
    fmt.Println("Number of words:", len(words))

    csv := "apple,banana,,cherry"
    fmt.Printf("Split CSV: %q\n", strings.Split(csv, ","))
    fmt.Printf("Fields (whitespace-aware): %q\n", strings.Fields("  lots   of    spaces  here  "))

    joined := strings.Join(words, "-")
    fmt.Println("Joined with '-':", joined)

    fmt.Println("\n=== Trimming ===")
    padded := "   trim me   "
    fmt.Printf("TrimSpace: %q\n", strings.TrimSpace(padded))
    fmt.Printf("Trim '*' from \"***hi***\": %q\n", strings.Trim("***hi***", "*"))
    fmt.Printf("TrimLeft '*' from \"***hi***\": %q\n", strings.TrimLeft("***hi***", "*"))
    fmt.Printf("TrimRight '*' from \"***hi***\": %q\n", strings.TrimRight("***hi***", "*"))
    fmt.Printf("TrimPrefix: %q\n", strings.TrimPrefix("filename.go", "file"))
    fmt.Printf("TrimSuffix: %q\n", strings.TrimSuffix("filename.go", ".go"))

    fmt.Println("\n=== Replacing and Casing ===")
    fmt.Println("Replace first \"o\":", strings.Replace(s, "o", "0", 1))
    fmt.Println("ReplaceAll \"o\":", strings.ReplaceAll(s, "o", "0"))
    fmt.Println("ToUpper:", strings.ToUpper(s))
    fmt.Println("ToLower:", strings.ToLower(s))

    fmt.Println("\n=== Repeat and Count ===")
    fmt.Println("Repeat \"ab\" 3 times:", strings.Repeat("ab", 3))
    fmt.Println("Count of \"o\" in s:", strings.Count(s, "o"))
    fmt.Println("Count of \"\" in \"abc\" (edge case):", strings.Count("abc", ""))
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Searching ===
Contains "fox": true
HasPrefix "The": true
HasSuffix "dog": true
Index of "brown": 10
Index of "cat" (not found): -1

=== Splitting and Joining ===
Split by space: ["The" "quick" "brown" "fox" "jumps" "over" "the" "lazy" "dog"]
Number of words: 9
Split CSV: ["apple" "banana" "" "cherry"]
Fields (whitespace-aware): ["lots" "of" "spaces" "here"]
Joined with '-': The-quick-brown-fox-jumps-over-the-lazy-dog

=== Trimming ===
TrimSpace: "trim me"
Trim '*' from "***hi***": "hi"
TrimLeft '*' from "***hi***": "hi***"
TrimRight '*' from "***hi***": "***hi"
TrimPrefix: "name.go"
TrimSuffix: "filename"

=== Replacing and Casing ===
Replace first "o": The quick br0wn fox jumps over the lazy dog
ReplaceAll "o": The quick br0wn f0x jumps 0ver the lazy d0g
ToUpper: THE QUICK BROWN FOX JUMPS OVER THE LAZY DOG
ToLower: the quick brown fox jumps over the lazy dog

=== Repeat and Count ===
Repeat "ab" 3 times: ababab
Count of "o" in s: 4
Count of "" in "abc" (edge case): 4
```

**Learning Objectives:**
- ✅ Use the core searching functions: Contains, HasPrefix, HasSuffix, Index
- ✅ Split and join strings, understanding Split vs Fields
- ✅ Trim whitespace and arbitrary characters from strings
- ✅ Replace, uppercase/lowercase, repeat, and count substrings

---

## Exercise 2: byte vs rune Fundamentals

**Objective:** See exactly what indexing a string gives you, and how byte and rune differ

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level9-exercise2
cd ~/projects/level9-exercise2
go mod init level9.example/exercise2
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "unicode/utf8"
)

func main() {
    s := "héllo"

    fmt.Println("=== Indexing a String Gives a Byte, Not a Character ===")
    fmt.Printf("s[0] = %v (type byte, prints as %c)\n", s[0], s[0])
    fmt.Printf("s[1] = %v (this is HALF of the 2-byte é, not a character!)\n", s[1])
    fmt.Printf("s[2] = %v (the other half of é)\n", s[2])
    fmt.Printf("s[3] = %v (type byte, prints as %c)\n", s[3], s[3])

    fmt.Println("\n=== len() Counts Bytes, Not Characters ===")
    fmt.Printf("len(%q) = %d bytes\n", s, len(s))
    fmt.Printf("utf8.RuneCountInString(%q) = %d runes (characters)\n", s, utf8.RuneCountInString(s))

    fmt.Println("\n=== byte is uint8, rune is int32 ===")
    var b byte = 'A'
    var r rune = 'A'
    fmt.Printf("byte 'A' = %d (type %T)\n", b, b)
    fmt.Printf("rune 'A' = %d (type %T)\n", r, r)

    fmt.Println("\n=== A Rune Beyond ASCII ===")
    var euro rune = '€'
    fmt.Printf("rune '€' = %d (type %T)\n", euro, euro)
    fmt.Printf("as a string: %s, byte length when encoded: %d\n", string(euro), len(string(euro)))

    fmt.Println("\n=== Iterating Bytes vs Runes ===")
    fmt.Println("Byte-by-byte (raw UTF-8 bytes):")
    for i := 0; i < len(s); i++ {
        fmt.Printf("  byte[%d] = %d\n", i, s[i])
    }

    fmt.Println("Rune-by-rune (decoded characters, via range):")
    for i, r := range s {
        fmt.Printf("  rune at byte-offset %d = %c (%d)\n", i, r, r)
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
=== Indexing a String Gives a Byte, Not a Character ===
s[0] = 104 (type byte, prints as h)
s[1] = 195 (this is HALF of the 2-byte é, not a character!)
s[2] = 169 (the other half of é)
s[3] = 108 (type byte, prints as l)

=== len() Counts Bytes, Not Characters ===
len("héllo") = 6 bytes
utf8.RuneCountInString("héllo") = 5 runes (characters)

=== byte is uint8, rune is int32 ===
byte 'A' = 65 (type uint8)
rune 'A' = 65 (type int32)

=== A Rune Beyond ASCII ===
rune '€' = 8364 (type int32)
as a string: €, byte length when encoded: 3

=== Iterating Bytes vs Runes ===
Byte-by-byte (raw UTF-8 bytes):
  byte[0] = 104
  byte[1] = 195
  byte[2] = 169
  byte[3] = 108
  byte[4] = 108
  byte[5] = 111
Rune-by-rune (decoded characters, via range):
  rune at byte-offset 0 = h (104)
  rune at byte-offset 1 = é (233)
  rune at byte-offset 3 = l (108)
  rune at byte-offset 4 = l (108)
  rune at byte-offset 5 = o (111)
```

**Learning Objectives:**
- ✅ See that s[i] returns a byte, and that this byte can be a meaningless fragment of a multi-byte character
- ✅ Understand that len() counts bytes while utf8.RuneCountInString counts characters
- ✅ Confirm byte is uint8 and rune is int32 using %T
- ✅ Compare byte-by-byte and rune-by-rune iteration over the same string

---

## Exercise 3: Converting Between string, []byte, and []rune

**Objective:** Practice every conversion and understand when each is needed

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level9-exercise3
cd ~/projects/level9-exercise3
go mod init level9.example/exercise3
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

func main() {
    s := "héllo"

    fmt.Println("=== string to []byte ===")
    b := []byte(s)
    fmt.Printf("[]byte(%q) = %v (len=%d)\n", s, b, len(b))

    fmt.Println("\n=== string to []rune ===")
    r := []rune(s)
    fmt.Printf("[]rune(%q) = %v (len=%d)\n", s, r, len(r))
    fmt.Printf("As characters: %c\n", r)

    fmt.Println("\n=== []byte back to string ===")
    backFromBytes := string(b)
    fmt.Printf("string(b) = %q (unchanged round trip)\n", backFromBytes)

    fmt.Println("\n=== []rune back to string ===")
    backFromRunes := string(r)
    fmt.Printf("string(r) = %q (unchanged round trip)\n", backFromRunes)

    fmt.Println("\n=== Why Convert to []byte: Mutating Bytes ===")
    // Strings are immutable - s[0] = 'H' would not compile.
    // Convert to []byte first, mutate the copy, convert back.
    greeting := "hello"
    gb := []byte(greeting)
    gb[0] = 'H'
    fmt.Printf("original string (untouched): %q\n", greeting)
    fmt.Printf("mutated []byte -> string:    %q\n", string(gb))

    fmt.Println("\n=== Why Convert to []rune: Safe Character-Level Editing ===")
    // If we tried gb[0] on a multi-byte string, we'd corrupt the UTF-8 encoding.
    // []rune lets us edit by character safely.
    word := "café"
    rw := []rune(word)
    rw[3] = 'e' // replace the 4th CHARACTER (é), not the 4th byte
    fmt.Printf("original: %q -> edited by rune: %q\n", word, string(rw))

    fmt.Println("\n=== Single Character (rune) to string ===")
    var singleRune rune = 'é'
    fmt.Printf("string(rune) = %q\n", string(singleRune))

    fmt.Println("\n=== Indexing []rune vs Indexing string Directly ===")
    fmt.Printf("[]rune(s)[1] = %c (correct: the 2nd CHARACTER)\n", r[1])
    fmt.Printf("s[1] = %d (WRONG for characters: just the 2nd BYTE)\n", s[1])
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== string to []byte ===
[]byte("héllo") = [104 195 169 108 108 111] (len=6)

=== string to []rune ===
[]rune("héllo") = [104 233 108 108 111] (len=5)
As characters: [h é l l o]

=== []byte back to string ===
string(b) = "héllo" (unchanged round trip)

=== []rune back to string ===
string(r) = "héllo" (unchanged round trip)

=== Why Convert to []byte: Mutating Bytes ===
original string (untouched): "hello"
mutated []byte -> string:    "Hello"

=== Why Convert to []rune: Safe Character-Level Editing ===
original: "café" -> edited by rune: "cafe"

=== Single Character (rune) to string ===
string(rune) = "é"

=== Indexing []rune vs Indexing string Directly ===
[]rune(s)[1] = é (correct: the 2nd CHARACTER)
s[1] = 195 (WRONG for characters: just the 2nd BYTE)
```

**Learning Objectives:**
- ✅ Convert string to []byte and []rune, and back again
- ✅ Understand that []byte(s) has the same length as len(s), while []rune(s) has the rune count
- ✅ Mutate a copy of string data safely via []byte and []rune, without touching the original
- ✅ See directly why []rune indexing is correct for characters while string indexing is not

---

## Exercise 4: strings.Builder for Efficient Concatenation

**Objective:** Measure why naive += concatenation is inefficient and confirm strings.Builder fixes it

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level9-exercise4
cd ~/projects/level9-exercise4
go mod init level9.example/exercise4
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "strings"
    "time"
)

func concatNaive(n int) string {
    result := ""
    for i := 0; i < n; i++ {
        result += "x"
    }
    return result
}

func concatBuilder(n int) string {
    var b strings.Builder
    for i := 0; i < n; i++ {
        b.WriteString("x")
    }
    return b.String()
}

func main() {
    const n = 50000

    fmt.Println("=== Naive += Concatenation ===")
    start := time.Now()
    naiveResult := concatNaive(n)
    naiveElapsed := time.Since(start)
    fmt.Printf("Built a string of length %d\n", len(naiveResult))
    fmt.Printf("Time taken: %v\n", naiveElapsed)

    fmt.Println("\n=== strings.Builder Concatenation ===")
    start = time.Now()
    builderResult := concatBuilder(n)
    builderElapsed := time.Since(start)
    fmt.Printf("Built a string of length %d\n", len(builderResult))
    fmt.Printf("Time taken: %v\n", builderElapsed)

    fmt.Println("\n=== Comparison ===")
    fmt.Printf("Results identical: %v\n", naiveResult == builderResult)
    if builderElapsed > 0 {
        fmt.Printf("strings.Builder was ~%.1fx faster\n", float64(naiveElapsed)/float64(builderElapsed))
    }

    fmt.Println("\n=== Why: strings are immutable ===")
    fmt.Println("Every += on a string allocates a brand-new string and copies")
    fmt.Println("everything so far into it. Building a string of length N one")
    fmt.Println("character at a time this way does roughly 1+2+3+...+N byte")
    fmt.Println("copies -- O(n^2) total work. strings.Builder grows an internal")
    fmt.Println("[]byte buffer (doubling capacity as needed) and only converts")
    fmt.Println("to a string once at the end, so it's O(n).")

    fmt.Println("\n=== Builder Also Supports WriteByte and WriteRune ===")
    var b strings.Builder
    b.WriteString("Go")
    b.WriteByte(' ')
    b.WriteRune('世')
    fmt.Println(b.String())
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Naive += Concatenation ===
Built a string of length 50000
Time taken: 109.648458ms

=== strings.Builder Concatenation ===
Built a string of length 50000
Time taken: 134.292µs

=== Comparison ===
Results identical: true
strings.Builder was ~816.5x faster

=== Why: strings are immutable ===
Every += on a string allocates a brand-new string and copies
everything so far into it. Building a string of length N one
character at a time this way does roughly 1+2+3+...+N byte
copies -- O(n^2) total work. strings.Builder grows an internal
[]byte buffer (doubling capacity as needed) and only converts
to a string once at the end, so it's O(n).

=== Builder Also Supports WriteByte and WriteRune ===
Go 世
```

**Note:** The exact timings and speedup factor will vary depending on your machine and Go version - what matters is that `strings.Builder` is dramatically and reliably faster, and the gap grows as `n` grows. Try changing `n` to `200000` and see the naive version slow down much more than proportionally.

**Learning Objectives:**
- ✅ Understand why naive string concatenation in a loop is O(n²)
- ✅ Use strings.Builder's WriteString, WriteByte, and WriteRune methods
- ✅ Measure a real, visible performance difference between the two approaches

---

## Exercise 5: strconv in a Strings Context

**Objective:** Bridge Level 3's strconv functions with strings package operations

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level9-exercise5
cd ~/projects/level9-exercise5
go mod init level9.example/exercise5
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "strconv"
    "strings"
)

func main() {
    fmt.Println("=== Recap from Level 3: strconv Basics ===")
    n, err := strconv.Atoi("42")
    fmt.Println("Atoi(\"42\") ->", n, err)
    s := strconv.Itoa(99)
    fmt.Println("Itoa(99) ->", s)

    fmt.Println("\n=== Combining strconv With the strings Package ===")
    csvLine := "10,25,7,42,3"
    parts := strings.Split(csvLine, ",")
    fmt.Printf("Split %q into %v\n", csvLine, parts)

    sum := 0
    var numbers []int
    for _, part := range parts {
        val, err := strconv.Atoi(part)
        if err != nil {
            fmt.Println("skipping invalid number:", part)
            continue
        }
        numbers = append(numbers, val)
        sum += val
    }
    fmt.Println("Parsed numbers:", numbers)
    fmt.Println("Sum:", sum)

    fmt.Println("\n=== Building a String Back From Numbers ===")
    var asStrings []string
    for _, val := range numbers {
        asStrings = append(asStrings, strconv.Itoa(val))
    }
    rejoined := strings.Join(asStrings, "-")
    fmt.Println("Rejoined with '-':", rejoined)

    fmt.Println("\n=== Handling a Bad Number Gracefully ===")
    badLine := "10,oops,30"
    badParts := strings.Split(badLine, ",")
    for _, part := range badParts {
        val, err := strconv.Atoi(part)
        if err != nil {
            fmt.Printf("%q is not a valid number: %v\n", part, err)
            continue
        }
        fmt.Printf("%q parsed as %d\n", part, val)
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
=== Recap from Level 3: strconv Basics ===
Atoi("42") -> 42 <nil>
Itoa(99) -> 99

=== Combining strconv With the strings Package ===
Split "10,25,7,42,3" into [10 25 7 42 3]
Parsed numbers: [10 25 7 42 3]
Sum: 87

=== Building a String Back From Numbers ===
Rejoined with '-': 10-25-7-42-3

=== Handling a Bad Number Gracefully ===
"10" parsed as 10
"oops" is not a valid number: strconv.Atoi: parsing "oops": invalid syntax
"30" parsed as 30
```

**Learning Objectives:**
- ✅ Recall strconv.Atoi and strconv.Itoa from Level 3
- ✅ Combine strconv with strings.Split and strings.Join in a realistic pipeline
- ✅ Handle strconv.Atoi errors gracefully instead of ignoring them

---

## Exercise 6: Unicode-Aware Classification

**Objective:** Classify runes correctly using the unicode and unicode/utf8 packages

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level9-exercise6
cd ~/projects/level9-exercise6
go mod init level9.example/exercise6
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "unicode"
    "unicode/utf8"
)

func main() {
    s := "Go 101: Héllo, 世界! 42"

    fmt.Println("=== Classifying Each Rune ===")
    for _, r := range s {
        switch {
        case unicode.IsLetter(r):
            fmt.Printf("%q is a letter\n", r)
        case unicode.IsDigit(r):
            fmt.Printf("%q is a digit\n", r)
        case unicode.IsSpace(r):
            fmt.Printf("%q is whitespace\n", r)
        case unicode.IsPunct(r):
            fmt.Printf("%q is punctuation\n", r)
        default:
            fmt.Printf("%q is something else\n", r)
        }
    }

    fmt.Println("\n=== Counting Categories ===")
    var letters, digits, spaces, other int
    for _, r := range s {
        switch {
        case unicode.IsLetter(r):
            letters++
        case unicode.IsDigit(r):
            digits++
        case unicode.IsSpace(r):
            spaces++
        default:
            other++
        }
    }
    fmt.Printf("letters=%d digits=%d spaces=%d other=%d\n", letters, digits, spaces, other)

    fmt.Println("\n=== IsUpper / IsLower ===")
    for _, r := range "Hé9" {
        fmt.Printf("%q: IsUpper=%v IsLower=%v\n", r, unicode.IsUpper(r), unicode.IsLower(r))
    }

    fmt.Println("\n=== unicode/utf8: RuneCountInString ===")
    fmt.Printf("%q has %d bytes but %d runes\n", s, len(s), utf8.RuneCountInString(s))

    fmt.Println("\n=== unicode/utf8: DecodeRuneInString (manual decode loop) ===")
    str := "日本語"
    for i := 0; i < len(str); {
        r, size := utf8.DecodeRuneInString(str[i:])
        fmt.Printf("byte offset %d: rune=%c (%d), encoded in %d bytes\n", i, r, r, size)
        i += size
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
=== Classifying Each Rune ===
'G' is a letter
'o' is a letter
' ' is whitespace
'1' is a digit
'0' is a digit
'1' is a digit
':' is punctuation
' ' is whitespace
'H' is a letter
'é' is a letter
'l' is a letter
'l' is a letter
'o' is a letter
',' is punctuation
' ' is whitespace
'世' is a letter
'界' is a letter
'!' is punctuation
' ' is whitespace
'4' is a digit
'2' is a digit

=== Counting Categories ===
letters=9 digits=5 spaces=4 other=3

=== IsUpper / IsLower ===
'H': IsUpper=true IsLower=false
'é': IsUpper=false IsLower=true
'9': IsUpper=false IsLower=false

=== unicode/utf8: RuneCountInString ===
"Go 101: Héllo, 世界! 42" has 26 bytes but 21 runes

=== unicode/utf8: DecodeRuneInString (manual decode loop) ===
byte offset 0: rune=日 (26085), encoded in 3 bytes
byte offset 3: rune=本 (26412), encoded in 3 bytes
byte offset 6: rune=語 (35486), encoded in 3 bytes
```

**Learning Objectives:**
- ✅ Classify runes with unicode.IsLetter/IsDigit/IsSpace/IsPunct correctly across languages
- ✅ Use unicode.IsUpper/IsLower to check letter case
- ✅ Use utf8.RuneCountInString to count characters instead of bytes
- ✅ Understand how utf8.DecodeRuneInString decodes one rune at a time (what range does internally)

---

## Exercise 7: String Reversal - Byte-Based (Wrong) vs Rune-Based (Correct)

**Objective:** Prove, with real multi-byte output, why a byte-based reversal is broken and a rune-based one is correct

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level9-exercise7
cd ~/projects/level9-exercise7
go mod init level9.example/exercise7
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import "fmt"

// reverseBytesWRONG reverses a string byte by byte.
// This is WRONG for any string containing multi-byte UTF-8 characters:
// it scrambles the byte order of multi-byte runes into invalid UTF-8.
func reverseBytesWRONG(s string) string {
    b := []byte(s)
    for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
        b[i], b[j] = b[j], b[i]
    }
    return string(b)
}

// reverseRunesCorrect reverses a string by decoding it into runes first,
// reversing the runes, then re-encoding to a string. This preserves
// multi-byte characters correctly.
func reverseRunesCorrect(s string) string {
    r := []rune(s)
    for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
        r[i], r[j] = r[j], r[i]
    }
    return string(r)
}

func main() {
    fmt.Println("=== ASCII String: Both Approaches Agree ===")
    ascii := "Hello"
    fmt.Printf("original:       %q\n", ascii)
    fmt.Printf("byte-reversed:  %q\n", reverseBytesWRONG(ascii))
    fmt.Printf("rune-reversed:  %q\n", reverseRunesCorrect(ascii))

    fmt.Println("\n=== Multi-Byte String: \"héllo\" ===")
    s1 := "héllo"
    fmt.Printf("original:       %q\n", s1)
    fmt.Printf("byte-reversed:  %q  <-- BROKEN: é's 2 bytes got split apart\n", reverseBytesWRONG(s1))
    fmt.Printf("rune-reversed:  %q  <-- CORRECT\n", reverseRunesCorrect(s1))

    fmt.Println("\n=== Multi-Byte String: \"日本語\" (Japanese, 3 bytes per rune) ===")
    s2 := "日本語"
    fmt.Printf("original:       %q\n", s2)
    fmt.Printf("byte-reversed:  %q  <-- BROKEN: not even valid UTF-8 anymore\n", reverseBytesWRONG(s2))
    fmt.Printf("rune-reversed:  %q  <-- CORRECT: %s\n", reverseRunesCorrect(s2), reverseRunesCorrect(s2))

    fmt.Println("\n=== Verifying Correctness Programmatically ===")
    expected := "語本日"
    got := reverseRunesCorrect(s2)
    fmt.Printf("expected reversal of %q: %q\n", s2, expected)
    fmt.Printf("got:                     %q\n", got)
    fmt.Printf("match: %v\n", got == expected)
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== ASCII String: Both Approaches Agree ===
original:       "Hello"
byte-reversed:  "olleH"
rune-reversed:  "olleH"

=== Multi-Byte String: "héllo" ===
original:       "héllo"
byte-reversed:  "oll\xa9\xc3h"  <-- BROKEN: é's 2 bytes got split apart
rune-reversed:  "olléh"  <-- CORRECT

=== Multi-Byte String: "日本語" (Japanese, 3 bytes per rune) ===
original:       "日本語"
byte-reversed:  "\x9e\xaa謜楗\xe6"  <-- BROKEN: not even valid UTF-8 anymore
rune-reversed:  "語本日"  <-- CORRECT: 語本日

=== Verifying Correctness Programmatically ===
expected reversal of "日本語": "語本日"
got:                     "語本日"
match: true
```

**Why this matters:** Notice that both functions produce identical, correct-looking output for `"Hello"` - an ASCII-only reversal example can hide a badly broken implementation. The moment a multi-byte character shows up, `reverseBytesWRONG` produces garbage (`%q` even shows raw `\xNN` escapes - a sign the bytes no longer form valid UTF-8), while `reverseRunesCorrect` reverses the actual characters correctly.

**Learning Objectives:**
- ✅ See a byte-based reversal visibly corrupt multi-byte UTF-8 text
- ✅ Implement a correct rune-based reversal
- ✅ Verify correctness programmatically, not just by eyeballing output
- ✅ Internalize why ASCII-only test cases are not enough to trust string code

---

## Exercise 8: Palindrome Checker

**Objective:** Build a rune-aware palindrome checker that ignores case, spaces, and punctuation

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level9-exercise8
cd ~/projects/level9-exercise8
go mod init level9.example/exercise8
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "strings"
    "unicode"
)

// isPalindrome checks whether s reads the same forwards and backwards,
// ignoring case, spaces, and punctuation. It works rune-by-rune so it's
// safe for multi-byte Unicode characters.
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

func main() {
    tests := []string{
        "racecar",
        "hello",
        "A man a plan a canal Panama",
        "No 'x' in Nixon",
        "Was it a car or a cat I saw?",
        "level",
        "Go",
        "level",
        "été", // French, contains é - rune-aware check
    }

    fmt.Println("=== Palindrome Checks ===")
    for _, t := range tests {
        fmt.Printf("%-35q -> %v\n", t, isPalindrome(t))
    }

    fmt.Println("\n=== Why Rune-Aware Matters ===")
    var cleanedEte []rune
    for _, r := range "été" {
        if unicode.IsLetter(r) {
            cleanedEte = append(cleanedEte, unicode.ToLower(r))
        }
    }
    fmt.Printf("%q cleaned to lowercase runes: %v (as characters: %c)\n", "été", cleanedEte, cleanedEte)
    fmt.Println(strings.Repeat("-", 20))
    fmt.Println("A byte-based approach would compare individual UTF-8 bytes of")
    fmt.Println("'é' out of order with the rest of the string and misjudge palindromes")
    fmt.Println("containing multi-byte characters.")
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Palindrome Checks ===
"racecar"                           -> true
"hello"                             -> false
"A man a plan a canal Panama"       -> true
"No 'x' in Nixon"                   -> true
"Was it a car or a cat I saw?"      -> true
"level"                             -> true
"Go"                                -> false
"level"                             -> true
"été"                               -> true

=== Why Rune-Aware Matters ===
"été" cleaned to lowercase runes: [233 116 233] (as characters: [é t é])
--------------------
A byte-based approach would compare individual UTF-8 bytes of
'é' out of order with the rest of the string and misjudge palindromes
containing multi-byte characters.
```

**Learning Objectives:**
- ✅ Build a palindrome checker that ignores case, spaces, and punctuation
- ✅ Use unicode.IsLetter/IsDigit and unicode.ToLower together to normalize input
- ✅ Confirm the checker works correctly on a string containing a multi-byte character

---

## Exercise 9: Word Frequency Counting

**Objective:** Apply Level 6's map-based frequency-counting pattern to real text

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level9-exercise9
cd ~/projects/level9-exercise9
go mod init level9.example/exercise9
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "sort"
    "strings"
)

func main() {
    text := "the quick brown fox jumps over the lazy dog the fox runs the dog barks"

    fmt.Println("=== Splitting Into Words ===")
    words := strings.Fields(text)
    fmt.Println("Words:", words)
    fmt.Println("Total word count:", len(words))

    fmt.Println("\n=== Counting Word Frequency (map pattern from Level 6) ===")
    counts := make(map[string]int)
    for _, w := range words {
        counts[strings.ToLower(w)]++
    }

    // Sort keys for stable, repeatable output (map iteration order is
    // randomized, as covered in Level 6).
    keys := make([]string, 0, len(counts))
    for k := range counts {
        keys = append(keys, k)
    }
    sort.Strings(keys)

    for _, k := range keys {
        fmt.Printf("%-8s %d\n", k, counts[k])
    }

    fmt.Println("\n=== Finding the Most Common Word ===")
    mostCommon := ""
    highest := 0
    for _, k := range keys {
        if counts[k] > highest {
            mostCommon = k
            highest = counts[k]
        }
    }
    fmt.Printf("Most common word: %q (appears %d times)\n", mostCommon, highest)

    fmt.Println("\n=== Unique Word Count ===")
    fmt.Println("Unique words:", len(counts))
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Splitting Into Words ===
Words: [the quick brown fox jumps over the lazy dog the fox runs the dog barks]
Total word count: 15

=== Counting Word Frequency (map pattern from Level 6) ===
barks    1
brown    1
dog      2
fox      2
jumps    1
lazy     1
over     1
quick    1
runs     1
the      4

=== Finding the Most Common Word ===
Most common word: "the" (appears 4 times)

=== Unique Word Count ===
Unique words: 10
```

**Learning Objectives:**
- ✅ Tokenize text into words with strings.Fields
- ✅ Reuse Level 6's map[string]int frequency-counting pattern
- ✅ Sort map keys for stable, repeatable output
- ✅ Find the most frequent entry in a frequency map

---

## Exercise 10: Comprehensive Practice — Text Analyzer

**Objective:** Combine everything from this level into a small but complete text analysis tool

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level9-exercise10
cd ~/projects/level9-exercise10
go mod init level9.example/exercise10
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "sort"
    "strings"
    "unicode"
    "unicode/utf8"
)

type textStats struct {
    byteCount   int
    runeCount   int
    wordCount   int
    letterCount int
    digitCount  int
    spaceCount  int
    mostCommon  string
    mostCommonN int
    longestWord string
    wordFreq    map[string]int
}

func analyze(text string) textStats {
    stats := textStats{
        byteCount: len(text),
        runeCount: utf8.RuneCountInString(text),
        wordFreq:  make(map[string]int),
    }

    for _, r := range text {
        switch {
        case unicode.IsLetter(r):
            stats.letterCount++
        case unicode.IsDigit(r):
            stats.digitCount++
        case unicode.IsSpace(r):
            stats.spaceCount++
        }
    }

    words := strings.Fields(text)
    stats.wordCount = len(words)

    for _, w := range words {
        cleaned := strings.ToLower(strings.Trim(w, ".,!?;:\"'"))
        if cleaned == "" {
            continue
        }
        stats.wordFreq[cleaned]++
        if len([]rune(cleaned)) > len([]rune(stats.longestWord)) {
            stats.longestWord = cleaned
        }
    }

    for word, count := range stats.wordFreq {
        if count > stats.mostCommonN {
            stats.mostCommon = word
            stats.mostCommonN = count
        } else if count == stats.mostCommonN && word < stats.mostCommon {
            // tie-break alphabetically for stable output
            stats.mostCommon = word
        }
    }

    return stats
}

func main() {
    text := "Go is fun. Go is fast. Go is simple, and Go is powerful! " +
        "Programming in Go feels productive."

    fmt.Println("=== Text Analyzer ===")
    fmt.Printf("Input: %q\n\n", text)

    stats := analyze(text)

    fmt.Println("--- Basic Counts ---")
    fmt.Printf("Bytes:   %d\n", stats.byteCount)
    fmt.Printf("Runes:   %d\n", stats.runeCount)
    fmt.Printf("Words:   %d\n", stats.wordCount)
    fmt.Printf("Letters: %d\n", stats.letterCount)
    fmt.Printf("Digits:  %d\n", stats.digitCount)
    fmt.Printf("Spaces:  %d\n", stats.spaceCount)

    fmt.Println("\n--- Word Frequency (sorted) ---")
    keys := make([]string, 0, len(stats.wordFreq))
    for k := range stats.wordFreq {
        keys = append(keys, k)
    }
    sort.Strings(keys)
    for _, k := range keys {
        fmt.Printf("%-10s %d\n", k, stats.wordFreq[k])
    }

    fmt.Println("\n--- Highlights ---")
    fmt.Printf("Most common word: %q (%d times)\n", stats.mostCommon, stats.mostCommonN)
    fmt.Printf("Longest word:     %q (%d characters)\n", stats.longestWord, len([]rune(stats.longestWord)))
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Text Analyzer ===
Input: "Go is fun. Go is fast. Go is simple, and Go is powerful! Programming in Go feels productive."

--- Basic Counts ---
Bytes:   92
Runes:   92
Words:   18
Letters: 70
Digits:  0
Spaces:  17

--- Word Frequency (sorted) ---
and        1
fast       1
feels      1
fun        1
go         5
in         1
is         4
powerful   1
productive 1
programming 1
simple     1

--- Highlights ---
Most common word: "go" (5 times)
Longest word:     "programming" (11 characters)
```

**Learning Objectives:**
- ✅ Combine byte/rune counting, Unicode classification, and word frequency into one program
- ✅ Design a struct to hold multiple related statistics
- ✅ Find both "most common" and "longest" from the same frequency map
- ✅ Apply strings.Trim to strip surrounding punctuation before counting words

---

## Bonus Challenges

### Challenge 1: Unicode-Aware Title Case Converter

Write a function that capitalizes the first letter of every word and lowercases the rest, correctly handling accented letters (and correctly leaving scripts with no case, like Japanese, unchanged).

```bash
mkdir -p ~/projects/level9-bonus1
cd ~/projects/level9-bonus1
go mod init level9.example/bonus1
```

**Hints:**
- Split into words with `strings.Fields`, convert each to `[]rune`
- Use `unicode.ToUpper` on the first rune, `unicode.ToLower` on a lowercased copy of the rest
- Test with a string that mixes ASCII, accented Latin letters, and CJK characters - CJK has no case, so those characters should pass through unchanged

**Expected behavior (verified):**
```
"the quick brown fox"          -> "The Quick Brown Fox"
"café society in münchen"      -> "Café Society In München"
"日本語 テスト"                      -> "日本語 テスト"
```

### Challenge 2: Caesar Cipher

Write a Caesar cipher that shifts each letter by N positions in the alphabet, preserving case, leaving non-letters untouched, and correctly wrapping around from 'z' back to 'a' (and handling negative shifts for decryption).

```bash
mkdir -p ~/projects/level9-bonus2
cd ~/projects/level9-bonus2
go mod init level9.example/bonus2
```

**Hints:**
- Convert the string to `[]rune` and work rune-by-rune
- Normalize the shift with `shift = ((shift % 26) + 26) % 26` so negative shifts work correctly
- Handle `'a'-'z'` and `'A'-'Z'` as two separate ranges so case is preserved
- Verify your cipher round-trips: `decrypt(encrypt(msg, n), n) == msg`

**Expected behavior (verified):**
```
Original:  "Hello, World! Go 101."
Shift +3:  "Khoor, Zruog! Jr 101."
Shift -3:  "Hello, World! Go 101."
caesarShift("xyz", 3) = "abc"   (wraps around correctly)
caesarShift("abc", -1) = "zab"  (negative shift wraps correctly)
```

### Challenge 3: CSV-Line Splitter With Quoted Fields

Write a simplified CSV line splitter that splits on commas but treats commas inside double-quoted fields as part of the field, not a separator (you do not need to handle escaped `""` quotes inside a field).

```bash
mkdir -p ~/projects/level9-bonus3
cd ~/projects/level9-bonus3
go mod init level9.example/bonus3
```

**Hints:**
- Walk the line rune by rune, tracking a boolean `inQuotes` flag that flips on every `"`
- Use a `strings.Builder` to accumulate the current field
- Only split on `,` when `inQuotes` is false

**Expected behavior (verified):**
```
input:  "Smith, Bob",25,"Product Manager"
fields: ["Smith, Bob" "25" "Product Manager"]

input:  Carol,,"New York, NY"
fields: ["Carol" "" "New York, NY"]
```

---

## What You've Learned

After completing these 10 exercises, you can:

✅ Use the strings package's core searching, splitting, trimming, replacing, and casing functions
✅ Explain exactly what byte and rune are, and what string indexing actually returns
✅ Convert correctly between string, []byte, and []rune, and know when each is needed
✅ Explain and demonstrate why naive string concatenation is O(n²), and fix it with strings.Builder
✅ Combine strconv with the strings package to parse and rebuild delimited text
✅ Classify text with the unicode and unicode/utf8 packages, correctly, across languages
✅ Reverse a string safely (rune-based), and prove a byte-based reversal is broken
✅ Build a rune-aware palindrome checker and a word-frequency counter
✅ Combine these skills into a small, real text analyzer

---

## Next Level

Level 10: Maps
- Map declaration, zero values, and the comma-ok idiom
- Adding, updating, and deleting keys
- Maps as sets and frequency counters
- Nested maps and maps of structs

Great work! You're mastering Go! 🚀
