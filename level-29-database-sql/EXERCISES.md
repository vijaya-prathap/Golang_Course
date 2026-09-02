# Level 29: Database/SQL - Exercises

Complete all exercises in order. Each exercise builds on previous knowledge.

**Before you start:** every exercise below needs the `modernc.org/sqlite` driver, pinned at `v1.57.0`. The `go get` step in Exercise 1 fetches it - this is the only exercise set in the course (along with Level 28: Gin) that needs internet access, and only for that one `go get`. Every program you run afterward talks only to an in-memory SQLite database - no network, no external server.

---

## Exercise 1: Connecting and Creating a Schema

**Objective:** Open a database connection, verify it, and create a table

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level29-exercise1
cd ~/projects/level29-exercise1
go mod init level29.example/exercise1
go get modernc.org/sqlite@v1.57.0
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "database/sql"
    "fmt"
    "log"

    _ "modernc.org/sqlite"
)

func main() {
    fmt.Println("=== Opening the Database ===")
    db, err := sql.Open("sqlite", "file::memory:?cache=shared")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()
    db.SetMaxOpenConns(1) // in-memory SQLite needs exactly one connection
    fmt.Println("sql.Open succeeded (no connection made yet)")

    fmt.Println("\n=== Verifying the Connection ===")
    if err := db.Ping(); err != nil {
        log.Fatal(err)
    }
    fmt.Println("Ping succeeded - connection is live")

    fmt.Println("\n=== Creating the Schema ===")
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
    fmt.Println("Table 'books' created")

    fmt.Println("\n=== Confirming the Table Exists ===")
    var name string
    err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='books'`).Scan(&name)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Found table:", name)
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

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

**Learning Objectives:**
- ✅ Import a driver with a blank import alongside database/sql
- ✅ Understand sql.Open is lazy and db.Ping() verifies connectivity eagerly
- ✅ Create a table with db.Exec
- ✅ Understand why in-memory SQLite needs the shared-cache DSN + SetMaxOpenConns(1)

---

## Exercise 2: Inserting Data and Querying Multiple Rows Into Structs

**Objective:** Insert several rows, then read them back into a slice of structs

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level29-exercise2
cd ~/projects/level29-exercise2
go mod init level29.example/exercise2
go get modernc.org/sqlite@v1.57.0
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "database/sql"
    "fmt"
    "log"

    _ "modernc.org/sqlite"
)

type Book struct {
    ID     int
    Title  string
    Author string
    Year   int
}

func main() {
    db, err := sql.Open("sqlite", "file::memory:?cache=shared")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()
    db.SetMaxOpenConns(1)

    if _, err := db.Exec(`CREATE TABLE books (
        id     INTEGER PRIMARY KEY AUTOINCREMENT,
        title  TEXT NOT NULL,
        author TEXT NOT NULL,
        year   INTEGER NOT NULL
    );`); err != nil {
        log.Fatal(err)
    }

    fmt.Println("=== Inserting Books ===")
    books := []Book{
        {Title: "The Go Programming Language", Author: "Donovan & Kernighan", Year: 2015},
        {Title: "Clean Code", Author: "Robert C. Martin", Year: 2008},
        {Title: "The Pragmatic Programmer", Author: "Hunt & Thomas", Year: 1999},
    }

    for _, b := range books {
        result, err := db.Exec(
            `INSERT INTO books (title, author, year) VALUES (?, ?, ?)`,
            b.Title, b.Author, b.Year,
        )
        if err != nil {
            log.Fatal(err)
        }
        id, err := result.LastInsertId()
        if err != nil {
            log.Fatal(err)
        }
        fmt.Printf("Inserted %q with id=%d\n", b.Title, id)
    }

    fmt.Println("\n=== Querying All Books ===")
    rows, err := db.Query(`SELECT id, title, author, year FROM books ORDER BY id`)
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()

    var results []Book
    for rows.Next() {
        var b Book
        if err := rows.Scan(&b.ID, &b.Title, &b.Author, &b.Year); err != nil {
            log.Fatal(err)
        }
        results = append(results, b)
    }
    if err := rows.Err(); err != nil {
        log.Fatal(err)
    }

    for _, b := range results {
        fmt.Printf("#%d: %s by %s (%d)\n", b.ID, b.Title, b.Author, b.Year)
    }
    fmt.Println("\nTotal books:", len(results))
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Inserting Books ===
Inserted "The Go Programming Language" with id=1
Inserted "Clean Code" with id=2
Inserted "The Pragmatic Programmer" with id=3

=== Querying All Books ===
#1: The Go Programming Language by Donovan & Kernighan (2015)
#2: Clean Code by Robert C. Martin (2008)
#3: The Pragmatic Programmer by Hunt & Thomas (1999)

Total books: 3
```

**Learning Objectives:**
- ✅ Insert rows with db.Exec and read back LastInsertId
- ✅ Query multiple rows with db.Query and rows.Next()/rows.Scan()
- ✅ Always defer rows.Close() and check rows.Err() after the loop

---

## Exercise 3: Querying a Single Row and Handling sql.ErrNoRows

**Objective:** Look up one row by id, and correctly handle the "not found" case

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level29-exercise3
cd ~/projects/level29-exercise3
go mod init level29.example/exercise3
go get modernc.org/sqlite@v1.57.0
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "database/sql"
    "errors"
    "fmt"
    "log"

    _ "modernc.org/sqlite"
)

