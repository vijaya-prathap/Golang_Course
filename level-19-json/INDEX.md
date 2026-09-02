# Level 19: JSON - INDEX

Welcome to **Level 19: JSON**! This is where your programs start speaking the format the rest of the world uses to exchange structured data.

---

## 📖 What You'll Learn

- ✅ What JSON is, and why encoding/json is Go's standard tool for it
- ✅ json.Marshal (struct → JSON) and json.Unmarshal (JSON → struct)
- ✅ Struct tags: rename, exclude (json:"-"), and omitempty
- ✅ Nested structs and slices of structs
- ✅ Pretty-printing with json.MarshalIndent
- ✅ Dynamic JSON with map[string]interface{} and type assertions
- ✅ json.RawMessage for deferred/partial parsing
- ✅ Streaming JSON with json.Decoder and json.Encoder
- ✅ Custom marshaling with MarshalJSON/UnmarshalJSON
- ✅ Real json.SyntaxError and json.UnmarshalTypeError messages

---

## 🗂️ Level 19 Materials

### 1. **START_HERE.txt** (Begin Here!)
Quick overview and welcome message. Start here first!

### 2. **README.md** (Main Theory)
Comprehensive guide to:
- Marshal, Unmarshal, and struct tags
- Nested structs, pretty-printing, dynamic JSON
- RawMessage, streaming, and custom marshaling
- Real, verified error messages
- Best practices
- Common mistakes

**Read Time:** 50-65 minutes
**When:** Start your session

### 3. **EXERCISES.md** (Hands-On Practice)
10 detailed, step-by-step exercises:
1. Basic struct-to-JSON marshal
2. JSON-to-struct unmarshal round trip
3. Struct tags - rename, exclude, omitempty
4. Nested structs with a slice of sub-structs
5. Pretty-printing with MarshalIndent
6. Dynamic JSON with map[string]interface{}
7. json.RawMessage for partial parsing
8. Streaming with json.Decoder/json.Encoder
9. Custom MarshalJSON/UnmarshalJSON
10. Comprehensive practice (config file load/save)

Plus 3 bonus challenges: a JSON-backed key-value store, JSON Lines conversion, and polymorphic JSON parsing.

**Time Commitment:** 5-6 hours
**When:** After reading README

### 4. **STUDY_GUIDE.md** (Visual Learning)
- Struct tag syntax reference
- Go type ↔ JSON type mapping table
- Marshal/Unmarshal round-trip diagram
- Exported vs unexported visibility reminder
- Dynamic JSON navigation diagram (interface{} → type assertion → value)
- Common mistakes table

**Read Time:** 30-40 minutes
**When:** Use while doing exercises

### 5. **QUICK_REFERENCE.md** (Cheat Sheet)
- Essential syntax reference
- Struct tag cheat sheet
- Dynamic JSON, RawMessage, and streaming snippets
- Common mistakes table

**Read Time:** 10-15 minutes
**When:** Quick lookup during practice

---

## 🎯 Recommended Learning Path

### Day 1: Marshal, Unmarshal, and Tags (2.5 hours)
1. Read **START_HERE.txt** (10 min)
2. Read **README.md** sections 1-4 (40 min)
3. Complete Exercises 1-3 (1 hour 40 min)

### Day 2: Nesting and Pretty-Printing (1.5 hours)
1. Read **README.md** sections 5-6 (20 min)
2. Complete Exercises 4-5 (1 hour 10 min)

### Day 3: Dynamic JSON and RawMessage (2 hours)
1. Read **README.md** sections 7-8 (25 min)
2. Complete Exercises 6-7 (1 hour 35 min)

### Day 4: Streaming and Custom Marshaling (2 hours)
1. Read **README.md** sections 9-10 (25 min)
2. Complete Exercises 8-9 (1 hour 35 min)

### Day 5: Errors and Comprehensive Practice (2 hours)
1. Read **README.md** sections 11-13 (25 min)
2. Complete Exercise 10 (1 hour 35 min)

### Day 6: Consolidation (1.5 hours)
1. Try the bonus challenges
2. Review with **QUICK_REFERENCE.md** and **STUDY_GUIDE.md**

---

## 💡 Key Concepts At A Glance

### Marshal and Unmarshal
```go
data, err := json.Marshal(v)     // Go -> JSON bytes
err := json.Unmarshal(data, &v)  // JSON bytes -> Go (needs a pointer!)
```

### Struct Tags
```go
type User struct {
    Name string `json:"name"`
    Pass string `json:"-"`
    Bio  string `json:"bio,omitempty"`
}
```

### Dynamic JSON
```go
var v map[string]interface{}
json.Unmarshal(data, &v)
n, ok := v["count"].(float64) // JSON numbers are always float64
```

### Custom Marshaling
```go
func (t MyType) MarshalJSON() ([]byte, error) { ... }
func (t *MyType) UnmarshalJSON(data []byte) error { ... }
```

---

## ✅ Prerequisites

Make sure you've completed **Level 18: File Handling**

You need:
- ✅ Comfortable with os/bufio file reading and writing
- ✅ Comfortable with exported/unexported fields (Level 2)
- ✅ Comfortable with struct tags at a basic level (Level 13 preview)
- ✅ Comfortable with type assertions (Level 15)
- ✅ Comfortable with error values and error types (Level 16)

---

## 🎓 Learning Objectives

By the end of Level 19, you'll be able to:

- ✅ Marshal and unmarshal structs, slices of structs, and nested structs
- ✅ Control JSON shape with struct tags (rename, exclude, omitempty)
- ✅ Pretty-print JSON for humans with MarshalIndent
- ✅ Parse JSON of unknown shape safely with map[string]interface{}
- ✅ Defer parsing part of a document with json.RawMessage
- ✅ Stream JSON with json.Decoder/json.Encoder
- ✅ Implement custom MarshalJSON/UnmarshalJSON for a type
- ✅ Recognize and handle real json.SyntaxError/json.UnmarshalTypeError output

---

## 📊 Statistics

- **Main Theory:** README.md covering Marshal/Unmarshal, tags, nesting, dynamic JSON, RawMessage, streaming, custom marshaling, and real errors
- **Exercises:** 10 detailed exercises + 3 bonus challenges
- **Study Guide:** type mapping table, round-trip diagram, dynamic JSON navigation diagram
- **Quick Reference:** One-page cheat sheet
- **Practice Time:** 5-6 hours
- **Difficulty:** ⭐⭐⭐ (Intermediate-Advanced)

---

## 🚀 How to Use These Materials

### For Reading
1. Open README.md for comprehensive theory
2. Use STUDY_GUIDE.md alongside for the type mapping table and diagrams
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
3. Practice the dynamic-JSON type-assertion dance until it's automatic

---

## 🆘 Common Questions

**Q: Do I need a third-party JSON library?**
A: No. `encoding/json` in the standard library handles the vast majority of real-world JSON work. Only reach for something else once you've measured a specific, real performance problem.

**Q: Why did a field just disappear from my JSON output?**
A: Either it's unexported (lowercase) and invisible to encoding/json, or it has `omitempty` and happens to hold its zero value. Both are covered in this level's Common Mistakes.

**Q: Why is my number showing up as a float when I expected an int?**
A: When you unmarshal into `interface{}` (e.g. `map[string]interface{}`), all JSON numbers decode as `float64` - there's no way for encoding/json to know you wanted an int. Convert explicitly, or unmarshal into a struct with a typed `int` field instead.

**Q: When should I use map[string]interface{} vs a struct?**
A: Use a struct whenever you know the shape - it's safer and clearer. Reach for `map[string]interface{}` only when the shape is genuinely unknown or varies at runtime.

**Q: What's the difference between json.Unmarshal and json.Decoder?**
A: `Unmarshal` works on a complete `[]byte` already in memory. `Decoder` reads from an `io.Reader` (a file, a network connection, a buffer) and can decode a stream of multiple JSON values one at a time - more efficient for large or incremental data.

---

## 🎯 Before Moving to Level 20

Make sure you can answer these questions:

- [ ] What happens to an unexported field when you marshal a struct?
- [ ] What are the four things a `json` struct tag can do?
- [ ] What Go type does a JSON number become inside `interface{}`?
- [ ] When would you reach for json.RawMessage instead of a plain struct?
- [ ] What's the difference between json.Marshal/Unmarshal and json.Encoder/Decoder?
- [ ] What must be true about a receiver for UnmarshalJSON to work correctly?

---

## 📚 Recommended Reading Order

For Maximum Learning:

1. **START_HERE.txt** (10 min)
   - Gets you oriented

2. **README.md** Sections 1-4 (35 min)
   - What JSON is, Marshal, struct tags

3. **README.md** Sections 5-8 (35 min)
   - Nesting, pretty-printing, dynamic JSON, RawMessage

4. **EXERCISES.md** Exercises 1-6 (3 hours)
   - Get hands-on with the core mechanics and dynamic JSON

5. **README.md** Sections 9-13 (35 min)
   - Streaming, custom marshaling, errors, best practices, mistakes

6. **STUDY_GUIDE.md** (35 min)
   - Study the type mapping table and navigation diagrams

7. **EXERCISES.md** Exercises 7-10 + Bonus (2.5+ hours)
   - Deep practice with RawMessage, streaming, custom types, and a real config workflow

8. **QUICK_REFERENCE.md** (10 min)
   - Create your mental cheat sheet

---

## 🏆 Success Indicators

You've mastered Level 19 when:

- ✅ You reach for struct tags automatically to shape your JSON output
- ✅ You can explain the exported-field rule without hesitation
- ✅ You never assert a JSON number as `int` when it came through `interface{}`
- ✅ You can decide confidently between a plain struct, RawMessage, and map[string]interface{}
- ✅ You've completed 8+ exercises
- ✅ You can explain JSON marshaling to someone else

---

## 🚀 What's Next?

After Level 19, you're ready for:

**Level 20: Goroutines**
- Concurrent execution with the go keyword
- How goroutines differ from OS threads
- The foundation for channels, select, and sync (Levels 21-23)

---

## 💬 Key Takeaway

> **json.Marshal and json.Unmarshal handle 90% of real work; struct tags are your control surface, and map[string]interface{} plus RawMessage are your escape hatches for the rest.**

---

## 📖 Quick Links

- [START_HERE.txt](./START_HERE.txt) - Welcome & Overview
- [README.md](./README.md) - Full Theory
- [EXERCISES.md](./EXERCISES.md) - Hands-On Exercises
- [STUDY_GUIDE.md](./STUDY_GUIDE.md) - Visual Patterns
- [QUICK_REFERENCE.md](./QUICK_REFERENCE.md) - Cheat Sheet

---

Happy learning! Level 19 gives your programs a common language to speak with the rest of the world! 🎉

*Estimated time to complete Level 19: 5-6 hours*
*Difficulty: ⭐⭐⭐ (Intermediate-Advanced)*
*Next Level: Level 20 - Goroutines*
