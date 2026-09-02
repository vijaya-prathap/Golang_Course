# Level 18: Study Guide & Visual Reference

## 📚 Learning Path

### Week 1: File Handling Fundamentals
```
Day 1:  io.Reader / io.Writer, opening and reading files
Day 2:  Writing files and understanding permissions
Day 3:  Buffered I/O (bufio.Scanner, bufio.Reader/Writer)
Day 4:  Appending to files
Day 5:  Directories - creating, listing, joining paths
Day 6:  Walking directory trees, temp files
Day 7:  Metadata (os.Stat) and safe deletion
```

### Week 2: Practice & Application
```
Day 1:  Exercises 1-3 (read/write, missing files, permissions)
Day 2:  Exercises 4-6 (bufio.Scanner, buffered I/O, appending)
Day 3:  Exercises 7-8 (directories, WalkDir)
Day 4:  Exercises 9-10 (temp files, comprehensive log rotator)
Day 5:  Bonus challenges
Day 6-7: Practice & review
```

---

## 🎯 io.Reader / io.Writer: The Universal Interface

```
        ┌───────────────────────────┐
        │   io.Reader interface     │
        │   Read(p []byte)          │
        │     (n int, err error)    │
        └───────────────────────────┘
                    ▲
     implemented by │
        ┌───────────┼────────────────┬─────────────────┐
        │           │                │                  │
   *os.File   bufio.Reader     bytes.Buffer       net.Conn / os.Stdin
   (a file)   (wraps a Reader) (in-memory bytes)  (sockets, stdin)

        ┌───────────────────────────┐
        │   io.Writer interface     │
        │   Write(p []byte)         │
        │     (n int, err error)    │
        └───────────────────────────┘
                    ▲
     implemented by │
        ┌───────────┼────────────────┬─────────────────┐
        │           │                │                  │
   *os.File   bufio.Writer     bytes.Buffer       net.Conn / os.Stdout
   (a file)   (wraps a Writer) (in-memory bytes)  (sockets, stdout)
```

Everything in this level either **is** one of these two interfaces or **wraps** one. That's why `fmt.Fprintf(w, ...)` works identically whether `w` is a file, a buffer, or the network - it only needs `Write`.

---

## 🚦 File Open Flags (for os.OpenFile)

| Flag | Meaning |
|------|---------|
| `os.O_RDONLY` | open for reading only |
| `os.O_WRONLY` | open for writing only |
| `os.O_RDWR` | open for reading and writing |
| `os.O_APPEND` | writes go to the end of the file |
| `os.O_CREATE` | create the file if it doesn't exist |
| `os.O_TRUNC` | truncate the file to zero length on open |
| `os.O_EXCL` | (with O_CREATE) fail if the file already exists |

### Common Combinations

| Combination | Result |
|-------------|--------|
| `os.O_RDONLY` | plain read-only open (same as `os.Open`) |
| `os.O_WRONLY\|os.O_CREATE\|os.O_TRUNC` | write-only, create if missing, wipe existing content (same as `os.Create`) |
| `os.O_WRONLY\|os.O_CREATE\|os.O_APPEND` | write-only, create if missing, **append** without wiping (the log-file pattern) |
| `os.O_RDWR\|os.O_CREATE` | read and write, create if missing, keep existing content, cursor starts at 0 |
| `os.O_WRONLY\|os.O_CREATE\|os.O_EXCL` | write-only, but **fail** if the file already exists (safe "create new, never overwrite") |

```go
// The log-file pattern from Exercise 6:
file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
```

---

## 🔢 File Permission Octal Breakdown

```
      0   6    4    4
      │   │    │    │
      │   │    │    └── other:  4(r) + 0(w) + 0(x) = 4  -> r--
      │   │    └─────── group:  4(r) + 0(w) + 0(x) = 4  -> r--
      │   └──────────── owner:  4(r) + 2(w) + 0(x) = 6  -> rw-
      └──────────────── marks this literal as octal (base 8) in Go source

  bit values:  read = 4   write = 2   execute = 1  (sum the ones that apply)
```

| Octal | Owner | Group | Other | Typical Use |
|-------|-------|-------|-------|-------------|
| `0600` | rw- | --- | --- | private file, owner only |
| `0644` | rw- | r-- | r-- | ordinary file, world-readable |
| `0666` | rw- | rw- | rw- | os.Create's requested default (narrowed by umask) |
| `0700` | rwx | --- | --- | private directory, owner only |
| `0755` | rwx | r-x | r-x | ordinary directory, world-readable/traversable |
| `0777` | rwx | rwx | rwx | fully open (rarely appropriate) |

**Verified on this machine:** `os.Create` requested the default `0666`, but the OS umask narrowed it to `-rw-r--r--` (0644) - which is exactly why `data.txt` and `public.txt` (explicitly `0644`) printed identical permissions in Exercise 3.

---

## 🧰 os vs bufio: Choosing the Right Tool

| Task | Convenience (whole file) | Streaming / fine control |
|------|---------------------------|---------------------------|
| Read a file | `os.ReadFile(path)` | `os.Open` + `bufio.NewScanner` or `bufio.NewReader` |
| Write a file | `os.WriteFile(path, data, perm)` | `os.Create` + `bufio.NewWriter` (remember `Flush()`) |
| Append to a file | *(no convenience function - use OpenFile)* | `os.OpenFile(path, O_APPEND\|O_CREATE\|O_WRONLY, perm)` |
| Line-by-line read | *(not applicable)* | `bufio.NewScanner(file)` + `scanner.Scan()`/`Text()` |
| Read until a delimiter | *(not applicable)* | `bufio.NewReader(file).ReadString(delim)` |

