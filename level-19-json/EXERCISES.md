# Level 19: JSON - Exercises

Complete all exercises in order. Each exercise builds on previous knowledge.

---

## Exercise 1: Basic Struct-to-JSON Marshal

**Objective:** Convert a struct (and a slice of structs) to JSON with json.Marshal

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level19-exercise1
cd ~/projects/level19-exercise1
go mod init level19.example/exercise1
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "encoding/json"
    "fmt"
)

type Book struct {
    Title    string
    Author   string
    Pages    int
    InStock  bool
}

func main() {
    fmt.Println("=== Marshal a Single Struct ===")
    b := Book{Title: "The Go Programming Language", Author: "Donovan & Kernighan", Pages: 380, InStock: true}
    data, err := json.Marshal(b)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println(string(data))

    fmt.Println("\n=== Marshal a Slice of Structs ===")
    books := []Book{
        {Title: "The Go Programming Language", Author: "Donovan & Kernighan", Pages: 380, InStock: true},
        {Title: "Learning Go", Author: "Jon Bodner", Pages: 375, InStock: false},
    }
    data2, err := json.Marshal(books)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println(string(data2))

    fmt.Println("\n=== Type of Marshal's Result ===")
    fmt.Printf("%T\n", data)
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Marshal a Single Struct ===
{"Title":"The Go Programming Language","Author":"Donovan & Kernighan","Pages":380,"InStock":true}

=== Marshal a Slice of Structs ===
[{"Title":"The Go Programming Language","Author":"Donovan & Kernighan","Pages":380,"InStock":true},{"Title":"Learning Go","Author":"Jon Bodner","Pages":375,"InStock":false}]

=== Type of Marshal's Result ===
[]uint8
```

Notice `&` became `&` - `json.Marshal` escapes `<`, `>`, and `&` by default so the output is safe to embed in HTML. It's still perfectly valid JSON.

**Learning Objectives:**
- ✅ Use json.Marshal to convert a struct to JSON bytes
- ✅ Marshal a slice of structs into a JSON array
- ✅ Confirm Marshal's return type is []byte ([]uint8)

---

## Exercise 2: JSON-to-Struct Unmarshal Round Trip

**Objective:** Convert JSON back into a struct, and verify a full round trip

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level19-exercise2
cd ~/projects/level19-exercise2
go mod init level19.example/exercise2
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "encoding/json"
    "fmt"
)

type Book struct {
    Title   string
    Author  string
    Pages   int
    InStock bool
}

func main() {
    fmt.Println("=== Unmarshal JSON Into a Struct ===")
    data := `{"Title":"Learning Go","Author":"Jon Bodner","Pages":375,"InStock":false}`

    var b Book
    if err := json.Unmarshal([]byte(data), &b); err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Printf("%+v\n", b)

    fmt.Println("\n=== Round Trip: Struct -> JSON -> Struct ===")
    original := Book{Title: "Go in Action", Author: "William Kennedy", Pages: 300, InStock: true}

    marshaled, err := json.Marshal(original)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println("Marshaled:", string(marshaled))

    var roundTripped Book
    if err := json.Unmarshal(marshaled, &roundTripped); err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Printf("Round-tripped: %+v\n", roundTripped)
    fmt.Println("Equal to original:", original == roundTripped)

    fmt.Println("\n=== Unmarshal Into a Slice ===")
    listData := `[{"Title":"A","Author":"X","Pages":100,"InStock":true},{"Title":"B","Author":"Y","Pages":200,"InStock":false}]`
    var books []Book
    if err := json.Unmarshal([]byte(listData), &books); err != nil {
        fmt.Println("Error:", err)
        return
    }
    for _, book := range books {
        fmt.Printf("%+v\n", book)
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
=== Unmarshal JSON Into a Struct ===
{Title:Learning Go Author:Jon Bodner Pages:375 InStock:false}

=== Round Trip: Struct -> JSON -> Struct ===
Marshaled: {"Title":"Go in Action","Author":"William Kennedy","Pages":300,"InStock":true}
Round-tripped: {Title:Go in Action Author:William Kennedy Pages:300 InStock:true}
Equal to original: true

=== Unmarshal Into a Slice ===
{Title:A Author:X Pages:100 InStock:true}
{Title:B Author:Y Pages:200 InStock:false}
```

**Learning Objectives:**
- ✅ Use json.Unmarshal with a pointer to fill in a struct
- ✅ Verify a full struct -> JSON -> struct round trip
- ✅ Unmarshal a JSON array directly into a Go slice

---

## Exercise 3: Struct Tags - Rename, Exclude, omitempty

**Objective:** Control JSON field names, exclude sensitive fields, and skip zero values

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level19-exercise3
cd ~/projects/level19-exercise3
go mod init level19.example/exercise3
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "encoding/json"
    "fmt"
)

