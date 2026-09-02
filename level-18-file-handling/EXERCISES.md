# Level 18: File Handling - Exercises

Complete all exercises in order. Each exercise builds on previous knowledge.

Every exercise below creates and deletes files **only inside its own project directory** (the one you `cd` into in step 1). None of them touch anything outside that folder, and each cleans up after itself unless the exercise says otherwise.

---

## Exercise 1: Writing and Reading a File

**Objective:** Use os.WriteFile and os.ReadFile for whole-file I/O

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level18-exercise1
cd ~/projects/level18-exercise1
go mod init level18.example/exercise1
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "os"
)

func main() {
    fmt.Println("=== Writing a File ===")
    content := []byte("Hello, Go file handling!\nThis is line two.\n")
    err := os.WriteFile("greeting.txt", content, 0644)
    if err != nil {
        fmt.Println("Error writing file:", err)
        return
    }
    fmt.Println("Wrote greeting.txt")

    fmt.Println("\n=== Reading It Back ===")
    data, err := os.ReadFile("greeting.txt")
    if err != nil {
        fmt.Println("Error reading file:", err)
        return
    }
    fmt.Println("Contents:")
    fmt.Print(string(data))

    fmt.Println("\n=== Cleaning Up ===")
    err = os.Remove("greeting.txt")
    if err != nil {
        fmt.Println("Error removing file:", err)
        return
    }
    fmt.Println("Removed greeting.txt")
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Writing a File ===
Wrote greeting.txt

=== Reading It Back ===
Contents:
Hello, Go file handling!
This is line two.

=== Cleaning Up ===
Removed greeting.txt
```

**Learning Objectives:**
- ✅ Write a file in one call with os.WriteFile
- ✅ Read a whole file in one call with os.ReadFile
- ✅ Clean up a file the exercise created, with os.Remove

---

## Exercise 2: Handling a Nonexistent File

**Objective:** Correctly detect a missing file with os.IsNotExist

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level18-exercise2
cd ~/projects/level18-exercise2
go mod init level18.example/exercise2
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "os"
)

func main() {
    fmt.Println("=== Reading a File That Does Not Exist ===")
    _, err := os.ReadFile("does-not-exist.txt")
    if err != nil {
        fmt.Println("Error:", err)
        if os.IsNotExist(err) {
            fmt.Println("Confirmed: the file does not exist (os.IsNotExist)")
        } else {
            fmt.Println("Some other error occurred")
        }
    }

    fmt.Println("\n=== Same Check With os.Open ===")
    file, err := os.Open("also-missing.txt")
    if err != nil {
        fmt.Println("Error:", err)
        if os.IsNotExist(err) {
            fmt.Println("Confirmed: os.Open also reports IsNotExist")
        }
    } else {
        file.Close()
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
=== Reading a File That Does Not Exist ===
Error: open does-not-exist.txt: no such file or directory
Confirmed: the file does not exist (os.IsNotExist)

=== Same Check With os.Open ===
Error: open also-missing.txt: no such file or directory
Confirmed: os.Open also reports IsNotExist
```

This exercise never creates a file at all - both filenames are intentionally missing so you can see the real, unmodified error text Go produces (verified above; the wording is identical no matter which directory you run this in, since these are relative paths).

**Learning Objectives:**
- ✅ See the exact error text Go returns for a missing file
- ✅ Use os.IsNotExist to check *why* an operation failed, not just *that* it failed
- ✅ Understand os.Open and os.ReadFile report missing files the same way

---

## Exercise 3: File Permissions

**Objective:** Understand os.FileMode and set explicit permissions

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level18-exercise3
cd ~/projects/level18-exercise3
go mod init level18.example/exercise3
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "os"
)

