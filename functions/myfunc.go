package functions

import (
	"strconv"
	"strings"
)

// ProcessSingleQuotes handles the single quote logic
func ProcessSingleQuotes(line string) string {
	res := ""
	insideSingleQuote := false
	singleQuoteText := ""

	for i := 0; i < len(line); i++ {
		ch := line[i]
		if ch == '\'' {
			if insideSingleQuote {
				res += strings.TrimSpace(singleQuoteText) + "'"
				insideSingleQuote = false
				singleQuoteText = ""
			} else {
				insideSingleQuote = true
				singleQuoteText = ""
				res += "'"
			}
			continue
		}
		if insideSingleQuote {
			singleQuoteText += string(ch)
		} else {
			res += string(ch)
		}
	}

	if insideSingleQuote {
		res += singleQuoteText
	}

	return res
}

// Helper function to split the word and transformation marker
func splitWordAndMarker(word string) (string, string) {
	if strings.Contains(word, "(") && strings.Contains(word, ")") {
		openParenIndex := strings.Index(word, "(")
		closeParenIndex := strings.Index(word, ")")
		if openParenIndex < closeParenIndex {
			// Split into base word and marker
			baseWord := word[:openParenIndex]
			marker := word[openParenIndex:closeParenIndex+1]
			return strings.TrimSpace(baseWord), marker
		}
	}
	return word, ""
}

// ApplyParenthesesLogic processes transformations (up, low, cap, bin, hex)
func ApplyParenthesesLogic(res string) string {
	datafile := strings.Split(res, " ")

	// Process transformations with arguments
	for i := 0; i < len(datafile); i++ {
		word := datafile[i]

		// Separate word from any transformation marker
		baseWord, marker := splitWordAndMarker(word)

		if marker != "" {
			// Handle (bin) and (hex) markers first
			if strings.HasPrefix(marker, "(bin)") && i-1 >= 0 {
				if IsValidBinary(datafile[i-1]) {
					if binNumber, err := strconv.ParseInt(datafile[i-1], 2, 64); err == nil {
						datafile[i-1] = strconv.FormatInt(binNumber, 10)
					} else {
						datafile[i] = baseWord + " invalid bin format"
					}
				} else {
					datafile[i] = baseWord + " invalid bin format"
				}
				datafile[i] = "" // Remove the (bin) marker
				continue // Skip further processing for this marker

			} else if strings.HasPrefix(marker, "(hex)") && i-1 >= 0 {
				if IsValidHex(datafile[i-1]) {
					if hexNumber, err := strconv.ParseInt(datafile[i-1], 16, 64); err == nil {
						datafile[i-1] = strconv.FormatInt(hexNumber, 10)
					} else {
						datafile[i] = baseWord + " invalid hex format"
					}
				} else {
					datafile[i] = baseWord + " invalid hex format"
				}
				datafile[i] = "" // Remove the (hex) marker
				continue // Skip further processing for this marker
			}
		}
	}

	// Process transformations without arguments
	for i := 0; i < len(datafile); i++ {
		word := datafile[i]

		// Separate word from any transformation marker
		baseWord, marker := splitWordAndMarker(word)

		if marker != "" {
			if strings.HasPrefix(marker, "(") && strings.HasSuffix(marker, ")") {
				if strings.Contains(marker, ",") {
					parts := strings.Split(marker[1:len(marker)-1], ",")
					if len(parts) == 2 {
						transformationType := strings.TrimSpace(parts[0])
						argument := strings.TrimSpace(parts[1])

						// Convert argument to integer
						n, err := strconv.Atoi(argument)
						if err != nil || n < 0 {
							datafile[i] = baseWord + " invalid format"
							continue
						}

						// Apply transformation based on type
						if transformationType == "up" || transformationType == "low" || transformationType == "cap" {
							if i-n >= 0 {
								for j := i - n; j < i; j++ {
									switch transformationType {
									case "up":
										datafile[j] = ToUpper(datafile[j])
									case "low":
										datafile[j] = ToLower(datafile[j])
									case "cap":
										datafile[j] = Capitalize(datafile[j])
									}
								}
							} else {
								datafile[i] = baseWord + " invalid format"
							}
						}
						datafile[i] = "" // Remove the marker
						continue
					} else {
						datafile[i] = baseWord + marker
						continue
					}
				} else {
					// Handle cases with no argument, e.g., (cap), (low), (up)
					transformationType := strings.TrimSpace(marker[1 : len(marker)-1])
					if i > 0 {
						// Apply transformation to the single preceding word
						switch transformationType {
						case "up":
							datafile[i-1] = ToUpper(datafile[i-1])
						case "low":
							datafile[i-1] = ToLower(datafile[i-1])
						case "cap":
							datafile[i-1] = Capitalize(datafile[i-1])
						}
					}
					datafile[i] = "" // Remove the marker
				}
			} else {
				// If marker is not recognized, treat it as text
				datafile[i] = baseWord + marker
			}
		} else {
			datafile[i] = baseWord
		}
	}

	// Join the processed words into a single string and return
	return strings.Join(datafile, " ")
}