type Book struct {
    ID     int
    Title  string
    Author string
    Year   int
}

func getBookByID(db *sql.DB, id int) (Book, error) {
    var b Book
    row := db.QueryRow(`SELECT id, title, author, year FROM books WHERE id = ?`, id)
    err := row.Scan(&b.ID, &b.Title, &b.Author, &b.Year)
    return b, err
}

func main() {
    db, err := sql.Open("sqlite", "file::memory:?cache=shared")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()
    db.SetMaxOpenConns(1)

    if _, err := db.Exec(`CREATE TABLE books (
        id     INTEGER PRIMARY KEY AUTOINCREMENT,
        title  TEXT NOT NULL,
        author TEXT NOT NULL,
        year   INTEGER NOT NULL
    );`); err != nil {
        log.Fatal(err)
    }
    if _, err := db.Exec(
        `INSERT INTO books (title, author, year) VALUES (?, ?, ?)`,
        "Clean Code", "Robert C. Martin", 2008,
    ); err != nil {
        log.Fatal(err)
    }

    fmt.Println("=== Looking Up an Existing Book (id=1) ===")
    b, err := getBookByID(db, 1)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Found: %s by %s (%d)\n", b.Title, b.Author, b.Year)

    fmt.Println("\n=== Looking Up a Book That Doesn't Exist (id=99) ===")
    _, err = getBookByID(db, 99)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            fmt.Println("No book found with id=99 (sql.ErrNoRows)")
        } else {
            log.Fatal(err)
        }
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
=== Looking Up an Existing Book (id=1) ===
Found: Clean Code by Robert C. Martin (2008)

=== Looking Up a Book That Doesn't Exist (id=99) ===
No book found with id=99 (sql.ErrNoRows)
```

**Learning Objectives:**
- ✅ Use db.QueryRow + row.Scan for a single-row lookup
- ✅ Recognize and handle sql.ErrNoRows with errors.Is
- ✅ Understand that "not found" is a normal outcome, not a fatal error

---

## Exercise 4: Preventing SQL Injection - Parameterized Queries

**Objective:** Prove the safe, parameterized query pattern holds even against hostile-looking input

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level29-exercise4
cd ~/projects/level29-exercise4
go mod init level29.example/exercise4
go get modernc.org/sqlite@v1.57.0
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "database/sql"
    "fmt"
    "log"

    _ "modernc.org/sqlite"
)

// searchByAuthor is the SAFE pattern: the user-supplied value is passed as a
// query parameter, never concatenated into the SQL string.
func searchByAuthor(db *sql.DB, author string) ([]string, error) {
    rows, err := db.Query(`SELECT title FROM books WHERE author = ?`, author)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var titles []string
    for rows.Next() {
        var title string
        if err := rows.Scan(&title); err != nil {
            return nil, err
        }
        titles = append(titles, title)
    }
    return titles, rows.Err()
}

func main() {
    db, err := sql.Open("sqlite", "file::memory:?cache=shared")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()
    db.SetMaxOpenConns(1)

    if _, err := db.Exec(`CREATE TABLE books (
        id     INTEGER PRIMARY KEY AUTOINCREMENT,
        title  TEXT NOT NULL,
        author TEXT NOT NULL,
        year   INTEGER NOT NULL
    );`); err != nil {
        log.Fatal(err)
    }
    seed := []struct {
        title, author string
        year          int
    }{
        {"Clean Code", "Robert C. Martin", 2008},
        {"The Clean Coder", "Robert C. Martin", 2011},
        {"Refactoring", "Martin Fowler", 1999},
    }
    for _, s := range seed {
        if _, err := db.Exec(`INSERT INTO books (title, author, year) VALUES (?, ?, ?)`, s.title, s.author, s.year); err != nil {
            log.Fatal(err)
        }
    }

    fmt.Println("=== Safe Lookup With a Normal Value ===")
    titles, err := searchByAuthor(db, "Robert C. Martin")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Books by Robert C. Martin:", titles)

    fmt.Println("\n=== Safe Lookup With a Hostile-Looking Value ===")
    // This string LOOKS like a SQL injection attempt. Because it travels as a
    // bound parameter (the ? placeholder), the driver sends it to SQLite as
    // pure data, not as SQL syntax. It is compared literally to the author
    // column - it does NOT alter the query's structure or leak all rows.
    hostileInput := "Robert C. Martin' OR '1'='1"
    titles, err = searchByAuthor(db, hostileInput)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Books by %q: %v\n", hostileInput, titles)
    fmt.Println("(No rows match because the whole string is treated as one literal author name - the query was never restructured.)")

    fmt.Println("\n=== Why This Matters ===")
    fmt.Println("NEVER build SQL with fmt.Sprintf or string concatenation of user input.")
    fmt.Println("ALWAYS use ? placeholders and pass values as separate Query/Exec arguments.")
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Safe Lookup With a Normal Value ===
Books by Robert C. Martin: [Clean Code The Clean Coder]

=== Safe Lookup With a Hostile-Looking Value ===
Books by "Robert C. Martin' OR '1'='1": []
(No rows match because the whole string is treated as one literal author name - the query was never restructured.)

=== Why This Matters ===
NEVER build SQL with fmt.Sprintf or string concatenation of user input.
ALWAYS use ? placeholders and pass values as separate Query/Exec arguments.
```