func main() {
    fmt.Println("=== os.Create Makes an Empty File (Truncates if it Exists) ===")
    file, err := os.Create("data.txt")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    file.WriteString("created with os.Create\n")
    file.Close()

    info, err := os.Stat("data.txt")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Printf("data.txt permissions: %v\n", info.Mode().Perm())

    fmt.Println("\n=== os.WriteFile With Explicit Permissions ===")
    err = os.WriteFile("public.txt", []byte("readable by everyone\n"), 0644)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    info, _ = os.Stat("public.txt")
    fmt.Printf("public.txt permissions: %v\n", info.Mode().Perm())

    err = os.WriteFile("private.txt", []byte("owner only\n"), 0600)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    info, _ = os.Stat("private.txt")
    fmt.Printf("private.txt permissions: %v\n", info.Mode().Perm())

    fmt.Println("\n=== Cleaning Up ===")
    for _, name := range []string{"data.txt", "public.txt", "private.txt"} {
        os.Remove(name)
    }
    fmt.Println("Removed all files")
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== os.Create Makes an Empty File (Truncates if it Exists) ===
data.txt permissions: -rw-r--r--

=== os.WriteFile With Explicit Permissions ===
public.txt permissions: -rw-r--r--
private.txt permissions: -rw-------

=== Cleaning Up ===
Removed all files
```

Note: `os.Create` always requests `0666`, but the operating system's umask narrows it - on this machine that produced `0644`, which is why it matches `public.txt` above. Your exact umask may differ, but `private.txt` (explicitly `0600`) will always come out `-rw-------` since it was set directly.

**Learning Objectives:**
- ✅ Understand os.FileMode octal permission bits (owner/group/other)
- ✅ See that os.Create's permissions are subject to the OS umask
- ✅ Set exact permissions explicitly with os.WriteFile

---

## Exercise 4: Reading Line by Line With bufio.Scanner

**Objective:** Process a multi-line text file one line at a time

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level18-exercise4
cd ~/projects/level18-exercise4
go mod init level18.example/exercise4
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "bufio"
    "fmt"
    "os"
)

func main() {
    fmt.Println("=== Creating a Multi-Line File ===")
    lines := "apple\nbanana\ncherry\ndate\n"
    err := os.WriteFile("fruits.txt", []byte(lines), 0644)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println("Wrote fruits.txt")

    fmt.Println("\n=== Reading Line by Line With bufio.Scanner ===")
    file, err := os.Open("fruits.txt")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    lineNum := 1
    count := 0
    for scanner.Scan() {
        fmt.Printf("%d: %s\n", lineNum, scanner.Text())
        lineNum++
        count++
    }
    if err := scanner.Err(); err != nil {
        fmt.Println("Scanner error:", err)
    }
    fmt.Println("Total lines:", count)

    fmt.Println("\n=== Cleaning Up ===")
    os.Remove("fruits.txt")
    fmt.Println("Removed fruits.txt")
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Creating a Multi-Line File ===
Wrote fruits.txt

=== Reading Line by Line With bufio.Scanner ===
1: apple
2: banana
3: cherry
4: date
Total lines: 4

=== Cleaning Up ===
Removed fruits.txt
```

**Learning Objectives:**
- ✅ Create a small text file to process
- ✅ Use bufio.Scanner to read a file line by line
- ✅ Check scanner.Err() after the loop to catch mid-read failures

---

## Exercise 5: Buffered Reading and Writing

**Objective:** Use bufio.NewWriter and bufio.NewReader for finer control

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level18-exercise5
cd ~/projects/level18-exercise5
go mod init level18.example/exercise5
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "bufio"
    "fmt"
    "io"
    "os"
)

