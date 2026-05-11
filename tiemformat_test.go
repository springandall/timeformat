package timeformat

import (
	"testing"
	"time"
)

func TestOfPattern_Valid(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		wantErr bool
	}{
		{"year", "yyyy", false},
		{"month", "MM", false},
		{"day", "dd", false},
		{"hour", "HH", false},
		{"minute", "mm", false},
		{"second", "ss", false},
		{"nano 2 digits", "SS", false},
		{"nano 9 digits", "SSSSSSSSS", false},
		{"full datetime", "yyyy-MM-dd HH:mm:ss", false},
		{"with nano", "yyyy-MM-dd HH:mm:ss.SSS", false},
		{"with literal", "yyyy年MM月dd日", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := OfPattern(tt.pattern)
			if (err != nil) != tt.wantErr {
				t.Errorf("OfPattern(%q) error = %v, wantErr = %v", tt.pattern, err, tt.wantErr)
			}
		})
	}
}

func TestOfPattern_Invalid(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
	}{
		{"year not 4", "yyy"},
		{"month not 2", "MMM"},
		{"day not 2", "ddd"},
		{"hour not 2", "HHH"},
		{"minute not 2", "mmm"},
		{"second not 2", "sss"},
		{"nano 1 digit", "S"},
		{"nano 10 digits", "SSSSSSSSSS"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := OfPattern(tt.pattern)
			if err == nil {
				t.Errorf("OfPattern(%q) expected error, got nil", tt.pattern)
			}
		})
	}
}

func TestFormat(t *testing.T) {
	loc := time.UTC
	tm := time.Date(2024, 6, 15, 14, 30, 45, 123456789, loc)

	tests := []struct {
		name    string
		pattern string
		want    string
	}{
		{"year", "yyyy", "2024"},
		{"month", "MM", "06"},
		{"day", "dd", "15"},
		{"hour", "HH", "14"},
		{"minute", "mm", "30"},
		{"second", "ss", "45"},
		{"nano 3 digits", "SSS", "123"},
		{"nano 6 digits", "SSSSSS", "123456"},
		{"nano 9 digits", "SSSSSSSSS", "123456789"},
		{"full datetime", "yyyy-MM-dd HH:mm:ss", "2024-06-15 14:30:45"},
		{"with nano", "yyyy-MM-dd HH:mm:ss.SSS", "2024-06-15 14:30:45.123"},
		{"with literals", "yyyy年MM月dd日", "2024年06月15日"},
		{"only literal", "T=", "T="},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := OfPattern(tt.pattern)
			if err != nil {
				t.Fatalf("OfPattern(%q) failed: %v", tt.pattern, err)
			}
			got := p.Format(tm)
			if got != tt.want {
				t.Errorf("Format(%q) = %q, want %q", tt.pattern, got, tt.want)
			}
		})
	}
}