**Learning Objectives:**
- ✅ Always pass user input as a bound `?` parameter, never string-built
- ✅ See that a hostile-looking string is treated as pure literal data
- ✅ State the one rule that prevents SQL injection, from memory

---

## Exercise 5: Prepared Statements

**Objective:** Prepare a statement once and reuse it for many inserts and queries

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level29-exercise5
cd ~/projects/level29-exercise5
go mod init level29.example/exercise5
go get modernc.org/sqlite@v1.57.0
```

2. Create `main.go`:

```bash
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
    db.SetMaxOpenConns(1)

    if _, err := db.Exec(`CREATE TABLE books (
        id     INTEGER PRIMARY KEY AUTOINCREMENT,
        title  TEXT NOT NULL,
        author TEXT NOT NULL,
        year   INTEGER NOT NULL
    );`); err != nil {
        log.Fatal(err)
    }

    fmt.Println("=== Preparing an INSERT Statement Once ===")
    stmt, err := db.Prepare(`INSERT INTO books (title, author, year) VALUES (?, ?, ?)`)
    if err != nil {
        log.Fatal(err)
    }
    defer stmt.Close()
    fmt.Println("Statement prepared")

    fmt.Println("\n=== Reusing It Many Times ===")
    books := []struct {
        title, author string
        year          int
    }{
        {"Clean Code", "Robert C. Martin", 2008},
        {"Refactoring", "Martin Fowler", 1999},
        {"The Mythical Man-Month", "Fred Brooks", 1975},
        {"Design Patterns", "Gang of Four", 1994},
    }
    for _, b := range books {
        result, err := stmt.Exec(b.title, b.author, b.year)
        if err != nil {
            log.Fatal(err)
        }
        id, _ := result.LastInsertId()
        fmt.Printf("Inserted id=%d: %s\n", id, b.title)
    }

    fmt.Println("\n=== Preparing a SELECT Statement Once ===")
    lookup, err := db.Prepare(`SELECT title, author, year FROM books WHERE year > ?`)
    if err != nil {
        log.Fatal(err)
    }
    defer lookup.Close()

    fmt.Println("\n=== Reusing the SELECT for Different Parameters ===")
    for _, minYear := range []int{1990, 2000} {
        fmt.Printf("Books after %d:\n", minYear)
        rows, err := lookup.Query(minYear)
        if err != nil {
            log.Fatal(err)
        }
        for rows.Next() {
            var title, author string
            var year int
            if err := rows.Scan(&title, &author, &year); err != nil {
                log.Fatal(err)
            }
            fmt.Printf("  - %s by %s (%d)\n", title, author, year)
        }
        if err := rows.Err(); err != nil {
            log.Fatal(err)
        }
        rows.Close()
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
=== Preparing an INSERT Statement Once ===
Statement prepared

=== Reusing It Many Times ===
Inserted id=1: Clean Code
Inserted id=2: Refactoring
Inserted id=3: The Mythical Man-Month
Inserted id=4: Design Patterns

=== Preparing a SELECT Statement Once ===

=== Reusing the SELECT for Different Parameters ===
Books after 1990:
  - Clean Code by Robert C. Martin (2008)
  - Refactoring by Martin Fowler (1999)
  - Design Patterns by Gang of Four (1994)
Books after 2000:
  - Clean Code by Robert C. Martin (2008)
```

**Learning Objectives:**
- ✅ Prepare a statement once with db.Prepare
- ✅ Reuse a *sql.Stmt across many Exec/Query calls with different parameters
- ✅ Know when prepared statements help (repeated statement text) vs. when they don't (one-off queries)

---

## Exercise 6: A Successful Transaction

**Objective:** Commit a multi-step write as one atomic unit

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level29-exercise6
cd ~/projects/level29-exercise6
go mod init level29.example/exercise6
go get modernc.org/sqlite@v1.57.0
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "database/sql"
    "fmt"
    "log"

    _ "modernc.org/sqlite"
)

func printState(db *sql.DB, label string) {
    var stock int
    if err := db.QueryRow(`SELECT stock FROM books WHERE id = 1`).Scan(&stock); err != nil {
        log.Fatal(err)
    }
    var orderCount int
    if err := db.QueryRow(`SELECT COUNT(*) FROM orders`).Scan(&orderCount); err != nil {
        log.Fatal(err)
    }
    fmt.Printf("[%s] stock=%d, orders=%d\n", label, stock, orderCount)
}

func main() {
    db, err := sql.Open("sqlite", "file::memory:?cache=shared")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()
    db.SetMaxOpenConns(1)

    if _, err := db.Exec(`CREATE TABLE books (
        id     INTEGER PRIMARY KEY AUTOINCREMENT,
        title  TEXT NOT NULL,
        stock  INTEGER NOT NULL
    );`); err != nil {
        log.Fatal(err)
    }
    if _, err := db.Exec(`CREATE TABLE orders (
        id      INTEGER PRIMARY KEY AUTOINCREMENT,
        book_id INTEGER NOT NULL,
        quantity INTEGER NOT NULL
    );`); err != nil {
        log.Fatal(err)
    }
    if _, err := db.Exec(`INSERT INTO books (title, stock) VALUES (?, ?)`, "Clean Code", 5); err != nil {
        log.Fatal(err)
    }

    printState(db, "before")

    fmt.Println("\n=== Placing an Order for 2 Copies (Successful Transaction) ===")
    tx, err := db.Begin()
    if err != nil {
        log.Fatal(err)
    }
    committed := false
    defer func() {
        if !committed {
            tx.Rollback()
        }
    }()

    if _, err := tx.Exec(`UPDATE books SET stock = stock - ? WHERE id = 1`, 2); err != nil {
        log.Fatal(err)
    }
    if _, err := tx.Exec(`INSERT INTO orders (book_id, quantity) VALUES (?, ?)`, 1, 2); err != nil {
        log.Fatal(err)
    }

    if err := tx.Commit(); err != nil {
        log.Fatal(err)
    }
    committed = true
    fmt.Println("Transaction committed: stock decremented AND order recorded")

    printState(db, "after commit")
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
[before] stock=5, orders=0

=== Placing an Order for 2 Copies (Successful Transaction) ===
Transaction committed: stock decremented AND order recorded
[after commit] stock=3, orders=1
```

**Learning Objectives:**
- ✅ Start a transaction with db.Begin()
- ✅ Run multiple writes as one atomic unit with tx.Exec
- ✅ Use the defer-rollback-if-not-committed pattern and commit with tx.Commit()

---

## Exercise 7: A Failed Transaction (Rollback)

**Objective:** Verify that a failure partway through a transaction undoes everything

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level29-exercise7
cd ~/projects/level29-exercise7
go mod init level29.example/exercise7
go get modernc.org/sqlite@v1.57.0
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "database/sql"
    "fmt"
    "log"

    _ "modernc.org/sqlite"
)