func main() {
    fmt.Println("=== Writing With bufio.Writer ===")
    file, err := os.Create("buffered.txt")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    writer := bufio.NewWriter(file)
    for i := 1; i <= 5; i++ {
        fmt.Fprintf(writer, "line %d\n", i)
    }
    // Data sits in the writer's buffer until Flush is called (or the buffer fills up)
    if err := writer.Flush(); err != nil {
        fmt.Println("Flush error:", err)
    }
    file.Close()
    fmt.Println("Wrote 5 lines through a buffered writer")

    fmt.Println("\n=== Reading With bufio.NewReader ===")
    file, err = os.Open("buffered.txt")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    reader := bufio.NewReader(file)
    for {
        line, err := reader.ReadString('\n')
        if len(line) > 0 {
            fmt.Print("Read: ", line)
        }
        if err == io.EOF {
            fmt.Println("Reached end of file")
            break
        }
        if err != nil {
            fmt.Println("Read error:", err)
            break
        }
    }
    file.Close()

    fmt.Println("\n=== Cleaning Up ===")
    os.Remove("buffered.txt")
    fmt.Println("Removed buffered.txt")
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Writing With bufio.Writer ===
Wrote 5 lines through a buffered writer

=== Reading With bufio.NewReader ===
Read: line 1
Read: line 2
Read: line 3
Read: line 4
Read: line 5
Reached end of file

=== Cleaning Up ===
Removed buffered.txt
```

**Learning Objectives:**
- ✅ Write through a bufio.Writer and remember to Flush()
- ✅ Read through a bufio.Reader with ReadString until io.EOF
- ✅ Understand why forgetting Flush() can silently lose data

---

## Exercise 6: Appending Across Multiple Writes

**Objective:** Build up a file's contents over several separate writes

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level18-exercise6
cd ~/projects/level18-exercise6
go mod init level18.example/exercise6
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "os"
)

func appendLine(path, line string) error {
    file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    defer file.Close()
    _, err = file.WriteString(line + "\n")
    return err
}

func main() {
    fmt.Println("=== Appending to a Log File Across Multiple Writes ===")
    logFile := "activity.log"

    events := []string{
        "user logged in",
        "user viewed dashboard",
        "user logged out",
    }

    for _, event := range events {
        if err := appendLine(logFile, event); err != nil {
            fmt.Println("Error appending:", err)
            return
        }
        fmt.Println("Appended:", event)
    }

    fmt.Println("\n=== Final Log Contents ===")
    data, err := os.ReadFile(logFile)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Print(string(data))

    fmt.Println("=== Cleaning Up ===")
    os.Remove(logFile)
    fmt.Println("Removed", logFile)
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Appending to a Log File Across Multiple Writes ===
Appended: user logged in
Appended: user viewed dashboard
Appended: user logged out

=== Final Log Contents ===
user logged in
user viewed dashboard
user logged out
=== Cleaning Up ===
Removed activity.log
```

**Learning Objectives:**
- ✅ Open a file with O_APPEND|O_CREATE|O_WRONLY
- ✅ See that each open-write-close cycle adds to the file instead of erasing it
- ✅ Reuse a small helper function across multiple calls

---

## Exercise 7: Building and Listing a Directory Tree

**Objective:** Create nested directories and list them portably

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level18-exercise7
cd ~/projects/level18-exercise7
go mod init level18.example/exercise7
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "os"
    "path/filepath"
)

