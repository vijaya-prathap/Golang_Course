# Level 29: Database/SQL - Complete Guide

## Introduction

Welcome to Level 29! You've mastered the Gin framework (Level 28) - routing HTTP requests, binding JSON, and returning responses. Every API you've built so far has kept its data in memory, which means it vanishes the moment the process restarts. It's time to fix that: **database/sql**, Go's standard library for talking to relational databases.

`database/sql` is unusual among standard library packages in one important way: it is a generic *interface* over databases, not a database itself, and not even a driver. It defines `sql.DB`, `sql.Rows`, `sql.Tx`, and the rest of the vocabulary you'll use for the remainder of this course - but it cannot open a single connection on its own. You always pair it with a **driver** package, imported purely for its side effects (`_ "modernc.org/sqlite"`), which registers itself with `database/sql` and does the actual talking to a real database engine.

**A note on this level's setup:** this course reaches offline into your machine for almost everything, but this level - like Level 28: Gin Framework - is a deliberate, clearly-flagged exception. Go's standard library alone cannot reach any real database; you need a driver module, and driver modules are fetched from the internet the same way Gin was. This level uses **`modernc.org/sqlite` (pinned at `v1.57.0`)**, a pure-Go SQLite driver with no cgo and no external database server required. It was chosen specifically because, once fetched, every exercise in this level runs against an **in-memory SQLite database** - zero network access, zero server process, zero configuration beyond `go get` - which keeps the rest of the level exactly as self-contained as every other one. All code, output, and behavior in this level's materials was written, run, and captured for real (`go run`, `go test`) against this exact driver version on Go 1.26.5 - there is nothing hand-computed or simulated here.

---

## Table of Contents