**Rule of thumb:** reach for `os.ReadFile`/`os.WriteFile` first - they're one line and handle open/close for you. Drop to `os.Open`/`os.Create` plus `bufio` only when the file might be large, you need line-by-line access, or you need to append.

---

## 🗂️ Directory Operations at a Glance

```
os.Mkdir("a/b", 0755)       -> ERROR if "a" doesn't already exist
os.MkdirAll("a/b", 0755)    -> creates "a" AND "b", no error if already there

os.ReadDir("dir")           -> []fs.DirEntry, already sorted by name

filepath.Join("a", "b", "c.txt")   -> "a/b/c.txt" (portable separator)
"a" + "/" + "b" + "/" + "c.txt"    -> fragile, avoid this

filepath.WalkDir("root", fn)  -> visits every file & dir under root, lexical order
```

---

## 🌡️ Temp Files/Directories: The Safe Scratch Pattern

```go
tmpFile, err := os.CreateTemp(dir, "prefix-*.ext")
// dir == ""  -> OS-wide temp folder
// dir == "." -> inside the CURRENT directory (keeps things self-contained)

tmpDir, err := os.MkdirTemp(dir, "prefix-*")

// Pair creation with cleanup IMMEDIATELY:
defer os.Remove(tmpFile.Name())
defer os.RemoveAll(tmpDir)
```

The trailing `*` in the pattern is replaced with random characters, guaranteeing the name doesn't collide with another run.

---

## 🚨 Common Mistakes

### Mistake 1: Leaking File Descriptors

```go
// ❌ WRONG - never closed
file, _ := os.Open("data.txt")

// ✅ RIGHT
file, err := os.Open("data.txt")
if err != nil { return err }
defer file.Close()
```

### Mistake 2: String-Concatenating Paths

```go
// ❌ WRONG
path := dir + "/" + name

// ✅ RIGHT
path := filepath.Join(dir, name)
```

### Mistake 3: Skipping os.IsNotExist

```go
// ❌ WRONG - treats "not created yet" as fatal
data, err := os.ReadFile(path)
if err != nil { panic(err) }

// ✅ RIGHT
if err != nil && !os.IsNotExist(err) { panic(err) }
```

### Mistake 4: os.RemoveAll on an Unchecked Path

```go
// ❌ DANGEROUS if basePath could be ""
os.RemoveAll(filepath.Join(basePath, "cache"))

// ✅ RIGHT - validate first
if basePath == "" { return errors.New("basePath required") }
```

### Mistake 5: Forgetting bufio.Writer.Flush()

```go
writer := bufio.NewWriter(file)
writer.WriteString("data")
writer.Flush() // don't forget this before closing!
file.Close()
```

---

## 📈 Progression Summary

### Understanding Level 18

Level 18 teaches how Go programs persist and retrieve data from the filesystem:

1. **The interface** - io.Reader/io.Writer underlie everything
2. **Whole-file convenience** - os.ReadFile/os.WriteFile for small files
3. **Streaming control** - os.Open/os.Create + bufio for large files and line-by-line work
4. **Appending** - os.OpenFile with the right flag combination
5. **Directories** - building, listing, walking, and safely deleting trees
6. **Scratch space** - temp files/directories that clean up easily
7. **Metadata** - os.Stat and the os.IsNotExist idiom

### Prerequisites for Level 19

Before moving to Level 19 (JSON), you need:

- ✅ Comfortable reading and writing whole files
- ✅ Can correctly check for a missing file with os.IsNotExist
- ✅ Understand file permissions well enough to pick a sensible mode
- ✅ Comfortable with bufio.Scanner for line-by-line reading
- ✅ Can build portable paths with filepath.Join
- ✅ Understand the difference between os.Remove and os.RemoveAll

### Ready for Level 19?

Level 19 builds directly on this level - you'll take the bytes you now know how to read from (and write to) files, and learn to encode/decode them as **JSON**, Go's most common data-interchange format:
- Marshaling Go structs to JSON and back
- Struct tags to control field names
- Reading/writing JSON files using the exact os.ReadFile/os.WriteFile patterns from this level

---

## ✅ Checklist Before Level 19

- [ ] Can read and write a whole file with os.ReadFile/os.WriteFile
- [ ] Can explain what os.IsNotExist checks and why it's better than string-matching an error
- [ ] Can read the octal digits of a permission like 0644 and say what each digit means
- [ ] Can read a file line by line with bufio.Scanner
- [ ] Know why os.OpenFile with O_APPEND is needed instead of os.WriteFile for logs
- [ ] Can build a path with filepath.Join instead of string concatenation
- [ ] Can walk a directory tree with filepath.WalkDir
- [ ] Understand why os.RemoveAll needs extra care
- [ ] Completed 8+ exercises

---

## 💡 Key Takeaways

### The Interface
`io.Reader`/`io.Writer` are two one-method interfaces that almost everything file-related implements or wraps.

### The Convenience/Control Tradeoff
`os.ReadFile`/`os.WriteFile` for small files done in one call; `os.Open`/`os.Create` plus `bufio` when you need streaming, line-by-line access, or appending.

### The Error Idiom
Don't check "does this file exist?" before acting - try the operation and check `os.IsNotExist(err)` on the result.

### The Danger Zone
`os.RemoveAll` and raw string path concatenation are the two things in this level most likely to bite you in production. Always build paths with `filepath.Join`, and always validate a path before recursively deleting it.

---

## 📚 Next Level

Level 19: JSON
- Encoding Go values to JSON and decoding JSON into Go values
- Struct tags for controlling field names
- Working with nested and dynamic JSON data

You've got file handling down! Keep going! 🚀
