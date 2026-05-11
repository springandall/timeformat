package timeformat

import (
	"errors"
	"time"
	"unicode/utf8"
)

func OfPattern(pattern string) (*TimeParser, error) {
	fields, err := formatter(pattern)
	if err != nil {
		return nil, err
	}
	return &TimeParser{fields}, nil
}

func Format(t time.Time, pattern string) (string, error) {
	p, err := OfPattern(pattern)
	if err != nil {
		return "", err
	}
	return p.Format(t), nil
}

func Parse(value, pattern string, loc *time.Location) (time.Time, error) {
	p, err := OfPattern(pattern)
	if err != nil {
		return time.Time{}, err
	}
	return p.Parse(value, loc)
}

type TimeType rune

const (
	Literal      TimeType = 0
	Year         TimeType = 'y'
	Month        TimeType = 'M'
	Day          TimeType = 'd'
	Hour         TimeType = 'H'
	Minute       TimeType = 'm'
	Second       TimeType = 's'
	NanoOfSecond TimeType = 'S'
)

var TyMap = map[TimeType]struct{}{
	Year:         {},
	Month:        {},
	Day:          {},
	Hour:         {},
	Minute:       {},
	Second:       {},
	NanoOfSecond: {},
}

type fieldKind int

const (
	kindLiteral fieldKind = iota
	kindYear
	kindMonth
	kindDay
	kindHour
	kindMinute
	kindSecond
	kindNano
)

type fieldInfo struct {
	kind   fieldKind
	cursor int
	length int
	text   string
}

type DateArgs struct {
	Year  int
	Month time.Month
	Day   int
	Hour  int
	Min   int
	Sec   int
	Nsec  int
	Loc   *time.Location
}

func (d *DateArgs) time() time.Time {
	return time.Date(d.Year, d.Month, d.Day, d.Hour, d.Min, d.Sec, d.Nsec, d.Loc)
}

type TimeParser struct {
	fields []fieldInfo
}

func (p *TimeParser) Format(t time.Time) string {
	buf := make([]byte, 0, 64)
	for _, f := range p.fields {
		switch f.kind {
		case kindLiteral:
			buf = append(buf, f.text...)
		case kindYear:
			buf = append4Digits(buf, t.Year())
		case kindMonth:
			buf = append2Digits(buf, int(t.Month()))
		case kindDay:
			buf = append2Digits(buf, t.Day())
		case kindHour:
			buf = append2Digits(buf, t.Hour())
		case kindMinute:
			buf = append2Digits(buf, t.Minute())
		case kindSecond:
			buf = append2Digits(buf, t.Second())
		case kindNano:
			buf = appendNano(buf, t.Nanosecond(), f.length)
		}
	}
	return string(buf)
}

func (p *TimeParser) Parse(layout string, loc *time.Location) (time.Time, error) {
	var args DateArgs
	for _, f := range p.fields {
		switch f.kind {
		case kindLiteral:
		case kindYear:
			v, err := atoiStr(layout, f.cursor, f.length)
			if err != nil {
				return time.Time{}, err
			}
			args.Year = v
		case kindMonth:
			v, err := atoiStr(layout, f.cursor, f.length)
			if err != nil {
				return time.Time{}, err
			}
			if v < 1 || v > 12 {
				return time.Time{}, errors.New("month out of range")
			}
			args.Month = time.Month(v)
		case kindDay:
			v, err := atoiStr(layout, f.cursor, f.length)
			if err != nil {
				return time.Time{}, err
			}
			if v < 1 || v > 31 {
				return time.Time{}, errors.New("day out of range")
			}
			args.Day = v
		case kindHour:
			v, err := atoiStr(layout, f.cursor, f.length)
			if err != nil {
				return time.Time{}, err
			}
			if v < 0 || v > 23 {
				return time.Time{}, errors.New("hour out of range")
			}
			args.Hour = v
		case kindMinute:
			v, err := atoiStr(layout, f.cursor, f.length)
			if err != nil {
				return time.Time{}, err
			}
			if v < 0 || v > 59 {
				return time.Time{}, errors.New("minute out of range")
			}
			args.Min = v
		case kindSecond:
			v, err := atoiStr(layout, f.cursor, f.length)
			if err != nil {
				return time.Time{}, err
			}
			if v < 0 || v > 59 {
				return time.Time{}, errors.New("second out of range")
			}
			args.Sec = v
		case kindNano:
			v, err := atoiStr(layout, f.cursor, f.length)
			if err != nil {
				return time.Time{}, err
			}
			for i := f.length; i < 9; i++ {
				v *= 10
			}
			if v < 0 || v > 999999999 {
				return time.Time{}, errors.New("nano out of range")
			}
			args.Nsec = v
		}
	}
	if loc == nil {
		loc = time.Local
	}
	args.Loc = loc
	return args.time(), nil
}