func ToUpper(s string) string {
	var res []rune
	for _, i := range s {
		if i >= 'a' && i <= 'z' {
			res = append(res, i-32)
		} else {
			res = append(res, i)
		}
	}
	return string(res)
}
func Capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(string(s[0])) + strings.ToLower(s[1:])
}

func ToLower(s string) string {
	var res []rune
	for _, i := range s {
		if i >= 'A' && i <= 'Z' {
			res = append(res, i+32)
		} else {
			res = append(res, i)
		}
	}
	return string(res)
}

func IsPunctuation(ch byte) bool {
	return ch == '.' || ch == ',' || ch == '!' || ch == '?' || ch == ':' || ch == ';'
}

func IsValidBinary(s string) bool {
    for _, c := range s {
        if c != '0' && c != '1' {
            return false
        }
    }
    return len(s) > 0 // Also ensure it's not an empty string
}

func IsValidHex(s string) bool {
    _, err := strconv.ParseInt(s, 16, 64)
    return err == nil
}



func IsAlphanumeric(ch byte) bool {
	// Check if the character is a letter (A-Z or a-z) or a digit (0-9)
	return (ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z') || (ch >= '0' && ch <= '9')
}


// ReplaceAWithAn replaces 'a' with 'an' if the next word begins with a vowel or 'h'.
// It also converts 'an' back to 'a' if the next word does not start with a vowel or 'h'.
func ReplaceAWithAn(text string) string {
	words := strings.Fields(text) // Split the text into words
	for i := 0; i < len(words)-1; i++ {
		// Check if the current word is 'a'
		if words[i] == "a" {
			// Check if the next word starts with a vowel or 'h'
			if startsWithVowelOrH(words[i+1]) {
				// Replace 'a' with 'an'
				words[i] = "an"
			}
		} else if words[i] == "an" {
			// Check if the next word does NOT start with a vowel or 'h'
			if !startsWithVowelOrH(words[i+1]) {
				// Revert 'an' back to 'a'
				words[i] = "a"
			}
		}else if words[i] == "AN" {
			// Check if the next word does NOT start with a vowel or 'h'
			if !startsWithVowelOrH(words[i+1]) {
				// Revert 'an' back to 'a'
				words[i] = "A"
			}
		}else if words[i] == "A" {
			// Check if the next word does NOT start with a vowel or 'h'
			if !startsWithVowelOrH(words[i+1]) {
				// Revert 'an' back to 'a'
				words[i] = "AN"
			}
		}
	}
	return strings.Join(words, " ")
}

// Helper function to check if a word starts with a vowel or 'h'
func startsWithVowelOrH(word string) bool {
	if len(word) == 0 {
		return false
	}
	firstChar := strings.ToLower(string(word[0])) // Convert the first character to lowercase
	return firstChar == "a" || firstChar == "e" || firstChar == "i" || firstChar == "o" || firstChar == "u" || firstChar == "h"
}