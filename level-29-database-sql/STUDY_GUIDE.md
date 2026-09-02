# Level 29: Study Guide & Visual Reference

## 📚 Learning Path

### Week 1: database/sql Fundamentals
```
Day 1:  The database/sql package, drivers, and blank imports
Day 2:  sql.Open, lazy connections, and sql.DB as a pool
Day 3:  Schema setup with db.Exec; querying multiple rows into structs
Day 4:  Querying a single row and sql.ErrNoRows
Day 5:  Writing data: INSERT/UPDATE/DELETE and sql.Result
Day 6:  Parameterized queries and prepared statements
Day 7:  Transactions: commit and rollback
```

### Week 2: Practice & Application
```
Day 1:  Exercises 1-3 (connecting, schema, querying rows and a single row)
Day 2:  Exercises 4-5 (parameterized queries, prepared statements)
Day 3:  Exercises 6-7 (successful and rolled-back transactions)
Day 4:  Exercises 8-9 (context-aware queries, connection pool settings)
Day 5:  Exercise 10 (comprehensive tested bookstore data layer)
Day 6:  Bonus challenges (migrations, pool exhaustion, repository wrapper)
Day 7:  Review & consolidation
```

---

## ⚠️ A Note on This Level's Verification

Every piece of code, every `go run` output, and the `go test -v` output in this level's materials was written, executed, and captured for real against `modernc.org/sqlite v1.57.0` on Go 1.26.5. This level (along with Level 28: Gin Framework) is the one place in the course where fetching a module from the internet was required first - a single `go get` in Exercise 1 - but every exercise afterward runs entirely against an in-memory database, with zero network or server dependency. Nothing here is simulated or hand-computed.

---

## 🧱 sql.DB: A Connection Pool, Not a Connection

```
                         *sql.DB
                    (one Go value, shared
                     safely across goroutines)
                             │
              ┌──────────────┼──────────────┐
              │              │              │
        connection 1   connection 2   connection 3   ... up to MaxOpenConns
        (idle or          (idle or       (idle or
         in use)           in use)        in use)

Query() / Exec() / QueryRow() each:
  1. borrow an idle connection from the pool (or open a new one)
  2. run the statement on it
  3. return the connection to the pool when the *Rows/*Stmt is closed
     (or immediately, for Exec/QueryRow)

You almost NEVER touch an individual connection directly - the pool
is opaque, and that opacity is the whole point.
```

**The SQLite in-memory exception:** a plain `:memory:` DSN makes every new pooled connection open its own brand-new, empty database - the pool's normal "any connection is as good as any other" assumption breaks. Every exercise in this level uses `file::memory:?cache=shared` plus `db.SetMaxOpenConns(1)` to route every call through the *same* single connection, sidestepping the issue entirely. A real client/server database (Postgres, MySQL) never has this problem.

---

## 🎯 Query vs. QueryRow vs. Exec

| Method | Use for | Returns | Cursor management |
|---|---|---|---|
| `db.Query` | Zero or more rows expected | `*sql.Rows`, `error` | You: `defer rows.Close()`, loop `rows.Next()`, check `rows.Err()` |
| `db.QueryRow` | Exactly zero-or-one row expected | `*sql.Row` (no separate error return) | Handled for you - just call `.Scan(...)` and check its returned error |
| `db.Exec` | `INSERT`/`UPDATE`/`DELETE`/DDL - no rows returned | `sql.Result`, `error` | None needed - `LastInsertId()`/`RowsAffected()` on the result |

```
Do I expect rows back?
├─ No (INSERT/UPDATE/DELETE/CREATE TABLE) ──────────► db.Exec
└─ Yes
   ├─ Exactly zero-or-one row (lookup by key) ──────► db.QueryRow
   └─ Zero-to-many rows (a list) ───────────────────► db.Query
```

Every `*Context` twin (`QueryContext`, `QueryRowContext`, `ExecContext`) follows the exact same three-way split - just with a `context.Context` as the first argument.

---

## 🔁 The rows.Next() / Scan() / Close() / Err() Lifecycle

```
rows, err := db.Query(...)
if err != nil { return err }        ①  check the Query error FIRST
defer rows.Close()                  ②  ALWAYS defer Close immediately after ①

for rows.Next() {                   ③  Next() advances the cursor,
                                        returns false when exhausted
                                        OR when an error occurred
    var x T
    if err := rows.Scan(&x); err != nil {
        return err                  ④  Scan can fail per-row (type mismatch, etc.)
    }
    results = append(results, x)
}

if err := rows.Err(); err != nil {  ⑤  the ONLY way to tell "finished normally"
    return err                          apart from "stopped early due to an error"
}
```

```
Next() returned false
        │
        ├─ rows.Err() == nil  → iteration finished normally, all rows consumed
        └─ rows.Err() != nil  → something went wrong mid-stream (connection
                                  dropped, etc.) - you may have gotten a PARTIAL
                                  result set without knowing it, if you skip ⑤
```

---

## 🔀 Transaction Commit/Rollback Flow

```
                         db.Begin()
                              │
                              ▼
                    ┌───────────────────┐
                    │   tx *sql.Tx       │
                    │  (holds ONE        │
                    │   connection for   │
                    │   its whole life)  │
                    └─────────┬──────────┘
                              │
              defer func() {                    ← set up BEFORE any tx.Exec
                  if !committed {
                      tx.Rollback()              ← the safety net
                  }
              }()
                              │
                 ┌────────────┴────────────┐
                 │                         │
          tx.Exec(...) #1            tx.Exec(...) #2
                 │                         │
         success │                success  │  FAILURE
                 ▼                         ▼        │
          tx.Commit()              (return err)     │
                 │                         │         │
          committed = true         deferred Rollback runs
                 │                  undoes EVERYTHING done
                 ▼                  inside this tx, including
        both statements STICK       statement #1's success
```

