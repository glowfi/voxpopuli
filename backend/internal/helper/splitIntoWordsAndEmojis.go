package helper

import (
	"strings"
)

func SplitIntoWordsAndEmojis(input string) []string {
	var result []string
	var currentSegment strings.Builder
	inEmoji := false

	for i := 0; i < len(input); i++ {
		if input[i] == ':' {
			if inEmoji { // End of an emoji
				currentSegment.WriteRune(':') // Include the closing ':'
				result = append(result, currentSegment.String())
				currentSegment.Reset()
				inEmoji = false
			} else { // Start of an emoji
				if currentSegment.Len() > 0 { // Add the accumulated word before the emoji
					result = append(result, currentSegment.String())
					currentSegment.Reset()
				}
				currentSegment.WriteRune(':') // Include the opening ':'
				inEmoji = true
			}
		} else {
			if inEmoji {
				currentSegment.WriteRune(rune(input[i])) // Continue building the emoji
			} else {
				currentSegment.WriteRune(rune(input[i])) // Continue building the word
			}
		}
	}

	// Add any remaining segment after the loop
	if currentSegment.Len() > 0 {
		result = append(result, currentSegment.String())
	}

	return result
}
