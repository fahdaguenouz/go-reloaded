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

// ApplyParenthesesLogic processes (up, low, cap, bin, hex) transformations
// but leaves the text as is if it's part of a regular sentence.
func ApplyParenthesesLogic(res string) string {
	datafile := strings.Fields(res)
	for i := 0; i < len(datafile); i++ {
		word := datafile[i]

		// Process transformations (up, low, cap) only if they have a closing parenthesis
		if (strings.HasPrefix(word, "(up") || strings.HasPrefix(word, "(low") || strings.HasPrefix(word, "(cap")) &&
			strings.Contains(word, ")") {

			// Extract the number part
			numStart := strings.Index(word, ",")
			numEnd := strings.Index(word, ")")
			if numStart != -1 && numEnd != -1 {
				number, err := strconv.Atoi(word[numStart+1 : numEnd])
				if err != nil {
					continue // Treat as invalid instruction, ignore this word
				}

				// Process the words before the transformation instruction
				count := 0
				for j := i - 1; j >= 0 && count < number; j-- {
					if datafile[j] != "" {
						count++
						if strings.HasPrefix(word, "(up") {
							datafile[j] = ToUpper(datafile[j])
						} else if strings.HasPrefix(word, "(low") {
							datafile[j] = ToLower(datafile[j])
						} else if strings.HasPrefix(word, "(cap") {
							// Capitalize first letter, lower the rest
							datafile[j] = ToUpper(string(datafile[j][0])) + ToLower(datafile[j][1:])
						}
					}
				}
				// Remove the transformation marker after applying it
				datafile[i] = ""
			}

		} else if strings.HasPrefix(word, "(bin") && strings.Contains(word, ")") && i-1 >= 0 && IsValidBinary(datafile[i-1]) {
			binNumber, err := strconv.ParseInt(datafile[i-1], 2, 64)
			if err == nil {
				datafile[i-1] = strconv.FormatInt(binNumber, 10)
				datafile[i] = ""
			}
		} else if strings.HasPrefix(word, "(hex") && strings.Contains(word, ")") && i-1 >= 0 && IsValidHex(datafile[i-1]) {
			hexNumber, err := strconv.ParseInt(datafile[i-1], 16, 64)
			if err == nil {
				datafile[i-1] = strconv.FormatInt(hexNumber, 10)
				datafile[i] = ""
			}
		}
	}

	// Return the processed result as a single string
	return strings.Join(datafile, " ")
}

func ToUpper(s string) string {
	var res []rune
	for _, i := range s {
		if i >= 'a' && i <= 'z' {
			res = append(res, i-32)
		}else if i=='é' {
			res = append(res, 'É')
		} else {
			res = append(res, i)
		}
	}
	return string(res)
}

func ToLower(s string) string {
	var res []rune
	for _, i := range s {
		if i >= 'A' && i <= 'Z' {
			res = append(res, i+32)
		}else if i=='é' {
			res = append(res, 'É')
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