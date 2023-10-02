package bytes_diff_visualizer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVisualizeByteDiffWorksWithStrings(t *testing.T) {
	actual := []byte("Hello, World!")
	expected := []byte("Hello, Go!")

	result := VisualizeByteDiff(actual, expected)

	if len(result) == 0 {
		t.Errorf("Expected a non-empty result")
	}

	expectedLines := []string{
		"Expected (bytes 0-20), hexadecimal:                         | Printable characters:",
		"48 65 6c 6c 6f 2c 20 47 6f 21                               | Hello, Go!          ",
		"",
		"Actual (bytes 0-20), hexadecimal:                           | Printable characters:",
		"48 65 6c 6c 6f 2c 20 57 6f 72 6c 64 21                      | Hello, World!       ",
	}

	for i, expectedLine := range expectedLines {
		if i >= len(result) {
			t.Fatalf("Expected %v lines, but only got %v", len(expectedLines), len(result))
		}

		assert.Equal(t, expectedLine, result[i])
	}
}

func TestVisualizeByteDiffWorksWithNonPrintableCharacters(t *testing.T) {
	actual := []byte("blob\000header")
	expected := []byte("blob\000\000header") // Has an extra null byte

	result := VisualizeByteDiff(actual, expected)

	if len(result) == 0 {
		t.Errorf("Expected a non-empty result")
	}

	expectedLines := []string{
		"Expected (bytes 0-20), hexadecimal:                         | Printable characters:",
		"62 6c 6f 62 00 00 68 65 61 64 65 72                         | blob..header        ",
		"",
		"Actual (bytes 0-20), hexadecimal:                           | Printable characters:",
		"62 6c 6f 62 00 68 65 61 64 65 72                            | blob.header         ",
	}

	for i, expectedLine := range expectedLines {
		if i >= len(result) {
			t.Fatalf("Expected %v lines, but only got %v", len(expectedLines), len(result))
		}

		assert.Equal(t, expectedLine, result[i])
	}
}

// func TestVisualizeByteDiffWorksWithLongerSequences(t *testing.T) {
// 	expected := []byte("1234567890123457890123457890abcd")
// 	actual := []byte("1234567890123457890123457890efgh")

// 	result := VisualizeByteDiff(actual, expected)

// 	if len(result) == 0 {
// 		t.Errorf("Expected a non-empty result")
// 	}

// 	expectedLines := []string{
// 		"Expected (bytes 18-38), hexadecimal:                        | Printable characters:",
// 		"31 32 33 34 35 36 37 38 39 30 31 32 33 34 35 37 38 39 30 61 | 1234567890123457890a",
// 		"",
// 		"Actual (bytes 18-38), hexadecimal:                          | Printable characters:",
// 		"62 6c 6f 62 00 68 65 61 64 65 72                            | blob.header         ",
// 	}

// 	for i, expectedLine := range expectedLines {
// 		if i >= len(result) {
// 			t.Fatalf("Expected %v lines, but only got %v", len(expectedLines), len(result))
// 		}

// 		assert.Equal(t, expectedLine, result[i])
// 	}
// }