func printState(db *sql.DB, label string) {
    var stock int
    if err := db.QueryRow(`SELECT stock FROM books WHERE id = 1`).Scan(&stock); err != nil {
        log.Fatal(err)
    }
    var orderCount int
    if err := db.QueryRow(`SELECT COUNT(*) FROM orders`).Scan(&orderCount); err != nil {
        log.Fatal(err)
    }
    fmt.Printf("[%s] stock=%d, orders=%d\n", label, stock, orderCount)
}

// placeOrder runs entirely inside one transaction. If any step fails, the
// deferred Rollback runs BEFORE this function returns, so by the time the
// caller looks at the database again, the transaction is fully closed out.
func placeOrder(db *sql.DB, bookID, quantity int) error {
    tx, err := db.Begin()
    if err != nil {
        return err
    }
    committed := false
    defer func() {
        if !committed {
            tx.Rollback()
        }
    }()

    if _, err := tx.Exec(`UPDATE books SET stock = stock - ? WHERE id = ?`, quantity, bookID); err != nil {
        return err
    }
    if _, err := tx.Exec(`INSERT INTO orders (book_id, quantity) VALUES (?, ?)`, bookID, quantity); err != nil {
        return err
    }

    if err := tx.Commit(); err != nil {
        return err
    }
    committed = true
    return nil
}