1. [The database/sql Package](#the-databasesql-package)
2. [Connecting: sql.Open, Lazy Connections, and the Pool](#connecting-sqlopen-lazy-connections-and-the-pool)
3. [Schema Setup: Creating Tables With db.Exec](#schema-setup-creating-tables-with-dbexec)
4. [Querying Multiple Rows](#querying-multiple-rows)
5. [Querying a Single Row](#querying-a-single-row)
6. [Writing Data: INSERT, UPDATE, DELETE](#writing-data-insert-update-delete)
7. [Preventing SQL Injection: Parameterized Queries](#preventing-sql-injection-parameterized-queries)
8. [Prepared Statements](#prepared-statements)
9. [Transactions](#transactions)
10. [Context-Aware Queries](#context-aware-queries)
11. [Best Practices](#best-practices)
12. [Common Mistakes](#common-mistakes)

---

## The database/sql Package

`database/sql` defines the shape every Go database interaction follows, no matter which database sits underneath: `sql.DB` for a connection pool, `sql.Rows`/`sql.Row` for query results, `sql.Result` for write outcomes, `sql.Tx` for transactions. None of that code knows how to speak SQLite's, Postgres's, or MySQL's wire protocol - that's the driver's job.

A driver registers itself with `database/sql` through a **blank import** - imported only for its `init()` side effect, never referenced by name in your code:

```go
import (
    "database/sql"

    _ "modernc.org/sqlite" // registers the "sqlite" driver name
)
```

The underscore (`_`) tells Go "import this for its side effects, I'm not calling anything on it directly." Internally, the driver calls `sql.Register("sqlite", ...)` in its own `init()` function, which is what makes the string `"sqlite"` a valid first argument to `sql.Open` later. If you forget this import, `sql.Open("sqlite", ...)` compiles fine but fails at runtime with `sql: unknown driver "sqlite" (forgotten import?)`.

This design is why switching from SQLite to Postgres in a real project is (mostly) a one-line change: swap the driver import and the DSN, and almost all of your `database/sql` code - `Query`, `Exec`, `Scan`, transactions - stays exactly the same.

---

## Connecting: sql.Open, Lazy Connections, and the Pool

```go
db, err := sql.Open("sqlite", "file::memory:?cache=shared")
if err != nil {
    log.Fatal(err)
}
defer db.Close()
```

**`sql.Open` is lazy.** It validates its arguments and returns a `*sql.DB` immediately - it does **not** open a network connection, does **not** verify credentials, and does **not** confirm the database even exists. The first real connection attempt happens on the first query, exec, or explicit ping. This means `sql.Open` returning `nil` error tells you almost nothing about whether the database is actually reachable.

To verify connectivity eagerly - typically right after startup, before you start accepting traffic - call `db.Ping()`:

```go
if err := db.Ping(); err != nil {
    log.Fatal(err) // NOW you know if the database is unreachable
}
```

**`sql.DB` is a connection *pool*, not a single connection.** This is the single most important mental model shift coming from other languages: every `Query`, `QueryRow`, and `Exec` call borrows a connection from the pool, uses it, and returns it - you almost never see or manage an individual connection yourself. `database/sql` opens new connections as needed (up to a configurable limit, Section 11) and reuses idle ones automatically.

### The SQLite In-Memory Gotcha

This pooling model has a real, verified consequence for `:memory:` SQLite databases that every exercise in this level works around: **the DSN `:memory:` alone creates a brand-new, separate, empty in-memory database for every new connection the pool opens.** Two connections opened from the same `*sql.DB` with plain `:memory:` are *not* looking at the same data - the second one sees an empty database with none of your tables.

The fix, used throughout this level, is two-part:

```go
db, err := sql.Open("sqlite", "file::memory:?cache=shared")
db.SetMaxOpenConns(1)
```

- `file::memory:?cache=shared` puts SQLite in **shared-cache mode**, so multiple connections *can* see the same in-memory database.
- `db.SetMaxOpenConns(1)` then pins the pool to exactly one connection, which is what a single-writer, file-based/in-memory database like SQLite wants anyway - and it sidesteps a second, subtler problem where an uncommitted write transaction on one shared-cache connection can block a query issued on another.

This isn't a workaround you'll need for a real client/server database like Postgres or MySQL (they don't have this "new connection, new empty database" behavior), but it's essential, correct, real-world knowledge for exactly the kind of embedded/in-memory SQLite usage this level relies on. Section 11 covers connection pool settings, including `SetMaxOpenConns`, in full.

---

## Schema Setup: Creating Tables With db.Exec

DDL (Data Definition Language - `CREATE TABLE`, `ALTER TABLE`, and friends) goes through the same `db.Exec` method you'll use for writes later. This level uses one small schema throughout - a `books` table:

```go
schema := `
CREATE TABLE books (
    id     INTEGER PRIMARY KEY AUTOINCREMENT,
    title  TEXT NOT NULL,
    author TEXT NOT NULL,
    year   INTEGER NOT NULL
);`

if _, err := db.Exec(schema); err != nil {
    log.Fatal(err)
}
```

`db.Exec` returns a `sql.Result` and an `error`. For DDL, the `sql.Result` isn't useful (there's no "row" that was inserted), so it's idiomatic to discard it with `_` and just check the error, exactly as above.

Real captured output from Exercise 1, which opens a database, pings it, creates this schema, and confirms the table exists via SQLite's own `sqlite_master` catalog table:

```
=== Opening the Database ===
sql.Open succeeded (no connection made yet)

=== Verifying the Connection ===
Ping succeeded - connection is live

=== Creating the Schema ===
Table 'books' created

=== Confirming the Table Exists ===
Found table: books
```

---

## Querying Multiple Rows

`db.Query` runs a SQL query expected to return zero or more rows. The result is a `*sql.Rows` cursor you iterate manually:

```go
type Book struct {
    ID     int
    Title  string
    Author string
    Year   int
}

rows, err := db.Query(`SELECT id, title, author, year FROM books ORDER BY id`)
if err != nil {
    log.Fatal(err)
}
defer rows.Close() // ALWAYS - even if you return early on error below

var books []Book
for rows.Next() {
    var b Book
    if err := rows.Scan(&b.ID, &b.Title, &b.Author, &b.Year); err != nil {
        log.Fatal(err)
    }
    books = append(books, b)
}
if err := rows.Err(); err != nil { // checks for errors DURING iteration
    log.Fatal(err)
}
```

Four rules, every single time you call `db.Query`:

1. **Check the error from `Query` itself** before touching `rows`.
2. **`defer rows.Close()` immediately** after confirming `err == nil` - forgetting this leaks the underlying connection back to the pool late (or never), starving other callers.
3. **`rows.Scan(&dest...)` inside the `for rows.Next()` loop**, with one pointer argument per selected column, in order.
4. **Check `rows.Err()` after the loop** - `rows.Next()` returns `false` both when iteration finishes normally *and* when an error occurred mid-stream, and `rows.Err()` is the only way to tell them apart.

---

## Querying a Single Row

When you expect exactly zero or one row - a lookup by primary key, for example - `db.QueryRow` is simpler than `db.Query`: no cursor, no `rows.Next()`, no `rows.Close()` (it's handled for you internally).

```go
var b Book
err := db.QueryRow(
    `SELECT id, title, author, year FROM books WHERE id = ?`, id,
).Scan(&b.ID, &b.Title, &b.Author, &b.Year)
```

The special case: **when no row matches, `Scan` returns `sql.ErrNoRows`** - not a generic "not found" error, and not a nil `Book` with no error. You must check for it explicitly:

```go
if errors.Is(err, sql.ErrNoRows) {
    // no book with that id - this is often NOT a fatal error in your program,
    // just "not found"
} else if err != nil {
    log.Fatal(err) // a real problem: bad SQL, connection lost, etc.
}
```

Real captured output, looking up an existing id and then a missing one:

```
=== Looking Up an Existing Book (id=1) ===
Found: Clean Code by Robert C. Martin (2008)

=== Looking Up a Book That Doesn't Exist (id=99) ===
No book found with id=99 (sql.ErrNoRows)
```

Forgetting to check for `sql.ErrNoRows` is one of the most common real-world `database/sql` bugs - the error gets logged as if it were a server outage, when it's really just "zero results," which is often a perfectly normal outcome.

---

## Writing Data: INSERT, UPDATE, DELETE

`db.Exec` handles all three write statements and returns a `sql.Result` with two useful methods:

```go
result, err := db.Exec(
    `INSERT INTO books (title, author, year) VALUES (?, ?, ?)`,
    "Clean Code", "Robert C. Martin", 2008,
)
if err != nil {
    log.Fatal(err)
}

id, err := result.LastInsertId() // the auto-generated id of the new row
rows, err := result.RowsAffected() // how many rows were touched
```

- **`LastInsertId()`** is meaningful after an `INSERT` into a table with an auto-increment primary key (not all databases support it the same way - SQLite and MySQL do, Postgres traditionally wants `RETURNING id` instead).
- **`RowsAffected()`** works after `INSERT`, `UPDATE`, or `DELETE`, and is the standard way to detect "I tried to update/delete something that didn't exist" - a `RowsAffected() == 0` after an `UPDATE ... WHERE id = ?` almost always means that id wasn't there.

```go
result, _ := db.Exec(`UPDATE books SET year = ? WHERE id = ?`, 2009, 1)
affected, _ := result.RowsAffected()
if affected == 0 {
    // nothing had that id - treat this like a "not found"
}
```

---

## Preventing SQL Injection: Parameterized Queries

**Rule, no exceptions: never build a SQL string by concatenating or `fmt.Sprintf`-ing user input into it.** Always use `?` placeholders and pass the real values as separate arguments to `Query`/`QueryRow`/`Exec`. The driver sends the SQL text and the parameter values to the database *separately* - a parameter can never be interpreted as SQL syntax, no matter what it contains.

```go
// ✅ SAFE - the value travels as a bound parameter, never as SQL text
rows, err := db.Query(`SELECT title FROM books WHERE author = ?`, userInput)

// ❌ NEVER DO THIS - string-built SQL lets user input control the query's structure
query := fmt.Sprintf(`SELECT title FROM books WHERE author = '%s'`, userInput)
rows, err := db.Query(query)
```

This isn't a theoretical concern demonstrated with a fake exploit - the safe pattern is simply the *only* pattern this course ever uses, everywhere, including every exercise in this level. To see why it holds even against hostile-looking input, Exercise 4 passes a string that *looks* like an injection attempt straight through the safe, parameterized `searchByAuthor` function:

```
=== Safe Lookup With a Normal Value ===
Books by Robert C. Martin: [Clean Code The Clean Coder]

=== Safe Lookup With a Hostile-Looking Value ===
Books by "Robert C. Martin' OR '1'='1": []
(No rows match because the whole string is treated as one literal author name - the query was never restructured.)
```

The entire hostile-looking string is compared, literally, character-for-character, against the `author` column - it never gets a chance to alter the query's shape, because it was never part of the SQL text in the first place.

---

## Prepared Statements

`db.Prepare` compiles a SQL statement once and returns a `*sql.Stmt` you can `Exec`/`Query` repeatedly with different parameters, skipping re-parsing/re-planning on every call:

```go
stmt, err := db.Prepare(`INSERT INTO books (title, author, year) VALUES (?, ?, ?)`)
if err != nil {
    log.Fatal(err)
}
defer stmt.Close()

for _, b := range books {
    if _, err := stmt.Exec(b.title, b.author, b.year); err != nil {
        log.Fatal(err)
    }
}
```

Prepared statements help most when the **same** statement text runs **many** times in a loop or across many requests - batch inserts, a lookup query called on every HTTP request, and so on. For a one-off query run exactly once, `db.Query`/`db.Exec` already prepare-and-execute internally, so a separate `Prepare` call buys you nothing. `*sql.Stmt` is itself safe for concurrent use and, like `*sql.DB`, draws from the connection pool - always `defer stmt.Close()` when you're done with it, just like `rows.Close()`.

---

## Transactions

A transaction groups multiple statements into one all-or-nothing unit: either every statement succeeds and the whole group is committed, or something fails and the whole group is rolled back as if none of it happened.

```go
tx, err := db.Begin()
if err != nil {
    log.Fatal(err)
}
committed := false
defer func() {
    if !committed {
        tx.Rollback() // safety net - runs if we return before Commit for ANY reason
    }
}()

if _, err := tx.Exec(`UPDATE books SET stock = stock - ? WHERE id = ?`, 2, 1); err != nil {
    return err // deferred Rollback cleans up
}
if _, err := tx.Exec(`INSERT INTO orders (book_id, quantity) VALUES (?, ?)`, 1, 2); err != nil {
    return err // deferred Rollback cleans up
}

if err := tx.Commit(); err != nil {
    return err
}
committed = true
```

The **defer-rollback-if-not-committed** pattern above is the standard, idiomatic shape for every transaction you write: declare a `committed` flag, defer a rollback that only fires when the flag is still `false`, and set the flag to `true` only immediately after a successful `Commit()`. Calling `Rollback()` on an already-committed transaction is a harmless no-op error you can safely ignore here - `database/sql` tracks whether the transaction is still open.

**Verified successful transaction** (Exercise 6) - decrementing stock and recording an order together:

```
[before] stock=5, orders=0

=== Placing an Order for 2 Copies (Successful Transaction) ===
Transaction committed: stock decremented AND order recorded
[after commit] stock=3, orders=1
```

**Verified failed transaction** (Exercise 7) - a `CHECK (stock >= 0)` constraint rejects an order for more books than are in stock, and the whole transaction rolls back, leaving *both* statements' worth of change completely undone:

```
[before] stock=5, orders=0

=== Attempting to Order 10 Copies When Only 5 Are In Stock (Failing Transaction) ===
Transaction failed and was rolled back: constraint failed: CHECK constraint failed: stock >= 0 (275)
[after] stock=5, orders=0

Stock and order count are UNCHANGED from before - the partially-applied UPDATE never survived the rollback.
```

Notice `stock` is `5` both before and after - not `-5`, not partially decremented. That's the entire point of a transaction: the failed `UPDATE` inside it never became visible on its own.

---

## Context-Aware Queries

Every `database/sql` method has a `*Context` twin that accepts a `context.Context` (Level 24) as its first argument: `QueryContext`, `QueryRowContext`, `ExecContext`, `PrepareContext`, `BeginTx`. Passing a context lets a caller cancel a slow query or enforce a timeout, exactly like any other context-aware operation:

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

var title string
err := db.QueryRowContext(ctx, `SELECT title FROM books WHERE id = ?`, 1).Scan(&title)
```

Verified: an already-expired context makes the query fail immediately with `context.DeadlineExceeded`, without ever reaching the database engine:

```
=== Query With a Generous Timeout (Succeeds) ===
Found: Clean Code

=== Query With an Already-Expired Context (Fails) ===
Query canceled: context deadline exceeded

=== ExecContext Respects Cancellation Too ===
Insert canceled: context canceled
```

In a real HTTP handler (Level 28's territory), you'd pass `r.Context()` straight through to every database call: if the client disconnects or the request's deadline elapses, every in-flight query tied to that context stops immediately instead of continuing to do useless work.

---

## Best Practices

### 1. Always Check Every Error

`database/sql` returns an `error` from nearly every method. A dropped error here is invisible right up until it silently corrupts data or hides a real outage.

### 2. Always Close What You Open

```go
rows, err := db.Query(...)
defer rows.Close()

stmt, err := db.Prepare(...)
defer stmt.Close()
```

An unclosed `*sql.Rows` holds its underlying connection out of the pool until it's garbage collected - which might be much later, or never, under sustained load.

### 3. Always Use Parameterized Queries

`?` placeholders, every time, for every value that isn't a compile-time constant you wrote yourself. See Section 7.

### 4. Use Transactions for Multi-Step Writes

If two or more `Exec` calls need to succeed or fail together as one logical unit, wrap them in a `db.Begin()`/`tx.Commit()`/`tx.Rollback()` block - never rely on "well, they'll probably both succeed."

### 5. Set Sensible Connection Pool Limits

```go
db.SetMaxOpenConns(25)
db.SetMaxIdleConns(25)
db.SetConnMaxLifetime(5 * time.Minute)
```

The right numbers depend on your database server's own connection limits and your workload - but leaving these at their defaults (unlimited open connections) is a common way to accidentally overwhelm a real database server under load. In-memory SQLite is the one case in this level where the answer is a hard `SetMaxOpenConns(1)` (Section 2).

---

## Common Mistakes

### Mistake 1: Forgetting rows.Close()

```go
// ❌ WRONG - leaks the connection back to the pool
rows, _ := db.Query(`SELECT * FROM books`)
for rows.Next() { /* ... */ }

// ✅ RIGHT
rows, err := db.Query(`SELECT * FROM books`)
if err != nil {
    log.Fatal(err)
}
defer rows.Close()
```

### Mistake 2: Not Checking rows.Err()

```go
// ❌ WRONG - a mid-iteration error looks identical to "no more rows"
for rows.Next() {
    rows.Scan(&x)
}
// silently ignores a possible error!

// ✅ RIGHT
for rows.Next() {
    rows.Scan(&x)
}
if err := rows.Err(); err != nil {
    log.Fatal(err)
}
```

### Mistake 3: String-Concatenating SQL

```go
// ❌ NEVER
query := "SELECT * FROM books WHERE title = '" + userInput + "'"

// ✅ ALWAYS
db.Query(`SELECT * FROM books WHERE title = ?`, userInput)
```

### Mistake 4: Forgetting to Check sql.ErrNoRows

```go
// ❌ WRONG - treats "not found" as a fatal server error
err := db.QueryRow(`SELECT * FROM books WHERE id = ?`, id).Scan(&b)
if err != nil {
    log.Fatal(err) // crashes the program just because id didn't exist!
}

// ✅ RIGHT
err := db.QueryRow(`SELECT * FROM books WHERE id = ?`, id).Scan(&b)
if errors.Is(err, sql.ErrNoRows) {
    // handle "not found" as a normal case
} else if err != nil {
    log.Fatal(err) // a REAL error
}
```

### Mistake 5: Forgetting to Commit or Rollback a Transaction

```go
// ❌ WRONG - the transaction is left open, holding a connection, forever
tx, _ := db.Begin()
tx.Exec(`UPDATE books SET stock = stock - 1 WHERE id = ?`, 1)
// ... no Commit(), no Rollback() ...

// ✅ RIGHT - the defer-rollback-if-not-committed pattern (Section 9)
tx, err := db.Begin()
committed := false
defer func() {
    if !committed {
        tx.Rollback()
    }
}()
// ... tx.Exec calls ...
if err := tx.Commit(); err != nil {
    return err
}
committed = true
```

---

## Summary

**The Core Model:**
- `database/sql` is an interface over drivers - always paired with a blank-imported driver like `_ "modernc.org/sqlite"`
- `sql.Open` is lazy; `sql.DB` is a connection *pool*, not one connection; `db.Ping()` verifies connectivity eagerly

**Reading Data:**
- `db.Query` + `rows.Next()`/`rows.Scan()`/`defer rows.Close()`/`rows.Err()` for many rows
- `db.QueryRow` + `row.Scan()` for exactly zero-or-one row, watching for `sql.ErrNoRows`

**Writing Data:**
- `db.Exec` for `INSERT`/`UPDATE`/`DELETE`, returning a `sql.Result` with `LastInsertId()`/`RowsAffected()`
- Always parameterized with `?` placeholders - never string-built

**Structure:**
- `db.Prepare` for a statement reused many times
- `db.Begin()`/`tx.Commit()`/`tx.Rollback()` with the defer-rollback-if-not-committed pattern for multi-step writes
- `*Context` variants (`QueryContext`, `ExecContext`, ...) for cancellation and timeouts

---

## Next Steps

You now understand:
- ✅ How database/sql, drivers, and connection pools fit together
- ✅ Querying single and multiple rows into Go structs
- ✅ Writing data safely with parameterized queries
- ✅ Prepared statements and transactions, including a verified rollback
- ✅ Context-aware database calls bridging to Level 24

**Next level:** Level 30 - Repository-Service-Handler
- Structuring a real application into repository, service, and handler layers
- Wrapping the `*sql.DB` access you just learned behind a clean repository interface
- Connecting that layered architecture to the Gin handlers from Level 28

You can now make your programs remember things! Keep going! 🚀
