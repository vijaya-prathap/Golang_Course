# Level 9: Strings & Runes - INDEX

Welcome to **Level 9: Strings & Runes**! This is where you go beneath the surface of Go's most-used type and learn exactly how text is represented, converted, and manipulated correctly.

---

## 📖 What You'll Learn

- ✅ Why a Go string is an immutable, UTF-8 encoded sequence of bytes
- ✅ The precise difference between byte (uint8) and rune (int32)
- ✅ Converting between string, []byte, and []rune - and when each is needed
- ✅ The strings package: searching, splitting, trimming, replacing, casing
- ✅ strings.Builder and why naive += concatenation is O(n²)
- ✅ strconv in a strings context (bridging back to Level 3)
- ✅ Unicode-aware operations with unicode and unicode/utf8
- ✅ Correct (rune-safe) string reversal, palindrome checking, and word/character counting

---

## 🗂️ Level 9 Materials

### 1. **START_HERE.txt** (Begin Here!)
Quick overview and welcome message. Start here first!

### 2. **README.md** (Main Theory)
Comprehensive guide to:
- Strings as immutable UTF-8 byte sequences
- byte vs rune in depth
- Every conversion between string, []byte, and []rune
- The strings package tour
- strings.Builder and efficient concatenation
- strconv recap in a strings context
- Unicode-aware operations
- Common string manipulation tasks (reversal, palindromes, counting)
- Best practices
- Common mistakes

**Read Time:** 50-65 minutes
**When:** Start your session

### 3. **EXERCISES.md** (Hands-On Practice)
10 detailed, step-by-step exercises:
1. The strings package tour
2. byte vs rune fundamentals
3. Converting between string, []byte, and []rune
4. strings.Builder for efficient concatenation
5. strconv in a strings context
6. Unicode-aware classification
7. String reversal - byte-based (wrong) vs rune-based (correct)
8. Palindrome checker
9. Word frequency counting
10. Comprehensive practice (text analyzer)

Plus 3 bonus challenges: title-case converter, Caesar cipher, CSV-line splitter with quoted fields.

**Time Commitment:** 4-5 hours
**When:** After reading README

### 4. **STUDY_GUIDE.md** (Visual Learning)
- String/byte/rune conceptual diagram
- byte vs rune comparison table
- Conversion cheat sheet (string ↔ []byte ↔ []rune)
- Wrong vs right reversal comparison
- strings package function reference table
- unicode / unicode/utf8 reference table
- Common mistakes guide

**Read Time:** 30-40 minutes
**When:** Use while doing exercises

### 5. **QUICK_REFERENCE.md** (Cheat Sheet)
- Essential syntax reference
- strings package syntax by category
- strings.Builder pattern
- Common mistakes table

**Read Time:** 10-15 minutes
**When:** Quick lookup during practice

---

## 🎯 Recommended Learning Path

### Day 1: Strings, Bytes, and Runes (2 hours)
1. Read **START_HERE.txt** (10 min)
2. Read **README.md** sections 1-3 (35 min)
3. Complete Exercises 1-3 (1 hour 15 min)

### Day 2: Efficient Building & strconv (1.5 hours)
1. Read **README.md** sections 4-6 (25 min)
2. Complete Exercises 4-5 (1 hour 5 min)

### Day 3: Unicode-Aware Text Processing (2 hours)
1. Read **README.md** sections 7-8 (25 min)
2. Complete Exercises 6-8 (1.5 hours)

### Day 4: Real Programs & Consolidation (1.5 hours)
1. Read **README.md** sections 9-10 (15 min)
2. Complete Exercises 9-10 (1 hour 15 min)

### Day 5: Consolidation (1 hour)
1. Try the bonus challenges
2. Review with **QUICK_REFERENCE.md**
3. Make sure the byte vs rune distinction feels automatic

---

## 💡 Key Concepts At A Glance

### Strings Are Immutable UTF-8 Bytes
```go
s := "héllo"
len(s)                        // 6 (bytes)
s[0] = 'H'                    // compile error - strings can't be mutated
```

### byte vs rune
```go
var b byte = 'A'    // uint8  - one raw UTF-8 code unit
var r rune = 'A'    // int32  - one decoded Unicode code point
```

### Conversions
```go
b := []byte(s)      // string -> []byte
r := []rune(s)      // string -> []rune
s2 := string(r)      // []rune -> string
```

### strings.Builder
```go
var b strings.Builder
b.WriteString("Go")
result := b.String()   // O(n), unlike repeated += which is O(n^2)
```

### Correct Reversal
```go
r := []rune(s)                              // decode to characters first
for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
    r[i], r[j] = r[j], r[i]
}
return string(r)                            // safe for multi-byte text
```

---

## ✅ Prerequisites

Make sure you've completed **Level 8: Slices**

You need:
- ✅ Comfort with slices, indexing, and `len`/`cap`
- ✅ Comfort with `for`/`range` loops (Level 6), including range over strings
- ✅ Comfort with maps for frequency counting (introduced informally in Level 6)
- ✅ Basic familiarity with strconv.Atoi/Itoa (Level 3)

---

## 🎓 Learning Objectives

By the end of Level 9, you'll be able to:

- ✅ Explain why Go strings are immutable, UTF-8 encoded byte sequences
- ✅ State the exact difference between byte and rune, and why s[i] gives a byte
- ✅ Convert correctly between string, []byte, and []rune, and know when each is needed
- ✅ Use the core strings package functions fluently
- ✅ Explain why naive concatenation is O(n²) and use strings.Builder to fix it
- ✅ Use unicode and unicode/utf8 for language-agnostic text classification
- ✅ Implement a rune-safe string reversal and explain why a byte-based one breaks
- ✅ Build a palindrome checker and a word-frequency counter

---

## 📊 Statistics

- **Main Theory:** README.md covering strings, byte/rune, conversions, the strings package, Builder, strconv, Unicode, and manipulation tasks
- **Exercises:** 10 detailed exercises + 3 bonus challenges
- **Study Guide:** conceptual diagrams, function reference tables, wrong-vs-right comparisons
- **Quick Reference:** One-page cheat sheet
- **Practice Time:** 4-5 hours
- **Difficulty:** ⭐⭐⭐ (Intermediate)

---

## 🚀 How to Use These Materials

### For Reading
1. Open README.md for comprehensive theory
2. Use STUDY_GUIDE.md alongside for diagrams and reference tables
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
3. Practice the byte vs rune distinction until it's automatic

---

## 🆘 Common Questions

**Q: Why does s[0] not give me the first character?**
A: `s[i]` indexes into the string's underlying bytes. For ASCII text a byte and a character are the same thing, but for anything else (accents, emoji, CJK), one character can span multiple bytes, so `s[i]` may only be a fragment of one.

**Q: Why isn't len(s) the number of characters?**
A: `len()` returns the number of bytes in the string's UTF-8 encoding. Use `utf8.RuneCountInString(s)` (or count while ranging) for the character count.

**Q: When should I convert to []byte vs []rune?**
A: Use `[]byte` for raw/binary work or when you'll pass the data to a byte-oriented API. Use `[]rune` when you need to index, edit, or reorder by character. If you're only iterating read-only, skip both and just use `range`.

**Q: Why is my string-building loop slow?**
A: If you're using `+=` inside a loop, each iteration allocates a new string and copies everything built so far - that's O(n²) work. Use `strings.Builder` instead.

**Q: Why did my "reverse a string" function break on non-English input?**
A: It almost certainly reversed bytes instead of runes. Convert to `[]rune` first, reverse that, then convert back to a string.

---

## 🎯 Before Moving to Level 10

Make sure you can answer these questions:

- [ ] Why is a Go string an immutable, UTF-8 encoded byte sequence?
- [ ] What's the exact difference between byte and rune?
- [ ] When do you convert to []byte vs []rune, and why does each conversion allocate?
- [ ] Why is naive += string concatenation O(n²), and how does strings.Builder fix it?
- [ ] Why does a byte-based string reversal break on multi-byte characters?

---

## 📚 Recommended Reading Order

For Maximum Learning:

1. **START_HERE.txt** (10 min)
   - Gets you oriented

2. **README.md** Sections 1-3 (30 min)
   - Immutability, byte vs rune, conversions

3. **README.md** Sections 4-6 (30 min)
   - The strings package, Builder, strconv recap

4. **EXERCISES.md** Exercises 1-5 (2 hours)
   - Get hands-on with the strings package, byte/rune, conversions, and Builder

5. **STUDY_GUIDE.md** (35 min)
   - Study the conceptual diagrams and reference tables

6. **EXERCISES.md** Exercises 6-10 (2+ hours)
   - Deep practice with Unicode, reversal, palindromes, and a real text analyzer

7. **QUICK_REFERENCE.md** (10 min)
   - Create your mental cheat sheet

---

## 🏆 Success Indicators

You've mastered Level 9 when:

- ✅ You instinctively reach for []rune (or range) instead of raw indexing when working with characters
- ✅ You can explain byte vs rune to someone else without hesitation
- ✅ You never write a naive += concatenation loop for large or unbounded input
- ✅ You default to the unicode package instead of hardcoded ASCII ranges
- ✅ You can implement and explain a correct, rune-safe string reversal
- ✅ You've completed 8+ exercises
- ✅ You can explain strings to someone else

---

## 🚀 What's Next?

After Level 9, you're ready for:

**Level 10: Maps**
- Map declaration, zero values, and the comma-ok idiom
- Adding, updating, and deleting keys
- Maps as sets and frequency counters
- Nested maps and maps of structs

---

## 💬 Key Takeaway

> **A Go string is immutable UTF-8 bytes, not an array of characters - byte gives you a code unit, rune gives you a character, and knowing which one you're holding is the whole game.**

---

## 📖 Quick Links

- [START_HERE.txt](./START_HERE.txt) - Welcome & Overview
- [README.md](./README.md) - Full Theory
- [EXERCISES.md](./EXERCISES.md) - Hands-On Exercises
- [STUDY_GUIDE.md](./STUDY_GUIDE.md) - Visual Patterns
- [QUICK_REFERENCE.md](./QUICK_REFERENCE.md) - Cheat Sheet

---

Happy learning! Level 9 gives you the precision to handle real-world text correctly! 🎉

*Estimated time to complete Level 9: 4-5 hours*
*Difficulty: ⭐⭐⭐ (Intermediate)*
*Next Level: Level 10 - Maps*
