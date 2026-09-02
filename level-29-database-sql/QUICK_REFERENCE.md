# Level 29: Quick Reference Card

## 🚀 Quick Start (5 Minutes)

```bash
# Setup
mkdir myapp && cd myapp
go mod init myapp
go get modernc.org/sqlite@v1.57.0

# Create main.go
cat > main.go << 'EOF'
package main

import (
    "database/sql"
    "fmt"
    "log"

    _ "modernc.org/sqlite"
)

func main() {
    db, err := sql.Open("sqlite", "file::memory:?cache=shared")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()
    db.SetMaxOpenConns(1) // required for in-memory SQLite

    db.Exec(`CREATE TABLE books (id INTEGER PRIMARY KEY, title TEXT)`)
    db.Exec(`INSERT INTO books (title) VALUES (?)`, "Clean Code")

    var title string
    db.QueryRow(`SELECT title FROM books WHERE id = ?`, 1).Scan(&title)
    fmt.Println(title)
}
EOF

# Run
go run main.go
```

---

## 📋 Connecting

```go
import (
    "database/sql"

    _ "modernc.org/sqlite" // blank import registers the driver
)

db, err := sql.Open("sqlite", "file::memory:?cache=shared") // lazy - no real connection yet
defer db.Close()
db.SetMaxOpenConns(1) // in-memory SQLite needs exactly one connection

if err := db.Ping(); err != nil { /* verify connectivity eagerly */ }
```

---

## 🎯 Query vs. QueryRow vs. Exec

```go
// Many rows expected
rows, err := db.Query(`SELECT id, title FROM books`)
defer rows.Close()
for rows.Next() {
    rows.Scan(&id, &title)
}
if err := rows.Err(); err != nil { /* ... */ }

// Zero-or-one row expected
var title string
err := db.QueryRow(`SELECT title FROM books WHERE id = ?`, id).Scan(&title)
if errors.Is(err, sql.ErrNoRows) { /* not found - normal case */ }

// No rows returned (INSERT/UPDATE/DELETE/DDL)
result, err := db.Exec(`INSERT INTO books (title) VALUES (?)`, title)
id, _ := result.LastInsertId()
affected, _ := result.RowsAffected()
```

---

## 🔒 Parameterized Queries - Always

```go
// ✅ SAFE
db.Query(`SELECT * FROM books WHERE author = ?`, userInput)

// ❌ NEVER
db.Query(fmt.Sprintf(`SELECT * FROM books WHERE author = '%s'`, userInput))
```

---

## 🧰 Prepared Statements

```go
stmt, err := db.Prepare(`INSERT INTO books (title, year) VALUES (?, ?)`)
defer stmt.Close()

for _, b := range books {
    stmt.Exec(b.Title, b.Year)
}
```

---

## 🔀 Transactions - defer-rollback-if-not-committed

```go
tx, err := db.Begin()
committed := false
defer func() {
    if !committed {
        tx.Rollback()
    }
}()

if _, err := tx.Exec(`UPDATE ...`); err != nil {
    return err
}
if _, err := tx.Exec(`INSERT ...`); err != nil {
    return err
}

if err := tx.Commit(); err != nil {
    return err
}
committed = true
```

---

## ⏱️ Context-Aware Calls

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

db.QueryRowContext(ctx, `SELECT ...`, id).Scan(&x)
db.ExecContext(ctx, `INSERT ...`, args...)
db.QueryContext(ctx, `SELECT ...`)
```

---

## 🏊 Connection Pool Settings

```go
db.SetMaxOpenConns(25)              // max total connections
db.SetMaxIdleConns(25)              // max kept idle for reuse
db.SetConnMaxLifetime(5 * time.Minute)
db.SetConnMaxIdleTime(1 * time.Minute)

stats := db.Stats() // OpenConnections, InUse, Idle, MaxOpenConnections, ...
```

---

## ⚠️ Common Mistakes

| Mistake | Wrong | Right |
|---------|-------|-------|
| Forgetting rows.Close() | loop over `rows` with no Close | `defer rows.Close()` right after the Query error check |
| Not checking rows.Err() | stop at `for rows.Next() {}` | check `rows.Err()` after the loop too |
| String-concatenated SQL | `"...WHERE x='" + input + "'"` | `db.Query("...WHERE x=?", input)` |
| Ignoring sql.ErrNoRows | `log.Fatal(err)` on any Scan error | `errors.Is(err, sql.ErrNoRows)` handled as "not found" |
| No Commit/Rollback | `db.Begin()` then nothing | defer-rollback-if-not-committed, every transaction |
| Plain `:memory:` DSN | `sql.Open("sqlite", ":memory:")` | `file::memory:?cache=shared` + `SetMaxOpenConns(1)` |

---

## 🎓 Before Next Level

Can you:
- [ ] Explain why sql.Open doesn't guarantee a working connection?
- [ ] Write the rows.Next()/Scan()/Close()/Err() loop from memory?
- [ ] Handle sql.ErrNoRows as a normal, non-fatal case?
- [ ] State the one rule that prevents SQL injection?
- [ ] Write the defer-rollback-if-not-committed transaction pattern?
- [ ] Pass a context.Context into a *Context database call?
- [ ] Explain why in-memory SQLite needs SetMaxOpenConns(1)?

If YES → You're ready for Level 30!

---

## 📚 Next Level

Level 30: Repository-Service-Handler
- Structuring a real application into repository, service, and handler layers
- Turning free CRUD functions into a proper repository interface
- Wiring that repository into Gin (Level 28) handlers through a service layer

You've got database/sql down! 💪
