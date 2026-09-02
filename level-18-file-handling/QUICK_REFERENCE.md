# Level 18: Quick Reference Card

## 🚀 Quick Start (5 Minutes)

```bash
# Setup
mkdir myapp && cd myapp
go mod init myapp

# Create main.go
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "os"
)

func main() {
    err := os.WriteFile("hello.txt", []byte("Hello, files!\n"), 0644)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }

    data, err := os.ReadFile("hello.txt")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Print(string(data))

    os.Remove("hello.txt")
}
EOF

# Run
go run main.go
```

---

## 📋 Reading Files

```go
data, err := os.ReadFile("path.txt")        // whole file, one call

file, err := os.Open("path.txt")            // manual control
if err != nil { /* handle */ }
defer file.Close()
```

---

## 📝 Writing Files

```go
err := os.WriteFile("path.txt", []byte("data"), 0644) // whole file, truncates

file, err := os.Create("path.txt")          // manual control, truncates
if err != nil { /* handle */ }
defer file.Close()
file.WriteString("data")
```

---

## 🔎 Detecting a Missing File

```go
_, err := os.ReadFile("path.txt")
if err != nil {
    if os.IsNotExist(err) {
        // file simply doesn't exist
    } else {
        // some other problem
    }
}
```

---

## 🔢 Permissions (os.FileMode)

```go
0600  // rw-------  owner read/write only
0644  // rw-r--r--  owner read/write, everyone else read-only
0755  // rwxr-xr-x  typical for directories
```

---

## 📚 Buffered I/O

```go
// Line by line
scanner := bufio.NewScanner(file)
for scanner.Scan() {
    line := scanner.Text()
}
if err := scanner.Err(); err != nil { /* handle */ }

// Buffered writer - MUST Flush()
writer := bufio.NewWriter(file)
writer.WriteString("data")
writer.Flush()

// Buffered reader
reader := bufio.NewReader(file)
line, err := reader.ReadString('\n')
```

---

## ➕ Appending

```go
file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
if err != nil { /* handle */ }
defer file.Close()
file.WriteString("new line\n")
```

---

## 🗂️ Directories

```go
os.Mkdir("dir", 0755)               // one directory, parent must exist
os.MkdirAll("a/b/c", 0755)          // all missing parents too

entries, err := os.ReadDir("dir")   // sorted []fs.DirEntry
for _, e := range entries {
    fmt.Println(e.Name(), e.IsDir())
}

path := filepath.Join("a", "b", "c.txt")  // ALWAYS use this, never "+/+"

filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
    if err != nil { return err }
    // ...
    return nil
})
```

---

## 🌡️ Temp Files/Directories

```go
tmpFile, err := os.CreateTemp(".", "prefix-*.txt") // "." = stay in current dir
defer os.Remove(tmpFile.Name())

tmpDir, err := os.MkdirTemp(".", "prefix-*")
defer os.RemoveAll(tmpDir)
```

---

## 📊 Metadata

```go
info, err := os.Stat("path")
if err != nil {
    if os.IsNotExist(err) { /* doesn't exist */ }
    return
}
info.Size()      // int64 bytes
info.ModTime()   // time.Time
info.IsDir()     // bool
info.Mode().Perm() // os.FileMode
```

---

## 🗑️ Deleting

```go
os.Remove("file.txt")      // one file, or one EMPTY directory
os.RemoveAll("dir")        // path + everything inside, recursively - be careful!
```

---

## ⚠️ Common Mistakes

| Mistake | Wrong | Right |
|---------|-------|-------|
| Leaked file descriptor | opening without `defer Close()` | `defer file.Close()` right after a successful open |
| Path concatenation | `dir + "/" + name` | `filepath.Join(dir, name)` |
| Treating missing file as fatal | `if err != nil { panic(err) }` | check `os.IsNotExist(err)` first |
| Overwriting instead of appending | `os.WriteFile` called repeatedly | `os.OpenFile` with `O_APPEND\|O_CREATE\|O_WRONLY` |
| Lost buffered writes | closing without `Flush()` | always `writer.Flush()` before `file.Close()` |
| Dangerous deletion | `os.RemoveAll` on an unchecked/empty path | validate the path before calling `os.RemoveAll` |

---

## 🎓 Before Next Level

Can you:
- [ ] Read and write a whole file in one call?
- [ ] Correctly check for a missing file with os.IsNotExist?
- [ ] Explain what 0644 means, digit by digit?
- [ ] Read a file line by line with bufio.Scanner?
- [ ] Append to a file without truncating it?
- [ ] Build a path with filepath.Join and walk a tree with filepath.WalkDir?
- [ ] Explain the danger of os.RemoveAll?

If YES → You're ready for Level 19!

---

## 📚 Next Level

Level 19: JSON
- Encoding Go values to JSON and decoding JSON into Go values
- Struct tags for controlling field names
- Working with nested and dynamic JSON data

You've got file handling down! 💪
