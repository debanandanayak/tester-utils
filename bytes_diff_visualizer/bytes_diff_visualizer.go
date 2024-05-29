package bytes_diff_visualizer

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/fatih/color"
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

	leftHeader := PadRight(fmt.Sprintf("Expected (bytes %v-%v), hexadecimal:", byteRangeStart, byteRangeEnd), " ", 60)
	fmt.Fprintf(linesBuffer, "%s| ASCII:\n", leftHeader)

	for i := byteRangeStart; i < intmin(byteRangeEnd, len(expected)); i += byteCountPerLine {
		end := intmin(i+byteCountPerLine, len(expected))

		var bytesAsAscii, bytesAsHex string
		if firstDiffIndex >= i && firstDiffIndex < end {
			bytesAsHex = formatBytesAsHex(expected[i:firstDiffIndex]) + " " + colorizeString(color.FgHiGreen, formatBytesAsHex(expected[firstDiffIndex:firstDiffIndex+1])) + " " + formatBytesAsHex(expected[firstDiffIndex+1:end])
			bytesAsHex = PadRight(bytesAsHex, " ", 69)
			bytesAsAscii = formatBytesAsAscii(expected[i:firstDiffIndex]) + colorizeString(color.FgHiGreen, formatBytesAsAscii(expected[firstDiffIndex:firstDiffIndex+1])) + formatBytesAsAscii(expected[firstDiffIndex+1:end])
		} else {
			bytesAsHex = PadRight(formatBytesAsHex(expected[i:end]), " ", 60)
			bytesAsAscii = formatBytesAsAscii(expected[i:end])
		}

		fmt.Fprintf(linesBuffer, "%v| %v\n", bytesAsHex, bytesAsAscii)
	}

	linesBuffer.Write([]byte("\n"))

	leftHeader = PadRight(fmt.Sprintf("Actual (bytes %v-%v), hexadecimal:", byteRangeStart, byteRangeEnd), " ", 60)
	fmt.Fprintf(linesBuffer, "%s| ASCII:\n", leftHeader)

	for i := byteRangeStart; i < intmin(byteRangeEnd, len(actual)); i += byteCountPerLine {
		end := intmin(i+byteCountPerLine, len(actual))

		var bytesAsAscii, bytesAsHex string
		if firstDiffIndex >= i && firstDiffIndex < end {
			bytesAsHex = formatBytesAsHex(actual[i:firstDiffIndex]) + " " + colorizeString(color.FgHiRed, formatBytesAsHex(actual[firstDiffIndex:firstDiffIndex+1])) + " " + formatBytesAsHex(actual[firstDiffIndex+1:end])
			bytesAsHex = PadRight(bytesAsHex, " ", 69)
			bytesAsAscii = formatBytesAsAscii(actual[i:firstDiffIndex]) + colorizeString(color.FgHiRed, formatBytesAsAscii(actual[firstDiffIndex:firstDiffIndex+1])) + formatBytesAsAscii(actual[firstDiffIndex+1:end])
		} else {
			bytesAsHex = PadRight(formatBytesAsHex(actual[i:end]), " ", 60)
			bytesAsAscii = (formatBytesAsAscii(actual[i:end]))
		}

		fmt.Fprintf(linesBuffer, "%v| %v\n", bytesAsHex, bytesAsAscii)
	}

	output := linesBuffer.String()
	if output[len(output)-1] == '\n' {
		output = output[:len(output)-1]
	}
	return strings.Split(output, "\n")
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

func formatHexColorized(value []byte, i int, firstDiffIndex int, end int) string {
	bytesAsHex := formatBytesAsHex(value[i:firstDiffIndex]) + " " + colorizeString(color.FgHiRed, formatBytesAsHex(value[firstDiffIndex:firstDiffIndex+1])) + " " + formatBytesAsHex(value[firstDiffIndex+1:end])
	return bytesAsHex
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

// func colorize(colorToUse color.Attribute, msg string) string {
// 	colorizedLine := color.New(colorToUse).SprintFunc()(msg)

// 	return colorizedLine
// }

func colorizeString(colorToUse color.Attribute, msg string) string {
	c := color.New(colorToUse)
	return c.Sprint(msg)
}

func PadRight(str, pad string, lenght int) string {
	for {
		str += pad
		if len(str) > lenght {
			return str[0:lenght]
		}
	}
}
func PadLeft(str, pad string, lenght int) string {
	for {
		str = pad + str
		if len(str) > lenght {
			return str[0:lenght]
		}
	}
}
