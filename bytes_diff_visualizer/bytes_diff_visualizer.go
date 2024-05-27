package bytes_diff_visualizer

import (
	"bytes"
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/fatih/color"
)

// VisualizeByteDiff visualizes the difference between two byte slices, returning lines to be presented to the user.
//
// The lines will include ANSI escape codes to colorize the output.
func VisualizeByteDiff(actual []byte, expected []byte) []string {
	byteDiffLines, row, col := visualizeByteDiff(actual, expected)

	// Always odd, so get the size of each half exactly, +1 for the sepearator line
	rowsInEachHalf := len(byteDiffLines)/2 + 1
	// The line that contains the first differing byte, in the "expected" half
	expectedLine := byteDiffLines[row]
	// The line that contains the first differing byte, in the "actual" half
	actualLine := byteDiffLines[row+rowsInEachHalf]
	// The start and end index of the hex representation of the differing byte
	// col is the index of the differing byte, so we have (col-1+1*2) bytes of hex representation before it, and (col*2) bytes of space
	startIdxHex, endIdxHex := col*3, col*3+2
	// We use the "|" separator to find the index of the differing byte in the ASCII representation
	// We add 2 to the index to skip the space and the separator
	idxAscii := strings.Index(expectedLine, "|") + 2 + col

	byteDiffLines[row] = expectedLine[:startIdxHex] + colorize(color.FgHiGreen, expectedLine[startIdxHex:endIdxHex]) + expectedLine[endIdxHex:idxAscii] + colorize(color.FgHiGreen, expectedLine[idxAscii:idxAscii+1]) + expectedLine[idxAscii+1:]

	byteDiffLines[row+rowsInEachHalf] = actualLine[:startIdxHex] + colorize(color.FgHiRed, actualLine[startIdxHex:endIdxHex]) + actualLine[endIdxHex:idxAscii] + colorize(color.FgHiRed, actualLine[idxAscii:idxAscii+1]) + actualLine[idxAscii+1:]

	return byteDiffLines
}

// visualizeByteDiff visualizes the difference between two byte slices, returning lines to be presented to the user.
// We also want to highlight the first differing byte in the output.
// For that reason we return the row and column of the first differing byte.
func visualizeByteDiff(actual []byte, expected []byte) ([]string, int, int) {
	// If both are exactly the same, return an empty slice
	if bytes.Equal(actual, expected) {
		return []string{}, -1, -1
	}

	rows := 1 // "Expected ... line"
	var firstDiffIndexRow, firstDiffIndexCol int
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

	fmt.Fprintf(tabWriter, "Expected (bytes %v-%v), hexadecimal:       \t ASCII:\n", byteRangeStart, byteRangeEnd)

	for i := byteRangeStart; i < intmin(byteRangeEnd, len(expected)); i += byteCountPerLine {
		end := intmin(i+byteCountPerLine, len(expected))

		if firstDiffIndex >= i && firstDiffIndex < end {
			firstDiffIndexCol = firstDiffIndex - i
			firstDiffIndexRow = rows
		} else {
			rows += 1
		}

		fmt.Fprintf(tabWriter, "%v\t %v\n", formatBytesAsHex(expected[i:end]), formatBytesAsAscii(expected[i:end]))
	}

	tabWriter.Write([]byte("\n"))
	fmt.Fprintf(tabWriter, "Actual (bytes %v-%v), hexadecimal:         \t ASCII:\n", byteRangeStart, byteRangeEnd)

	for i := byteRangeStart; i < intmin(byteRangeEnd, len(actual)); i += byteCountPerLine {
		end := intmin(i+byteCountPerLine, len(actual))
		fmt.Fprintf(tabWriter, "%v\t %v\n", formatBytesAsHex(actual[i:end]), formatBytesAsAscii(actual[i:end]))
	}

	tabWriter.Flush()

	output := linesBuffer.String()
	if output[len(output)-1] == '\n' {
		output = output[:len(output)-1]
	}
	return strings.Split(output, "\n"), firstDiffIndexRow, firstDiffIndexCol
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

func colorize(colorToUse color.Attribute, msg string) string {
	colorizedLine := color.New(colorToUse).SprintFunc()(msg)

	return colorizedLine
}