// Before: no tags - Go field names are used as-is
type UserBefore struct {
    FullName string
    Email    string
    Password string
    Nickname string
}

// After: tags rename, exclude, and omitempty
type UserAfter struct {
    FullName string `json:"full_name"`
    Email    string `json:"email"`
    Password string `json:"-"`
    Nickname string `json:"nickname,omitempty"`
}

func main() {
    fmt.Println("=== Before: No Struct Tags ===")
    before := UserBefore{FullName: "Alice Smith", Email: "alice@example.com", Password: "secret123", Nickname: ""}
    dataBefore, _ := json.Marshal(before)
    fmt.Println(string(dataBefore))

    fmt.Println("\n=== After: Renamed + Excluded + omitempty (empty Nickname) ===")
    after1 := UserAfter{FullName: "Alice Smith", Email: "alice@example.com", Password: "secret123", Nickname: ""}
    dataAfter1, _ := json.Marshal(after1)
    fmt.Println(string(dataAfter1))

    fmt.Println("\n=== After: omitempty With a Non-Empty Nickname ===")
    after2 := UserAfter{FullName: "Bob Jones", Email: "bob@example.com", Password: "hunter2", Nickname: "Bobby"}
    dataAfter2, _ := json.Marshal(after2)
    fmt.Println(string(dataAfter2))

    fmt.Println("\n=== omitempty With Various Zero Values ===")
    type Flags struct {
        Count   int    `json:"count,omitempty"`
        Label   string `json:"label,omitempty"`
        Enabled bool   `json:"enabled,omitempty"`
    }
    zero := Flags{Count: 0, Label: "", Enabled: false}
    dataZero, _ := json.Marshal(zero)
    fmt.Println("all zero values omitted:", string(dataZero))

    nonZero := Flags{Count: 5, Label: "x", Enabled: true}
    dataNonZero, _ := json.Marshal(nonZero)
    fmt.Println("all non-zero, all present:", string(dataNonZero))
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Before: No Struct Tags ===
{"FullName":"Alice Smith","Email":"alice@example.com","Password":"secret123","Nickname":""}

=== After: Renamed + Excluded + omitempty (empty Nickname) ===
{"full_name":"Alice Smith","email":"alice@example.com"}

=== After: omitempty With a Non-Empty Nickname ===
{"full_name":"Bob Jones","email":"bob@example.com","nickname":"Bobby"}

=== omitempty With Various Zero Values ===
all zero values omitted: {}
all non-zero, all present: {"count":5,"label":"x","enabled":true}
```

**Learning Objectives:**
- ✅ Rename JSON keys with json:"name"
- ✅ Exclude a field entirely with json:"-"
- ✅ Skip zero-valued fields with omitempty
- ✅ See the before/after difference on the exact same data

---

## Exercise 4: Nested Structs With a Slice of Sub-Structs

**Objective:** Model and (un)marshal a realistic multi-level JSON document

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level19-exercise4
cd ~/projects/level19-exercise4
go mod init level19.example/exercise4
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "encoding/json"
    "fmt"
)

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

func (o Order) Total() float64 {
    total := 0.0
    for _, item := range o.Items {
        total += float64(item.Quantity) * item.UnitPrice
    }
    return total
}

func main() {
    order := Order{
        OrderID:  5001,
        Customer: "Alice Smith",
        Items: []LineItem{
            {ProductName: "Mechanical Keyboard", Quantity: 1, UnitPrice: 89.99},
            {ProductName: "USB-C Cable", Quantity: 3, UnitPrice: 7.50},
        },
        Shipping: ShippingAddress{Street: "12 Rose Lane", City: "Lisbon", Zip: "1100-048"},
    }

    fmt.Println("=== Marshal Nested Struct ===")
    data, err := json.Marshal(order)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println(string(data))
    fmt.Printf("Computed total: %.2f\n", order.Total())

    fmt.Println("\n=== Unmarshal Back Into Nested Struct ===")
    var decoded Order
    if err := json.Unmarshal(data, &decoded); err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Printf("%+v\n", decoded)
    fmt.Println("Number of items:", len(decoded.Items))
    fmt.Println("Second item name:", decoded.Items[1].ProductName)
    fmt.Println("Shipping city:", decoded.Shipping.City)
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Marshal Nested Struct ===
{"order_id":5001,"customer":"Alice Smith","items":[{"product_name":"Mechanical Keyboard","quantity":1,"unit_price":89.99},{"product_name":"USB-C Cable","quantity":3,"unit_price":7.5}],"shipping":{"street":"12 Rose Lane","city":"Lisbon","zip":"1100-048"}}
Computed total: 112.49

=== Unmarshal Back Into Nested Struct ===
{OrderID:5001 Customer:Alice Smith Items:[{ProductName:Mechanical Keyboard Quantity:1 UnitPrice:89.99} {ProductName:USB-C Cable Quantity:3 UnitPrice:7.5}] Shipping:{Street:12 Rose Lane City:Lisbon Zip:1100-048}}
Number of items: 2
Second item name: USB-C Cable
Shipping city: Lisbon
```

**Learning Objectives:**
- ✅ Model a realistic multi-level JSON document with nested structs
- ✅ Marshal and unmarshal a slice of sub-structs (LineItem) inside a parent struct
- ✅ Reach through a decoded struct's slice and nested struct fields

---

## Exercise 5: Pretty-Printing With MarshalIndent

**Objective:** Produce human-readable, indented JSON

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level19-exercise5
cd ~/projects/level19-exercise5
go mod init level19.example/exercise5
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "encoding/json"
    "fmt"
)