func main() {
    fmt.Println("=== Building a Directory Tree With os.MkdirAll ===")
    root := "project"
    dirs := []string{
        filepath.Join(root, "src"),
        filepath.Join(root, "docs"),
        filepath.Join(root, "src", "utils"),
    }
    for _, d := range dirs {
        if err := os.MkdirAll(d, 0755); err != nil {
            fmt.Println("Error:", err)
            return
        }
        fmt.Println("Created:", d)
    }

    fmt.Println("\n=== Writing a File Using filepath.Join ===")
    readmePath := filepath.Join(root, "docs", "README.md")
    err := os.WriteFile(readmePath, []byte("# Project\n"), 0644)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println("Wrote:", readmePath)

    fmt.Println("\n=== Listing the Top-Level Directory With os.ReadDir ===")
    entries, err := os.ReadDir(root)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    for _, e := range entries {
        fmt.Printf("%s (dir=%v)\n", e.Name(), e.IsDir())
    }

    fmt.Println("\n=== Cleaning Up ===")
    err = os.RemoveAll(root)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println("Removed", root, "and everything inside it")
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Building a Directory Tree With os.MkdirAll ===
Created: project/src
Created: project/docs
Created: project/src/utils

=== Writing a File Using filepath.Join ===
Wrote: project/docs/README.md

=== Listing the Top-Level Directory With os.ReadDir ===
docs (dir=true)
src (dir=true)

=== Cleaning Up ===
Removed project and everything inside it
```

Note: `os.ReadDir` returns entries already sorted by name, which is why `docs` prints before `src` even though `src` was created first.

**Learning Objectives:**
- ✅ Use os.MkdirAll to create nested directories in one call
- ✅ Build every path with filepath.Join instead of string concatenation
- ✅ List directory contents with os.ReadDir
- ✅ Use os.RemoveAll to clean up an entire tree at once

---

## Exercise 8: Walking a Directory Tree

**Objective:** Recursively process files and compute a total size

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level18-exercise8
cd ~/projects/level18-exercise8
go mod init level18.example/exercise8
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "io/fs"
    "os"
    "path/filepath"
)

func main() {
    fmt.Println("=== Building a Small Tree of Files ===")
    root := "data"
    type fileSpec struct {
        path    string
        content string
    }
    files := []fileSpec{
        {filepath.Join(root, "a.txt"), "12345"},
        {filepath.Join(root, "sub", "b.txt"), "1234567890"},
        {filepath.Join(root, "sub", "deep", "c.txt"), "abc"},
    }
    for _, f := range files {
        if err := os.MkdirAll(filepath.Dir(f.path), 0755); err != nil {
            fmt.Println("Error:", err)
            return
        }
        if err := os.WriteFile(f.path, []byte(f.content), 0644); err != nil {
            fmt.Println("Error:", err)
            return
        }
    }
    fmt.Println("Created 3 files under", root)

    fmt.Println("\n=== Walking the Tree With filepath.WalkDir ===")
    var totalSize int64
    var fileCount int
    err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
        if err != nil {
            return err
        }
        if d.IsDir() {
            return nil
        }
        info, err := d.Info()
        if err != nil {
            return err
        }
        fmt.Printf("%s (%d bytes)\n", path, info.Size())
        totalSize += info.Size()
        fileCount++
        return nil
    })
    if err != nil {
        fmt.Println("Walk error:", err)
        return
    }
    fmt.Printf("\nTotal: %d files, %d bytes\n", fileCount, totalSize)

    fmt.Println("\n=== Cleaning Up ===")
    os.RemoveAll(root)
    fmt.Println("Removed", root)
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Building a Small Tree of Files ===
Created 3 files under data

=== Walking the Tree With filepath.WalkDir ===
data/a.txt (5 bytes)
data/sub/b.txt (10 bytes)
data/sub/deep/c.txt (3 bytes)

Total: 3 files, 18 bytes

=== Cleaning Up ===
Removed data
```

**Learning Objectives:**
- ✅ Use filepath.WalkDir to visit every file in a tree recursively
- ✅ Skip directory entries and accumulate stats for files only
- ✅ Combine WalkDir with fs.DirEntry.Info() to read file size

---

## Exercise 9: Temporary Files, Directories, and Metadata

**Objective:** Use os.CreateTemp/os.MkdirTemp for scratch space and inspect it with os.Stat

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level18-exercise9
cd ~/projects/level18-exercise9
go mod init level18.example/exercise9
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "os"
    "strings"
)

