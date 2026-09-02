# Level 18: File Handling - INDEX

Welcome to **Level 18: File Handling**! This is where your programs stop living only in memory and start reading from, and writing to, the filesystem.

---

## 📖 What You'll Learn

- ✅ io.Reader/io.Writer, the two interfaces underlying almost all of Go's I/O
- ✅ Reading and writing whole files, plus streaming large ones
- ✅ File permissions and what the octal digits mean
- ✅ Buffered I/O with bufio.Scanner, bufio.Reader, and bufio.Writer
- ✅ Appending to files without truncating them
- ✅ Building, listing, and walking directory trees
- ✅ Temporary files/directories and safe scratch-space patterns
- ✅ File metadata with os.Stat and safe deletion

---

## 🗂️ Level 18 Materials

### 1. **START_HERE.txt** (Begin Here!)
Quick overview and welcome message. Start here first!

### 2. **README.md** (Main Theory)
Comprehensive guide to:
- io.Reader/io.Writer as Go's universal I/O abstraction
- Opening, reading, writing, and appending to files
- File permissions (os.FileMode)
- Buffered I/O and why it matters
- Directories, temp files, metadata, and deletion
- Best practices
- Common mistakes

**Read Time:** 45-60 minutes
**When:** Start your session

### 3. **EXERCISES.md** (Hands-On Practice)
10 detailed, step-by-step exercises:
1. Writing and reading a file
2. Handling a nonexistent file
3. File permissions
4. Reading line by line with bufio.Scanner
5. Buffered reading and writing
6. Appending across multiple writes
7. Building and listing a directory tree
8. Walking a directory tree
9. Temporary files, directories, and metadata
10. Comprehensive practice (rotating log file)

**Time Commitment:** 5-6 hours
**When:** After reading README

### 4. **STUDY_GUIDE.md** (Visual Learning)
- io.Reader/io.Writer interface diagram
- File open flags table
- Permission octal breakdown table
- os vs bufio function comparison table
- Common mistakes guide

**Read Time:** 35-45 minutes
**When:** Use while doing exercises

### 5. **QUICK_REFERENCE.md** (Cheat Sheet)
- Essential syntax reference
- Permissions, buffered I/O, directories at a glance
- Common mistakes table

**Read Time:** 10-15 minutes
**When:** Quick lookup during practice

---

## 🎯 Recommended Learning Path

### Day 1: Reading, Writing & Errors (2 hours)
1. Read **START_HERE.txt** (10 min)
2. Read **README.md** sections 1-3 (35 min)
3. Complete Exercises 1-3 (1 hour 15 min)

### Day 2: Buffered I/O & Appending (2 hours)
1. Read **README.md** sections 4-5 (25 min)
2. Complete Exercises 4-6 (1.5 hours)

### Day 3: Directories & Trees (1.5 hours)
1. Read **README.md** section 6 (20 min)
2. Complete Exercises 7-8 (1 hour 10 min)

### Day 4: Temp Files, Metadata & Real-World Practice (2 hours)
1. Read **README.md** sections 7-9 (25 min)
2. Complete Exercises 9-10 (1.5 hours)

### Day 5: Consolidation (1 hour)
1. Try bonus challenges
2. Review with **QUICK_REFERENCE.md**
3. Make sure os.IsNotExist and filepath.Join feel automatic

---

## 💡 Key Concepts At A Glance

### io.Reader / io.Writer
```go
type Reader interface { Read(p []byte) (n int, err error) }
type Writer interface { Write(p []byte) (n int, err error) }
```

### Reading & Writing
```go
data, err := os.ReadFile("path.txt")
err = os.WriteFile("path.txt", data, 0644)
```

### The os.IsNotExist Idiom
```go
if _, err := os.Stat(path); err != nil {
    if os.IsNotExist(err) { /* not there yet */ }
}
```

### Buffered Line Reading
```go
scanner := bufio.NewScanner(file)
for scanner.Scan() { line := scanner.Text() }
```

### Portable Paths
```go
path := filepath.Join(dir, "sub", "file.txt") // never dir + "/" + "file.txt"
```

---

## ✅ Prerequisites

Make sure you've completed **Level 17: Packages & Modules**

You need:
- ✅ Comfort organizing code into packages and modules
- ✅ Comfort with defer (Level 11) and error handling (Level 16)
- ✅ Comfort with interfaces (Level 15) - io.Reader/io.Writer build directly on that idea

---

## 🎓 Learning Objectives

By the end of Level 18, you'll be able to:

- ✅ Explain why io.Reader/io.Writer underlie almost all of Go's I/O
- ✅ Read and write whole files, and stream large ones line by line
- ✅ Set and interpret file permissions
- ✅ Use buffered I/O correctly, including remembering to Flush()
- ✅ Append to files without destroying existing content
- ✅ Build directory trees, list them, and walk them recursively
- ✅ Create safe scratch files/directories with proper cleanup
- ✅ Check file existence/metadata with os.Stat and os.IsNotExist
- ✅ Delete files and directories safely, understanding the danger of os.RemoveAll