func TestParse(t *testing.T) {
	loc := time.UTC

	tests := []struct {
		name    string
		pattern string
		input   string
		want    time.Time
		wantErr bool
	}{
		{
			name:    "full datetime",
			pattern: "yyyy-MM-dd HH:mm:ss",
			input:   "2024-06-15 14:30:45",
			want:    time.Date(2024, 6, 15, 14, 30, 45, 0, loc),
		},
		{
			name:    "with nano",
			pattern: "yyyy-MM-dd HH:mm:ss.SSS",
			input:   "2024-06-15 14:30:45.123",
			want:    time.Date(2024, 6, 15, 14, 30, 45, 123000000, loc),
		},
		{
			name:    "nano full 9 digits",
			pattern: "yyyy-MM-dd HH:mm:ss.SSSSSSSSS",
			input:   "2024-06-15 14:30:45.123456789",
			want:    time.Date(2024, 6, 15, 14, 30, 45, 123456789, loc),
		},
		{
			name:    "with literals",
			pattern: "yyyy年MM月dd日",
			input:   "2024年06月15日",
			want:    time.Date(2024, 6, 15, 0, 0, 0, 0, loc),
		},
		{
			name:    "single digit month day",
			pattern: "yyyy-MM-dd",
			input:   "2024-06-05",
			want:    time.Date(2024, 6, 5, 0, 0, 0, 0, loc),
		},
		{
			name:    "date only",
			pattern: "yyyyMMdd",
			input:   "20240615",
			want:    time.Date(2024, 6, 15, 0, 0, 0, 0, loc),
		},
		{
			name:    "err invalid month",
			pattern: "yyyy-MM-dd",
			input:   "2024-13-15",
			wantErr: true,
		},
		{
			name:    "err invalid day",
			pattern: "yyyy-MM-dd",
			input:   "2024-06-32",
			wantErr: true,
		},
		{
			name:    "err invalid hour",
			pattern: "yyyy-MM-dd HH:mm:ss",
			input:   "2024-06-15 24:30:45",
			wantErr: true,
		},
		{
			name:    "err invalid minute",
			pattern: "yyyy-MM-dd HH:mm:ss",
			input:   "2024-06-15 14:60:45",
			wantErr: true,
		},
		{
			name:    "err invalid second",
			pattern: "yyyy-MM-dd HH:mm:ss",
			input:   "2024-06-15 14:30:60",
			wantErr: true,
		},
		{
			name:    "err input too short",
			pattern: "yyyy-MM-dd",
			input:   "2024-06-1",
			wantErr: true,
		},
		{
			name:    "err non-numeric",
			pattern: "yyyy",
			input:   "abcd",
			wantErr: true,
		},
		{
			name:    "err month zero",
			pattern: "yyyy-MM-dd",
			input:   "2024-00-15",
			wantErr: true,
		},
		{
			name:    "err day zero",
			pattern: "yyyy-MM-dd",
			input:   "2024-06-00",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := OfPattern(tt.pattern)
			if err != nil {
				t.Fatalf("OfPattern(%q) failed: %v", tt.pattern, err)
			}
			got, err := p.Parse(tt.input, loc)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse(%q) error = %v, wantErr = %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr && !got.Equal(tt.want) {
				t.Errorf("Parse(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestParse_ReturnsZeroTimeOnError(t *testing.T) {
	p, err := OfPattern("yyyy")
	if err != nil {
		t.Fatal(err)
	}
	got, err := p.Parse("abc", time.UTC)
	if err == nil {
		t.Fatal("expected error")
	}
	if !got.IsZero() {
		t.Errorf("expected zero time on error, got %v", got)
	}
}

func TestRoundtrip(t *testing.T) {
	loc := time.UTC
	original := time.Date(2024, 6, 15, 14, 30, 45, 987654321, loc)

	tests := []struct {
		pattern string
		trunc   time.Duration
	}{
		{"yyyy-MM-dd HH:mm:ss", time.Second},
		{"yyyy-MM-dd HH:mm:ss.SSS", time.Millisecond},
		{"yyyy-MM-dd HH:mm:ss.SSSSSS", time.Microsecond},
		{"yyyy-MM-dd HH:mm:ss.SSSSSSSSS", time.Nanosecond},
		{"yyyyMMdd", 24 * time.Hour},
	}

	for _, tt := range tests {
		t.Run(tt.pattern, func(t *testing.T) {
			p, err := OfPattern(tt.pattern)
			if err != nil {
				t.Fatalf("OfPattern(%q) failed: %v", tt.pattern, err)
			}
			formatted := p.Format(original)
			parsed, err := p.Parse(formatted, loc)
			if err != nil {
				t.Fatalf("Parse(%q) failed: %v", formatted, err)
			}

			want := original.Truncate(tt.trunc)
			if !parsed.Equal(want) {
				t.Errorf("roundtrip %q: got %v, want %v", tt.pattern, parsed, want)
			}
		})
	}
}

func TestOfPattern_EmptyPattern(t *testing.T) {
	p, err := OfPattern("")
	if err != nil {
		t.Fatalf("OfPattern(\"\") failed: %v", err)
	}
	if p.Format(time.Now()) != "" {
		t.Error("expected empty string from empty pattern")
	}
}

func TestFormat_NanoPadding(t *testing.T) {
	loc := time.UTC
	tests := []struct {
		nsec    int
		pattern string
		want    string
	}{
		{0, "SSS", "000"},
		{1, "SSS", "000"},
		{1000, "SSS", "000"},
		{1000000, "SSS", "001"},
		{123456789, "SSSSSSSSS", "123456789"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			tm := time.Date(2024, 1, 1, 0, 0, 0, tt.nsec, loc)
			p, err := OfPattern(tt.pattern)
			if err != nil {
				t.Fatal(err)
			}
			got := p.Format(tm)
			if got != tt.want {
				t.Errorf("Format = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestTyMap(t *testing.T) {
	expected := map[TimeType]struct{}{
		Year:         {},
		Month:        {},
		Day:          {},
		Hour:         {},
		Minute:       {},
		Second:       {},
		NanoOfSecond: {},
	}
	if len(TyMap) != len(expected) {
		t.Errorf("TyMap length = %d, want %d", len(TyMap), len(expected))
	}
	for k, v := range expected {
		if _, ok := TyMap[k]; !ok {
			t.Errorf("TyMap missing key %c", k)
		}
		_ = v
	}
}

func TestFormatFunc(t *testing.T) {
	tm := time.Date(2024, 6, 15, 14, 30, 45, 0, time.UTC)
	got, _ := Format(tm, "yyyy-MM-dd")
	want := "2024-06-15"
	if got != want {
		t.Errorf("Format() = %q, want %q", got, want)
	}
}

func TestParseFunc(t *testing.T) {
	got, err := Parse("2024-06-15", "yyyy-MM-dd", time.UTC)
	if err != nil {
		t.Fatalf("Parse() failed: %v", err)
	}
	want := time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC)
	if !got.Equal(want) {
		t.Errorf("Parse() = %v, want %v", got, want)
	}
}

func TestParseFunc_InvalidPattern(t *testing.T) {
	_, err := Parse("2024", "yyy", time.UTC)
	if err == nil {
		t.Error("Parse() with invalid pattern should error")
	}
}

func TestParseFunc_InvalidValue(t *testing.T) {
	_, err := Parse("abc", "yyyy", time.UTC)
	if err == nil {
		t.Error("Parse() with invalid value should error")
	}
}

func TestParse_NilLocation(t *testing.T) {
	got, err := Parse("2024-06-15", "yyyy-MM-dd", nil)
	if err != nil {
		t.Fatalf("Parse() failed: %v", err)
	}
	if got.IsZero() {
		t.Error("expected non-zero time")
	}
}

var benchSink time.Time

func BenchmarkParse_Custom(b *testing.B) {
	p, err := OfPattern("yyyy-MM-dd HH:mm:ss")
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchSink, _ = p.Parse("2024-06-15 14:30:45", time.UTC)
	}
}

func BenchmarkParse_Stdlib(b *testing.B) {
	for i := 0; i < b.N; i++ {
		benchSink, _ = time.Parse("2006-01-02 15:04:05", "2024-06-15 14:30:45")
	}
}

func BenchmarkParse_CustomCached(b *testing.B) {
	p, err := OfPattern("yyyy-MM-dd HH:mm:ss")
	if err != nil {
		b.Fatal(err)
	}
	loc := time.UTC
	val := "2024-06-15 14:30:45"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchSink, _ = p.Parse(val, loc)
	}
}

func BenchmarkParse_StdlibCached(b *testing.B) {
	layout := "2006-01-02 15:04:05"
	val := "2024-06-15 14:30:45"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchSink, _ = time.Parse(layout, val)
	}
}

func BenchmarkParse_CustomWithNano(b *testing.B) {
	p, err := OfPattern("yyyy-MM-dd HH:mm:ss.SSSSSSSSS")
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchSink, _ = p.Parse("2024-06-15 14:30:45.123456789", time.UTC)
	}
}

func BenchmarkParse_StdlibWithNano(b *testing.B) {
	for i := 0; i < b.N; i++ {
		benchSink, _ = time.Parse("2006-01-02 15:04:05.000000000", "2024-06-15 14:30:45.123456789")
	}
}

func BenchmarkFormat_Custom(b *testing.B) {
	p, err := OfPattern("yyyy-MM-dd HH:mm:ss")
	if err != nil {
		b.Fatal(err)
	}
	tm := time.Date(2024, 6, 15, 14, 30, 45, 0, time.UTC)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchSink2 = p.Format(tm)
	}
}

var benchSink2 string

func BenchmarkFormat_Stdlib(b *testing.B) {
	tm := time.Date(2024, 6, 15, 14, 30, 45, 0, time.UTC)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchSink2 = tm.Format("2006-01-02 15:04:05")
	}
}
