package timeformat

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func OfPattern(pattern string) (*TimeParser, error) {
	formats, err := formatter(pattern)
	if err != nil {
		return nil, err
	}
	return &TimeParser{formats}, nil
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

type DateArgs struct {
	Year  int
	Month time.Month
	Day   int
	Hour  int
	Min   int
	Sec   int
	Nsec  int

	Loc *time.Location
}

func (d *DateArgs) time() time.Time {
	return time.Date(
		d.Year,
		d.Month,
		d.Day,
		d.Hour,
		d.Min,
		d.Sec,
		d.Nsec,
		d.Loc,
	)
}

type TimeFormat interface {
	Format(time time.Time) string
	Parse(layout []rune, args *DateArgs) error
}

type YearFormat struct {
	cursor int
	length int
}

func (y *YearFormat) Format(time time.Time) string {
	return fmt.Sprintf("%04d", time.Year())
}

func (y *YearFormat) Parse(layout []rune, args *DateArgs) error {
	if len(layout[y.cursor:]) < y.length {
		return errors.New("length out of range")
	}
	year := layout[y.cursor : y.cursor+y.length]
	yearNum, err := strconv.Atoi(string(year))
	if err != nil {
		return err
	}
	args.Year = yearNum
	return nil
}

type MonthFormat struct {
	cursor int
	length int
}

func (m *MonthFormat) Format(time time.Time) string {
	return fmt.Sprintf("%02d", time.Month())
}

func (m *MonthFormat) Parse(layout []rune, args *DateArgs) error {
	if len(layout[m.cursor:]) < m.length {
		return errors.New("length out of range")
	}
	month := layout[m.cursor : m.cursor+m.length]
	monthNum, err := strconv.Atoi(string(month))
	if err != nil {
		return err
	}
	if monthNum <= 0 || monthNum > 12 {
		return errors.New("month out of range")
	}
	args.Month = time.Month(monthNum)
	return nil
}

type DayFormat struct {
	cursor int
	length int
}

func (d *DayFormat) Format(time time.Time) string {
	return fmt.Sprintf("%02d", time.Day())

}

func (d *DayFormat) Parse(layout []rune, args *DateArgs) error {
	if len(layout[d.cursor:]) < d.length {
		return errors.New("length out of range")
	}
	day := layout[d.cursor : d.cursor+d.length]
	dayNum, err := strconv.Atoi(string(day))
	if err != nil {
		return err
	}
	if dayNum <= 0 || dayNum > 31 {
		return errors.New("day out of range")
	}
	args.Day = dayNum
	return nil
}

type HourFormat struct {
	cursor int
	length int
}

func (h *HourFormat) Format(time time.Time) string {
	return fmt.Sprintf("%02d", time.Hour())

}

func (h *HourFormat) Parse(layout []rune, args *DateArgs) error {
	if len(layout[h.cursor:]) < h.length {
		return errors.New("length out of range")
	}
	hour := layout[h.cursor : h.cursor+h.length]
	hourNum, err := strconv.Atoi(string(hour))
	if err != nil {
		return err
	}
	if hourNum < 0 || hourNum > 23 {
		return errors.New("hour out of range")
	}
	args.Hour = hourNum
	return nil
}

type MinuteFormat struct {
	cursor int
	length int
}

func (m *MinuteFormat) Format(time time.Time) string {
	return fmt.Sprintf("%02d", time.Minute())

}

func (m *MinuteFormat) Parse(layout []rune, args *DateArgs) error {
	if len(layout[m.cursor:]) < m.length {
		return errors.New("length out of range")
	}
	minute := layout[m.cursor : m.cursor+m.length]
	minuteNum, err := strconv.Atoi(string(minute))
	if err != nil {
		return err
	}
	if minuteNum < 0 || minuteNum > 59 {
		return errors.New("minute out of range")
	}
	args.Min = minuteNum
	return nil
}

type SecondFormat struct {
	cursor int
	length int
}

func (s *SecondFormat) Format(time time.Time) string {
	return fmt.Sprintf("%02d", time.Second())

}

func (s *SecondFormat) Parse(layout []rune, args *DateArgs) error {
	if len(layout[s.cursor:]) < s.length {
		return errors.New("length out of range")
	}
	second := layout[s.cursor : s.cursor+s.length]
	secondNum, err := strconv.Atoi(string(second))
	if err != nil {
		return err
	}
	if secondNum < 0 || secondNum > 59 {
		return errors.New("second out of range")
	}
	args.Sec = secondNum
	return nil
}

type NanoOfSecondFormat struct {
	cursor int
	length int
}

func (n *NanoOfSecondFormat) Format(time time.Time) string {
	s := fmt.Sprintf("%09d", time.Nanosecond())
	return s[:n.length]
}

func (n *NanoOfSecondFormat) Parse(layout []rune, args *DateArgs) error {
	if len(layout[n.cursor:]) < n.length {
		return errors.New("length out of range")
	}
	nano := layout[n.cursor : n.cursor+n.length]
	nanoNum, err := strconv.Atoi(string(nano))
	if err != nil {
		return err
	}
	for i := n.length; i < 9; i++ {
		nanoNum *= 10
	}
	if nanoNum < 0 || nanoNum > 999999999 {
		return errors.New("nanoNum out of range")
	}
	args.Nsec = nanoNum
	return nil
}

type LiteralFormat struct {
	text []rune
}

func (l *LiteralFormat) Format(time time.Time) string {
	return string(l.text)
}

func (l *LiteralFormat) Parse(_ []rune, _ *DateArgs) error {
	return nil
}

type TimeParser struct {
	formatter []TimeFormat
}

func (t *TimeParser) Format(time time.Time) string {
	var builder strings.Builder
	for _, f := range t.formatter {
		builder.WriteString(f.Format(time))
	}
	return builder.String()
}
func (t *TimeParser) Parse(format string, loc *time.Location) (time.Time, error) {
	var args DateArgs
	v := []rune(format)
	for _, f := range t.formatter {
		err := f.Parse(v, &args)
		if err != nil {
			return time.Time{}, err
		}
	}
	if loc == nil {
		loc = time.Local
	}
	args.Loc = loc
	return args.time(), nil
}

type typeRule struct {
	minLen    int
	maxLen    int
	newFormat func(rlayout []rune, cursor, length int) TimeFormat
	errMsg    string
}

var typeRules = map[TimeType]typeRule{
	Literal:      {1, -1, func(rl []rune, c, l int) TimeFormat { return &LiteralFormat{text: rl[c : c+l]} }, ""},
	Year:         {4, 4, func(_ []rune, c, l int) TimeFormat { return &YearFormat{c, l} }, "year out of range"},
	Month:        {2, 2, func(_ []rune, c, l int) TimeFormat { return &MonthFormat{c, l} }, "month out of range"},
	Day:          {2, 2, func(_ []rune, c, l int) TimeFormat { return &DayFormat{c, l} }, "day out of range"},
	Hour:         {2, 2, func(_ []rune, c, l int) TimeFormat { return &HourFormat{c, l} }, "hour out of range"},
	Minute:       {2, 2, func(_ []rune, c, l int) TimeFormat { return &MinuteFormat{c, l} }, "minute out of range"},
	Second:       {2, 2, func(_ []rune, c, l int) TimeFormat { return &SecondFormat{c, l} }, "second out of range"},
	NanoOfSecond: {2, 9, func(_ []rune, c, l int) TimeFormat { return &NanoOfSecondFormat{c, l} }, "nano out of range"},
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

func formatter(layoutS string) ([]TimeFormat, error) {
	layout := []rune(layoutS)
	var formats []TimeFormat
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

		rule := typeRules[ty]
		if length < rule.minLen || (rule.maxLen != -1 && length > rule.maxLen) {
			return nil, errors.New(rule.errMsg)
		}
		formats = append(formats, rule.newFormat(layout, i, length))

		i += length
	}
	return formats, nil
}
