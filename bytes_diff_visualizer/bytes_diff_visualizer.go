package bytes_diff_visualizer

import (
	"bytes"
	"fmt"
	"strings"
	"text/tabwriter"
)

// VisualizeByteDiff visualizes the difference between two byte slices, returning lines to be presented to the user.
//
// The lines will include ANSI escape codes to colorize the output.
func VisualizeByteDiff(actual []byte, expected []byte) []string {
	// If both are exactly the same, return an empty slice
	if bytes.Equal(actual, expected) {
		return []string{}
	}

	// Find the index of the first differing byte
	var firstDiffIndex int = -1

	for i := 0; i < len(actual) && i < len(expected); i++ {
		if actual[i] != expected[i] {
			firstDiffIndex = i
			break
		}
	}

	// If we got here, then the first differing byte is the first byte of the longer slice
	if firstDiffIndex == -1 {
		if len(actual) > len(expected) {
			firstDiffIndex = len(expected)
		} else {
			firstDiffIndex = len(actual)
		}
	}

	byteCountToDisplay := 20
	byteRangeStart := intmax(0, firstDiffIndex-(byteCountToDisplay/2))
	byteRangeEnd := byteRangeStart + byteCountToDisplay

	linesBuffer := bytes.NewBuffer([]byte{})
	tabWriter := tabwriter.NewWriter(linesBuffer, 20, 0, 1, ' ', tabwriter.Debug)

	tabWriter.Write([]byte(fmt.Sprintf("Expected (bytes %v-%v), hexadecimal:\t Printable characters:\n", byteRangeStart, byteRangeEnd)))
	tabWriter.Write([]byte(fmt.Sprintf("%v\t %v\n", formatBytesAsHex(expected[byteRangeStart:intmin(byteRangeEnd, len(expected))], byteCountToDisplay), formatBytesAsAscii(expected[byteRangeStart:intmin(byteRangeEnd, len(expected))], byteCountToDisplay))))
	tabWriter.Write([]byte("\n"))
	tabWriter.Write([]byte(fmt.Sprintf("Actual (bytes %v-%v), hexadecimal:\t Printable characters:\n", byteRangeStart, byteRangeEnd)))
	tabWriter.Write([]byte(fmt.Sprintf("%v\t %v\n", formatBytesAsHex(actual[byteRangeStart:intmin(byteRangeEnd, len(actual))], byteCountToDisplay), formatBytesAsAscii(actual[byteRangeStart:intmin(byteRangeEnd, len(actual))], byteCountToDisplay))))
	tabWriter.Flush()

	return strings.Split(string(linesBuffer.Bytes()), "\n")
}

func formatBytesAsHexAndAscii(value []byte, expectedCount int) string {
	return fmt.Sprintf("%v | %v", formatBytesAsHex(value, expectedCount), formatBytesAsAscii(value, expectedCount))
}

func formatBytesAsAscii(value []byte, expectedCount int) string {
	var asciiRepresentations []string

	for i := 0; i < expectedCount; i++ {
		if i >= len(value) {
			// Pad with spaces if we're out of bytes
			asciiRepresentations = append(asciiRepresentations, " ")
		} else if value[i] < 32 || value[i] > 126 {
			// If the byte is not printable, replace it with a dot
			asciiRepresentations = append(asciiRepresentations, ".")
		} else {
			asciiRepresentations = append(asciiRepresentations, string(value[i]))
		}
	}

	return strings.Join(asciiRepresentations, "")
}

func formatBytesAsHex(value []byte, expectedCount int) string {
	var hexadecimalRepresentations []string

	for i := 0; i < expectedCount; i++ {
		if i >= len(value) {
			hexadecimalRepresentations = append(hexadecimalRepresentations, "  ")
		} else {
			hexadecimalRepresentations = append(hexadecimalRepresentations, fmt.Sprintf("%02x", value[i]))
		}
	}

	return strings.Join(hexadecimalRepresentations, " ")
}

func intmax(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

func intmin(a int, b int) int {
	if a < b {
		return a
	}
	return b
}