type Server struct {
    Name    string   `json:"name"`
    Port    int      `json:"port"`
    Tags    []string `json:"tags"`
    Healthy bool     `json:"healthy"`
}

func main() {
    s := Server{Name: "api-01", Port: 8080, Tags: []string{"prod", "us-east"}, Healthy: true}

    fmt.Println("=== Compact (json.Marshal) ===")
    compact, err := json.Marshal(s)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println(string(compact))

    fmt.Println("\n=== Pretty (json.MarshalIndent, 4-space) ===")
    pretty, err := json.MarshalIndent(s, "", "    ")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println(string(pretty))

    fmt.Println("\n=== Pretty With a Prefix ===")
    prefixed, err := json.MarshalIndent(s, ">> ", "  ")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println(string(prefixed))

    fmt.Println("\n=== Pretty-Printing a Slice ===")
    servers := []Server{
        {Name: "api-01", Port: 8080, Tags: []string{"prod"}, Healthy: true},
        {Name: "api-02", Port: 8081, Tags: []string{"staging"}, Healthy: false},
    }
    prettyList, err := json.MarshalIndent(servers, "", "    ")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println(string(prettyList))
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Compact (json.Marshal) ===
{"name":"api-01","port":8080,"tags":["prod","us-east"],"healthy":true}

=== Pretty (json.MarshalIndent, 4-space) ===
{
    "name": "api-01",
    "port": 8080,
    "tags": [
        "prod",
        "us-east"
    ],
    "healthy": true
}

=== Pretty With a Prefix ===
{
>>   "name": "api-01",
>>   "port": 8080,
>>   "tags": [
>>     "prod",
>>     "us-east"
>>   ],
>>   "healthy": true
>> }

=== Pretty-Printing a Slice ===
[
    {
        "name": "api-01",
        "port": 8080,
        "tags": [
            "prod"
        ],
        "healthy": true
    },
    {
        "name": "api-02",
        "port": 8081,
        "tags": [
            "staging"
        ],
        "healthy": false
    }
]
```

**Learning Objectives:**
- ✅ Compare compact Marshal output to indented MarshalIndent output
- ✅ Understand the prefix argument applies to every line after the first
- ✅ Pretty-print a slice of structs

---

## Exercise 6: Dynamic JSON With map[string]interface{}

**Objective:** Parse JSON of unknown/dynamic shape and extract values with type assertions

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level19-exercise6
cd ~/projects/level19-exercise6
go mod init level19.example/exercise6
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "encoding/json"
    "fmt"
)

func main() {
    data := `{
        "id": 42,
        "name": "Widget",
        "price": 19.99,
        "in_stock": true,
        "categories": ["hardware", "tools"],
        "metadata": {"weight_kg": 1.5, "manufacturer": "Acme"}
    }`

    fmt.Println("=== Unmarshal Into map[string]interface{} ===")
    var result map[string]interface{}
    if err := json.Unmarshal([]byte(data), &result); err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Printf("Go type of \"id\": %T\n", result["id"])
    fmt.Printf("Go type of \"in_stock\": %T\n", result["in_stock"])
    fmt.Printf("Go type of \"categories\": %T\n", result["categories"])
    fmt.Printf("Go type of \"metadata\": %T\n", result["metadata"])

    fmt.Println("\n=== Extracting Scalar Values With Type Assertions ===")
    id, ok := result["id"].(float64)
    fmt.Printf("id = %v (ok=%v, as int: %d)\n", id, ok, int(id))

    name, ok := result["name"].(string)
    fmt.Printf("name = %q (ok=%v)\n", name, ok)

    inStock, ok := result["in_stock"].(bool)
    fmt.Printf("in_stock = %v (ok=%v)\n", inStock, ok)

    fmt.Println("\n=== Extracting a Nested Array ===")
    categories, ok := result["categories"].([]interface{})
    fmt.Println("categories ok:", ok)
    for i, c := range categories {
        s, _ := c.(string)
        fmt.Printf("  [%d] %s\n", i, s)
    }

    fmt.Println("\n=== Extracting a Nested Object ===")
    metadata, ok := result["metadata"].(map[string]interface{})
    fmt.Println("metadata ok:", ok)
    weight, _ := metadata["weight_kg"].(float64)
    manufacturer, _ := metadata["manufacturer"].(string)
    fmt.Printf("weight_kg = %v, manufacturer = %s\n", weight, manufacturer)

    fmt.Println("\n=== Failed Type Assertion (Safe Two-Value Form) ===")
    badAge, ok := result["id"].(string)
    fmt.Printf("id as string = %q, ok = %v\n", badAge, ok)

    fmt.Println("\n=== Missing Key ===")
    missing, ok := result["does_not_exist"]
    fmt.Printf("missing = %v, ok = %v\n", missing, ok)
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Unmarshal Into map[string]interface{} ===
Go type of "id": float64
Go type of "in_stock": bool
Go type of "categories": []interface {}
Go type of "metadata": map[string]interface {}

=== Extracting Scalar Values With Type Assertions ===
id = 42 (ok=true, as int: 42)
name = "Widget" (ok=true)
in_stock = true (ok=true)

=== Extracting a Nested Array ===
categories ok: true
  [0] hardware
  [1] tools

=== Extracting a Nested Object ===
metadata ok: true
weight_kg = 1.5, manufacturer = Acme

=== Failed Type Assertion (Safe Two-Value Form) ===
id as string = "", ok = false

=== Missing Key ===
missing = <nil>, ok = false
```

**Learning Objectives:**
- ✅ Unmarshal JSON of unknown shape into map[string]interface{}
- ✅ Confirm JSON numbers always decode as float64
- ✅ Use the safe two-value type assertion to extract scalars, arrays, and nested objects
- ✅ Handle a failed assertion and a missing key without panicking

---

## Exercise 7: json.RawMessage for Partial Parsing

**Objective:** Parse an outer JSON shape immediately, deferring a sub-object until its type is known

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level19-exercise7
cd ~/projects/level19-exercise7
go mod init level19.example/exercise7
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "encoding/json"
    "fmt"
)