---

## 📊 Statistics

- **Main Theory:** README.md covering I/O interfaces, reading, writing, buffering, directories, and cleanup
- **Exercises:** 10 detailed exercises + 3 bonus challenges
- **Study Guide:** interface diagram, flags table, permissions table, pattern library
- **Quick Reference:** One-page cheat sheet
- **Practice Time:** 5-6 hours
- **Difficulty:** ⭐⭐⭐ (Intermediate-Advanced)

---

## 🚀 How to Use These Materials

### For Reading
1. Open README.md for comprehensive theory
2. Use STUDY_GUIDE.md alongside for the flags/permissions tables and diagrams
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
3. Practice the os.IsNotExist and filepath.Join idioms until they're automatic

---

## 🆘 Common Questions

**Q: Why does the error from a missing file look so plain, like "open x.txt: no such file or directory"?**
A: That's the real, unmodified error Go's `os` package returns - it wraps the operating system's own error. Instead of matching that text, use `os.IsNotExist(err)` to check the underlying cause portably.

**Q: When should I use os.ReadFile vs os.Open?**
A: `os.ReadFile` when the whole file comfortably fits in memory and you just want its bytes. `os.Open` (plus `bufio`) when the file could be large, or you need to process it line by line rather than all at once.

**Q: Why do I need filepath.Join instead of just concatenating strings with "/"?**
A: `filepath.Join` uses the correct separator for the current OS and cleans up redundant slashes. Hardcoding `"/"` breaks on Windows and is fragile even on Unix-like systems.

**Q: What's the real danger with os.RemoveAll?**
A: It recursively deletes a path and everything inside it with no confirmation. If the path is built incorrectly (e.g., from an empty variable), it can delete far more than intended. Always validate the path first.

**Q: Do I need to call Flush() on every bufio.Writer?**
A: Yes - buffered data isn't guaranteed to reach the file until you call `Flush()` (or the internal buffer fills on its own). Always flush before closing the file.

---

## 🎯 Before Moving to Level 19

Make sure you can answer these questions:

- [ ] What are io.Reader and io.Writer, and why do so many types implement them?
- [ ] What's the difference between os.ReadFile and os.Open?
- [ ] How do you correctly detect that a file doesn't exist?
- [ ] What do the three octal digits in a permission like 0644 each mean?
- [ ] Why does appending require os.OpenFile instead of os.WriteFile?
- [ ] Why is os.RemoveAll dangerous, and how do you use it safely?

---

## 📚 Recommended Reading Order

For Maximum Learning:

1. **START_HERE.txt** (10 min)
   - Gets you oriented

2. **README.md** Sections 1-4 (30 min)
   - The Reader/Writer interfaces, reading, writing, buffering

3. **README.md** Sections 5-9 (35 min)
   - Appending, directories, temp files, metadata, deletion

4. **EXERCISES.md** Exercises 1-5 (2.5 hours)
   - Get hands-on with reading, writing, permissions, and buffered I/O

5. **STUDY_GUIDE.md** (35 min)
   - Study the flags table, permissions table, and pattern library

6. **EXERCISES.md** Exercises 6-10 (2.5+ hours)
   - Deep practice with appending, directories, temp files, and a real log rotator

7. **QUICK_REFERENCE.md** (10 min)
   - Create your mental cheat sheet

---

## 🏆 Success Indicators

You've mastered Level 18 when:

- ✅ You reach for os.ReadFile/os.WriteFile by default for small files
- ✅ You can explain why os.IsNotExist is safer than matching error text
- ✅ You never string-concatenate a path when filepath.Join is available
- ✅ You remember to Flush() a bufio.Writer before closing
- ✅ You treat os.RemoveAll with real caution
- ✅ You've completed 8+ exercises
- ✅ You can explain file handling to someone else

---

## 🚀 What's Next?

After Level 18, you're ready for:

**Level 19: JSON**
- Encoding Go values to JSON and decoding JSON into Go values
- Struct tags for controlling field names
- Working with nested and dynamic JSON data

---

## 💬 Key Takeaway

> **Almost every file operation in Go boils down to io.Reader/io.Writer underneath - learn the convenience functions (os.ReadFile/WriteFile) for the common case, and drop to os.Open/Create plus bufio when you need streaming, appending, or line-by-line control.**

---

## 📖 Quick Links

- [START_HERE.txt](./START_HERE.txt) - Welcome & Overview
- [README.md](./README.md) - Full Theory
- [EXERCISES.md](./EXERCISES.md) - Hands-On Exercises
- [STUDY_GUIDE.md](./STUDY_GUIDE.md) - Visual Patterns
- [QUICK_REFERENCE.md](./QUICK_REFERENCE.md) - Cheat Sheet

---

Happy learning! Level 18 gives your programs the power to persist data beyond a single run! 🎉

*Estimated time to complete Level 18: 5-6 hours*
*Difficulty: ⭐⭐⭐ (Intermediate-Advanced)*
*Next Level: Level 19 - JSON*
