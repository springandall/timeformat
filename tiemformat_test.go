package timeformat

import (
	"testing"
	"time"
)

func TestFormat(t *testing.T) {
	tm := time.Date(2024, 6, 15, 14, 30, 45, 123456789, time.UTC)

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
		{"empty", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Format(tt.pattern, tm)
			if err != nil {
				t.Fatalf("Format(%q) failed: %v", tt.pattern, err)
			}
			if got != tt.want {
				t.Errorf("Format(%q) = %q, want %q", tt.pattern, got, tt.want)
			}
		})
	}
}

func TestFormat_InvalidPattern(t *testing.T) {
	tm := time.Now()
	invalid := []string{"yyy", "MMM", "ddd", "HHH", "mmm", "sss", "S", "SSSSSSSSSS"}
	for _, p := range invalid {
		_, err := Format(p, tm)
		if err == nil {
			t.Errorf("Format(%q) expected error", p)
		}
	}
}

func TestFormat_Boundary(t *testing.T) {
	loc := time.UTC
	tests := []struct {
		name    string
		tm      time.Time
		pattern string
		want    string
	}{
		{"zero time", time.Time{}, "yyyy-MM-dd HH:mm:ss", "0001-01-01 00:00:00"},
		{"leap year", time.Date(2024, 2, 29, 0, 0, 0, 0, loc), "yyyy-MM-dd", "2024-02-29"},
		{"max year", time.Date(9999, 12, 31, 23, 59, 59, 999999999, loc), "yyyy-MM-dd HH:mm:ss.SSSSSSSSS", "9999-12-31 23:59:59.999999999"},
		{"first second 2020", time.Date(2020, 1, 1, 0, 0, 0, 0, loc), "yyyy-MM-dd HH:mm:ss", "2020-01-01 00:00:00"},
		{"no nano", time.Date(2024, 6, 15, 14, 30, 45, 0, loc), "yyyy-MM-dd HH:mm:ss.SSS", "2024-06-15 14:30:45.000"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Format(tt.pattern, tt.tm)
			if err != nil {
				t.Fatalf("Format(%q) failed: %v", tt.pattern, err)
			}
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
			name:    "date only",
			pattern: "yyyyMMdd",
			input:   "20240615",
			want:    time.Date(2024, 6, 15, 0, 0, 0, 0, loc),
		},
		{
			name:    "leap year",
			pattern: "yyyy-MM-dd",
			input:   "2024-02-29",
			want:    time.Date(2024, 2, 29, 0, 0, 0, 0, loc),
		},
		{
			name:    "max year",
			pattern: "yyyyMMdd",
			input:   "99991231",
			want:    time.Date(9999, 12, 31, 0, 0, 0, 0, loc),
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
		{
			name:    "err hour negative",
			pattern: "yyyy-MM-dd HH:mm:ss",
			input:   "2024-06-15 -1:30:45",
			wantErr: true,
		},
		{
			name:    "err minute negative",
			pattern: "yyyy-MM-dd HH:mm:ss",
			input:   "2024-06-15 14:-1:45",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.pattern, tt.input, loc)
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

func TestParse_ExtraChars(t *testing.T) {
	loc := time.UTC
	got, err := Parse("yyyy-MM-dd", "2024-06-15extra", loc)
	if err != nil {
		t.Fatalf("Parse() failed with extra chars: %v", err)
	}
	want := time.Date(2024, 6, 15, 0, 0, 0, 0, loc)
	if !got.Equal(want) {
		t.Errorf("Parse() = %v, want %v", got, want)
	}
}

func TestParse_ReturnsZeroTimeOnError(t *testing.T) {
	got, err := Parse("yyyy", "abc", time.UTC)
	if err == nil {
		t.Fatal("expected error")
	}
	if !got.IsZero() {
		t.Errorf("expected zero time on error, got %v", got)
	}
}

func TestParse_NilLocation(t *testing.T) {
	got, err := Parse("yyyy-MM-dd", "2024-06-15", nil)
	if err != nil {
		t.Fatalf("Parse() failed: %v", err)
	}
	if got.IsZero() {
		t.Error("expected non-zero time")
	}
	if got.Location() != time.UTC {
		t.Errorf("expected UTC location, got %v", got.Location())
	}
}

func TestParse_Location(t *testing.T) {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		t.Fatal(err)
	}
	got, err := Parse("yyyy-MM-dd HH:mm:ss", "2024-06-15 14:30:45", loc)
	if err != nil {
		t.Fatalf("Parse() failed: %v", err)
	}
	if got.Location() != loc {
		t.Errorf("expected %v location, got %v", loc, got.Location())
	}
	if got.Hour() != 14 {
		t.Errorf("expected hour 14, got %d", got.Hour())
	}
}

func TestParse_InvalidPattern(t *testing.T) {
	_, err := Parse("yyy", "2024", time.UTC)
	if err == nil {
		t.Error("Parse() with invalid pattern should error")
	}
}

func TestParse_EmptyPattern(t *testing.T) {
	_, err := Parse("", "", time.UTC)
	if err != nil {
		t.Fatal(err)
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
		{"yyyy年MM月dd日 HH:mm:ss", time.Second},
	}

	for _, tt := range tests {
		t.Run(tt.pattern, func(t *testing.T) {
			formatted, err := Format(tt.pattern, original)
			if err != nil {
				t.Fatalf("Format(%q) failed: %v", tt.pattern, err)
			}
			parsed, err := Parse(tt.pattern, formatted, loc)
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

func TestFormat_NanoPadding(t *testing.T) {
	loc := time.UTC
	tests := []struct {
		name    string
		nsec    int
		pattern string
		want    string
	}{
		{"zero nano", 0, "SSS", "000"},
		{"1 ns", 1, "SSS", "000"},
		{"1 us", 1000, "SSS", "000"},
		{"1 ms", 1000000, "SSS", "001"},
		{"full 9 digits", 123456789, "SSSSSSSSS", "123456789"},
		{"2 digits trunc", 1000000, "SS", "00"},
		{"6 digits trunc", 1, "SSSSSS", "000000"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tm := time.Date(2024, 1, 1, 0, 0, 0, tt.nsec, loc)
			got, err := Format(tt.pattern, tm)
			if err != nil {
				t.Fatal(err)
			}
			if got != tt.want {
				t.Errorf("Format = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestParse_NanoPadding(t *testing.T) {
	loc := time.UTC
	tests := []struct {
		name     string
		pattern  string
		input    string
		wantNsec int
		wantErr  bool
	}{
		{"2 digits", "SS", "12", 120000000, false},
		{"3 digits", "SSS", "123", 123000000, false},
		{"6 digits", "SSSSSS", "123456", 123456000, false},
		{"9 digits", "SSSSSSSSS", "123456789", 123456789, false},
		{"2 zero pad", "SS", "00", 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.pattern, tt.input, loc)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Parse() error = %v, wantErr = %v", err, tt.wantErr)
			}
			if !tt.wantErr && got.Nanosecond() != tt.wantNsec {
				t.Errorf("nanosecond = %d, want %d", got.Nanosecond(), tt.wantNsec)
			}
		})
	}
}

var benchSink time.Time
var benchSink2 string

func BenchmarkParse_Custom(b *testing.B) {
	for i := 0; i < b.N; i++ {
		benchSink, _ = Parse("yyyy-MM-dd HH:mm:ss", "2024-06-15 14:30:45", time.UTC)
	}
}

func BenchmarkParse_Stdlib(b *testing.B) {
	for i := 0; i < b.N; i++ {
		benchSink, _ = time.Parse("2006-01-02 15:04:05", "2024-06-15 14:30:45")
	}
}

func BenchmarkParse_CustomSimple(b *testing.B) {
	for i := 0; i < b.N; i++ {
		benchSink, _ = Parse("yyyyMMdd", "20240615", time.UTC)
	}
}

func BenchmarkParse_StdlibSimple(b *testing.B) {
	for i := 0; i < b.N; i++ {
		benchSink, _ = time.Parse("20060102", "20240615")
	}
}

func BenchmarkParse_CustomWithNano(b *testing.B) {
	for i := 0; i < b.N; i++ {
		benchSink, _ = Parse("yyyy-MM-dd HH:mm:ss.SSSSSSSSS", "2024-06-15 14:30:45.123456789", time.UTC)
	}
}

func BenchmarkParse_StdlibWithNano(b *testing.B) {
	for i := 0; i < b.N; i++ {
		benchSink, _ = time.Parse("2006-01-02 15:04:05.000000000", "2024-06-15 14:30:45.123456789")
	}
}

func BenchmarkParse_CustomUnicode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		benchSink, _ = Parse("yyyy年MM月dd日 HH:mm:ss", "2024年06月15日 14:30:45", time.UTC)
	}
}

func BenchmarkFormat_Custom(b *testing.B) {
	tm := time.Date(2024, 6, 15, 14, 30, 45, 0, time.UTC)
	for i := 0; i < b.N; i++ {
		benchSink2, _ = Format("yyyy-MM-dd HH:mm:ss", tm)
	}
}

func BenchmarkFormat_Stdlib(b *testing.B) {
	tm := time.Date(2024, 6, 15, 14, 30, 45, 0, time.UTC)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchSink2 = tm.Format("2006-01-02 15:04:05")
	}
}

func BenchmarkFormat_CustomSimple(b *testing.B) {
	tm := time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC)
	for i := 0; i < b.N; i++ {
		benchSink2, _ = Format("yyyyMMdd", tm)
	}
}

func BenchmarkFormat_CustomWithNano(b *testing.B) {
	tm := time.Date(2024, 6, 15, 14, 30, 45, 123456789, time.UTC)
	for i := 0; i < b.N; i++ {
		benchSink2, _ = Format("yyyy-MM-dd HH:mm:ss.SSSSSSSSS", tm)
	}
}

func BenchmarkFormat_StdlibWithNano(b *testing.B) {
	tm := time.Date(2024, 6, 15, 14, 30, 45, 123456789, time.UTC)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchSink2 = tm.Format("2006-01-02 15:04:05.000000000")
	}
}

func BenchmarkFormat_CustomUnicode(b *testing.B) {
	tm := time.Date(2024, 6, 15, 14, 30, 45, 0, time.UTC)
	for i := 0; i < b.N; i++ {
		benchSink2, _ = Format("yyyy年MM月dd日 HH:mm:ss", tm)
	}
}