type Notification struct {
    ID      int             `json:"id"`
    Kind    string          `json:"kind"`
    Details json.RawMessage `json:"details"`
}

type EmailDetails struct {
    To      string `json:"to"`
    Subject string `json:"subject"`
}

type SMSDetails struct {
    PhoneNumber string `json:"phone_number"`
}

func main() {
    raw := `[
        {"id": 1, "kind": "email", "details": {"to": "alice@example.com", "subject": "Welcome"}},
        {"id": 2, "kind": "sms", "details": {"phone_number": "+15551234567"}}
    ]`

    fmt.Println("=== Stage 1: Parse Outer Shape, Defer Details ===")
    var notifications []Notification
    if err := json.Unmarshal([]byte(raw), &notifications); err != nil {
        fmt.Println("Error:", err)
        return
    }
    for _, n := range notifications {
        fmt.Printf("id=%d kind=%s rawDetails=%s\n", n.ID, n.Kind, string(n.Details))
    }

    fmt.Println("\n=== Stage 2: Parse Details Based on Kind ===")
    for _, n := range notifications {
        switch n.Kind {
        case "email":
            var e EmailDetails
            if err := json.Unmarshal(n.Details, &e); err != nil {
                fmt.Println("Error:", err)
                continue
            }
            fmt.Printf("Email #%d -> %+v\n", n.ID, e)
        case "sms":
            var s SMSDetails
            if err := json.Unmarshal(n.Details, &s); err != nil {
                fmt.Println("Error:", err)
                continue
            }
            fmt.Printf("SMS #%d -> %+v\n", n.ID, s)
        }
    }

    fmt.Println("\n=== Marshaling a RawMessage Back Out ===")
    out, err := json.Marshal(notifications[0])
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println(string(out))
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Stage 1: Parse Outer Shape, Defer Details ===
id=1 kind=email rawDetails={"to": "alice@example.com", "subject": "Welcome"}
id=2 kind=sms rawDetails={"phone_number": "+15551234567"}

=== Stage 2: Parse Details Based on Kind ===
Email #1 -> {To:alice@example.com Subject:Welcome}
SMS #2 -> {PhoneNumber:+15551234567}

=== Marshaling a RawMessage Back Out ===
{"id":1,"kind":"email","details":{"to":"alice@example.com","subject":"Welcome"}}
```

**Learning Objectives:**
- ✅ Use json.RawMessage to defer parsing part of a JSON document
- ✅ Choose which struct to decode a raw sub-object into, based on a discriminator field
- ✅ Confirm a RawMessage marshals cleanly back to its original bytes

---

## Exercise 8: Streaming With json.Decoder and json.Encoder

**Objective:** Read and write a stream of JSON values without loading everything into one []byte

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level19-exercise8
cd ~/projects/level19-exercise8
go mod init level19.example/exercise8
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "strings"
)

type LogEntry struct {
    Level   string `json:"level"`
    Message string `json:"message"`
}

func main() {
    fmt.Println("=== Decoder: Reading a Stream of JSON Values ===")
    stream := `{"level":"info","message":"server started"}
{"level":"warn","message":"disk usage high"}
{"level":"error","message":"connection refused"}`

    reader := strings.NewReader(stream)
    dec := json.NewDecoder(reader)

    var entries []LogEntry
    for dec.More() {
        var entry LogEntry
        if err := dec.Decode(&entry); err != nil {
            fmt.Println("Error:", err)
            return
        }
        entries = append(entries, entry)
    }
    for _, e := range entries {
        fmt.Printf("[%s] %s\n", e.Level, e.Message)
    }

    fmt.Println("\n=== Encoder: Writing a Stream of JSON Values ===")
    var buf bytes.Buffer
    enc := json.NewEncoder(&buf)
    for _, e := range entries {
        if err := enc.Encode(e); err != nil {
            fmt.Println("Error:", err)
            return
        }
    }
    fmt.Print(buf.String())

    fmt.Println("=== Encoder With SetIndent ===")
    var prettyBuf bytes.Buffer
    prettyEnc := json.NewEncoder(&prettyBuf)
    prettyEnc.SetIndent("", "  ")
    if err := prettyEnc.Encode(entries[0]); err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Print(prettyBuf.String())
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Decoder: Reading a Stream of JSON Values ===
[info] server started
[warn] disk usage high
[error] connection refused

=== Encoder: Writing a Stream of JSON Values ===
{"level":"info","message":"server started"}
{"level":"warn","message":"disk usage high"}
{"level":"error","message":"connection refused"}
=== Encoder With SetIndent ===
{
  "level": "info",
  "message": "server started"
}
```

**Learning Objectives:**
- ✅ Use json.Decoder with dec.More() to read multiple JSON values from a single stream
- ✅ Use json.Encoder to write multiple JSON values to an io.Writer
- ✅ Use SetIndent to pretty-print encoder output

---

## Exercise 9: Custom MarshalJSON/UnmarshalJSON

**Objective:** Implement custom JSON formatting for an enum-like type and a custom date type

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level19-exercise9
cd ~/projects/level19-exercise9
go mod init level19.example/exercise9
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "encoding/json"
    "fmt"
    "strings"
    "time"
)

