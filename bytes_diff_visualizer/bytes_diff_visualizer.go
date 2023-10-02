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

	totalByteCountToDisplay := 100
	byteCountPerLine := 20
	byteRangeStart := intmax(0, firstDiffIndex-(totalByteCountToDisplay/2))
	byteRangeEnd := intmin(byteRangeStart+totalByteCountToDisplay, intmax(len(actual), len(expected)))

	linesBuffer := bytes.NewBuffer([]byte{})
	tabWriter := tabwriter.NewWriter(linesBuffer, 20, 0, 1, ' ', tabwriter.Debug)

	tabWriter.Write([]byte(fmt.Sprintf("Expected (bytes %v-%v), hexadecimal:       \t ASCII:\n", byteRangeStart, byteRangeEnd)))

	for i := byteRangeStart; i < intmin(byteRangeEnd, len(expected)); i += byteCountPerLine {
		end := intmin(i+byteCountPerLine, len(expected))
		tabWriter.Write([]byte(fmt.Sprintf("%v\t %v\n", formatBytesAsHex(expected[i:end]), formatBytesAsAscii(expected[i:end]))))
	}

	tabWriter.Write([]byte("\n"))
	tabWriter.Write([]byte(fmt.Sprintf("Actual (bytes %v-%v), hexadecimal:         \t ASCII:\n", byteRangeStart, byteRangeEnd)))

	for i := byteRangeStart; i < intmin(byteRangeEnd, len(actual)); i += byteCountPerLine {
		end := intmin(i+byteCountPerLine, len(actual))
		tabWriter.Write([]byte(fmt.Sprintf("%v\t %v\n", formatBytesAsHex(actual[i:end]), formatBytesAsAscii(actual[i:end]))))
	}

	tabWriter.Flush()

	return strings.Split(string(linesBuffer.Bytes()), "\n")
}

func formatBytesAsAscii(value []byte) string {
	var asciiRepresentations []string

	for _, b := range value {
		if b < 32 || b > 126 {
			// If the byte is not printable, replace it with a dot
			asciiRepresentations = append(asciiRepresentations, ".")
		} else {
			asciiRepresentations = append(asciiRepresentations, string(b))
		}
	}

	return strings.Join(asciiRepresentations, "")
}

func formatBytesAsHex(value []byte) string {
	var hexadecimalRepresentations []string

	for _, b := range value {
		hexadecimalRepresentations = append(hexadecimalRepresentations, fmt.Sprintf("%02x", b))
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
