# ini-reader

Memory-efficient INI file parser for Go. The primary purpose of this package is to read large INI files without loading
the entire file into memory.

# Installing

```bash
go get github.com/eugeniypetrov/ini-reader
```

# Usage

```go
content := `
; comment
[section1]
key=value

[section2]
foo="bar" ; comment
`
r := ini.NewReader(strings.NewReader(content))
for r.Next() {
    s := r.Section()
    log.Println(s.Name, s.Properties)
}

if err := r.Err(); err != nil {
    panic(err)
}
```