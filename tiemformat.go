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

type timeType rune

const (
	literal      timeType = 0
	year         timeType = 'y'
	month        timeType = 'M'
	day          timeType = 'd'
	hour         timeType = 'H'
	minute       timeType = 'm'
	second       timeType = 's'
	nanoOfSecond timeType = 'S'
)

var tyMap = map[timeType]struct{}{
	year:         {},
	month:        {},
	day:          {},
	hour:         {},
	minute:       {},
	second:       {},
	nanoOfSecond: {},
}

type fieldInfo struct {
	kind   timeType
	cursor int
	length int
	text   string
}

type dateArgs struct {
	Year  int
	Month time.Month
	Day   int
	Hour  int
	Min   int
	Sec   int
	Nsec  int
	Loc   *time.Location
}

func (d *dateArgs) time() time.Time {
	return time.Date(d.Year, d.Month, d.Day, d.Hour, d.Min, d.Sec, d.Nsec, d.Loc)
}

type TimeParser struct {
	fields []fieldInfo
}

func (p *TimeParser) Format(t time.Time) string {
	buf := make([]byte, 0, 64)
	for _, f := range p.fields {
		switch f.kind {
		case literal:
			buf = append(buf, f.text...)
		case year:
			buf = append4Digits(buf, t.Year())
		case month:
			buf = append2Digits(buf, int(t.Month()))
		case day:
			buf = append2Digits(buf, t.Day())
		case hour:
			buf = append2Digits(buf, t.Hour())
		case minute:
			buf = append2Digits(buf, t.Minute())
		case second:
			buf = append2Digits(buf, t.Second())
		case nanoOfSecond:
			buf = appendNano(buf, t.Nanosecond(), f.length)
		}
	}
	return string(buf)
}

func (p *TimeParser) Parse(layout string, loc *time.Location) (time.Time, error) {
	var args dateArgs
	for _, f := range p.fields {
		switch f.kind {
		case literal:
		case year:
			v, err := atoiStr(layout, f.cursor, f.length)
			if err != nil {
				return time.Time{}, err
			}
			args.Year = v
		case month:
			v, err := atoiStr(layout, f.cursor, f.length)
			if err != nil {
				return time.Time{}, err
			}
			if v < 1 || v > 12 {
				return time.Time{}, errors.New("month out of range")
			}
			args.Month = time.Month(v)
		case day:
			v, err := atoiStr(layout, f.cursor, f.length)
			if err != nil {
				return time.Time{}, err
			}
			if v < 1 || v > 31 {
				return time.Time{}, errors.New("day out of range")
			}
			args.Day = v
		case hour:
			v, err := atoiStr(layout, f.cursor, f.length)
			if err != nil {
				return time.Time{}, err
			}
			if v < 0 || v > 23 {
				return time.Time{}, errors.New("hour out of range")
			}
			args.Hour = v
		case minute:
			v, err := atoiStr(layout, f.cursor, f.length)
			if err != nil {
				return time.Time{}, err
			}
			if v < 0 || v > 59 {
				return time.Time{}, errors.New("minute out of range")
			}
			args.Min = v
		case second:
			v, err := atoiStr(layout, f.cursor, f.length)
			if err != nil {
				return time.Time{}, err
			}
			if v < 0 || v > 59 {
				return time.Time{}, errors.New("second out of range")
			}
			args.Sec = v
		case nanoOfSecond:
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
		loc = time.UTC
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

var fieldRules = map[timeType]fieldRule{
	literal: {1, -1, func(rl []rune, c, l, _ int) fieldInfo { return fieldInfo{kind: literal, text: string(rl[c : c+l])} }, ""},
	year:    {4, 4, func(_ []rune, _ int, l, bc int) fieldInfo { return fieldInfo{kind: year, cursor: bc, length: l} }, "year out of range"},
	month:   {2, 2, func(_ []rune, _ int, l, bc int) fieldInfo { return fieldInfo{kind: month, cursor: bc, length: l} }, "month out of range"},
	day:     {2, 2, func(_ []rune, _ int, l, bc int) fieldInfo { return fieldInfo{kind: day, cursor: bc, length: l} }, "day out of range"},
	hour:    {2, 2, func(_ []rune, _ int, l, bc int) fieldInfo { return fieldInfo{kind: hour, cursor: bc, length: l} }, "hour out of range"},
	minute:  {2, 2, func(_ []rune, _ int, l, bc int) fieldInfo { return fieldInfo{kind: minute, cursor: bc, length: l} }, "minute out of range"},
	second:  {2, 2, func(_ []rune, _ int, l, bc int) fieldInfo { return fieldInfo{kind: second, cursor: bc, length: l} }, "second out of range"},
	nanoOfSecond: {2, 9, func(_ []rune, _ int, l, bc int) fieldInfo {
		return fieldInfo{kind: nanoOfSecond, cursor: bc, length: l}
	}, "nano out of range"},
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
		if _, ok := tyMap[timeType(r)]; ok {
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
		ty := timeType(c)

		var length int
		if _, ok := tyMap[ty]; ok {
			length = countToken(layout, i, c)
		} else {
			length = countLiteral(layout, i)
			ty = literal
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
