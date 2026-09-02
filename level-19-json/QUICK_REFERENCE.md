# Level 19: Quick Reference Card

## 🚀 Quick Start (5 Minutes)

```bash
# Setup
mkdir myapp && cd myapp
go mod init myapp

# Create main.go
cat > main.go << 'EOF'
package main

import (
    "encoding/json"
    "fmt"
)

type Person struct {
    Name string `json:"name"`
    Age  int    `json:"age,omitempty"`
}

func main() {
    p := Person{Name: "Alice", Age: 30}
    data, _ := json.Marshal(p)
    fmt.Println(string(data))

    var decoded Person
    json.Unmarshal(data, &decoded)
    fmt.Printf("%+v\n", decoded)
}
EOF

# Run
go run main.go
```

---

## 📋 Marshal / Unmarshal

```go
data, err := json.Marshal(v)          // Go value -> []byte
err := json.Unmarshal(data, &v)       // []byte -> Go value (NEEDS a pointer!)

data, err := json.MarshalIndent(v, "", "  ")  // pretty-printed []byte
```

---

## 🏷️ Struct Tag Cheat Sheet

```go
type User struct {
    Name  string `json:"name"`            // renamed
    Email string `json:"-"`               // excluded entirely
    Age   int    `json:"age,omitempty"`   // renamed + omit if zero
    Bio   string `json:",omitempty"`      // same name, omit if zero
}
```

| Behavior | Tag |
|----------|-----|
| Rename | `json:"newname"` |
| Exclude | `json:"-"` |
| Omit if zero | `json:"name,omitempty"` |
| No tag | Go field name used as-is |

**Reminder:** only exported (capitalized) fields are ever (un)marshaled.

---

## 🗺️ interface{} Type Map (Dynamic JSON)

```go
var v map[string]interface{}
json.Unmarshal(data, &v)
```

| JSON | Go type inside interface{} |
|------|------------------------------|
| object | `map[string]interface{}` |
| array | `[]interface{}` |
| string | `string` |
| number | `float64` (always!) |
| bool | `bool` |
| null | `nil` |

```go
s, ok := v["key"].(string)             // safe two-value assertion
n, ok := v["count"].(float64)          // NOT int!
arr, ok := v["items"].([]interface{})
obj, ok := v["nested"].(map[string]interface{})
```

---

## 🧩 RawMessage: Deferred Parsing

```go
type Envelope struct {
    Kind string          `json:"kind"`
    Data json.RawMessage `json:"data"`
}

var e Envelope
json.Unmarshal(data, &e)

switch e.Kind {
case "a":
    var x TypeA
    json.Unmarshal(e.Data, &x)
case "b":
    var y TypeB
    json.Unmarshal(e.Data, &y)
}
```

---

## 🌊 Streaming: Decoder / Encoder

```go
// Reading multiple values from a stream
dec := json.NewDecoder(reader)
for dec.More() {
    var v MyType
    dec.Decode(&v)
}

// Writing multiple values to a stream
enc := json.NewEncoder(writer)
enc.Encode(v1)
enc.Encode(v2) // each call adds one value + newline (JSON Lines format)

enc.SetIndent("", "  ")     // pretty-print each encoded value
enc.SetEscapeHTML(false)    // stop escaping <, >, &
dec.DisallowUnknownFields() // error on unknown JSON keys
```

---

## 🎨 Custom Marshaling

```go
func (t MyType) MarshalJSON() ([]byte, error) {
    return json.Marshal(t.String()) // value receiver - read only
}

func (t *MyType) UnmarshalJSON(data []byte) error {
    var s string
    if err := json.Unmarshal(data, &s); err != nil {
        return err
    }
    // parse s into *t ...
    return nil // pointer receiver - must mutate in place
}
```

---

## ⚠️ Errors to Recognize

```go
*json.SyntaxError          // malformed JSON text
    // e.g. "invalid character '}' looking for beginning of object key string"

*json.UnmarshalTypeError   // JSON value doesn't match target Go type
    // e.g. "json: cannot unmarshal string into Go struct field Config.port of type int"
    // .Field, .Value, .Type carry structured detail
```

```go
if err := json.Unmarshal(data, &v); err != nil {
    var se *json.SyntaxError
    var ute *json.UnmarshalTypeError
    switch {
    case errors.As(err, &se):
        // malformed JSON
    case errors.As(err, &ute):
        // type mismatch: ute.Field, ute.Value, ute.Type
    default:
        // other error
    }
}
```

---

## ⚠️ Common Mistakes

| Mistake | Wrong | Right |
|---------|-------|-------|
| Unexported field | `type T struct{ name string }` never in JSON | Capitalize: `Name string` |
| omitempty surprise | hides a real `0`/`""`/`false` | drop `omitempty` when zero is meaningful |
| Ignored error | `data, _ := json.Marshal(v)` | check `err` every time |
| Wrong assertion | `v["n"].(int)` on interface{} data | `v["n"].(float64)` |
| Missing pointer | `json.Unmarshal(data, c)` | `json.Unmarshal(data, &c)` |
| Assuming map order | relying on marshaled map key order | use a slice/struct if order matters |

---

## 🎓 Before Next Level

Can you:
- [ ] Marshal and unmarshal a nested struct with a slice of sub-structs?
- [ ] Explain all four struct tag behaviors from memory?
- [ ] Navigate map[string]interface{} with safe type assertions?
- [ ] Explain why JSON numbers become float64 in interface{}?
- [ ] Use json.RawMessage to defer parsing based on a discriminator field?
- [ ] Stream JSON with json.Decoder/json.Encoder?
- [ ] Implement MarshalJSON/UnmarshalJSON on a custom type?

If YES → You're ready for Level 20!

---

## 📚 Next Level

Level 20: Goroutines
- Concurrent execution with the go keyword
- How goroutines differ from OS threads
- The foundation for channels, select, and sync

You've got JSON down! 💪
