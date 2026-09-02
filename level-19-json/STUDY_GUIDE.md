# Level 19: Study Guide & Visual Reference

## 📚 Learning Path

### Week 1: Marshal, Unmarshal, and Tags
```
Day 1:  What JSON is, json.Marshal basics
Day 2:  Exported vs unexported fields (bridge to Level 2)
Day 3:  Struct tags - rename, exclude, omitempty
Day 4:  json.Unmarshal, missing/extra field behavior
Day 5:  Nested structs and slices of structs
Day 6:  MarshalIndent pretty-printing
Day 7:  Review and consolidate
```

### Week 2: Dynamic JSON, Streaming, and Custom Types
```
Day 1:  map[string]interface{} and type assertions (bridge to Level 15)
Day 2:  json.RawMessage for deferred/partial parsing
Day 3:  json.Decoder / json.Encoder streaming
Day 4:  Custom MarshalJSON/UnmarshalJSON
Day 5:  Real error messages: SyntaxError, UnmarshalTypeError
Day 6:  Exercises 1-6
Day 7:  Exercises 7-10 + bonus challenges
```

---

## 🎯 Struct Tag Syntax Reference

```
`json:"<name>,<option1>,<option2>,..."`
        │       │
        │       └─ options: omitempty (only one that matters in practice)
        └───────── the JSON key to use instead of the Go field name
```

| Tag | Meaning |
|-----|---------|
| *(no tag)* | JSON key = Go field name, exactly as written |
| `json:"name"` | JSON key = `"name"` |
| `json:"-"` | Field is never marshaled or unmarshaled - completely excluded |
| `json:"-,"` | Special case: JSON key is literally `"-"` (rare) |
| `json:",omitempty"` | Keep the Go field name, omit if zero value |
| `json:"name,omitempty"` | Rename to `"name"`, omit if zero value |

```
     ┌── field name in JSON output
     │        ┌── comma separator
     │        │        ┌── option
     ▼        ▼        ▼
`json:"full_name,omitempty"`
```

