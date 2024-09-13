package functions

import (
	"fmt"
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

// Add spaces around parentheses if they are attached to words
func addSpacesAroundParentheses(input string) string {
	var result strings.Builder
	inParenthesis := false

	for i, r := range input {
		// Check if we are encountering an opening parenthesis
		if r == '(' {
			// Add a space before if it's attached to a word
			if i > 0 && !strings.ContainsRune(" (", rune(input[i-1])) {
				result.WriteRune(' ')
			}
			inParenthesis = true
		} else if r == ')' {
			// Add a space after if it's attached to a word
			if inParenthesis && (i+1 < len(input) && !strings.ContainsRune(" )", rune(input[i+1]))) {
				result.WriteRune(' ')
			}
			inParenthesis = false
		}
		result.WriteRune(r)
	}

	// If we end inside parentheses, add a space at the end
	if inParenthesis {
		result.WriteRune(' ')
	}

	return result.String()
}

func processBinaryHex(datafile []string, i int) []string {
	previousWord := strings.TrimSpace(datafile[i-1])
	marker := strings.TrimSpace(datafile[i])

	if strings.Contains(marker, "bin") && i-1 >= 0 {
		if IsValidBinary(previousWord) {
			if binNumber, err := strconv.ParseInt(previousWord, 2, 64); err == nil {
				datafile[i-1] = strconv.FormatInt(binNumber, 10)
			}
		}
		datafile[i] = "" // Remove the (bin) marker
	} else if strings.Contains(marker, "hex") && i-1 >= 0 {
		if IsValidHex(previousWord) {
			if hexNumber, err := strconv.ParseInt(previousWord, 16, 64); err == nil {
				datafile[i-1] = strconv.FormatInt(hexNumber, 10)
			}
		}
		datafile[i] = "" // Remove the (hex) marker
	}

	return datafile
}


// Apply transformations with arguments like (up, 2), (low, 3), etc.
func processTransformationsWithArgs(datafile []string, transformationType string, n int, i int) []string {
	// Ensure there are enough previous words to apply the transformation
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
	}
	// Remove the marker after processing the transformation
	datafile[i] = ""
	return datafile
}
// Process transformations without arguments
func processTransformationsWithoutArgs(datafile []string, marker string, i int) []string {
	transformationType := strings.TrimSpace(marker[1 : len(marker)-1])
	if i > 0 {
		switch transformationType {
		case "up":
			datafile[i-1] = ToUpper(datafile[i-1])
		case "low":
			datafile[i-1] = ToLower(datafile[i-1])
		case "cap":
			datafile[i-1] = Capitalize(datafile[i-1])
		}
	}
	datafile[i] = "" // Remove the marker after applying the transformation
	return datafile
}





// Apply transformations based on markers
func ApplyParenthesesLogic(res string) string {
	res = addSpacesAroundParentheses(res)
	var datafile []string
	word := ""

	// Split the string into individual words and markers
	for _, r := range res {
		if r == '(' {
			if len(word) > 0 {
				datafile = append(datafile, word)
				word = ""
			}
			word += string(r)
		} else if r == ')' {
			word += string(r)
			datafile = append(datafile, word)
			word = ""
		} else if r == ' ' {
			if len(word) > 0 {
				datafile = append(datafile, word)
				word = ""
			}
		} else {
			word += string(r)
		}
	}

	if len(word) > 0 {
		datafile = append(datafile, word)
	}

	// Handle transformations sequentially
	for i := 0; i < len(datafile); i++ {
		word := strings.TrimSpace(datafile[i])
		fmt.Println(word)
		if word != "" {
			if word == "(bin)" || word == "(hex)" {
				datafile = processBinaryHex(datafile, i)
			} else if strings.HasPrefix(word, "(") && strings.HasSuffix(word, ")") {
				marker := word

				// Ensure the marker has a valid format
				if strings.Contains(marker, ",") {
					parts := strings.Split(marker[1:len(marker)-1], ",")
					if len(parts) == 2 {
						transformationType := strings.TrimSpace(parts[0])
						argument := strings.TrimSpace(parts[1])

						// Check if argument is a valid number
						if n, err := strconv.Atoi(argument); err == nil && n >= 0 {
							
							// Process transformations with valid arguments
							datafile = processTransformationsWithArgs(datafile, transformationType, n, i)
						} else {
							datafile[i] = "invalid format"
						}
					} else {
						datafile[i] = "invalid format"
					}
				} else {
					datafile = processTransformationsWithoutArgs(datafile, marker, i)
				}
			}
		}
	}

	return strings.TrimSpace(strings.Join(datafile, " "))
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