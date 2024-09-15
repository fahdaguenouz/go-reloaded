package functions

import (

	"strconv"
	"strings"
)

// ProcessSingleQuotes handles the single quote logic
func ProcessSingleQuotes(text string) string {
	result := ""
	first := false
	next := false
	for i, char := range text {
		if char == '\'' {
			// if single cote
			if i == 0 {
				// if its first char
				result += string(char)
				first = true
				// Check if the next char is a space.
				if text[i+1] == ' ' {
					next = true
				}
			} else if i == len(text)-1 {
				// If the single quote is the last character.
				// Remove the last character from result if it is a space.
				if text[i-1] == ' ' {
					result = result[:len(result)-1]
				}
				result += "'"
			} else {
				// your logical code
				if IsAlphanumeric(text[i-1]) && IsAlphanumeric(text[i+1]) {
					result += "'"
					continue
				}
				// If the single quote is not the first one.
				if !first {
					// first cote
					if result[len(result)-1] != ' ' {
						result += " "
					}
					result += "'"
					if text[i+1] == ' ' {
						next = true
					}
					first = true
				} else {
					// last cote
					first = false
					if result[len(result)-1] == ' ' {
						result = result[:len(result)-1]
					}
					ispunc := strings.Contains(",;:.!? ", string(text[i+1]))
					result += "'"
					if !ispunc {
						result += " "
					}
				}
			}
		} else {
			// if not single cote
			if next {
				next = false
				continue
			}
			result += string(char)
		}
	}
	return result
}
// Add spaces around parentheses if they are attached to words
func addSpacesAroundParentheses(input string) string {
	result := ""
	inParenthesis := false

	for i := 0; i < len(input); i++ {
		r := rune(input[i])
		// Check if we are encountering an opening parenthesis
		if r == '(' {
			
			// Add a space before if it's attached to a word
			if i > 0 && !strings.ContainsRune(" (", rune(input[i-1])) {
				result += " "
			}
			inParenthesis = true
		} else if r == ')' {
			// Add a space after if it's attached to a word
			result += string(r)
			if inParenthesis && (i+1 < len(input) && !strings.ContainsRune(" )", rune(input[i+1]))) {
				result += " "
				
				}
				inParenthesis = false
				continue
				}
				
				// Add the current character to the result
				result += string(r)
				
				
	}

	// If we end inside parentheses, add a space at the end
	if inParenthesis {
		result += " "
	}

	return result
}

