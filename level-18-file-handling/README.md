# Level 18: File Handling - Complete Guide

## Introduction

Welcome to Level 18! You've learned how to organize code into packages and modules (Level 17). Now it's time to make your programs talk to the outside world by reading and writing **files**.

Every Go program that touches a file - config files, logs, CSVs, saved data - goes through the same small set of building blocks: `os` for opening/creating/removing files and directories, `io`/`bufio` for reading and writing their contents, and `path/filepath` for building paths that work on every operating system. This level walks through all of them, verified end to end with real, runnable code.

---

## Table of Contents

1. [io.Reader and io.Writer: Go's Universal I/O Abstraction](#ioreader-and-iowriter-gos-universal-io-abstraction)
2. [Opening and Reading Files](#opening-and-reading-files)
3. [Writing Files](#writing-files)
4. [Buffered I/O](#buffered-io)
5. [Appending to Files](#appending-to-files)
6. [Working With Directories](#working-with-directories)
7. [Temporary Files and Directories](#temporary-files-and-directories)
8. [Checking File Existence and Metadata](#checking-file-existence-and-metadata)
9. [Deleting Files and Directories](#deleting-files-and-directories)
10. [Best Practices](#best-practices)
11. [Common Mistakes](#common-mistakes)

---

## io.Reader and io.Writer: Go's Universal I/O Abstraction

Back in Level 15 you learned that a Go interface is satisfied implicitly - any type with the right methods qualifies. Almost the entire standard library's I/O story rests on two tiny interfaces built exactly that way:

```go
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Writer interface {
    Write(p []byte) (n int, err error)
}
```

That's it. `Read` fills the given byte slice with data and reports how many bytes it read; `Write` consumes bytes from the given slice and reports how many it wrote. Files (`*os.File`), network connections, in-memory buffers (`bytes.Buffer`), compressors, and even `os.Stdin`/`os.Stdout` all implement these two interfaces.

Why this matters for this level: almost everything you'll use - `os.File`, `bufio.Scanner`, `bufio.Reader`, `bufio.Writer` - either **is** an `io.Reader`/`io.Writer` or **wraps** one. Once you understand that a file is "just something you can Read from or Write to," the rest of this level is about convenience functions built on top of that idea.

```go
var r io.Reader = someFile   // *os.File satisfies io.Reader
var w io.Writer = someFile   // *os.File satisfies io.Writer
```

This is also why functions like `io.Copy(dst io.Writer, src io.Reader)` can copy from a file to a network connection, or from stdin to a file, without caring what the concrete types are.

---

## Opening and Reading Files

Go gives you two levels of control for reading files: a whole-file convenience function, and a lower-level handle you manage yourself.

### os.ReadFile: The Convenience Way

For small-to-medium files, `os.ReadFile` reads the entire file into memory in one call - open, read, close, all handled for you:

```go
data, err := os.ReadFile("config.txt")
if err != nil {
    // handle error
}
fmt.Println(string(data))
```

### os.Open: The Manual Way

When you need more control - reading line by line, reading a huge file without loading it all into memory, or seeking around - use `os.Open` to get an `*os.File` handle:

```go
file, err := os.Open("config.txt")
if err != nil {
    // handle error
}
defer file.Close()
```

### defer for Closing (Connects to Level 11/16)

You met `defer` back in Level 11 and used it for cleanup in Level 16's error handling. File handles are exactly the resource `defer` was made for: open the file, immediately `defer file.Close()`, and the file is guaranteed to close when the function returns - no matter which `return` statement fires or whether a later step fails.

```go
file, err := os.Open("config.txt")
if err != nil {
    return err
}
defer file.Close() // guaranteed to run when this function exits
```

### The Real Error From a Missing File

Opening (or reading) a file that doesn't exist doesn't panic - it returns an `error` you check like any other. The verified, real error text from this machine looks like:

```
Error: open does-not-exist.txt: no such file or directory
```

Rather than string-matching that message, use `os.IsNotExist` to check *why* the operation failed:

```go
_, err := os.ReadFile("does-not-exist.txt")
if err != nil {
    if os.IsNotExist(err) {
        fmt.Println("the file simply isn't there")
    } else {
        fmt.Println("some other problem:", err)
    }
}
```

`os.IsNotExist` works the same way whether the error came from `os.Open`, `os.ReadFile`, or `os.Stat` - it inspects the underlying error rather than comparing strings, so it's the correct, portable way to make this check.

---

## Writing Files

Just like reading, Go gives you a whole-file convenience function and a lower-level handle.

### os.WriteFile: The Convenience Way

Creates the file (or truncates it if it exists) and writes the given bytes, all in one call:

```go
err := os.WriteFile("output.txt", []byte("hello\n"), 0644)
if err != nil {
    // handle error
}
```

### os.Create: The Manual Way

Creates a new file for writing (truncating it if it already exists) and returns an `*os.File` you write to and must close yourself:

```go
file, err := os.Create("output.txt")
if err != nil {
    // handle error
}
defer file.Close()
file.WriteString("hello\n")
```

### File Permissions (os.FileMode)

Both `os.WriteFile` and `os.OpenFile` take a permission argument like `0644`. This is an `os.FileMode` value written in **octal** (base 8), and it directly mirrors Unix file permission bits: three groups of three bits for **owner**, **group**, and **other**, where each group is `read (4) + write (2) + execute (1)`.

```
0644
│└┴┴─ owner=6 (rw-), group=4 (r--), other=4 (r--)
└──── leading 0 marks this as an octal literal in Go source

  6 = 4(read) + 2(write)          -> rw-
  4 = 4(read)                     -> r--
  0 = nothing                     -> ---

0644 -> owner: read+write, group: read-only, other: read-only
0600 -> owner: read+write, group: nothing,   other: nothing
0755 -> owner: read+write+execute, group/other: read+execute (typical for directories)
```

Verified with `os.Stat(...).Mode().Perm()`:

```
data.txt permissions: -rw-r--r--   (created with os.Create, then Stat'd)
public.txt permissions: -rw-r--r-- (written with 0644)
private.txt permissions: -rw-------(written with 0600)
```

`os.Create` uses a fixed default (`0666`) that the operating system's umask then narrows - which is why the plain `os.Create` file above still ended up `0644` on this machine. When you need a *specific* permission, use `os.WriteFile` or `os.OpenFile` and pass the mode explicitly.

---

## Buffered I/O

### Why Buffering Matters

Every `Read`/`Write` call on an `*os.File` is a system call - a relatively expensive trip into the operating system kernel. If you read or write one byte (or one line) at a time directly against a file, you pay that cost over and over. **Buffered I/O** batches many small reads/writes into fewer, larger system calls by staging data through an in-memory buffer first.

### bufio.Scanner: Reading Line by Line

The most common way to process a text file line by line:

```go
file, err := os.Open("fruits.txt")
if err != nil {
    // handle error
}
defer file.Close()

scanner := bufio.NewScanner(file)
for scanner.Scan() {
    line := scanner.Text() // current line, without the trailing newline
    fmt.Println(line)
}
if err := scanner.Err(); err != nil {
    // handle a read error that happened mid-scan
}
```

`Scan()` advances to the next line and returns `false` when it runs out of input **or** hits an error - which is why you always check `scanner.Err()` after the loop to tell "ran out of lines" apart from "something went wrong."

### bufio.NewReader / bufio.NewWriter: More Control

`bufio.Scanner` is great for "give me the next line." When you need finer control - read until a specific delimiter, peek ahead, or write many small pieces efficiently - wrap the file in a `bufio.Reader` or `bufio.Writer`:

```go
reader := bufio.NewReader(file)
line, err := reader.ReadString('\n') // reads up to and including the delimiter

writer := bufio.NewWriter(file)
fmt.Fprintf(writer, "line %d\n", 1)
writer.Flush() // IMPORTANT: buffered data isn't written until you Flush (or the buffer fills)
```

The single most important rule with `bufio.Writer`: **data sits in memory until you call `Flush()`** (or the internal buffer fills up on its own). Forget to flush, and your last few writes can silently vanish when the program exits.

---

## Appending to Files

`os.WriteFile` and `os.Create` both **truncate** - they wipe out any existing content. To add to a file without destroying what's already there (think: log files), open it with explicit flags via `os.OpenFile`:

```go
file, err := os.OpenFile("activity.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
if err != nil {
    // handle error
}
defer file.Close()
file.WriteString("new event\n")
```

The three flags combined with `|` (bitwise OR) each add a capability:

- `os.O_APPEND` - writes go to the end of the file, not wherever the write cursor happens to be
- `os.O_CREATE` - create the file if it doesn't exist yet
- `os.O_WRONLY` - open for writing only (as opposed to `O_RDONLY` or `O_RDWR`)

Every call to `os.OpenFile` with these flags on the same path appends one more chunk - open, write, close, repeat - which is exactly the pattern a logger uses.

---

## Working With Directories

### Creating Directories: Mkdir vs MkdirAll

```go
os.Mkdir("data", 0755)      // creates ONE directory; fails if "data"'s parent doesn't exist
os.MkdirAll("data/sub/deep", 0755) // creates ALL missing parent directories too, no error if it already exists
```

Use `os.MkdirAll` almost by default - it's the "make sure this path exists" version and won't complain if part (or all) of the path is already there.

### Listing Contents: os.ReadDir

```go
entries, err := os.ReadDir("project")
if err != nil {
    // handle error
}
for _, e := range entries {
    fmt.Println(e.Name(), e.IsDir())
}
```

`os.ReadDir` returns entries already sorted by filename, and each entry is a lightweight `fs.DirEntry` (name, whether it's a directory, and enough to fetch full metadata via `e.Info()` without an extra `os.Stat` call for most uses).

### Building Paths: filepath.Join

**Never concatenate paths with a raw `"/"`.** Use `filepath.Join` instead:

```go
// ❌ Fragile - hardcodes "/", breaks on Windows, mishandles trailing slashes
path := dir + "/" + "file.txt"

// ✅ Portable - uses the right separator for the OS, cleans up extra slashes
path := filepath.Join(dir, "file.txt")
```

`filepath.Join` uses the correct separator for whatever OS the program runs on (`/` on Linux/macOS, `\` on Windows) and normalizes away doubled or trailing separators. There's no good reason to hand-build a path with string concatenation once `filepath.Join` exists.

### Walking a Directory Tree: filepath.WalkDir

To process every file (and subdirectory) under a root, recursively, use `filepath.WalkDir`:

```go
err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
    if err != nil {
        return err // something went wrong accessing this entry
    }
    if d.IsDir() {
        return nil // nothing to do for directories themselves
    }
    info, err := d.Info()
    if err != nil {
        return err
    }
    fmt.Println(path, info.Size())
    return nil
})
```

`filepath.WalkDir` visits every file and directory under `root` in a deterministic, lexical order, calling your function once per entry. (There's an older `filepath.Walk`, which uses `os.FileInfo` instead of the lighter `fs.DirEntry` - prefer `WalkDir`, it's faster because it avoids an extra `stat` call per entry in the common case.)

---

## Temporary Files and Directories

When you need scratch space - a place to stage data you don't want to name or clean up by hand - use `os.CreateTemp` and `os.MkdirTemp`. Both take a directory (empty string means "the OS's default temp directory") and a name pattern where a trailing `*` is replaced with random characters, guaranteeing a unique name:

```go
tmpFile, err := os.CreateTemp("", "myapp-*.txt")
if err != nil {
    // handle error
}
defer os.Remove(tmpFile.Name()) // clean up when done
defer tmpFile.Close()

tmpDir, err := os.MkdirTemp("", "myapp-*")
if err != nil {
    // handle error
}
defer os.RemoveAll(tmpDir) // clean up the whole directory when done
```

This is the right way to get a scratch file or directory: the name is guaranteed unique (no collisions with another run of your program), and pairing the creation with a `defer os.Remove`/`defer os.RemoveAll` right away means cleanup can never be forgotten later in the function.

---

## Checking File Existence and Metadata

`os.Stat` fetches metadata about a path - size, permissions, modification time, and whether it's a directory - without opening the file's contents:

```go
info, err := os.Stat("report.txt")
if err != nil {
    if os.IsNotExist(err) {
        fmt.Println("doesn't exist yet")
    }
    return
}
fmt.Println("Size:", info.Size())
fmt.Println("Modified:", info.ModTime())
fmt.Println("Is a directory:", info.IsDir())
fmt.Println("Permissions:", info.Mode().Perm())
```

The same `os.IsNotExist` check from the [reading section](#opening-and-reading-files) works here too - `os.Stat` on a missing path returns an error that satisfies it. This is the standard "does this file exist?" idiom in Go: **try the operation and check the error**, rather than checking existence first and then acting (which has a race condition between the check and the use anyway).

---

## Deleting Files and Directories

```go
os.Remove("file.txt")       // removes ONE file, or ONE empty directory - errors if the directory isn't empty
os.RemoveAll("some/dir")    // removes a path AND everything inside it, recursively - no error if it doesn't exist
```

`os.RemoveAll` is powerful and convenient - and that's exactly why it deserves caution. It will silently delete an entire directory tree with no confirmation and no undo. Always build the path with `filepath.Join` (never raw string concatenation - see [Working With Directories](#working-with-directories)) and double-check it points inside a directory you own and expect, especially if any part of the path comes from user input or a variable that could be empty (an empty string or `"."` passed to `os.RemoveAll` is a genuinely dangerous mistake).

---

## Best Practices

### 1. Always defer Close() Right After a Successful Open

```go
// ✅ Good - Close is guaranteed even if later code returns early
file, err := os.Open("data.txt")
if err != nil {
    return err
}
defer file.Close()
```

### 2. Check the Error From Close() When It Matters

For files you only read, a failed `Close()` rarely matters. For files you **wrote** to, a failed `Close()` can mean data never actually made it to disk (some errors, like a full disk, only surface on close/flush):

```go
file, err := os.Create("important.txt")
if err != nil {
    return err
}
defer func() {
    if cerr := file.Close(); cerr != nil {
        log.Println("warning: failed to close file:", cerr)
    }
}()
```

### 3. Prefer os.ReadFile/WriteFile for Small Files, Stream for Large Ones

```go
// ✅ Good for small config files, short text files, etc.
data, err := os.ReadFile("config.json")

// ✅ Good for large files - avoids loading gigabytes into memory at once
file, _ := os.Open("huge.csv")
defer file.Close()
scanner := bufio.NewScanner(file)
for scanner.Scan() { /* process one line at a time */ }
```

### 4. Use filepath.Join, Never Raw String Concatenation

```go
// ✅ Good
path := filepath.Join(dir, "sub", "file.txt")

// ❌ Bad
path := dir + "/sub/" + "file.txt"
```

### 5. Check os.IsNotExist Before Treating an Error as Fatal

```go
// ✅ Good - "file doesn't exist" is often a normal, expected case
if _, err := os.Stat(path); err != nil {
    if os.IsNotExist(err) {
        // create it, use a default, etc. - not necessarily an error
    } else {
        return err // a real problem: permissions, disk error, etc.
    }
}
```

### 6. Pair Temp File/Directory Creation With an Immediate defer for Cleanup

```go
// ✅ Good
tmpDir, err := os.MkdirTemp("", "scratch-*")
if err != nil {
    return err
}
defer os.RemoveAll(tmpDir)
```

---

## Common Mistakes

### Mistake 1: Forgetting to Close Files (Leaking File Descriptors)

```go
// ❌ WRONG - file is never closed; in a long-running program or loop,
// this exhausts the OS's file descriptor limit
file, _ := os.Open("data.txt")
data := make([]byte, 100)
file.Read(data)

// ✅ RIGHT - defer immediately after a successful open
file, err := os.Open("data.txt")
if err != nil {
    return err
}
defer file.Close()
```

### Mistake 2: Using "/" Instead of filepath.Join

```go
// ❌ WRONG - breaks on Windows, mishandles slashes
path := base + "/" + name

// ✅ RIGHT
path := filepath.Join(base, name)
```

### Mistake 3: Not Checking os.IsNotExist Before Treating an Error as Fatal

```go
// ❌ WRONG - crashes the program on the completely normal case
// of a config file that simply hasn't been created yet
data, err := os.ReadFile("config.txt")
if err != nil {
    panic(err)
}

// ✅ RIGHT
data, err := os.ReadFile("config.txt")
if err != nil {
    if os.IsNotExist(err) {
        data = []byte(defaultConfig)
    } else {
        panic(err) // some other, genuinely unexpected error
    }
}
```

### Mistake 4: os.RemoveAll on the Wrong Path

```go
// ❌ DANGEROUS - if basePath is ever empty (e.g. an unset variable),
// this deletes the current directory and everything under it
os.RemoveAll(filepath.Join(basePath, "cache"))

// ✅ RIGHT - validate the path is what you expect before deleting,
// and never let os.RemoveAll run on an empty or unchecked path
if basePath == "" {
    return errors.New("basePath must not be empty")
}
target := filepath.Join(basePath, "cache")
os.RemoveAll(target)
```

### Mistake 5: Forgetting to Flush a bufio.Writer

```go
// ❌ WRONG - the last writes may never reach disk
writer := bufio.NewWriter(file)
writer.WriteString("important data")
file.Close() // buffered data that was never Flushed can be lost

// ✅ RIGHT
writer := bufio.NewWriter(file)
writer.WriteString("important data")
writer.Flush() // always flush before closing
file.Close()
```

### Mistake 6: Using os.WriteFile/os.Create When You Meant to Append

```go
// ❌ WRONG - each call wipes out the previous contents
os.WriteFile("log.txt", []byte("event 1\n"), 0644)
os.WriteFile("log.txt", []byte("event 2\n"), 0644) // "event 1" is gone!

// ✅ RIGHT - use O_APPEND to add without truncating
file, _ := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
defer file.Close()
file.WriteString("event 2\n")
```

---

## Summary

**The Core Abstraction:**
- `io.Reader`/`io.Writer` - the two tiny interfaces almost everything in this level implements

**Reading:**
- `os.ReadFile` - whole file, one call
- `os.Open` + `defer Close()` - manual control, streaming
- `os.IsNotExist(err)` - the correct way to detect a missing file

**Writing:**
- `os.WriteFile` - whole file, one call, truncates
- `os.Create` - manual control, truncates
- `os.FileMode` (e.g. `0644`) - octal permission bits for owner/group/other

**Buffered I/O:**
- `bufio.Scanner` - line-by-line reading
- `bufio.NewReader`/`bufio.NewWriter` - finer control; remember to `Flush()`

**Appending:**
- `os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, perm)`

**Directories:**
- `os.Mkdir` vs `os.MkdirAll`, `os.ReadDir`, `filepath.Join`, `filepath.WalkDir`

**Temp Space:**
- `os.CreateTemp`/`os.MkdirTemp` - unique names, pair with an immediate `defer` cleanup

**Metadata & Cleanup:**
- `os.Stat` + `os.IsNotExist`, `info.Size()`/`ModTime()`/`IsDir()`
- `os.Remove` vs `os.RemoveAll` - the latter is powerful and dangerous

---

## Next Steps

You now understand:
- ✅ Why `io.Reader`/`io.Writer` underlie almost all of Go's I/O
- ✅ Reading and writing whole files, and streaming large ones
- ✅ File permissions and what the octal digits mean
- ✅ Buffered I/O and why it matters for performance
- ✅ Appending, directory trees, temp files, metadata, and safe deletion

**Next level:** Level 19 - JSON
- Encoding Go values to JSON and decoding JSON into Go values
- Struct tags for controlling field names
- Working with nested and dynamic JSON data

You can now make your programs remember things between runs! Keep going! 🚀
