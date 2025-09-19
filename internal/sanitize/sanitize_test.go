package sanitize

import (
	"testing"
)

func TestSanitizer_SanitizeText(t *testing.T) {
	s := NewSanitizer()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "clean text",
			input:    "This is clean text",
			expected: "This is clean text",
		},
		{
			name:     "text with null bytes",
			input:    "Text with\x00null bytes",
			expected: "Text withnull bytes",
		},
		{
			name:     "text with HTML",
			input:    "Text with <script>alert('xss')</script>",
			expected: "Text with &lt;script&gt;alert(&#39;xss&#39;)&lt;/script&gt;",
		},
		{
			name:     "text with control characters",
			input:    "Text\x01with\x02control\x03chars",
			expected: "Textwithcontrolchars",
		},
		{
			name:     "text with excessive whitespace",
			input:    "Text   with    multiple     spaces",
			expected: "Text with multiple spaces",
		},
		{
			name:     "text with line breaks",
			input:    "Text\nwith\nline\nbreaks",
			expected: "Text with line breaks",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := s.SanitizeText(tt.input)
			if result != tt.expected {
				t.Errorf("SanitizeText() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestSanitizer_SanitizeTitle(t *testing.T) {
	s := NewSanitizer()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "clean title",
			input:    "Clean Task Title",
			expected: "Clean Task Title",
		},
		{
			name:     "title with line breaks",
			input:    "Title\nwith\nline\nbreaks",
			expected: "Title with line breaks",
		},
		{
			name:     "title with tabs",
			input:    "Title\twith\ttabs",
			expected: "Title with tabs",
		},
		{
			name:     "title with excessive whitespace",
			input:    "  Title   with   spaces  ",
			expected: "Title with spaces",
		},
		{
			name:     "title with HTML",
			input:    "Title with <b>bold</b> text",
			expected: "Title with &lt;b&gt;bold&lt;/b&gt; text",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := s.SanitizeTitle(tt.input)
			if result != tt.expected {
				t.Errorf("SanitizeTitle() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestSanitizer_SanitizeLabel(t *testing.T) {
	s := NewSanitizer()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "clean label",
			input:    "bug",
			expected: "bug",
		},
		{
			name:     "label with spaces",
			input:    "high priority",
			expected: "high-priority",
		},
		{
			name:     "label with uppercase",
			input:    "Frontend",
			expected: "frontend",
		},
		{
			name:     "label with special characters",
			input:    "bug@fix!",
			expected: "bugfix",
		},
		{
			name:     "label with multiple hyphens",
			input:    "high---priority",
			expected: "high-priority",
		},
		{
			name:     "label with leading/trailing separators",
			input:    "-_bug_-",
			expected: "bug",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := s.SanitizeLabel(tt.input)
			if result != tt.expected {
				t.Errorf("SanitizeLabel() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestSanitizer_SanitizeTaskID(t *testing.T) {
	s := NewSanitizer()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "clean task ID",
			input:    "T1.2.3",
			expected: "T1.2.3",
		},
		{
			name:     "task ID without T prefix",
			input:    "1.2.3",
			expected: "1.2.3",
		},
		{
			name:     "task ID with spaces",
			input:    " T1.2.3 ",
			expected: "T1.2.3",
		},
		{
			name:     "task ID with invalid characters",
			input:    "T1.2.3@invalid",
			expected: "T1.2.3",
		},
		{
			name:     "task ID with null bytes",
			input:    "T1\x00.2.3",
			expected: "T1.2.3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := s.SanitizeTaskID(tt.input)
			if result != tt.expected {
				t.Errorf("SanitizeTaskID() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestSanitizer_SanitizeFilePath(t *testing.T) {
	s := NewSanitizer()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "clean path",
			input:    "/home/user/tasks",
			expected: "/home/user/tasks",
		},
		{
			name:     "path with traversal",
			input:    "/home/user/../../../etc/passwd",
			expected: "/home/user/etc/passwd",
		},
		{
			name:     "path with multiple slashes",
			input:    "/home//user///tasks",
			expected: "/home/user/tasks",
		},
		{
			name:     "path with null bytes",
			input:    "/home\x00/user/tasks",
			expected: "/home/user/tasks",
		},
		{
			name:     "path with spaces",
			input:    "  /home/user/tasks  ",
			expected: "/home/user/tasks",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := s.SanitizeFilePath(tt.input)
			if result != tt.expected {
				t.Errorf("SanitizeFilePath() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestSanitizer_RemoveUnsafeCharacters(t *testing.T) {
	s := NewSanitizer()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "clean text",
			input:    "This is safe text",
			expected: "This is safe text",
		},
		{
			name:     "script tags",
			input:    "Text with <script>alert('xss')</script> content",
			expected: "Text with  content",
		},
		{
			name:     "javascript protocol",
			input:    "Click javascript:alert('xss') here",
			expected: "Click  here",
		},
		{
			name:     "event handlers",
			input:    "Text with onload=alert('xss') handler",
			expected: "Text with  handler",
		},
		{
			name:     "hex escapes",
			input:    "Text with \\x41\\x42 escapes",
			expected: "Text with  escapes",
		},
		{
			name:     "null bytes",
			input:    "Text with\x00null bytes",
			expected: "Text withnull bytes",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := s.RemoveUnsafeCharacters(tt.input)
			if result != tt.expected {
				t.Errorf("RemoveUnsafeCharacters() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestSanitizer_TruncateString(t *testing.T) {
	s := NewSanitizer()

	tests := []struct {
		name      string
		input     string
		maxLength int
		expected  string
	}{
		{
			name:      "short string",
			input:     "Short",
			maxLength: 10,
			expected:  "Short",
		},
		{
			name:      "exact length",
			input:     "Exact",
			maxLength: 5,
			expected:  "Exact",
		},
		{
			name:      "truncate string",
			input:     "This is a long string",
			maxLength: 10,
			expected:  "This is a ",
		},
		{
			name:      "zero length",
			input:     "Any string",
			maxLength: 0,
			expected:  "",
		},
		{
			name:      "UTF-8 string",
			input:     "Hello 世界",
			maxLength: 7,
			expected:  "Hello 世",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := s.TruncateString(tt.input, tt.maxLength)
			if result != tt.expected {
				t.Errorf("TruncateString() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestSanitizer_SanitizeSlice(t *testing.T) {
	s := NewSanitizer()

	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "clean slice",
			input:    []string{"label1", "label2", "label3"},
			expected: []string{"label1", "label2", "label3"},
		},
		{
			name:     "slice with empty strings",
			input:    []string{"label1", "", "label3"},
			expected: []string{"label1", "label3"},
		},
		{
			name:     "slice with spaces",
			input:    []string{"  label1  ", "label2", "  label3  "},
			expected: []string{"label1", "label2", "label3"},
		},
		{
			name:     "slice with special characters",
			input:    []string{"label@1", "label#2", "label$3"},
			expected: []string{"label1", "label2", "label3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := s.SanitizeSlice(tt.input, s.SanitizeLabel)
			if len(result) != len(tt.expected) {
				t.Errorf("SanitizeSlice() length = %d, want %d", len(result), len(tt.expected))
				return
			}
			for i, expected := range tt.expected {
				if result[i] != expected {
					t.Errorf("SanitizeSlice()[%d] = %q, want %q", i, result[i], expected)
				}
			}
		})
	}
}

func TestSanitizer_IsValidUTF8(t *testing.T) {
	s := NewSanitizer()

	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "valid UTF-8",
			input:    "Hello 世界",
			expected: true,
		},
		{
			name:     "valid ASCII",
			input:    "Hello World",
			expected: true,
		},
		{
			name:     "empty string",
			input:    "",
			expected: true,
		},
		{
			name:     "invalid UTF-8",
			input:    string([]byte{0xff, 0xfe, 0xfd}),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := s.IsValidUTF8(tt.input)
			if result != tt.expected {
				t.Errorf("IsValidUTF8() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestNormalizeWhitespace(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "single spaces",
			input:    "word1 word2 word3",
			expected: "word1 word2 word3",
		},
		{
			name:     "multiple spaces",
			input:    "word1   word2     word3",
			expected: "word1 word2 word3",
		},
		{
			name:     "mixed whitespace",
			input:    "word1\t\t\nword2\r\nword3",
			expected: "word1 word2 word3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeWhitespace(tt.input)
			if result != tt.expected {
				t.Errorf("normalizeWhitespace() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestNormalizeLineBreaks(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Unix line breaks",
			input:    "line1\nline2\nline3",
			expected: "line1\nline2\nline3",
		},
		{
			name:     "Windows line breaks",
			input:    "line1\r\nline2\r\nline3",
			expected: "line1\nline2\nline3",
		},
		{
			name:     "Old Mac line breaks",
			input:    "line1\rline2\rline3",
			expected: "line1\nline2\nline3",
		},
		{
			name:     "excessive line breaks",
			input:    "line1\n\n\n\nline2",
			expected: "line1\n\nline2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeLineBreaks(tt.input)
			if result != tt.expected {
				t.Errorf("normalizeLineBreaks() = %q, want %q", result, tt.expected)
			}
		})
	}
}