func processBinaryHex(datafile []string, i int) []string {
	
	previousWord := strings.TrimSpace(datafile[i-1])
	marker := strings.TrimSpace(datafile[i])
	

	if strings.Contains(marker, "bin") && i-1 >= 0 {
		if IsValidBinary(previousWord) {
			if binNumber, err := strconv.ParseInt(previousWord, 2, 64); err == nil {// verifier transform to base 2
				datafile[i-1] = strconv.FormatInt(binNumber, 10) // trnasform to decimal base 10
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

func isNumber(word string) bool {
	for i := 0; i < len(word); i++ {
		if !isDigit(word[i]) {
			return false
		}
	}
	return len(word) > 0
}
// Apply transformations with arguments like (up, 2), (low, 3), etc.
func processTransformationsWithArgs(datafile []string, transformationType string, n int, i int) []string {
	// Ensure there are enough previous words to apply the transformation
	if n > 0 { // Only apply if n is positive
		c := 0
		for j := i - 1; j >= 0 && c < n; j-- {
			if datafile[j] != ""  && !IsPunctuation(datafile[j][len(datafile[j])-1]) && !isNumber(datafile[j]){
				switch transformationType {
				case "up":
					datafile[j] = ToUpper(datafile[j])
				case "low":
					datafile[j] = ToLower(datafile[j])
				case "cap":
					datafile[j] = Capitalize(datafile[j])
				}
				c++
			}
		}
		// Remove the marker after processing the transformation
		datafile[i] = ""
		datafile[i+1] = ""
	}
	// If n is invalid, do not remove the marker
	return datafile
}

// Process transformations without arguments
func processTransformationsWithoutArgs(datafile []string, marker string, i int) []string {
	transformationType := strings.TrimSpace(marker[1 : len(marker)-1])
	// Traverse backwards to find the first non-punctuation word
	if i > 0 {
		j := i - 1
		for j >= 0 && IsPunctuation(datafile[j][len(datafile[j])-1]) || isNumber(datafile[j]) {
			j-- // Skip punctuation
		}
		
		// Apply the transformation to the first non-punctuation word
		if j >= 0 && datafile[j] != "" {
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
	
	datafile[i] = ""
	// Remove the marker after applying the transformation
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
		if word != "" {
			if word == "(bin)" || word == "(hex)" {
				datafile = processBinaryHex(datafile, i)
			} else if strings.HasPrefix(word, "(up") || strings.HasPrefix(word, "(cap") || strings.HasPrefix(word, "(low") {
				marker := word
		

				// Ensure that the mark  has a valid format
				if strings.Contains(marker, ",") && (i+1 < len(datafile) && strings.HasSuffix(datafile[i+1], ")") && len(datafile[i+1]) == 2) {
					nextPart := strings.TrimSpace(datafile[i+1])

					// Check if the next part contains only a number and a closing parenthesis
					if strings.HasSuffix(nextPart, ")") {
						argument := nextPart[:len(nextPart)-1] // Remove the closing parenthesis
						argument = strings.TrimSpace(argument)

						// Check if argument is a valid number
						if n, err := strconv.Atoi(argument); err == nil && n >= 0 {
							transformationType := strings.TrimSpace(marker[1 : len(marker)-1])
							// Process transformations with valid positive arguments
							if n >= 0 {
								datafile = processTransformationsWithArgs(datafile, transformationType, n, i)

								i++ // Skip the next index as we already processed it
							} else {
								// If the number is negative, keep the marker as text
								datafile[i] = marker
								datafile[i+1] = nextPart
							} // Skip the next index as we already processed it
						} else {
							datafile[i] = marker
							datafile[i+1] = nextPart
						}
					} else {
						datafile[i] = marker
					}
				} else if word == "(up)" || word == "(cap)" || word == "(low)" {
					datafile = processTransformationsWithoutArgs(datafile, marker, i)
				}
			}
		}
	}
	return strings.TrimSpace(strings.Join(datafile, " "))
}

// Process the input data to handle parentheses, spaces, and punctuation
func Ponctuation(data string) string {

	res := ""
	inPunctuationGroup := false

	for i := 0; i < len(data); i++ {
		ch := data[i]

		if ch == '\n' {
			res += "\n" // Preserve newlines in the input
			continue
		}

		if IsPunctuation(ch) {
			// Handle the case of a dot between two digits 21.2
			if ch == '.' && i > 0 && i < len(data)-1 && isDigit(data[i-1]) && isDigit(data[i+1]) {
				// If the dot is between two digits, just add the dot without spaces
				res += "."
			} else {
				// Trim any trailing space before punctuation
				if !inPunctuationGroup {
					// If not already in a punctuation group, add a space before punctuation
					if len(res) > 0 && res[len(res)-1] == ' ' {
						res = res[:len(res)-1]
					}
					inPunctuationGroup = true
				}
				res += string(ch) // Add punctuation
			} // Add punctuation
		} else {
			if inPunctuationGroup {
				// After a punctuation group, add a space if the next character is not punctuation or space
				if len(res) > 0 && !IsPunctuation(ch) && ch != ' ' {
					res += " "
				}
				inPunctuationGroup = false
			}

			// Handle spaces and normal characters
			if ch == ' ' {
				res += " "

			} else {
				res += string(ch) // Add normal characters
			}
		}
	}

	return res
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
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
	words := strings.Fields(text)
	 // Split the text into words
	for i := 0; i < len(words)-1; i++ {
		// Check if the current word is 'a'
		

		// Check if the current word is 'a' or 'A'
		if words[i] == "a" || words[i] == "A" || words[i]=="'a" ||words[i]=="'A" {
			// Check if the next word starts with a vowel or 'h'
			if startsWithVowelOrH(words[i+1]) && words[i]!="'a" &&words[i]!="'A"{
				// Replace 'a' with 'an' 
				if words[i] == "A" {
					words[i] = "An"
				} else {
					words[i] = "an"
				}
			}else if words[i+1] == "a'" || words[i+1] == "A'" {
				// 'a a' or 'A A'
				if words[i] == "'A" {
					words[i] = "'An"
				} else {
					words[i] = "'an"
				}
			}
		} else if words[i] == "an" || words[i] == "An" || words[i] == "AN" {
			// Check if the next word does NOT start with a vowel or 'h'
			if !startsWithVowelOrH(words[i+1]) {
				// Revert 'an' back to 'a' 
				if words[i] == "AN" {
					words[i] = "A"
				} else if words[i] == "An" {
					words[i] = "A"
				} else {
					words[i] = "a"
				}
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
