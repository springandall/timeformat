# timeformat

使用 `yyyy-MM-dd` 风格模式来格式化和解析时间的 Go 库（代替 Go 标准库的参考时间 `2006-01-02 15:04:05`）。

## 用法

```go
package main

import (
    "fmt"
    "time"
    "github.com/springandall/timeformat"
)

func main() {
    // 格式化
    s, _ := timeformat.Format(time.Now(), "yyyy-MM-dd HH:mm:ss")
    fmt.Println(s) // 2024-06-15 14:30:45

    // 解析
    t, _ := timeformat.Parse("2024-06-15 14:30:45", "yyyy-MM-dd HH:mm:ss", time.UTC)
    fmt.Println(t)

    // 带纳秒
    s, _ = timeformat.Format(time.Now(), "yyyy-MM-dd HH:mm:ss.SSS")
    fmt.Println(s) // 2024-06-15 14:30:45.123
	
	// 带纳秒
	s, _ = timeformat.Format(time.Now(), "yyyy-MM-dd HH:mm:ss SSS")
	fmt.Println(s) // 2024-06-15 14:30:45 123
}
```

## 模式符号

| 符号    | 含义         | 长度 |
|---------|-------------|------|
| `yyyy`  | 年（4位）    | 4    |
| `MM`    | 月（2位）    | 2    |
| `dd`    | 日（2位）    | 2    |
| `HH`    | 时（2位）    | 2    |
| `mm`    | 分（2位）    | 2    |
| `ss`    | 秒（2位）    | 2    |
| `SS`    | 纳秒（2-9位）| 2-9  |

其他字符作为字面量原样输出。

> 注意：`SSS` 支持不带 `.` 的纳秒，`yyyy-MM-dd HH:mm:ss.SSS` 和 `yyyy-MM-dd HH:mm:ssSSS` 都是有效的写法。

## API

- `OfPattern(pattern)` — 返回 `*TimeParser`
- `Format(t time.Time, pattern)` — 格式化时间为字符串
- `Parse(value, pattern, loc)` — 解析字符串为 `time.Time`
- `TimeParser.Format(t)` — 使用缓存好的解析器格式化
- `TimeParser.Parse(value, loc)` — 使用缓存好的解析器解析