func main() {
    fmt.Println("=== Scratch File With os.CreateTemp ===")
    // "." keeps the temp file inside THIS project directory instead of the
    // OS-wide temp folder, so this exercise stays fully self-contained.
    tmpFile, err := os.CreateTemp(".", "level18-*.txt")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    name := tmpFile.Name()
    fmt.Println("Created a temp file with a unique, auto-generated name")
    fmt.Println("Name matches pattern 'level18-*.txt':", strings.HasPrefix(name, "./level18-") && strings.HasSuffix(name, ".txt"))
    _, err = tmpFile.WriteString("scratch data\n")
    if err != nil {
        fmt.Println("Error writing:", err)
        return
    }
    tmpFile.Close()

    fmt.Println("\n=== Checking Metadata With os.Stat ===")
    info, err := os.Stat(name)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println("Size:", info.Size(), "bytes")
    fmt.Println("IsDir:", info.IsDir())
    fmt.Println("ModTime recorded:", !info.ModTime().IsZero())

    fmt.Println("\n=== Cleaning Up the Temp File ===")
    os.Remove(name)
    fmt.Println("Removed temp file")

    fmt.Println("\n=== Scratch Directory With os.MkdirTemp ===")
    tmpDir, err := os.MkdirTemp(".", "level18-dir-*")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println("Created a temp directory with a unique, auto-generated name")
    info, err = os.Stat(tmpDir)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println("IsDir:", info.IsDir())

    fmt.Println("\n=== Cleaning Up the Temp Directory ===")
    os.RemoveAll(tmpDir)
    fmt.Println("Removed temp directory")
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Scratch File With os.CreateTemp ===
Created a temp file with a unique, auto-generated name
Name matches pattern 'level18-*.txt': true

=== Checking Metadata With os.Stat ===
Size: 13 bytes
IsDir: false
ModTime recorded: true

=== Cleaning Up the Temp File ===
Removed temp file

=== Scratch Directory With os.MkdirTemp ===
Created a temp directory with a unique, auto-generated name
IsDir: true

=== Cleaning Up the Temp Directory ===
Removed temp directory
```

**Learning Objectives:**
- ✅ Use os.CreateTemp/os.MkdirTemp to get a guaranteed-unique scratch name
- ✅ Keep temp files inside the project's own sandbox by passing "." as the directory
- ✅ Inspect metadata (size, IsDir, ModTime) with os.Stat

---

## Exercise 10: Comprehensive Practice — Rotating Log File

**Objective:** Combine appending, os.Stat, and renaming into a small log rotator

**Instructions:**

1. Create working directory:

```bash
mkdir -p ~/projects/level18-exercise10
cd ~/projects/level18-exercise10
go mod init level18.example/exercise10
```

2. Create `main.go`:

```bash
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "os"
)

const maxLogSize = 100 // bytes; kept small on purpose so rotation is easy to trigger

func writeLog(path, message string) error {
    info, err := os.Stat(path)
    if err == nil && info.Size() >= maxLogSize {
        rotated := path + ".old"
        os.Remove(rotated) // drop any previous rotation
        if err := os.Rename(path, rotated); err != nil {
            return err
        }
        fmt.Println("Rotated log ->", rotated)
    } else if err != nil && !os.IsNotExist(err) {
        return err
    }

    file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    defer file.Close()
    _, err = file.WriteString(message + "\n")
    return err
}

func main() {
    fmt.Println("=== Simulating a Rotating Log File ===")
    logPath := "app.log"

    messages := []string{
        "server started",
        "listening on port 8080",
        "handled request GET /health",
        "handled request GET /users",
        "handled request POST /users",
        "handled request GET /users/42",
        "server shutting down",
    }

    for i, msg := range messages {
        if err := writeLog(logPath, msg); err != nil {
            fmt.Println("Error:", err)
            return
        }
        info, _ := os.Stat(logPath)
        fmt.Printf("Wrote message %d (%q), app.log is now %d bytes\n", i+1, msg, info.Size())
    }

    fmt.Println("\n=== Final State ===")
    for _, name := range []string{"app.log", "app.log.old"} {
        if info, err := os.Stat(name); err == nil {
            fmt.Printf("%s exists, %d bytes\n", name, info.Size())
        } else {
            fmt.Printf("%s does not exist\n", name)
        }
    }

    fmt.Println("\n=== Cleaning Up ===")
    os.Remove("app.log")
    os.Remove("app.log.old")
    fmt.Println("Removed app.log and app.log.old")
}
EOF
```

3. Run the program:

```bash
go run main.go
```

**Expected Output:**

```
=== Simulating a Rotating Log File ===
Wrote message 1 ("server started"), app.log is now 15 bytes
Wrote message 2 ("listening on port 8080"), app.log is now 38 bytes
Wrote message 3 ("handled request GET /health"), app.log is now 66 bytes
Wrote message 4 ("handled request GET /users"), app.log is now 93 bytes
Wrote message 5 ("handled request POST /users"), app.log is now 121 bytes
Rotated log -> app.log.old
Wrote message 6 ("handled request GET /users/42"), app.log is now 30 bytes
Wrote message 7 ("server shutting down"), app.log is now 51 bytes

