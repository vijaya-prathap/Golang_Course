# Level 29: Database/SQL - INDEX

Welcome to **Level 29: Database/SQL**! This is where your programs stop losing their data the moment they restart - Go's standard library way to talk to a real relational database.

**Setup note:** this level (along with Level 28: Gin Framework) is the one place in the course that needs internet access once, to fetch the `modernc.org/sqlite` driver (pinned at `v1.57.0`). After that single `go get`, every exercise runs against an in-memory SQLite database with zero network or server dependency - exactly like every other level.

---

## 📖 What You'll Learn

- ✅ database/sql as an interface over drivers, not a database itself
- ✅ sql.Open's laziness, sql.DB as a connection pool, and db.Ping()
- ✅ Creating a schema with db.Exec
- ✅ Querying multiple rows into structs, and a single row with sql.ErrNoRows
- ✅ Writing data safely with parameterized queries - proven against hostile input
- ✅ Prepared statements for repeated SQL
- ✅ Transactions, verified both committed and rolled back
- ✅ Context-aware queries bridging to Level 24
- ✅ Connection pool configuration and db.Stats()

---

## 🗂️ Level 29 Materials

### 1. **START_HERE.txt** (Begin Here!)
Quick overview and welcome message. Start here first!

### 2. **README.md** (Main Theory)
Comprehensive guide to:
- database/sql, drivers, and blank imports
- Connecting, the connection pool, and the SQLite in-memory gotcha
- Schema setup, querying rows and a single row
- Writing data, parameterized queries, prepared statements
- Transactions and context-aware queries
- Best practices
- Common mistakes

**Read Time:** 55-70 minutes
**When:** Start your session

### 3. **EXERCISES.md** (Hands-On Practice)
10 detailed, step-by-step exercises:
1. Connecting and creating a schema
2. Inserting data and querying multiple rows into structs
3. Querying a single row and handling sql.ErrNoRows
4. Preventing SQL injection - parameterized queries
5. Prepared statements
6. A successful transaction
7. A failed transaction (rollback)
8. Context-aware queries
9. Connection pool settings
10. Comprehensive practice - a tested bookstore data layer

Plus 3 bonus challenges: a migration runner, connection pool exhaustion, and a repository wrapper struct previewing Level 30.

**Time Commitment:** 5-6 hours
**When:** After reading README

### 4. **STUDY_GUIDE.md** (Visual Learning)
- sql.DB-as-connection-pool diagram
- Query vs. QueryRow vs. Exec decision table
- The rows.Next()/Scan()/Close()/Err() lifecycle diagram
- Transaction commit/rollback flow diagram
- Parameterized query safety illustration

**Read Time:** 35-45 minutes
**When:** Use while doing exercises

### 5. **QUICK_REFERENCE.md** (Cheat Sheet)
- Essential syntax reference
- Query/QueryRow/Exec at a glance
- Transaction pattern and pool settings
- Common mistakes table

**Read Time:** 10-15 minutes
**When:** Quick lookup during practice

---

## 🎯 Recommended Learning Path

### Day 1: Connecting and Reading Data (2.5 hours)
1. Read **START_HERE.txt** (10 min)
2. Read **README.md** sections 1-5 (35 min)
3. Complete Exercises 1-3 (1 hour 45 min)

### Day 2: Writing Data Safely (2 hours)
1. Read **README.md** sections 6-8 (25 min)
2. Complete Exercises 4-5 (1.5 hours)

### Day 3: Transactions (1.5 hours)
1. Read **README.md** section 9 (15 min)
2. Complete Exercises 6-7 (1 hour 15 min)

### Day 4: Context and Pools (1.5 hours)
1. Read **README.md** sections 10-12 (20 min)
2. Complete Exercises 8-9 (1 hour 10 min)

### Day 5: Comprehensive Practice & Consolidation (1.5 hours)
1. Complete Exercise 10 (1 hour)
2. Try bonus challenges
3. Review with **QUICK_REFERENCE.md**

---

## 💡 Key Concepts At A Glance

### Connecting
```go
db, err := sql.Open("sqlite", "file::memory:?cache=shared")
db.SetMaxOpenConns(1) // required for in-memory SQLite
db.Ping()             // verify connectivity eagerly
```

### Reading Data
```go
rows, err := db.Query(`SELECT ...`)   // many rows
defer rows.Close()
for rows.Next() { rows.Scan(&x) }
rows.Err()

db.QueryRow(`SELECT ...`, id).Scan(&x) // zero-or-one row, watch sql.ErrNoRows
```

### Writing Data Safely
```go
result, err := db.Exec(`INSERT INTO books (title) VALUES (?)`, title) // parameterized!
```

### Transactions
```go
tx, _ := db.Begin()
defer func() { if !committed { tx.Rollback() } }()
// ... tx.Exec ...
tx.Commit()
committed = true
```

---

## ✅ Prerequisites

Make sure you've completed **Level 28: Gin Framework**

You need:
- ✅ Comfort building HTTP handlers and routes with Gin
- ✅ Comfort with context.Context (Level 24) for cancellation and timeouts
- ✅ Comfort with error handling (Level 16) and errors.Is/errors.As
- ✅ Comfort writing tests with the testing package (Level 26)

---

## 🎓 Learning Objectives

By the end of Level 29, you'll be able to:

- ✅ Explain why database/sql needs a separately-imported driver
- ✅ Open a connection pool and verify it eagerly with db.Ping()
- ✅ Query multiple rows into structs and a single row with sql.ErrNoRows handling
- ✅ Write every value into SQL through a parameterized `?` placeholder, without exception
- ✅ Reuse prepared statements for repeated SQL
- ✅ Write a transaction using the defer-rollback-if-not-committed pattern
- ✅ Pass a context.Context into a *Context database call for cancellation/timeouts
- ✅ Configure a connection pool's limits and read db.Stats()

---

## 📊 Statistics

- **Main Theory:** README.md covering drivers, pooling, reading/writing data, transactions, and context
- **Exercises:** 10 detailed exercises + 3 bonus challenges
- **Study Guide:** connection pool diagram, decision tables, lifecycle and transaction flow diagrams
- **Quick Reference:** One-page cheat sheet
- **Practice Time:** 5-6 hours
- **Difficulty:** ⭐⭐⭐⭐ (Advanced)

---

## 🚀 How to Use These Materials

### For Reading
1. Open README.md for comprehensive theory
2. Use STUDY_GUIDE.md alongside for the pool, lifecycle, and transaction diagrams
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
3. Practice the defer-rollback-if-not-committed transaction pattern until it's automatic

---

## 🆘 Common Questions

**Q: Why do I need a separate driver if database/sql is in the standard library?**
A: `database/sql` only defines the interface - `sql.DB`, `sql.Rows`, and so on. It has no code that actually speaks any database's wire protocol. A driver (here, `modernc.org/sqlite`, blank-imported) registers itself and does the real work.

**Q: Does sql.Open fail if the database is unreachable?**
A: No - `sql.Open` is lazy and just validates its arguments. Call `db.Ping()` right after if you need to know immediately whether the database is actually reachable.

**Q: Why does this level use `file::memory:?cache=shared` instead of plain `:memory:`?**
A: Because `sql.DB` is a connection pool. With plain `:memory:`, every new pooled connection gets its own separate, empty database. The shared-cache DSN plus `SetMaxOpenConns(1)` makes every call go through one connection, one consistent database.

**Q: What's the difference between db.Query and db.QueryRow?**
A: `db.Query` returns a cursor (`*sql.Rows`) for zero-to-many rows - you loop, Scan, Close, and check Err(). `db.QueryRow` is for exactly zero-or-one row - call `.Scan()` directly and check for `sql.ErrNoRows`.

**Q: How do I know a transaction actually rolled back correctly?**
A: Exercise 7 verifies it for real: a `CHECK` constraint rejects an UPDATE inside a transaction, the deferred `Rollback()` runs, and the book's stock count and order count both end up exactly where they started - not partially changed.

---

## 🎯 Before Moving to Level 30

Make sure you can answer these questions:

- [ ] Why doesn't a successful sql.Open guarantee the database is reachable?
- [ ] What's the difference between db.Query, db.QueryRow, and db.Exec?
- [ ] What is sql.ErrNoRows, and why is it usually not a fatal error?
- [ ] What's the one rule that prevents SQL injection?
- [ ] What does the defer-rollback-if-not-committed pattern protect you from?

---

## 📚 Recommended Reading Order

For Maximum Learning:

1. **START_HERE.txt** (10 min)
   - Gets you oriented

2. **README.md** Sections 1-6 (40 min)
   - The interface, connecting, schema, reading, and writing data

3. **README.md** Sections 7-12 (30 min)
   - Safety, prepared statements, transactions, context, best practices, mistakes

4. **EXERCISES.md** Exercises 1-5 (2.5 hours)
   - Get hands-on with connecting, reading, writing, and safe queries

5. **STUDY_GUIDE.md** (40 min)
   - Study the pool diagram and the rows lifecycle diagram

6. **EXERCISES.md** Exercises 6-10 (2.5+ hours)
   - Deep practice with transactions, context, pools, and the tested data layer

7. **QUICK_REFERENCE.md** (10 min)
   - Create your mental cheat sheet

---

## 🏆 Success Indicators

You've mastered Level 29 when:

- ✅ You reach for a parameterized `?` placeholder automatically, every time
- ✅ You can explain sql.Open's laziness and sql.DB's pooling without hesitation
- ✅ You never forget rows.Close() or checking rows.Err()
- ✅ You can write the defer-rollback-if-not-committed transaction pattern from memory
- ✅ You've personally seen a rollback leave data completely unchanged
- ✅ You've completed 8+ exercises, including Exercise 10's fully-tested data layer
- ✅ You can explain database/sql to someone else

---

## 🚀 What's Next?

After Level 29, you're ready for:

**Level 30: Repository-Service-Handler**
- Structuring a real application into repository, service, and handler layers
- Turning Exercise 10's free functions into a proper repository interface
- Wiring that repository into Gin (Level 28) handlers through a service layer

---

## 💬 Key Takeaway

> **database/sql is an interface over drivers, not a database itself - sql.Open is lazy, sql.DB is a pool, and the one rule with no exceptions is: every value goes through a `?` placeholder, never string concatenation.**

---

## 📖 Quick Links

- [START_HERE.txt](./START_HERE.txt) - Welcome & Overview
- [README.md](./README.md) - Full Theory
- [EXERCISES.md](./EXERCISES.md) - Hands-On Exercises
- [STUDY_GUIDE.md](./STUDY_GUIDE.md) - Visual Patterns
- [QUICK_REFERENCE.md](./QUICK_REFERENCE.md) - Cheat Sheet

---

Happy learning! Level 29 gives your programs the power to remember things, safely and atomically! 🎉

*Estimated time to complete Level 29: 5-6 hours*
*Difficulty: ⭐⭐⭐⭐ (Advanced)*
*Next Level: Level 30 - Repository-Service-Handler*
