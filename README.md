# timeformat

A Go library for formatting and parsing time using `yyyy-MM-dd` style patterns (instead of Go's reference time `2006-01-02 15:04:05`).

## Usage

```go
package main

import (
    "fmt"
    "time"
    "github.com/springandall/timeformat"
)

func main() {
    // Format
    s, _ := timeformat.Format("yyyy-MM-dd HH:mm:ss", time.Now())
    fmt.Println(s) // 2024-06-15 14:30:45

    // Parse
    t, _ := timeformat.Parse("yyyy-MM-dd HH:mm:ss", "2024-06-15 14:30:45", time.UTC)
    fmt.Println(t)

    // With nanosecond
    s, _ = timeformat.Format("yyyy-MM-dd HH:mm:ss.SSS", time.Now())
    fmt.Println(s) // 2024-06-15 14:30:45.123

    // With nanosecond
    s, _ = timeformat.Format("yyyy-MM-dd HH:mm:ss SSS", time.Now())
    fmt.Println(s) // 2024-06-15 14:30:45 123
}
```

## Pattern Symbols

| Symbol | Meaning             | Length |
|--------|---------------------|--------|
| `yyyy` | Year (4 digits)     | 4      |
| `MM`   | Month (2 digits)    | 2      |
| `dd`   | Day (2 digits)      | 2      |
| `HH`   | Hour (2 digits)     | 2      |
| `mm`   | Minute (2 digits)   | 2      |
| `ss`   | Second (2 digits)   | 2      |
| `SS`   | Nanosecond fraction | 2-9    |

Other characters are treated as literals.

> Note: `SSS` supports nanosecond without a dot separator — both `yyyy-MM-dd HH:mm:ss.SSS` and `yyyy-MM-dd HH:mm:ssSSS` are valid.

## API

- `Format(pattern, t)` — format time as string
- `Parse(pattern, value, loc)` — parse string into `time.Time`