**Verified successful commit** (Exercise 6): stock 5→3, orders 0→1, both changes visible together.

**Verified rollback** (Exercise 7): an UPDATE that would violate a `CHECK (stock >= 0)` constraint fails; the deferred `Rollback()` runs; stock and order count both end up **exactly where they started** - `5` and `0` - not partially applied.

---

## 🔒 Parameterized Query Safety, Illustrated

```
❌ STRING-BUILT SQL                    ✅ PARAMETERIZED SQL

  query := fmt.Sprintf(                  db.Query(
    "SELECT * FROM books "                 `SELECT * FROM books
    + "WHERE author = '%s'",                WHERE author = ?`,
    userInput,                              userInput,
  )                                       )
  db.Query(query)

  userInput becomes PART OF             userInput travels as a
  THE SQL TEXT ITSELF -                 separate, typed VALUE -
  quotes/operators in it can            the driver can NEVER
  change the query's structure          interpret it as SQL syntax


  Verified (Exercise 4): passing "Robert C. Martin' OR '1'='1"
  through the PARAMETERIZED path returns ZERO rows - the entire
  string is compared literally against the author column. It never
  gets a chance to alter the query.
```

**The rule, with no exceptions:** every value that isn't a literal you wrote yourself in the source code goes through a `?` placeholder. Always.

---

## 🚨 Common Mistakes

| Mistake | Wrong | Right |
|---|---|---|
| Forgetting rows.Close() | `rows, _ := db.Query(...)` then loop, no Close | `defer rows.Close()` right after checking the Query error |
| Not checking rows.Err() | loop `rows.Next()` and stop | check `rows.Err()` after the loop too |
| String-concatenating SQL | `"...WHERE x = '" + input + "'"` | `db.Query("...WHERE x = ?", input)` |
| Forgetting sql.ErrNoRows | treating "not found" as `log.Fatal` | `errors.Is(err, sql.ErrNoRows)` as a normal case |
| Forgetting Commit/Rollback | `db.Begin()` then nothing | the defer-rollback-if-not-committed pattern, every time |
| Unlimited pool on real DBs | leaving `SetMaxOpenConns` at its default (unlimited) | set a limit matched to your database server's capacity |
| Plain `:memory:` + multiple connections | `sql.Open("sqlite", ":memory:")` with default pool settings | `file::memory:?cache=shared` + `SetMaxOpenConns(1)` |

---

## 📈 Progression Summary

### Understanding Level 29

Level 29 gives your Gin handlers (Level 28) a place to actually keep data:

1. **The interface** - database/sql is driver-agnostic; the driver is a separate, blank-imported module
2. **The pool** - sql.DB manages connections for you; you interact with Query/Exec/QueryRow, not raw connections
3. **Reading** - Query for many rows, QueryRow for zero-or-one, always mindful of Close()/Err()/ErrNoRows
4. **Writing** - Exec + sql.Result, always parameterized, never string-built
5. **Structure** - prepared statements for repeated SQL, transactions for multi-step atomicity
6. **Cancellation** - the *Context family, bridging to Level 24

### Prerequisites for Level 30

Before moving to Level 30 (Repository-Service-Handler), you need:

- ✅ Comfortable opening a database and understanding sql.Open's laziness
- ✅ Comfortable querying multiple rows into structs and a single row with ErrNoRows handling
- ✅ Always using parameterized queries, without exception
- ✅ Comfortable writing a transaction with the defer-rollback-if-not-committed pattern
- ✅ Understand what sql.DB's connection pool is and how to configure it
- ✅ Completed Exercise 10's tested CRUD data layer

### Ready for Level 30?

Level 30 takes Exercise 10's free functions and formalizes them into a proper layered architecture:
- A **repository** layer wrapping `*sql.DB` behind an interface (Bonus Challenge 3's preview)
- A **service** layer holding business logic, depending on the repository interface, not the concrete database
- A **handler** layer (Gin, from Level 28) translating HTTP requests into service calls

---

## ✅ Checklist Before Level 30

- [ ] Can explain why sql.Open doesn't guarantee a live connection
- [ ] Can query multiple rows into a struct slice, with proper Close()/Err() handling
- [ ] Can query a single row and correctly handle sql.ErrNoRows
- [ ] Never writes SQL by string concatenation - always `?` placeholders
- [ ] Can write a transaction using the defer-rollback-if-not-committed pattern
- [ ] Has seen a real rollback leave data completely unchanged
- [ ] Can pass a context.Context into a *Context query/exec call
- [ ] Understands SetMaxOpenConns/SetMaxIdleConns and why in-memory SQLite needs SetMaxOpenConns(1)
- [ ] Completed 8+ exercises, including the fully-tested Exercise 10

---

## 💡 Key Takeaways

### The Interface
`database/sql` is a generic interface; the driver (blank-imported) does the real talking to the database.

### The Pool
`sql.DB` is a pool, not a connection - `sql.Open` is lazy, `db.Ping()` verifies eagerly.

### The Safety Rule
Parameterized `?` queries, every time, for every value that isn't a literal you wrote yourself.

### The Atomicity Tool
Transactions plus the defer-rollback-if-not-committed pattern turn multi-step writes into all-or-nothing units.

---

## 📚 Next Level

Level 30: Repository-Service-Handler
- Structuring a real application into repository, service, and handler layers
- Turning this level's free functions into a proper repository interface
- Wiring that repository into Gin (Level 28) handlers through a service layer

You've got database/sql down! Keep going! 🚀
