package timeformat

import (
	"errors"
	"time"
)

const (
	DateTime string = "yyyy-MM-dd HH:mm:ss"
	DateOnly        = "yyyy-MM-dd"
	TimeOnly        = "HH:mm:ss"
)

func Format(pattern string, t time.Time) (string, error) {
	buf := make([]byte, 0, len(pattern))
	for i := 0; i < len(pattern); {
		switch pattern[i] {
		case 'y':
			j := i + 1
			for j < len(pattern) && pattern[j] == 'y' {
				j++
			}
			if j-i != 4 {
				return "", errors.New("year must be 4 digits")
			}
			buf = append4Digits(buf, t.Year())
			i = j
		case 'M':
			j := i + 1
			for j < len(pattern) && pattern[j] == 'M' {
				j++
			}
			if j-i != 2 {
				return "", errors.New("month must be 2 digits")
			}
			buf = append2Digits(buf, int(t.Month()))
			i = j
		case 'd':
			j := i + 1
			for j < len(pattern) && pattern[j] == 'd' {
				j++
			}
			if j-i != 2 {
				return "", errors.New("day must be 2 digits")
			}
			buf = append2Digits(buf, t.Day())
			i = j
		case 'H':
			j := i + 1
			for j < len(pattern) && pattern[j] == 'H' {
				j++
			}
			if j-i != 2 {
				return "", errors.New("hour must be 2 digits")
			}
			buf = append2Digits(buf, t.Hour())
			i = j
		case 'm':
			j := i + 1
			for j < len(pattern) && pattern[j] == 'm' {
				j++
			}
			if j-i != 2 {
				return "", errors.New("minute must be 2 digits")
			}
			buf = append2Digits(buf, t.Minute())
			i = j
		case 's':
			j := i + 1
			for j < len(pattern) && pattern[j] == 's' {
				j++
			}
			if j-i != 2 {
				return "", errors.New("second must be 2 digits")
			}
			buf = append2Digits(buf, t.Second())
			i = j
		case 'S':
			j := i + 1
			for j < len(pattern) && pattern[j] == 'S' {
				j++
			}
			length := j - i
			if length < 2 || length > 9 {
				return "", errors.New("nano must be 2-9 digits")
			}
			buf = appendNano(buf, t.Nanosecond(), length)
			i = j
		default:
			j := i + 1
			for j < len(pattern) && !isTimeByte(pattern[j]) {
				j++
			}
			buf = append(buf, pattern[i:j]...)
			i = j
		}
	}
	return string(buf), nil
}

func Parse(pattern, value string, loc *time.Location) (time.Time, error) {
	var args dateArgs
	bytePos := 0
	for i := 0; i < len(pattern); {
		switch pattern[i] {
		case 'y':
			j := i + 1
			for j < len(pattern) && pattern[j] == 'y' {
				j++
			}
			length := j - i
			if length != 4 {
				return time.Time{}, errors.New("year must be 4 digits")
			}
			v, err := atoiStr(value, bytePos, length)
			if err != nil {
				return time.Time{}, err
			}
			args.Year = v
			bytePos += length
			i = j
		case 'M':
			j := i + 1
			for j < len(pattern) && pattern[j] == 'M' {
				j++
			}
			length := j - i
			if length != 2 {
				return time.Time{}, errors.New("month must be 2 digits")
			}
			v, err := atoiStr(value, bytePos, length)
			if err != nil {
				return time.Time{}, err
			}
			if v < 1 || v > 12 {
				return time.Time{}, errors.New("month out of range")
			}
			args.Month = time.Month(v)
			bytePos += length
			i = j
		case 'd':
			j := i + 1
			for j < len(pattern) && pattern[j] == 'd' {
				j++
			}
			length := j - i
			if length != 2 {
				return time.Time{}, errors.New("day must be 2 digits")
			}
			v, err := atoiStr(value, bytePos, length)
			if err != nil {
				return time.Time{}, err
			}
			if v < 1 || v > 31 {
				return time.Time{}, errors.New("day out of range")
			}
			args.Day = v
			bytePos += length
			i = j
		case 'H':
			j := i + 1
			for j < len(pattern) && pattern[j] == 'H' {
				j++
			}
			length := j - i
			if length != 2 {
				return time.Time{}, errors.New("hour must be 2 digits")
			}
			v, err := atoiStr(value, bytePos, length)
			if err != nil {
				return time.Time{}, err
			}
			if v < 0 || v > 23 {
				return time.Time{}, errors.New("hour out of range")
			}
			args.Hour = v
			bytePos += length
			i = j
		case 'm':
			j := i + 1
			for j < len(pattern) && pattern[j] == 'm' {
				j++
			}
			length := j - i
			if length != 2 {
				return time.Time{}, errors.New("minute must be 2 digits")
			}
			v, err := atoiStr(value, bytePos, length)
			if err != nil {
				return time.Time{}, err
			}
			if v < 0 || v > 59 {
				return time.Time{}, errors.New("minute out of range")
			}
			args.Min = v
			bytePos += length
			i = j
		case 's':
			j := i + 1
			for j < len(pattern) && pattern[j] == 's' {
				j++
			}
			length := j - i
			if length != 2 {
				return time.Time{}, errors.New("second must be 2 digits")
			}
			v, err := atoiStr(value, bytePos, length)
			if err != nil {
				return time.Time{}, err
			}
			if v < 0 || v > 59 {
				return time.Time{}, errors.New("second out of range")
			}
			args.Sec = v
			bytePos += length
			i = j
		case 'S':
			j := i + 1
			for j < len(pattern) && pattern[j] == 'S' {
				j++
			}
			length := j - i
			if length < 2 || length > 9 {
				return time.Time{}, errors.New("nano must be 2-9 digits")
			}
			v, err := atoiStr(value, bytePos, length)
			if err != nil {
				return time.Time{}, err
			}
			for k := length; k < 9; k++ {
				v *= 10
			}
			if v < 0 || v > 999999999 {
				return time.Time{}, errors.New("nano out of range")
			}
			args.Nsec = v
			bytePos += length
			i = j
		default:
			j := i + 1
			for j < len(pattern) && !isTimeByte(pattern[j]) {
				j++
			}
			bytePos += j - i
			i = j
		}
	}
	if loc == nil {
		loc = time.UTC
	}
	args.Loc = loc
	return args.time(), nil
}

func isTimeByte(c byte) bool {
	return c == 'y' || c == 'M' || c == 'd' || c == 'H' || c == 'm' || c == 's' || c == 'S'
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