**Reminder:** only **exported** (capitalized) struct fields are visible to `encoding/json` at all - a tag on an unexported field does nothing (Level 2's rule, still in force).

---

## 🗺️ Go Type ↔ JSON Type Mapping

| Go type | JSON type (Marshal) | JSON → Go (Unmarshal into concrete struct) | JSON → Go (Unmarshal into interface{}) |
|---------|----------------------|---------------------------------------------|------------------------------------------|
| `string` | string | string | `string` |
| `int`, `float64`, etc. | number | number | **always `float64`** |
| `bool` | boolean | boolean | `bool` |
| `[]T` | array | array | `[]interface{}` |
| `nil` slice | `null` | - | `nil` |
| `map[string]T` | object | object | `map[string]interface{}` |
| `nil` map/pointer | `null` | - | `nil` |
| `struct` | object | object | `map[string]interface{}` |
| `*T` (non-nil) | same as `T` | same as `T` | (pointer transparent) |
| `time.Time` | RFC 3339 string | RFC 3339 string | `string` |
| `null` (JSON) | - | zero value / no change | `nil` |

Verified round-trip of a struct with a mix of these types, showing exactly which zero values become `null` vs their normal zero form:

```
Populated: {"I":42,"F":3.14,"S":"hi","B":true,"Slice":[1,2,3],"NilSlice":null,"Map":{"a":1},"Ptr":5,"NilPtr":null,"T":"2026-09-01T12:00:00Z"}
Zero value: {"I":0,"F":0,"S":"","B":false,"Slice":null,"NilSlice":null,"Map":null,"Ptr":null,"NilPtr":null,"T":"0001-01-01T00:00:00Z"}
```

Note: a `nil` slice and an **empty, non-nil** slice (`[]int{}`) both marshal differently - `nil` becomes JSON `null`, but `[]int{}` becomes `[]`. Keep that distinction in mind if a consumer treats `null` and `[]` differently.

---

## 🔄 Marshal / Unmarshal Round-Trip Diagram

```
                     json.Marshal(v)
   Go value  ───────────────────────────►  []byte (JSON)
  (struct,                                  (compact text)
   slice,       ◄───────────────────────────
   map...)         json.Unmarshal(data, &v)


   Go value  ──────────────────────────►  []byte (JSON, indented)
                json.MarshalIndent(v, "", "  ")


  io.Reader  ──────────────────────────►  Go value
  (file,           dec := json.NewDecoder(r)
   network,        dec.Decode(&v)
   buffer)

   Go value  ──────────────────────────►  io.Writer
                enc := json.NewEncoder(w)     (file, network, buffer)
                enc.Encode(v)
```

**Key facts to remember:**
- `Marshal` needs a **value**; `Unmarshal` needs a **pointer** (it has to write into your variable)
- Both return an `error` as the second (or only, for Marshal, alongside the bytes) return value - always check it
- `Decoder`/`Encoder` do the same job as `Unmarshal`/`Marshal`, but stream through an `io.Reader`/`io.Writer` instead of holding one full `[]byte`

---

## 🔒 Exported vs Unexported: Visibility Reminder

```
type Account struct {
    Username string   // exported (capital U) -> visible to encoding/json
    password string   // unexported (lowercase p) -> INVISIBLE to encoding/json
    Balance  float64  // exported -> visible
}
```

Verified marshal of `Account{Username: "alice99", password: "supersecret", Balance: 42.5}`:

```
{"Username":"alice99","Balance":42.5}
```

`password` isn't empty in the output - it's **not there at all**, and there's no error. This is Level 2's exported/unexported rule resurfacing: reflection (what `encoding/json` uses under the hood) can only see exported identifiers from outside a type's own package.

```
                encoding/json uses reflection
                          │
                          ▼
        ┌─────────────────────────────────┐
        │  Can reflection see this field? │
        └─────────────────────────────────┘
             │                      │
         Exported                Unexported
     (Capitalized)              (lowercase)
             │                      │
             ▼                      ▼
     Included in JSON        SILENTLY SKIPPED
     (unless json:"-")       (no error, ever)
```

---

## 🧭 Dynamic JSON Navigation Diagram

```
JSON text
    │
    │ json.Unmarshal(data, &v)   where v is map[string]interface{}
    ▼
map[string]interface{}
    │
    │  v["key"]  →  interface{}   (concrete type hidden inside)
    ▼
interface{}
    │
    │  type assertion:  x, ok := v["key"].(ConcreteType)
    ▼
   ┌────────────┬──────────────────────────────┐
   │  ok=true   │  x is the concrete value,     │
   │            │  safe to use                  │
   ├────────────┼──────────────────────────────┤
   │  ok=false  │  x is the zero value of       │
   │            │  ConcreteType - assertion     │
   │            │  failed, but NO panic         │
   └────────────┴──────────────────────────────┘

Nesting repeats the same dance one level deeper:

  v["metadata"].(map[string]interface{})["weight_kg"].(float64)
  └──────────┬──────────────────────────┘└────────┬────────┘
      assert to unlock the nested map        assert the leaf value
```

This is Level 15's type-assertion mechanism, applied specifically to the `interface{}` values that `encoding/json` produces for JSON of unknown shape. The **single-value form** (`x := v["key"].(T)`) panics on a bad guess - always use the **two-value form** (`x, ok := ...`) for JSON you don't fully control.

---

## 🚨 Common Mistakes At a Glance

| Mistake | Wrong | Right |
|---------|-------|-------|
| Unexported field | `type T struct { name string }` never appears in JSON | Capitalize: `Name string` (+ tag if renaming) |
| omitempty on a meaningful zero | `Balance float64 \`json:"balance,omitempty"\`` hides a real $0 balance | Drop `omitempty` when zero is meaningful |
| Ignoring the error | `data, _ := json.Marshal(v)` | `data, err := json.Marshal(v); if err != nil { ... }` |
| Assuming map key order | Relying on a particular key order in marshaled JSON | Treat order as an implementation detail; use a slice or struct if order matters |
| Wrong assertion type for numbers | `v["count"].(int)` on interface{} data → **panics** | `v["count"].(float64)` - JSON numbers are always float64 in interface{} |
| Passing a value, not a pointer, to Unmarshal | `json.Unmarshal(data, c)` → compile error | `json.Unmarshal(data, &c)` |

---

## 📈 Progression Summary

### Understanding Level 19

Level 19 teaches how Go talks to the outside world in its most common data format:

1. **Marshal/Unmarshal** - the two core conversions, built on reflection
2. **Struct tags** - the control surface for JSON shape (rename, exclude, omitempty)
3. **Nesting** - structs and slices of structs compose automatically
4. **Dynamic JSON** - `map[string]interface{}` plus Level 15's type assertions for unknown shapes
5. **RawMessage** - deferred, two-stage parsing for polymorphic data
6. **Streaming** - `Decoder`/`Encoder` for `io.Reader`/`io.Writer`, connecting to Level 18's file I/O
7. **Custom marshaling** - `MarshalJSON`/`UnmarshalJSON` for full control over wire format

### Prerequisites for Level 20

Before moving to Level 20 (Goroutines), you need:

- ✅ Comfortable marshaling and unmarshaling structs, including nested ones
- ✅ Know all four struct tag behaviors (default, rename, exclude, omitempty) cold
- ✅ Can navigate `map[string]interface{}` with type assertions safely
- ✅ Understand when to reach for `json.RawMessage` vs a plain struct
- ✅ Comfortable with `json.Decoder`/`json.Encoder` for streaming
- ✅ Can implement `MarshalJSON`/`UnmarshalJSON` on a custom type
- ✅ Recognize `*json.SyntaxError` and `*json.UnmarshalTypeError` on sight

### Ready for Level 20?

Level 20 shifts from data format to program structure - your first taste of concurrency:
- Launching concurrent work with the `go` keyword
- Why goroutines are cheap compared to OS threads
- The foundation everything in Levels 21-23 (channels, select, sync) builds on

---

## ✅ Checklist Before Level 20

- [ ] Can marshal a struct and explain why unexported fields are skipped
- [ ] Can use all four struct tag forms and predict their JSON output
- [ ] Can unmarshal JSON into a struct and explain missing/extra field behavior
- [ ] Can model and round-trip a nested struct with a slice of sub-structs
- [ ] Can pretty-print JSON with MarshalIndent
- [ ] Can extract values from map[string]interface{} with safe type assertions
- [ ] Can use json.RawMessage to defer parsing based on a discriminator field
- [ ] Can stream JSON with json.Decoder/json.Encoder
- [ ] Can implement MarshalJSON/UnmarshalJSON on a custom type
- [ ] Completed 8+ exercises

---

## 💡 Key Takeaways

### The Core Pair
`json.Marshal` (Go → JSON bytes) and `json.Unmarshal` (JSON bytes → Go, needs a pointer) are 90% of what you'll use.

### The Control Surface
Struct tags are how you decouple your Go field names from the JSON your API or config file actually needs.

### The Escape Hatches
`map[string]interface{}` for truly unknown shapes, `json.RawMessage` for shapes you'll decide about later, and `MarshalJSON`/`UnmarshalJSON` for full manual control.

### The Errors
JSON almost always comes from outside your program - a file, a network call, a user. Check every error; `*json.SyntaxError` and `*json.UnmarshalTypeError` tell you exactly what went wrong.

---

## 📚 Next Level

Level 20: Goroutines
- Concurrent execution with the go keyword
- How goroutines differ from OS threads
- The foundation for channels, select, and sync

You've got JSON down! Keep going! 🚀
