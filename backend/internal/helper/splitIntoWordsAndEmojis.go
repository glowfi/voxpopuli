package helper

import (
	"fmt"
	"strings"

	"github.com/forPelevin/gomoji"
	"github.com/rivo/uniseg"
)

func IsCustomEmoji(input string) bool {
	if len(input) < 3 {
		return false
	}
	return input[0] == ':' && input[len(input)-1] == ':'
}

func SplitStringIntoCustomEmojisAndWords(input string) []string {
	var result []string
	var currentSegment strings.Builder
	inEmoji := false

	for i := 0; i < len(input); i++ {
		if input[i] == ':' {
			if inEmoji { // End of an emoji
				currentSegment.WriteByte(':') // Include the closing ':'
				result = append(result, currentSegment.String())
				currentSegment.Reset()
				inEmoji = false
			} else { // Start of an emoji
				if currentSegment.Len() > 0 { // Add the accumulated word before the emoji
					result = append(result, currentSegment.String())
					currentSegment.Reset()
				}
				currentSegment.WriteByte(':') // Include the opening ':'
				inEmoji = true
			}
		} else {
			currentSegment.WriteByte(input[i]) // Continue building the word or emoji
		}
	}

	// Add any remaining segment after the loop
	if currentSegment.Len() > 0 {
		result = append(result, currentSegment.String())
	}

	return result
}

func GetBestGuessedEmojiInfo(emojiStr string) (gomoji.Emoji, error) {
	emojis := gomoji.FindAll(emojiStr)

	if len(emojis) == 0 {
		return gomoji.Emoji{}, fmt.Errorf("no emojis present")
	}

	if len(emojis) > 1 {
		return gomoji.Emoji{}, fmt.Errorf("multiple emojis not allowed")
	}

	return emojis[0], nil
}

// isStandardEmoji checks if a rune is an emoji
// func isStandardEmoji(r rune) bool {
// 	return gomoji.ContainsEmoji(string(r))
// }

// splitByEmoji splits the input string into a slice of strings based on emojis
func SplitStringIntoStandardEmojisAndWords(input string) []string {
	var output []string

	gs := uniseg.NewGraphemes(input)

	var tmp strings.Builder

	for gs.Next() {
		grapheme := gs.Str()
		if gomoji.ContainsEmoji(grapheme) {
			if tmp.Len() > 0 {
				output = append(output, tmp.String())
				tmp.Reset()
			}
			output = append(output, grapheme)
		} else {
			tmp.Write([]byte(grapheme))
		}

	}

	if tmp.Len() > 0 {
		output = append(output, tmp.String())
	}

	return output
}