func main() {
    db, err := sql.Open("sqlite", "file::memory:?cache=shared")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()
    // In-memory SQLite needs exactly one connection: every new pooled
    // connection would otherwise be free to see a different, uncommitted
    // view of the shared-cache database.
    db.SetMaxOpenConns(1)

    if _, err := db.Exec(`CREATE TABLE books (
        id     INTEGER PRIMARY KEY AUTOINCREMENT,
        title  TEXT NOT NULL,
        stock  INTEGER NOT NULL CHECK (stock >= 0)
    );`); err != nil {
        log.Fatal(err)
    }
    if _, err := db.Exec(`CREATE TABLE orders (
        id       INTEGER PRIMARY KEY AUTOINCREMENT,
        book_id  INTEGER NOT NULL,
        quantity INTEGER NOT NULL
    );`); err != nil {
        log.Fatal(err)
    }
    if _, err := db.Exec(`INSERT INTO books (title, stock) VALUES (?, ?)`, "Clean Code", 5); err != nil {
        log.Fatal(err)
    }

    printState(db, "before")

    fmt.Println("\n=== Attempting to Order 10 Copies When Only 5 Are In Stock (Failing Transaction) ===")
    err = placeOrder(db, 1, 10)
    if err != nil {
        // The CHECK(stock >= 0) constraint rejects the UPDATE: stock would
        // go negative. placeOrder's deferred Rollback already undid
        // everything in this transaction before returning.
        fmt.Println("Transaction failed and was rolled back:", err)
    }

    printState(db, "after")
    fmt.Println("\nStock and order count are UNCHANGED from before - the partially-applied UPDATE never survived the rollback.")
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
[before] stock=5, orders=0

=== Attempting to Order 10 Copies When Only 5 Are In Stock (Failing Transaction) ===
Transaction failed and was rolled back: constraint failed: CHECK constraint failed: stock >= 0 (275)
[after] stock=5, orders=0

Stock and order count are UNCHANGED from before - the partially-applied UPDATE never survived the rollback.
```

**Learning Objectives:**
- ✅ See a real constraint violation trigger a transaction failure
- ✅ Verify the deferred Rollback undoes every statement in the transaction
- ✅ Understand why a transaction's own function boundary matters for the defer-rollback pattern

---

## Exercise 8: Context-Aware Queries

**Objective:** Use QueryContext/ExecContext with timeouts and cancellation

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level29-exercise8
cd ~/projects/level29-exercise8
go mod init level29.example/exercise8
go get modernc.org/sqlite@v1.57.0
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "context"
    "database/sql"
    "errors"
    "fmt"
    "log"
    "time"

    _ "modernc.org/sqlite"
)

func main() {
    db, err := sql.Open("sqlite", "file::memory:?cache=shared")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()
    db.SetMaxOpenConns(1)

    if _, err := db.Exec(`CREATE TABLE books (
        id     INTEGER PRIMARY KEY AUTOINCREMENT,
        title  TEXT NOT NULL,
        author TEXT NOT NULL,
        year   INTEGER NOT NULL
    );`); err != nil {
        log.Fatal(err)
    }
    if _, err := db.Exec(`INSERT INTO books (title, author, year) VALUES (?, ?, ?)`,
        "Clean Code", "Robert C. Martin", 2008); err != nil {
        log.Fatal(err)
    }

    fmt.Println("=== Query With a Generous Timeout (Succeeds) ===")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    var title string
    err = db.QueryRowContext(ctx, `SELECT title FROM books WHERE id = ?`, 1).Scan(&title)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Found:", title)

    fmt.Println("\n=== Query With an Already-Expired Context (Fails) ===")
    shortCtx, shortCancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
    defer shortCancel()
    // Give the deadline time to pass before we even issue the query.
    time.Sleep(10 * time.Millisecond)

    err = db.QueryRowContext(shortCtx, `SELECT title FROM books WHERE id = ?`, 1).Scan(&title)
    if err != nil {
        if errors.Is(err, context.DeadlineExceeded) {
            fmt.Println("Query canceled: context deadline exceeded")
        } else {
            fmt.Println("Query failed with a different error:", err)
        }
    } else {
        fmt.Println("Unexpectedly succeeded:", title)
    }

    fmt.Println("\n=== ExecContext Respects Cancellation Too ===")
    cancelCtx, cancelNow := context.WithCancel(context.Background())
    cancelNow() // cancel before the call
    _, err = db.ExecContext(cancelCtx, `INSERT INTO books (title, author, year) VALUES (?, ?, ?)`,
        "Too Late", "Nobody", 2024)
    if err != nil {
        if errors.Is(err, context.Canceled) {
            fmt.Println("Insert canceled: context canceled")
        } else {
            fmt.Println("Insert failed with a different error:", err)
        }
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
=== Query With a Generous Timeout (Succeeds) ===
Found: Clean Code

=== Query With an Already-Expired Context (Fails) ===
Query canceled: context deadline exceeded

=== ExecContext Respects Cancellation Too ===
Insert canceled: context canceled
```

**Learning Objectives:**
- ✅ Pass a context.Context to QueryRowContext/ExecContext
- ✅ See a context that has already expired fail a query with context.DeadlineExceeded
- ✅ See a canceled context stop an ExecContext with context.Canceled

---

## Exercise 9: Connection Pool Settings

**Objective:** Configure and observe the connection pool with db.Stats()

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level29-exercise9
cd ~/projects/level29-exercise9
go mod init level29.example/exercise9
go get modernc.org/sqlite@v1.57.0
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "database/sql"
    "fmt"
    "log"
    "time"

    _ "modernc.org/sqlite"
)

func printStats(db *sql.DB, label string) {
    s := db.Stats()
    fmt.Printf("[%s] OpenConnections=%d InUse=%d Idle=%d MaxOpenConnections=%d\n",
        label, s.OpenConnections, s.InUse, s.Idle, s.MaxOpenConnections)
}