func atoiStr(s string, start, length int) (int, error) {
	if len(s) < start+length {
		return 0, errors.New("input too short")
	}
	n := 0
	for i := 0; i < length; i++ {
		c := s[start+i]
		if c < '0' || c > '9' {
			return 0, errors.New("not a number")
		}
		n = n*10 + int(c-'0')
	}
	return n, nil
}

func append2Digits(buf []byte, v int) []byte {
	return append(buf, byte('0'+v/10), byte('0'+v%10))
}

func append4Digits(buf []byte, v int) []byte {
	return append(buf,
		byte('0'+v/1000),
		byte('0'+(v/100)%10),
		byte('0'+(v/10)%10),
		byte('0'+v%10),
	)
}

func appendNano(buf []byte, v, digits int) []byte {
	var tmp [9]byte
	for i := 8; i >= 0; i-- {
		tmp[i] = byte('0' + v%10)
		v /= 10
	}
	return append(buf, tmp[:digits]...)
}

type fieldRule struct {
	minLen   int
	maxLen   int
	newField func(rlayout []rune, cursor, length, byteCursor int) fieldInfo
	errMsg   string
}

var fieldRules = map[TimeType]fieldRule{
	Literal:      {1, -1, func(rl []rune, c, l, _ int) fieldInfo { return fieldInfo{kind: kindLiteral, text: string(rl[c : c+l])} }, ""},
	Year:         {4, 4, func(_ []rune, _ int, l, bc int) fieldInfo { return fieldInfo{kind: kindYear, cursor: bc, length: l} }, "year out of range"},
	Month:        {2, 2, func(_ []rune, _ int, l, bc int) fieldInfo { return fieldInfo{kind: kindMonth, cursor: bc, length: l} }, "month out of range"},
	Day:          {2, 2, func(_ []rune, _ int, l, bc int) fieldInfo { return fieldInfo{kind: kindDay, cursor: bc, length: l} }, "day out of range"},
	Hour:         {2, 2, func(_ []rune, _ int, l, bc int) fieldInfo { return fieldInfo{kind: kindHour, cursor: bc, length: l} }, "hour out of range"},
	Minute:       {2, 2, func(_ []rune, _ int, l, bc int) fieldInfo { return fieldInfo{kind: kindMinute, cursor: bc, length: l} }, "minute out of range"},
	Second:       {2, 2, func(_ []rune, _ int, l, bc int) fieldInfo { return fieldInfo{kind: kindSecond, cursor: bc, length: l} }, "second out of range"},
	NanoOfSecond: {2, 9, func(_ []rune, _ int, l, bc int) fieldInfo { return fieldInfo{kind: kindNano, cursor: bc, length: l} }, "nano out of range"},
}

func countToken(layout []rune, start int, c rune) int {
	count := 0
	for _, r := range layout[start:] {
		if r != c {
			break
		}
		count++
	}
	return count
}

func countLiteral(layout []rune, start int) int {
	count := 0
	for _, r := range layout[start:] {
		if _, ok := TyMap[TimeType(r)]; ok {
			break
		}
		count++
	}
	return count
}

func formatter(layoutS string) ([]fieldInfo, error) {
	layout := []rune(layoutS)
	var fields []fieldInfo
	byteCursor := 0
	for i := 0; i < len(layout); {
		c := layout[i]
		ty := TimeType(c)

		var length int
		if _, ok := TyMap[ty]; ok {
			length = countToken(layout, i, c)
		} else {
			length = countLiteral(layout, i)
			ty = Literal
		}

		rule := fieldRules[ty]
		if length < rule.minLen || (rule.maxLen != -1 && length > rule.maxLen) {
			return nil, errors.New(rule.errMsg)
		}
		fields = append(fields, rule.newField(layout, i, length, byteCursor))

		for j := i; j < i+length; j++ {
			byteCursor += utf8.RuneLen(layout[j])
		}

		i += length
	}
	return fields, nil
}