// Priority is an enum-like type that marshals as a lowercase string
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

// DateOnly wraps time.Time to marshal as YYYY-MM-DD instead of RFC3339
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

type Task struct {
    Title    string   `json:"title"`
    Priority Priority `json:"priority"`
    Due      DateOnly `json:"due"`
}

func main() {
    fmt.Println("=== Custom MarshalJSON: Enum-like Type ===")
    t := Task{
        Title:    "Ship the release",
        Priority: PriorityHigh,
        Due:      DateOnly{Time: time.Date(2026, 9, 15, 0, 0, 0, 0, time.UTC)},
    }
    data, err := json.Marshal(t)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println(string(data))

    fmt.Println("\n=== Custom UnmarshalJSON: Round Trip ===")
    var decoded Task
    if err := json.Unmarshal(data, &decoded); err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Printf("Title: %s\n", decoded.Title)
    fmt.Printf("Priority (Go value): %d\n", decoded.Priority)
    fmt.Printf("Priority (String()): %s\n", decoded.Priority)
    fmt.Printf("Due (formatted): %s\n", decoded.Due.Time.Format(dateLayout))

    fmt.Println("\n=== Invalid Priority String ===")
    bad := `{"title":"x","priority":"urgent","due":"2026-09-15"}`
    var badTask Task
    err = json.Unmarshal([]byte(bad), &badTask)
    fmt.Println("Error:", err)
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Custom MarshalJSON: Enum-like Type ===
{"title":"Ship the release","priority":"high","due":"2026-09-15"}

=== Custom UnmarshalJSON: Round Trip ===
Title: Ship the release
Priority (Go value): 2
Priority (String()): high
Due (formatted): 2026-09-15

=== Invalid Priority String ===
Error: invalid priority "urgent"
```

**Learning Objectives:**
- ✅ Implement MarshalJSON on a value receiver to control output format
- ✅ Implement UnmarshalJSON on a pointer receiver to control parsing and validation
- ✅ See custom marshaling apply automatically wherever the type appears in a struct
- ✅ Return a real error from UnmarshalJSON when the input is invalid

---

## Exercise 10: Comprehensive Practice — Config File Load/Save

**Objective:** Combine Level 18's file I/O with everything in this level: load a config file with defaults for missing fields, validate it, and save it back pretty-printed

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level19-exercise10
cd ~/projects/level19-exercise10
go mod init level19.example/exercise10
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "encoding/json"
    "errors"
    "fmt"
    "os"
)

type Config struct {
    Host    string `json:"host,omitempty"`
    Port    int    `json:"port,omitempty"`
    Debug   bool   `json:"debug,omitempty"`
    Workers int    `json:"workers,omitempty"`
}

func defaultConfig() Config {
    return Config{
        Host:    "localhost",
        Port:    8080,
        Debug:   false,
        Workers: 4,
    }
}

// applyDefaults fills any zero-valued fields with defaults after loading
func applyDefaults(c *Config) {
    defaults := defaultConfig()
    if c.Host == "" {
        c.Host = defaults.Host
    }
    if c.Port == 0 {
        c.Port = defaults.Port
    }
    if c.Workers == 0 {
        c.Workers = defaults.Workers
    }
    // Debug has no meaningful "unset" state (false is valid), so it's left as-is
}

func validate(c Config) error {
    if c.Port < 1 || c.Port > 65535 {
        return fmt.Errorf("invalid port %d: must be between 1 and 65535", c.Port)
    }
    if c.Workers < 1 {
        return fmt.Errorf("invalid workers %d: must be at least 1", c.Workers)
    }
    if c.Host == "" {
        return errors.New("host must not be empty")
    }
    return nil
}

func loadConfig(path string) (Config, error) {
    data, err := os.ReadFile(path)
    if errors.Is(err, os.ErrNotExist) {
        fmt.Println("Config file not found, using defaults")
        return defaultConfig(), nil
    }
    if err != nil {
        return Config{}, err
    }
    var c Config
    if err := json.Unmarshal(data, &c); err != nil {
        return Config{}, fmt.Errorf("parsing config: %w", err)
    }
    applyDefaults(&c)
    return c, nil
}

func saveConfig(path string, c Config) error {
    data, err := json.MarshalIndent(c, "", "  ")
    if err != nil {
        return err
    }
    return os.WriteFile(path, data, 0644)
}

func main() {
    path := "config.json"
    defer os.Remove(path)

    fmt.Println("=== Step 1: Load Config When File Doesn't Exist ===")
    cfg, err := loadConfig(path)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Printf("Loaded (defaults): %+v\n", cfg)

    if err := validate(cfg); err != nil {
        fmt.Println("Validation error:", err)
        return
    }
    fmt.Println("Validation passed")

    if err := saveConfig(path, cfg); err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println("Saved to", path)

    fmt.Println("\n=== Step 2: Simulate a Partial Config File ===")
    partial := `{"host": "api.example.com", "debug": true}`
    if err := os.WriteFile(path, []byte(partial), 0644); err != nil {
        fmt.Println("Error:", err)
        return
    }

    cfg2, err := loadConfig(path)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Printf("Loaded (partial + defaults filled in): %+v\n", cfg2)

    if err := validate(cfg2); err != nil {
        fmt.Println("Validation error:", err)
        return
    }
    fmt.Println("Validation passed")

    if err := saveConfig(path, cfg2); err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println("Saved to", path)

    finalContents, err := os.ReadFile(path)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println("\n=== Final config.json Contents ===")
    fmt.Println(string(finalContents))

    fmt.Println("=== Step 3: Invalid Config Fails Validation ===")
    invalid := `{"host": "api.example.com", "port": 99999, "workers": 4}`
    if err := os.WriteFile(path, []byte(invalid), 0644); err != nil {
        fmt.Println("Error:", err)
        return
    }
    cfg3, err := loadConfig(path)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    if err := validate(cfg3); err != nil {
        fmt.Println("Validation error:", err)
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
=== Step 1: Load Config When File Doesn't Exist ===
Config file not found, using defaults
Loaded (defaults): {Host:localhost Port:8080 Debug:false Workers:4}
Validation passed
Saved to config.json

=== Step 2: Simulate a Partial Config File ===
Loaded (partial + defaults filled in): {Host:api.example.com Port:8080 Debug:true Workers:4}
Validation passed
Saved to config.json

=== Final config.json Contents ===
{
  "host": "api.example.com",
  "port": 8080,
  "debug": true,
  "workers": 4
}
=== Step 3: Invalid Config Fails Validation ===
Validation error: invalid port 99999: must be between 1 and 65535
```

**Learning Objectives:**
- ✅ Combine os file I/O (Level 18) with json.Marshal/Unmarshal
- ✅ Apply defaults for fields omitted from a loaded config
- ✅ Validate decoded data before trusting it
- ✅ Save a config back to disk, pretty-printed, using os.WriteFile

---

## Bonus Challenges

### Challenge 1: JSON-Based Key-Value Store

Build a tiny persistent key-value store: `Set(key, value string)`, `Get(key string) (string, bool)`, `Save() error`, and `Load() error`, backed by a `map[string]string` marshaled to a JSON file on disk.

```bash
mkdir -p ~/projects/level19-bonus1
cd ~/projects/level19-bonus1
go mod init level19.example/bonus1
```

**Hints:**
- Store the map and the file path in a struct; `Save` calls `json.MarshalIndent` then `os.WriteFile`
- `Load` should treat a missing file (check with `os.IsNotExist` or `errors.Is(err, os.ErrNotExist)`) as "start with an empty store", not an error
- `os.WriteFile` after `os.ReadFile` round-trips cleanly - verify by reloading into a second instance

### Challenge 2: Slice of Structs to JSON Lines

Write a function that converts a `[]T` into "JSON Lines" format - one compact JSON object per line, no surrounding `[` `]` or commas - and a second function that reads JSON Lines back into a `[]T`.

```bash
mkdir -p ~/projects/level19-bonus2
cd ~/projects/level19-bonus2
go mod init level19.example/bonus2
```

**Hints:**
- `json.NewEncoder(w).Encode(v)` already appends a trailing newline after each value - call it once per item into the same `bytes.Buffer`
- Read it back with `json.NewDecoder(r)` and a `for dec.More()` loop, exactly like Exercise 8
- JSON Lines (`.jsonl`) is a real, common format for logs and bulk data exports

### Challenge 3: Polymorphic JSON Array (Two Shapes + a Discriminator)

Given a JSON array where each element is either `{"kind": "circle", "radius": ...}` or `{"kind": "rectangle", "width": ..., "height": ...}`, parse the whole array into a `[]interface{}` holding the correct concrete type (`Circle` or `Rectangle`) for each element.

```bash
mkdir -p ~/projects/level19-bonus3
cd ~/projects/level19-bonus3
go mod init level19.example/bonus3
```

**Hints:**
- Unmarshal the array into `[]json.RawMessage` first - each raw element is a shape you haven't decided how to parse yet
- Unmarshal each raw element into a small "probe" struct that has only the `Kind` field, to read the discriminator
- Switch on the probe's `Kind`, then unmarshal the *same* raw bytes again into the right concrete struct
- Return an error for an unrecognized `kind` instead of silently skipping it

---

## What You've Learned

After completing these 10 exercises, you can:

✅ Marshal structs and slices of structs to JSON, and unmarshal JSON back into them
✅ Use struct tags to rename, exclude, and conditionally omit fields
✅ Model and round-trip nested structs and slices of sub-structs
✅ Pretty-print JSON for humans with MarshalIndent
✅ Navigate dynamic/unknown JSON with map[string]interface{} and type assertions
✅ Defer parsing part of a document with json.RawMessage
✅ Stream JSON in and out with json.Decoder/json.Encoder
✅ Give a type full control over its own JSON format with MarshalJSON/UnmarshalJSON
✅ Recognize real json.SyntaxError and json.UnmarshalTypeError messages
✅ Combine file I/O and JSON into a realistic config load/validate/save workflow

---

## Next Level

Level 20: Goroutines
- Concurrent execution with the go keyword
- How goroutines differ from OS threads
- The foundation for channels, select, and sync (Levels 21-23)

Great work! Your programs can now read and write the world's most common data format! 🚀