=== Final State ===
app.log exists, 51 bytes
app.log.old exists, 121 bytes

=== Cleaning Up ===
Removed app.log and app.log.old
```

Rotation triggers right before message 6: after message 5 the log had grown to 121 bytes, which is over the (deliberately tiny) 100-byte limit, so `app.log` is renamed to `app.log.old` and a fresh `app.log` is started before message 6 is appended.

**Learning Objectives:**
- ✅ Combine os.OpenFile appending, os.Stat, and os.Rename in one workflow
- ✅ Trigger and observe a real rotation happening mid-run
- ✅ Practice the "check size before writing" pattern used by real loggers

---

## Bonus Challenges

### Challenge 1: Recursive Directory Size Calculator (Without WalkDir)

Write a function `dirSize(path string) (int64, error)` that computes the total size of everything under `path`, using `os.ReadDir` and manual recursion instead of `filepath.WalkDir`.

```bash
mkdir -p ~/projects/level18-bonus1
cd ~/projects/level18-bonus1
go mod init level18.example/bonus1
```

**Hints:**
- For each entry from `os.ReadDir(path)`, use `filepath.Join(path, entry.Name())` to build the child path (never string-concatenate)
- If the entry is a directory, recurse into it and add the result; if it's a file, use `entry.Info()` to get its size
- Build a small tree with `os.MkdirAll` and `os.WriteFile` first (inside your own project directory), verify your function's total against what you created, then clean up with `os.RemoveAll`

### Challenge 2: Duplicate File Finder by Content Hash

Given a directory of files, find groups of files with identical content by hashing each file with `crypto/sha256` and grouping by hash.

```bash
mkdir -p ~/projects/level18-bonus2
cd ~/projects/level18-bonus2
go mod init level18.example/bonus2
```

**Hints:**
- Use `filepath.WalkDir` to visit every file, `os.ReadFile` (or stream through `io.Copy` into a `sha256.New()` hasher for large files) to get its bytes
- `sha256.Sum256(data)` returns a `[32]byte`; convert it to a `string` (e.g. with `fmt.Sprintf("%x", hash)`) to use as a map key
- Build a `map[string][]string` from hash to matching file paths; any hash with more than one path is a duplicate group
- Create a few files with intentionally identical content to test this, then remove them when done

### Challenge 3: Text File Line Reverser

Read all lines from a file and write them back out to a *new* file in reverse order (last line first).

```bash
mkdir -p ~/projects/level18-bonus3
cd ~/projects/level18-bonus3
go mod init level18.example/bonus3
```

**Hints:**
- Use `bufio.Scanner` to collect every line into a `[]string`
- Reverse the slice in place with the two-index technique from Level 6's bonus challenges (`for left, right := 0, len(s)-1; left < right; left, right = left+1, right-1`)
- Join the reversed lines with `\n` and write them with `os.WriteFile` to a new filename (e.g. `reversed.txt`) inside the same project directory
- Remove both the input and output files when you're done inspecting them

---

## What You've Learned

After completing these 10 exercises, you can:

✅ Read and write whole files with os.ReadFile/os.WriteFile
✅ Correctly detect a missing file with os.IsNotExist instead of guessing from error text
✅ Understand and set file permissions with os.FileMode octal values
✅ Read a file line by line with bufio.Scanner, and use bufio.Reader/Writer for finer control
✅ Append to a file across multiple writes without truncating it
✅ Build, list, and portably join paths in a directory tree
✅ Recursively walk a directory tree and accumulate stats with filepath.WalkDir
✅ Create safe, unique scratch files/directories with os.CreateTemp/os.MkdirTemp
✅ Inspect file metadata with os.Stat and safely delete files/directories

---

## Next Level

Level 19: JSON
- Encoding Go values to JSON and decoding JSON into Go values
- Struct tags for controlling field names
- Working with nested and dynamic JSON data

Great work! Your programs can now persist data to disk! 🚀
