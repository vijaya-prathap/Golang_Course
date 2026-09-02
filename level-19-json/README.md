# Level 19: JSON - Complete Guide

## Introduction

Welcome to Level 19! You've mastered file handling (Level 18) - reading and writing bytes and text with `os` and `bufio`. Now it's time to give those bytes **structure**: JSON (JavaScript Object Notation), the format almost every API, config file, and log line on the planet speaks.

Go's standard library ships a complete, fast, well-designed JSON implementation in `encoding/json`. There is no need to reach for a third-party library for the vast majority of real work - `encoding/json` is the idiomatic, standard way to convert Go values to JSON (**marshaling**) and JSON back into Go values (**unmarshaling**). This level uses only the standard library.

This level also collects a lot of what you've already learned: exported vs unexported fields (Level 2), struct tags previewed in Level 13, type assertions from Level 15, and error types from Level 16 - all put to work for a single, very practical purpose.

---

## Table of Contents

1. [What Is JSON, and Why encoding/json?](#what-is-json-and-why-encodingjson)
2. [json.Marshal: Struct to JSON Bytes](#jsonmarshal-struct-to-json-bytes)
3. [Struct Tags for JSON](#struct-tags-for-json)
4. [json.Unmarshal: JSON Bytes to Struct](#jsonunmarshal-json-bytes-to-struct)
5. [Nested Structs and Slices of Structs](#nested-structs-and-slices-of-structs)
6. [Pretty-Printing With json.MarshalIndent](#pretty-printing-with-jsonmarshalindent)
7. [Dynamic JSON: map[string]interface{}](#dynamic-json-mapstringinterface)
8. [json.RawMessage: Deferred Parsing](#jsonrawmessage-deferred-parsing)
9. [Streaming JSON: json.Decoder and json.Encoder](#streaming-json-jsondecoder-and-jsonencoder)
10. [Custom Marshaling: MarshalJSON and UnmarshalJSON](#custom-marshaling-marshaljson-and-unmarshaljson)
11. [Common Error Cases](#common-error-cases)
12. [Best Practices](#best-practices)
13. [Common Mistakes](#common-mistakes)

---

## What Is JSON, and Why encoding/json?

JSON is a lightweight, text-based data format built from just a handful of shapes:

```
{"name": "Alice", "age": 30}     - object (like a Go map or struct)
["a", "b", "c"]                  - array (like a Go slice)
"hello"                          - string
42, 3.14                         - number (JSON has ONE number type)
true, false                      - boolean
null                             - absence of a value
```

Every mainstream language can read and write it, which makes it the default vocabulary for:

- **APIs** - REST services trade JSON requests and responses (Level 27 uses this heavily)
- **Config files** - human-readable settings, easy to hand-edit
- **Logs** - one JSON object per line is easy for machines to parse
- **Data interchange** - saving/loading structured data between programs, even between different languages

Go's standard library package for this is `encoding/json`. It converts between:

- Go values (structs, maps, slices, basic types) and JSON bytes - called **marshaling**
- JSON bytes and Go values - called **unmarshaling**

```go
import "encoding/json"
```

No external dependencies, no code generation required for the common case, and it is fast enough for the overwhelming majority of real programs. Reach for a third-party JSON library only once you've measured a real performance problem `encoding/json` can't solve - which is rare.

---

## json.Marshal: Struct to JSON Bytes

`json.Marshal` takes any Go value and returns its JSON representation as `[]byte`, plus an error.

```go
type Person struct {
    Name string
    Age  int
    City string
}

p := Person{Name: "Alice", Age: 30, City: "Lisbon"}
data, err := json.Marshal(p)
if err != nil {
    fmt.Println("error:", err)
    return
}
fmt.Println(string(data))
```

Verified output:

```
{"Name":"Alice","Age":30,"City":"Lisbon"}
```

`data` really is `[]byte` (verified with `%T`): `[]uint8` - `byte` is just an alias for `uint8`, as you'd expect from Level 3.

### Rule: Only Exported Fields Are Marshaled

This is the bridge to **Level 2's exported/unexported rule**. `encoding/json` uses reflection to inspect a struct's fields, and reflection can only see **exported** (capitalized) fields from another package - so `encoding/json` silently skips unexported ones. Verified:

```go
type Account struct {
    Username string
    password string // unexported - lowercase
    Balance  float64
}

a := Account{Username: "alice99", password: "supersecret", Balance: 42.5}
data, _ := json.Marshal(a)
fmt.Println(string(data))
```

```
{"Username":"alice99","Balance":42.5}
```

Notice `password` isn't just empty - it's **completely absent**, and no error is raised. This is a common source of "why is my field missing?!" bugs (see [Common Mistakes](#common-mistakes)).

### A Note on HTML Escaping

`json.Marshal` escapes certain characters (`<`, `>`, `&`) as Unicode escapes by default, to make JSON safe to embed directly in HTML `<script>` tags. Verified:

```go
b := Book{Title: "...", Author: "Donovan & Kernighan"}
data, _ := json.Marshal(b)
fmt.Println(string(data))
```

```
{"Title":"The Go Programming Language","Author":"Donovan & Kernighan","Pages":380,"InStock":true}
```

`&` became `&`. This is still valid, correct JSON (any JSON parser reads `&` back as `&`) - but if you're comparing raw output strings in tests, this can be surprising. To turn it off, use a `json.Encoder` with `SetEscapeHTML(false)` (see [Streaming JSON](#streaming-json-jsondecoder-and-jsonencoder)).

---

## Struct Tags for JSON

Level 13 previewed struct tags as backtick metadata read by libraries via reflection. `encoding/json` is the single most common consumer of struct tags in Go. Four things you can do with the `json` tag:

| Tag | Effect |
|-----|--------|
| `json:"name"` | Use `"name"` as the JSON key instead of the Go field name |
| `json:"-"` | Never include this field in JSON output or input |
| `json:"name,omitempty"` | Use `"name"`, but omit the field entirely if it holds its **zero value** |
| `json:",omitempty"` | Keep the Go field name, but still omit when zero-valued |

### Verified Before/After

```go
// Before: no tags - Go field names are used as-is
type UserBefore struct {
    FullName string
    Email    string
    Password string
    Nickname string
}
```

```
{"FullName":"Alice Smith","Email":"alice@example.com","Password":"secret123","Nickname":""}
```

```go
// After: tags rename, exclude, and omitempty
type UserAfter struct {
    FullName string `json:"full_name"`
    Email    string `json:"email"`
    Password string `json:"-"`
    Nickname string `json:"nickname,omitempty"`
}
```

With an empty `Nickname`:

```
{"full_name":"Alice Smith","email":"alice@example.com"}
```

With a non-empty `Nickname`:

```
{"full_name":"Bob Jones","email":"bob@example.com","nickname":"Bobby"}
```

Three things happened at once: `FullName` → `full_name` (renamed), `Password` vanished entirely (excluded), and `Nickname` only appears when it has a value (omitempty).

### omitempty and the Zero Value

`omitempty` doesn't mean "empty string" - it means **the field's zero value** (Level 3's zero values, all over again): `0` for numbers, `""` for strings, `false` for bools, `nil` for slices/maps/pointers. Verified:

```go
type Flags struct {
    Count   int    `json:"count,omitempty"`
    Label   string `json:"label,omitempty"`
    Enabled bool   `json:"enabled,omitempty"`
}

zero := Flags{Count: 0, Label: "", Enabled: false}
// {}

nonZero := Flags{Count: 5, Label: "x", Enabled: true}
// {"count":5,"label":"x","enabled":true}
```

A meaningful `false` or `0` (e.g., "explicitly disabled" or "zero balance") is **indistinguishable** from "not set" once `omitempty` is applied - see [Common Mistakes](#common-mistakes).

---

## json.Unmarshal: JSON Bytes to Struct

`json.Unmarshal` is the reverse of `Marshal`: it takes JSON bytes and a **pointer** to a Go value to fill in.

```go
var b Book
err := json.Unmarshal([]byte(data), &b)
```

The pointer matters - exactly like every mutating function since Level 12: `Unmarshal` needs to write into your variable, so it needs its address.

### Verified Round Trip

```go
original := Book{Title: "Go in Action", Author: "William Kennedy", Pages: 300, InStock: true}

marshaled, _ := json.Marshal(original)
var roundTripped Book
json.Unmarshal(marshaled, &roundTripped)

fmt.Println(original == roundTripped)
```

```
true
```

### Missing Fields: Zero Value

If the JSON is missing a field the struct has, that field keeps its **zero value** - no error. Verified:

```go
type Config struct {
    Host string
    Port int
}

missing := `{"Host": "localhost"}`
var c1 Config
json.Unmarshal([]byte(missing), &c1)
fmt.Printf("%+v\n", c1)
```

```
missing field result: {Host:localhost Port:0}
```

### Extra Fields: Silently Ignored

If the JSON has a field the struct doesn't declare, `Unmarshal` **silently ignores it** by default. Verified:

```go
extra := `{"Host": "localhost", "Port": 8080, "Timeout": 30}`
var c2 Config
json.Unmarshal([]byte(extra), &c2)
fmt.Printf("%+v\n", c2)
```

```
extra field result: {Host:localhost Port:8080}
```

`Timeout` simply disappears - no error, no warning. If you need to reject unknown fields (e.g., to catch typos in a config file), use a `json.Decoder` with `DisallowUnknownFields()` (covered in [Streaming JSON](#streaming-json-jsondecoder-and-jsonencoder)).

---

## Nested Structs and Slices of Structs

Real data isn't flat. `encoding/json` handles nested structs and slices of structs automatically, recursing all the way down.

```go
type LineItem struct {
    ProductName string  `json:"product_name"`
    Quantity    int     `json:"quantity"`
    UnitPrice   float64 `json:"unit_price"`
}

type ShippingAddress struct {
    Street string `json:"street"`
    City   string `json:"city"`
    Zip    string `json:"zip"`
}

type Order struct {
    OrderID  int             `json:"order_id"`
    Customer string          `json:"customer"`
    Items    []LineItem      `json:"items"`
    Shipping ShippingAddress `json:"shipping"`
}
```

Verified marshal of an `Order` with two `LineItem`s:

```
{"order_id":5001,"customer":"Alice Smith","items":[{"product_name":"Mechanical Keyboard","quantity":1,"unit_price":89.99},{"product_name":"USB-C Cable","quantity":3,"unit_price":7.5}],"shipping":{"street":"12 Rose Lane","city":"Lisbon","zip":"1100-048"}}
```

And unmarshaling it straight back:

```go
var decoded Order
json.Unmarshal(data, &decoded)
fmt.Println(decoded.Items[0].SKU) // reach through the slice, then the struct
```

Nested objects become nested structs; JSON arrays become Go slices of structs - one call handles the whole tree, no manual loop required.

---

## Pretty-Printing With json.MarshalIndent

`json.Marshal` produces compact, single-line JSON - great for wire transfer, hard to read. `json.MarshalIndent(v, prefix, indent)` adds line breaks and indentation.

```go
data, err := json.MarshalIndent(s, "", "    ")
```

Verified:

```
{
    "name": "api-01",
    "port": 8080,
    "tags": [
        "prod",
        "us-east"
    ],
    "healthy": true
}
```

The `prefix` argument is prepended to every line **after** the first (verified with `">> "` as prefix and `"  "` as indent - notice the closing brace also gets the prefix):

```
{
>>   "name": "api-01",
>>   "port": 8080,
...
>> }
```

In practice, `prefix` is almost always `""`; only `indent` (commonly `"  "` or `"    "`) is used. Use `MarshalIndent` for config files, CLI output, and anywhere a human will read the JSON; use `Marshal` for data that only machines will consume (APIs, storage).

---

## Dynamic JSON: map[string]interface{}

Sometimes you don't control the JSON's shape - a third-party API, a flexible config format, or JSON whose structure varies. Unmarshal it into `map[string]interface{}` instead of a fixed struct.

```go
var result map[string]interface{}
json.Unmarshal([]byte(data), &result)
```

Every JSON value type maps to a specific Go type inside the `interface{}`:

| JSON type | Go type inside `interface{}` |
|-----------|-------------------------------|
| object | `map[string]interface{}` |
| array | `[]interface{}` |
| string | `string` |
| number | `float64` (**always** - even for whole numbers like `42`) |
| boolean | `bool` |
| null | `nil` |

Verified with `%T`:

```
Go type of "id":         float64
Go type of "in_stock":   bool
Go type of "categories": []interface {}
Go type of "metadata":   map[string]interface {}
```

### The Type-Assertion Dance

This is the direct payoff of **Level 15's type assertions**: to get a usable value out of `interface{}`, assert its concrete type, using the safe two-value form so a wrong guess doesn't panic.

```go
id, ok := result["id"].(float64)       // 42, true
name, ok := result["name"].(string)    // "Widget", true
inStock, ok := result["in_stock"].(bool) // true, true

categories, ok := result["categories"].([]interface{})
for _, c := range categories {
    s, _ := c.(string) // assert each element individually
    fmt.Println(s)
}

metadata, ok := result["metadata"].(map[string]interface{})
weight, _ := metadata["weight_kg"].(float64)
```

Verified failure cases (both return the zero value and `ok == false`, never panic):

```go
badAge, ok := result["id"].(string)   // "" false  - id is really a float64
missing, ok := result["does_not_exist"] // <nil> false - key doesn't exist
```

**Rule:** JSON numbers are *always* `float64` when unmarshaled into `interface{}` - even `42` comes back as `float64(42)`, not `int(42)`. Convert explicitly (`int(id)`) if you need an integer.

---

## json.RawMessage: Deferred Parsing

`json.RawMessage` (defined as `type RawMessage []byte`) lets you parse a JSON document in two passes: decode the outer shape immediately, but keep a sub-object as raw, unparsed bytes to decode **later**, once you know what it actually is.

```go
type Notification struct {
    ID      int             `json:"id"`
    Kind    string          `json:"kind"`
    Details json.RawMessage `json:"details"` // parsed on demand
}
```

Stage 1 - parse the outer shape, leave `Details` as raw bytes:

```go
var notifications []Notification
json.Unmarshal([]byte(raw), &notifications)
fmt.Println(string(notifications[0].Details))
```

```
id=1 kind=email rawDetails={"to": "alice@example.com", "subject": "Welcome"}
id=2 kind=sms rawDetails={"phone_number": "+15551234567"}
```

Stage 2 - now that `Kind` is known, unmarshal `Details` into the right struct:

```go
switch n.Kind {
case "email":
    var e EmailDetails
    json.Unmarshal(n.Details, &e)
case "sms":
    var s SMSDetails
    json.Unmarshal(n.Details, &s)
}
```

Verified:

```
Email #1 -> {To:alice@example.com Subject:Welcome}
SMS #2 -> {PhoneNumber:+15551234567}
```

`RawMessage` also marshals cleanly back out - it implements `MarshalJSON` itself, re-emitting exactly the bytes it holds. This is the standard pattern for **polymorphic JSON**: a `kind`/`type` discriminator field plus a raw payload decoded based on that discriminator (see Bonus Challenge 3).

---

## Streaming JSON: json.Decoder and json.Encoder

`Marshal`/`Unmarshal` work on a complete `[]byte` in memory. For large data, or data arriving incrementally over an `io.Reader`/`io.Writer` (a file, a network connection, an in-memory buffer), `json.Decoder` and `json.Encoder` read and write JSON **as a stream**, without needing the whole thing in memory at once.

### Decoder: Reading Multiple Values

```go
reader := strings.NewReader(stream)
dec := json.NewDecoder(reader)

for dec.More() {
    var entry LogEntry
    if err := dec.Decode(&entry); err != nil {
        // handle error
    }
    // use entry
}
```

Verified against a stream of three back-to-back JSON objects with no separators:

```
[info] server started
[warn] disk usage high
[error] connection refused
```

`dec.More()` reports whether there's another JSON value left to decode - it's what makes reading a stream of unknown length possible.

### Encoder: Writing Multiple Values

```go
var buf bytes.Buffer
enc := json.NewEncoder(&buf)
for _, e := range entries {
    enc.Encode(e)
}
```

Verified output (each `Encode` call writes one JSON value followed by a newline - this is the basis of "JSON Lines" format, see Bonus Challenge 2):

```
{"level":"info","message":"server started"}
{"level":"warn","message":"disk usage high"}
{"level":"error","message":"connection refused"}
```

`SetIndent("", "  ")` pretty-prints each encoded value; `SetEscapeHTML(false)` turns off the `<`/`>`/`&` escaping mentioned earlier; `DisallowUnknownFields()` on a `Decoder` makes unknown JSON keys an **error** instead of silently ignoring them - useful for strict config parsing.

**When to reach for streaming:** reading/writing JSON to a file or network connection directly (Level 18's `os.File` satisfies both `io.Reader` and `io.Writer`), or a sequence of many JSON values where loading everything into one `[]byte` first would waste memory.

---

## Custom Marshaling: MarshalJSON and UnmarshalJSON

Sometimes the default field-by-field conversion isn't the representation you want. Implement two methods to take full control:

```go
type Marshaler interface {
    MarshalJSON() ([]byte, error)
}

type Unmarshaler interface {
    UnmarshalJSON([]byte) error
}
```

If a type has these methods, `encoding/json` calls them **instead of** its default behavior - anywhere that type appears, even nested inside other structs.

### Example: An Enum-like Type as a String

```go
type Priority int

const (
    PriorityLow Priority = iota
    PriorityMedium
    PriorityHigh
)

func (p Priority) String() string {
    switch p {
    case PriorityLow:
        return "low"
    case PriorityMedium:
        return "medium"
    case PriorityHigh:
        return "high"
    default:
        return "unknown"
    }
}

func (p Priority) MarshalJSON() ([]byte, error) {
    return json.Marshal(p.String())
}

func (p *Priority) UnmarshalJSON(data []byte) error {
    var s string
    if err := json.Unmarshal(data, &s); err != nil {
        return err
    }
    switch strings.ToLower(s) {
    case "low":
        *p = PriorityLow
    case "medium":
        *p = PriorityMedium
    case "high":
        *p = PriorityHigh
    default:
        return fmt.Errorf("invalid priority %q", s)
    }
    return nil
}
```

Note the receiver types: `MarshalJSON` uses a **value** receiver (Level 14) since it only reads; `UnmarshalJSON` must use a **pointer** receiver, since it needs to mutate the value in place - exactly the value-vs-pointer receiver reasoning from Level 14.

Verified marshal of a `Task{Priority: PriorityHigh}`:

```
{"title":"Ship the release","priority":"high","due":"2026-09-15"}
```

`priority` is the string `"high"`, not the underlying int `2` - custom marshaling ran instead of the default.

### Example: A Custom Date Format

`time.Time` already implements `MarshalJSON`/`UnmarshalJSON` (producing RFC 3339, e.g. `"2026-09-01T10:30:00Z"`). Wrap it in your own type to change the format:

```go
type DateOnly struct {
    time.Time
}

const dateLayout = "2006-01-02"

func (d DateOnly) MarshalJSON() ([]byte, error) {
    return json.Marshal(d.Time.Format(dateLayout))
}

func (d *DateOnly) UnmarshalJSON(data []byte) error {
    var s string
    if err := json.Unmarshal(data, &s); err != nil {
        return err
    }
    t, err := time.Parse(dateLayout, s)
    if err != nil {
        return err
    }
    d.Time = t
    return nil
}
```

Verified: a `DateOnly` wrapping `2026-09-15` marshals as `"2026-09-15"` (not the full RFC 3339 timestamp), and parses back correctly on `Unmarshal`. An invalid value correctly returns an error:

```
Error: invalid priority "urgent"
```

**Rule of thumb:** implement `MarshalJSON`/`UnmarshalJSON` when a type's *wire format* should differ from its *Go representation* - enums as strings, custom date formats, or values that need validation on the way in.

---

## Common Error Cases

Both `Marshal`/`Unmarshal` and `Decoder`/`Encoder` return `error` - Level 16's lesson that Go surfaces failure as an ordinary return value applies fully here. Two named error types matter most for JSON.

### json.SyntaxError: Malformed JSON

Returned when the input isn't valid JSON at all. Verified with a trailing comma (illegal in JSON):

```go
bad := `{"host": "localhost", "port": 8080,}`
var c Config
err := json.Unmarshal([]byte(bad), &c)
fmt.Printf("%v\n", err)
```

```
invalid character '}' looking for beginning of object key string
```

You can also assert the concrete type to inspect `Offset` - the byte position where parsing failed:

```go
if se, ok := err.(*json.SyntaxError); ok {
    fmt.Println(se.Offset) // 36
}
```

A couple of other verified variants worth recognizing on sight:

```
invalid character 'i' looking for beginning of object key string   // {invalid
unexpected end of JSON input                                       // "" (empty input)
```

### json.UnmarshalTypeError: Wrong Shape

Returned when the JSON is syntactically valid but doesn't match the target Go type - e.g., a JSON string where the struct expects a number. Verified:

```go
wrongType := `{"host": "localhost", "port": "not-a-number"}`
var c2 Config
err2 := json.Unmarshal([]byte(wrongType), &c2)
fmt.Printf("%v\n", err2)
```

```
json: cannot unmarshal string into Go struct field Config.port of type int
```

The concrete `*json.UnmarshalTypeError` carries structured detail (also verified):

```go
if ute, ok := err2.(*json.UnmarshalTypeError); ok {
    fmt.Println(ute.Field, ute.Value, ute.Type)
    // port string int
}
```

### Marshal Errors

`Marshal` fails too - most commonly on types it fundamentally cannot represent (channels, functions). Verified:

```go
_, err := json.Marshal(make(chan int))
fmt.Println(err)
```

```
json: unsupported type: chan int
```

**Always check these errors.** A JSON payload from a network call, a user-edited config file, or an external API is input you don't control - treat it the way Level 16 taught you to treat any fallible operation.

---

## Best Practices

### 1. Always Check the Error From Marshal and Unmarshal

```go
// ✅ Good
data, err := json.Marshal(v)
if err != nil {
    return fmt.Errorf("marshaling failed: %w", err)
}

// ❌ Risky - a failure is silently swallowed, data may be nil/garbage
data, _ := json.Marshal(v)
```

### 2. Use Struct Tags to Decouple Go Names From Wire Names

```go
// ✅ Good - idiomatic Go field names, idiomatic JSON keys
type User struct {
    FullName string `json:"full_name"`
}
```

Go convention is `PascalCase`; most JSON APIs use `snake_case` or `camelCase`. Tags let both sides stay idiomatic.

### 3. Use MarshalIndent Only for Human-Facing Output

```go
// ✅ Good - config files, CLI output, debug logs
data, _ := json.MarshalIndent(cfg, "", "  ")

// For APIs and storage, prefer compact Marshal - smaller, faster
data, _ := json.Marshal(cfg)
```

### 4. Prefer a Named Struct Over map[string]interface{} When the Shape Is Known

```go
// ✅ Good - compile-time safety, no type assertions needed
var order Order
json.Unmarshal(data, &order)

// ❌ Avoid when you already know the shape - error-prone, no safety
var order map[string]interface{}
json.Unmarshal(data, &order)
order["items"].([]interface{})[0].(map[string]interface{})["sku"].(string) // fragile
```

Reach for `map[string]interface{}` only when the shape is genuinely unknown or varies.

### 5. Use json.RawMessage for Polymorphic or Partially-Known JSON

Defer parsing a sub-object until a discriminator field (`kind`, `type`) tells you which struct to use - don't try to force every possible shape into one struct with lots of unused fields.

### 6. Use Decoder/Encoder for Files and Streams

```go
// ✅ Good - streams directly, no full-file []byte needed
f, _ := os.Open("data.json")
defer f.Close()
var v MyType
json.NewDecoder(f).Decode(&v)

// Less efficient for large files - reads everything into memory first
data, _ := os.ReadFile("data.json")
json.Unmarshal(data, &v)
```

### 7. Validate After Unmarshaling

`Unmarshal` only checks that the JSON *parses* and roughly matches your types - it does not enforce business rules (port ranges, required fields, non-negative counts). Validate explicitly afterward.

---

## Common Mistakes

### Mistake 1: Forgetting Fields Must Be Exported

```go
// ❌ WRONG - lowercase field is invisible to encoding/json
type User struct {
    name string // never appears in JSON output, never gets filled on input!
}

// ✅ RIGHT
type User struct {
    Name string `json:"name"` // exported, tag controls the JSON key
}
```

### Mistake 2: Misunderstanding omitempty

```go
// ❌ SURPRISE - a legitimately-zero value looks "unset" and disappears
type Account struct {
    Balance float64 `json:"balance,omitempty"`
}
a := Account{Balance: 0} // a real zero balance
// {} - Balance vanished! Can't tell "zero" from "not provided"

// ✅ RIGHT - drop omitempty when zero is a meaningful value
type Account struct {
    Balance float64 `json:"balance"`
}
```

### Mistake 3: Not Checking the Error From Marshal/Unmarshal

```go
// ❌ WRONG - errors silently ignored
data, _ := json.Marshal(v)
json.Unmarshal(input, &v)

// ✅ RIGHT
data, err := json.Marshal(v)
if err != nil { /* handle */ }
if err := json.Unmarshal(input, &v); err != nil { /* handle */ }
```

### Mistake 4: Assuming Map Key Order Is Preserved

```go
m := map[string]int{"z": 1, "a": 2, "m": 3}
data, _ := json.Marshal(m)
// encoding/json actually sorts map keys alphabetically when marshaling -
// but never rely on this; if order matters, use a slice of key/value pairs
// or a struct instead.
```

Go maps have randomized *iteration* order (Level 6/10), and while `encoding/json` happens to sort keys for deterministic output, treating map key order as meaningful is fragile - it's an implementation detail, not a documented guarantee for arbitrary consumers.

### Mistake 5: Expecting JSON Numbers to Unmarshal as int

```go
var v interface{}
json.Unmarshal([]byte(`{"count": 5}`), &v)
m := v.(map[string]interface{})
// ❌ WRONG - panics: interface conversion, interface {} is float64, not int
count := m["count"].(int)

// ✅ RIGHT - JSON numbers into interface{} are always float64
count := m["count"].(float64)
```

### Mistake 6: Forgetting Unmarshal Needs a Pointer

```go
var c Config
// ❌ WRONG - compile error: cannot use c (variable of type Config) as *Config
json.Unmarshal(data, c)

// ✅ RIGHT
json.Unmarshal(data, &c)
```

---

## Summary

**Marshal / Unmarshal:**
- `json.Marshal(v)` → `[]byte, error` - Go value to JSON
- `json.Unmarshal(data, &v)` → `error` - JSON to Go value (needs a pointer!)
- Only **exported** fields participate - unexported fields are silently skipped
- Missing JSON fields → zero value; extra JSON fields → silently ignored

**Struct Tags:**
- `json:"name"` renames a field
- `json:"-"` excludes a field entirely
- `json:"name,omitempty"` omits the field when it's the zero value

**Beyond the Basics:**
- Nested structs and slices of structs marshal/unmarshal recursively
- `json.MarshalIndent` pretty-prints for humans
- `map[string]interface{}` handles unknown/dynamic shapes (numbers come back as `float64`)
- `json.RawMessage` defers parsing a sub-object until you know its shape
- `json.Decoder`/`json.Encoder` stream JSON to/from an `io.Reader`/`io.Writer`
- `MarshalJSON`/`UnmarshalJSON` give a type full control over its own JSON format

**Errors:**
- `*json.SyntaxError` - malformed JSON text
- `*json.UnmarshalTypeError` - JSON value doesn't match the target Go type
- Always check the error - JSON usually comes from outside your program

---

## Next Steps

You now understand:
- ✅ Marshaling structs to JSON and unmarshaling JSON back to structs
- ✅ Controlling JSON output with struct tags (rename, exclude, omitempty)
- ✅ Nested structs, slices of structs, and pretty-printing
- ✅ Handling dynamic/unknown JSON with map[string]interface{} and type assertions
- ✅ Deferred parsing with json.RawMessage
- ✅ Streaming JSON with json.Decoder/json.Encoder
- ✅ Custom marshaling for enums, dates, and other special formats
- ✅ Recognizing real json.SyntaxError and json.UnmarshalTypeError messages

**Next level:** Level 20 - Goroutines
- Concurrent execution with the `go` keyword
- How goroutines differ from OS threads
- The foundation for everything in Go's concurrency story (channels, select, sync - Levels 21-23)

Your programs can now talk to the rest of the world in its native format. Next, they learn to do more than one thing at once! 🚀