func main() {
    db, err := sql.Open("sqlite", "file::memory:?cache=shared")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    fmt.Println("=== Default Pool Settings ===")
    printStats(db, "before any config")

    fmt.Println("\n=== Applying Sensible Limits ===")
    // In-memory SQLite must use exactly one connection - a second connection
    // would see its own separate database. Real server-backed databases
    // (Postgres, MySQL) would instead pick these based on server capacity.
    db.SetMaxOpenConns(1)
    db.SetMaxIdleConns(1)
    db.SetConnMaxLifetime(5 * time.Minute)
    db.SetConnMaxIdleTime(1 * time.Minute)
    fmt.Println("SetMaxOpenConns(1)")
    fmt.Println("SetMaxIdleConns(1)")
    fmt.Println("SetConnMaxLifetime(5m)")
    fmt.Println("SetConnMaxIdleTime(1m)")

    if _, err := db.Exec(`CREATE TABLE books (
        id     INTEGER PRIMARY KEY AUTOINCREMENT,
        title  TEXT NOT NULL
    );`); err != nil {
        log.Fatal(err)
    }

    printStats(db, "after CREATE TABLE (1 connection opened on demand)")

    if _, err := db.Exec(`INSERT INTO books (title) VALUES (?)`, "Clean Code"); err != nil {
        log.Fatal(err)
    }
    var count int
    if err := db.QueryRow(`SELECT COUNT(*) FROM books`).Scan(&count); err != nil {
        log.Fatal(err)
    }
    fmt.Println("\nBooks in table:", count)

    printStats(db, "after INSERT + SELECT (same connection reused)")
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Default Pool Settings ===
[before any config] OpenConnections=0 InUse=0 Idle=0 MaxOpenConnections=0

=== Applying Sensible Limits ===
SetMaxOpenConns(1)
SetMaxIdleConns(1)
SetConnMaxLifetime(5m)
SetConnMaxIdleTime(1m)
[after CREATE TABLE (1 connection opened on demand)] OpenConnections=1 InUse=0 Idle=1 MaxOpenConnections=1

Books in table: 1
[after INSERT + SELECT (same connection reused)] OpenConnections=1 InUse=0 Idle=1 MaxOpenConnections=1
```

**Learning Objectives:**
- ✅ Configure SetMaxOpenConns, SetMaxIdleConns, SetConnMaxLifetime, SetConnMaxIdleTime
- ✅ Read db.Stats() to see the pool's real, live state
- ✅ Understand that the same connection is reused across sequential calls once opened

---

## Exercise 10: Comprehensive Practice - A Tested Bookstore Data Layer

**Objective:** Build a small, fully-tested CRUD data layer, bridging directly into Level 30's repository pattern

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level29-exercise10
cd ~/projects/level29-exercise10
go mod init level29.example/exercise10
go get modernc.org/sqlite@v1.57.0
```

2. Create `bookstore.go`:

```bash
cat > bookstore.go << 'EOF'
package main

import (
    "database/sql"
    "errors"
    "fmt"

    _ "modernc.org/sqlite"
)

// Book is the domain model used throughout the bookstore data layer.
type Book struct {
    ID     int64
    Title  string
    Author string
    Year   int
}

// ErrBookNotFound is returned by Update/Delete when no row matches the id.
var ErrBookNotFound = errors.New("book not found")

// OpenBookstoreDB opens an in-memory SQLite database configured for
// single-connection use (required for :memory: databases) and creates the
// books table if it doesn't already exist.
func OpenBookstoreDB(dsn string) (*sql.DB, error) {
    db, err := sql.Open("sqlite", dsn)
    if err != nil {
        return nil, err
    }
    db.SetMaxOpenConns(1)

    if err := db.Ping(); err != nil {
        db.Close()
        return nil, err
    }

    schema := `
    CREATE TABLE IF NOT EXISTS books (
        id     INTEGER PRIMARY KEY AUTOINCREMENT,
        title  TEXT NOT NULL,
        author TEXT NOT NULL,
        year   INTEGER NOT NULL
    );`
    if _, err := db.Exec(schema); err != nil {
        db.Close()
        return nil, err
    }
    return db, nil
}

// CreateBook inserts a new book and returns its generated id.
func CreateBook(db *sql.DB, b Book) (int64, error) {
    result, err := db.Exec(
        `INSERT INTO books (title, author, year) VALUES (?, ?, ?)`,
        b.Title, b.Author, b.Year,
    )
    if err != nil {
        return 0, err
    }
    return result.LastInsertId()
}

// GetBook fetches a single book by id. Returns ErrBookNotFound (wrapping
// sql.ErrNoRows) when no row matches.
func GetBook(db *sql.DB, id int64) (Book, error) {
    var b Book
    err := db.QueryRow(
        `SELECT id, title, author, year FROM books WHERE id = ?`, id,
    ).Scan(&b.ID, &b.Title, &b.Author, &b.Year)
    if errors.Is(err, sql.ErrNoRows) {
        return Book{}, fmt.Errorf("get book %d: %w", id, ErrBookNotFound)
    }
    if err != nil {
        return Book{}, err
    }
    return b, nil
}

// ListBooks returns every book ordered by id.
func ListBooks(db *sql.DB) ([]Book, error) {
    rows, err := db.Query(`SELECT id, title, author, year FROM books ORDER BY id`)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var books []Book
    for rows.Next() {
        var b Book
        if err := rows.Scan(&b.ID, &b.Title, &b.Author, &b.Year); err != nil {
            return nil, err
        }
        books = append(books, b)
    }
    if err := rows.Err(); err != nil {
        return nil, err
    }
    return books, nil
}

// UpdateBook updates an existing book's fields by id. Returns
// ErrBookNotFound if no row had that id.
func UpdateBook(db *sql.DB, b Book) error {
    result, err := db.Exec(
        `UPDATE books SET title = ?, author = ?, year = ? WHERE id = ?`,
        b.Title, b.Author, b.Year, b.ID,
    )
    if err != nil {
        return err
    }
    affected, err := result.RowsAffected()
    if err != nil {
        return err
    }
    if affected == 0 {
        return fmt.Errorf("update book %d: %w", b.ID, ErrBookNotFound)
    }
    return nil
}

// DeleteBook removes a book by id. Returns ErrBookNotFound if no row had
// that id.
func DeleteBook(db *sql.DB, id int64) error {
    result, err := db.Exec(`DELETE FROM books WHERE id = ?`, id)
    if err != nil {
        return err
    }
    affected, err := result.RowsAffected()
    if err != nil {
        return err
    }
    if affected == 0 {
        return fmt.Errorf("delete book %d: %w", id, ErrBookNotFound)
    }
    return nil
}

func main() {
    db, err := OpenBookstoreDB("file::memory:?cache=shared")
    if err != nil {
        panic(err)
    }
    defer db.Close()

    id, err := CreateBook(db, Book{Title: "Clean Code", Author: "Robert C. Martin", Year: 2008})
    if err != nil {
        panic(err)
    }
    fmt.Println("Created book with id:", id)

    b, err := GetBook(db, id)
    if err != nil {
        panic(err)
    }
    fmt.Printf("Fetched: %+v\n", b)

    books, err := ListBooks(db)
    if err != nil {
        panic(err)
    }
    fmt.Println("Total books:", len(books))
}
EOF
```

3. Create `bookstore_test.go`:

```bash
cat > bookstore_test.go << 'EOF'
package main

import (
    "errors"
    "testing"
)

// Each test uses its own uniquely-named in-memory database (via the DSN)
// so tests never see each other's data, while still following the
// single-connection, shared-cache pattern SQLite's :memory: mode requires.
func TestBookstoreCRUD(t *testing.T) {
    db, err := OpenBookstoreDB("file:crud?mode=memory&cache=shared")
    if err != nil {
        t.Fatalf("OpenBookstoreDB: %v", err)
    }
    t.Cleanup(func() { db.Close() })

    var id int64

    t.Run("Create", func(t *testing.T) {
        id, err = CreateBook(db, Book{Title: "Clean Code", Author: "Robert C. Martin", Year: 2008})
        if err != nil {
            t.Fatalf("CreateBook: %v", err)
        }
        if id == 0 {
            t.Fatal("expected a non-zero id")
        }
    })

    t.Run("Get", func(t *testing.T) {
        b, err := GetBook(db, id)
        if err != nil {
            t.Fatalf("GetBook: %v", err)
        }
        if b.Title != "Clean Code" || b.Author != "Robert C. Martin" || b.Year != 2008 {
            t.Fatalf("unexpected book: %+v", b)
        }
    })

    t.Run("List", func(t *testing.T) {
        if _, err := CreateBook(db, Book{Title: "Refactoring", Author: "Martin Fowler", Year: 1999}); err != nil {
            t.Fatalf("CreateBook: %v", err)
        }
        books, err := ListBooks(db)
        if err != nil {
            t.Fatalf("ListBooks: %v", err)
        }
        if len(books) != 2 {
            t.Fatalf("expected 2 books, got %d", len(books))
        }
    })

    t.Run("Update", func(t *testing.T) {
        err := UpdateBook(db, Book{ID: id, Title: "Clean Code (2nd ed.)", Author: "Robert C. Martin", Year: 2008})
        if err != nil {
            t.Fatalf("UpdateBook: %v", err)
        }
        b, err := GetBook(db, id)
        if err != nil {
            t.Fatalf("GetBook: %v", err)
        }
        if b.Title != "Clean Code (2nd ed.)" {
            t.Fatalf("expected updated title, got %q", b.Title)
        }
    })

    t.Run("UpdateMissing", func(t *testing.T) {
        err := UpdateBook(db, Book{ID: 9999, Title: "Ghost", Author: "Nobody", Year: 2000})
        if !errors.Is(err, ErrBookNotFound) {
            t.Fatalf("expected ErrBookNotFound, got %v", err)
        }
    })

    t.Run("Delete", func(t *testing.T) {
        if err := DeleteBook(db, id); err != nil {
            t.Fatalf("DeleteBook: %v", err)
        }
        _, err := GetBook(db, id)
        if !errors.Is(err, ErrBookNotFound) {
            t.Fatalf("expected ErrBookNotFound after delete, got %v", err)
        }
    })

    t.Run("DeleteMissing", func(t *testing.T) {
        err := DeleteBook(db, 9999)
        if !errors.Is(err, ErrBookNotFound) {
            t.Fatalf("expected ErrBookNotFound, got %v", err)
        }
    })
}
EOF
```

4. Run the demo program, then the tests:

```bash
go run bookstore.go
go test -v ./...
```

**Expected Output (go run bookstore.go):**

```
Created book with id: 1
Fetched: {ID:1 Title:Clean Code Author:Robert C. Martin Year:2008}
Total books: 1
```

**Expected Output (go test -v ./...):**

```
=== RUN   TestBookstoreCRUD
=== RUN   TestBookstoreCRUD/Create
=== RUN   TestBookstoreCRUD/Get
=== RUN   TestBookstoreCRUD/List
=== RUN   TestBookstoreCRUD/Update
=== RUN   TestBookstoreCRUD/UpdateMissing
=== RUN   TestBookstoreCRUD/Delete
=== RUN   TestBookstoreCRUD/DeleteMissing
--- PASS: TestBookstoreCRUD (0.00s)
    --- PASS: TestBookstoreCRUD/Create (0.00s)
    --- PASS: TestBookstoreCRUD/Get (0.00s)
    --- PASS: TestBookstoreCRUD/List (0.00s)
    --- PASS: TestBookstoreCRUD/Update (0.00s)
    --- PASS: TestBookstoreCRUD/UpdateMissing (0.00s)
    --- PASS: TestBookstoreCRUD/Delete (0.00s)
    --- PASS: TestBookstoreCRUD/DeleteMissing (0.00s)
PASS
ok  	level29.example/exercise10	0.594s
```

(Your own `go test` run will show the same PASS lines; the exact tenths-of-a-second timing will vary run to run and machine to machine - that part is never worth comparing.)

**Learning Objectives:**
- ✅ Build a small repository-shaped data layer (Create/Get/List/Update/Delete) around *sql.DB
- ✅ Wrap sql.ErrNoRows into a package-level sentinel error with %w
- ✅ Use RowsAffected() to detect "update/delete matched nothing"
- ✅ Write real go test subtests (Level 26 callback) that exercise every function against a real in-memory database
- ✅ See exactly the shape Level 30 (Repository-Service-Handler) will formalize into its own layer

---

## Bonus Challenges

### Challenge 1: Migration Runner

Build a tiny migration system: a `[]string` slice of SQL statements applied in order, tracked in a `schema_migrations` table so re-running the program never re-applies a migration twice.

```bash
mkdir -p ~/projects/level29-bonus1
cd ~/projects/level29-bonus1
go mod init level29.example/bonus1
go get modernc.org/sqlite@v1.57.0
```

**Hints:**
- Create a `schema_migrations (version INTEGER PRIMARY KEY)` table first, if it doesn't exist
- For each migration, check whether its version number is already recorded before running it
- After a migration's SQL succeeds, insert its version into `schema_migrations` in the same pass
- Run your program twice in a row - the second run should apply zero new migrations

### Challenge 2: Connection Pool Exhaustion

Set `db.SetMaxOpenConns(1)` and launch several goroutines that each call `ExecContext` with a `context.WithTimeout`. Observe that only one goroutine's statement runs at a time - the rest wait for the single connection to free up - and that a short enough timeout can make a waiting goroutine fail with `context.DeadlineExceeded` before it ever gets a turn.

```bash
mkdir -p ~/projects/level29-bonus2
cd ~/projects/level29-bonus2
go mod init level29.example/bonus2
go get modernc.org/sqlite@v1.57.0
```

**Hints:**
- Use a `sync.WaitGroup` to launch several goroutines at once
- Give each goroutine its own `context.WithTimeout` and print how long its `ExecContext` call took
- Try a generous timeout first (everything succeeds, just serialized), then a very short one (some calls fail)

### Challenge 3: A Repository Wrapper Struct

Wrap `*sql.DB` in a `BookRepository` struct with methods instead of free functions - `func (r *BookRepository) Create(b Book) (int64, error)`, `Get`, `List`, `Update`, `Delete` - built on top of Exercise 10's SQL. This is exactly the shape Level 30 (Repository-Service-Handler) builds on.

```bash
mkdir -p ~/projects/level29-bonus3
cd ~/projects/level29-bonus3
go mod init level29.example/bonus3
go get modernc.org/sqlite@v1.57.0
```

**Hints:**
- `type BookRepository struct { db *sql.DB }` with a `NewBookRepository(db *sql.DB) *BookRepository` constructor
- Move Exercise 10's SQL directly into methods on `*BookRepository` - the SQL itself doesn't need to change, only where it lives
- Define a `BookStore` interface with the same method set, so a service layer (Level 30) can depend on the interface instead of the concrete struct

---

## What You've Learned

After completing these 10 exercises, you can:

✅ Connect to a database through a driver, understanding sql.Open's laziness and sql.DB's pooling
✅ Create a schema and understand the SQLite in-memory + connection-pool interaction
✅ Query multiple rows into structs and a single row with sql.ErrNoRows handling
✅ Write data safely with parameterized queries, proven against hostile-looking input
✅ Reuse prepared statements for repeated SQL
✅ Run multi-step writes atomically with transactions, verified both committed and rolled back
✅ Make queries cancelable and timeout-aware with context.Context
✅ Configure and observe a connection pool
✅ Build and fully test a small CRUD data layer

---

## Next Level

Level 30: Repository-Service-Handler
- Structuring a real application into repository, service, and handler layers
- Turning Exercise 10's free functions into a proper repository interface
- Wiring that repository into Gin (Level 28) handlers through a service layer

Great work! You can now make your programs remember things! 🚀
