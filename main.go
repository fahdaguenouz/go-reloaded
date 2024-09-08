package main

import (
	"fmt"
	"go-reloaded/functions"
	"io"
	"os"
	"strconv"
	"strings"
)

func main() {
	args := os.Args[1:]
	Inextention := ""
	Ouextention := ""

	if len(args) == 2 {
		input := args[0]
		output := args[1]
		Index := -1
		for i := len(input) - 1; i >= 0; i-- {
			if input[i] == '.' {
				Index = i
				break
			}
		}
		Inextention = input[Index+1:]

		for j := len(output) - 1; j >= 0; j-- {
			if output[j] == '.' {
				Index = j
				break
			}
		}
		Ouextention = output[Index+1:]

		if Inextention != "txt" || Ouextention != "txt" {
			fmt.Println("The files Extention must be : '.txt' ")
			os.Exit(1)
		} else {
			file, err := os.Open(input)
			if err != nil {
				fmt.Println("Error :", err)
				return
			}
			defer file.Close()
			data, err := io.ReadAll(file)
			if err != nil {
				fmt.Println("Error: ", err)
				return
			}
			res := ""
			parenth := false
			temp := ""
			groupPonc := false
			for i, ch := range data {

				if ch == '(' {
					parenth = true
					temp = " ("
				} else if ch == ')' {
					parenth = false
					if !strings.Contains(temp, ",") {
						temp += ",1"
					}
					temp += ")"
					res += temp + " "
					temp = ""
				} else {
					if parenth {
						// Inside parentheses, skip spaces and append characters
						if ch != ' ' {
							temp += string(ch)
						}
					} else {
						// Handling punctuation and spacing rules
						if functions.IsPunctuation(ch) {
							// Ensure only one punctuation group is handled
							if !groupPonc {
								// Trim any trailing space in res before punctuation
								if len(res) > 0 && res[len(res)-1] == ' ' {
									res = res[:len(res)-1]
								}
			
								res += string(ch)
			
								// Continue handling the punctuation group
								for i+1 < len(data) && functions.IsPunctuation(data[i+1]) {
									i++
									res += string(data[i])
								}
			
								// Add space if next character is not punctuation or space
								if i+1 < len(data) && data[i+1] != ' ' && !functions.IsPunctuation(data[i+1]) {
									res += " "
								}
			
								groupPonc = true
							}
						} else {
							groupPonc = false // Reset after handling punctuation
			
							// Handle spaces and normal characters
							if ch == ' ' {
								if len(res) > 0 && res[len(res)-1] != ' ' && !functions.IsPunctuation(res[len(res)-1]) {
									res += " "
								}
							} else {
								res += string(ch)
							}
						}
					}

				}

			}

			datafile := strings.Fields(string(res))
			for _, v := range datafile {
				fmt.Printf("%s\n", v)
			}
			for i := 0; i < len(datafile); i++ {

				word := datafile[i]
				if strings.HasPrefix(word, "(up") || strings.HasPrefix(word, "(low") || strings.HasPrefix(word, "(cap") {
					// Extract the number from the argument, e.g., (up,2) -> number = 2
					numStart := strings.Index(word, ",") + 1
					numEnd := strings.Index(word, ")")
					number, err := strconv.Atoi(word[numStart:numEnd])
					if err != nil {
						fmt.Println("Invalid input for number:", err)
						return
					}

					// Apply changes to the previous 'number' words
					for j := 0; j < number; j++ {
						if i-j-1 < 0 {
							break // To Verify the negativ enumber in the arguments
						}

						if strings.HasPrefix(word, "(up") {
							datafile[i-j-1] = functions.ToUpper(datafile[i-j-1])
						} else if strings.HasPrefix(word, "(low") {
							datafile[i-j-1] = functions.ToLower(datafile[i-j-1])
						} else if strings.HasPrefix(word, "(cap") {
							datafile[i-j-1] = functions.ToUpper(string(datafile[i-j-1][0])) + functions.ToLower(datafile[i-j-1][1:])
						}
					}

					if strings.HasPrefix(word, "(bin") {
						numStart := strings.Index(word, ",") + 1
						numEnd := strings.Index(word, ")")
						number, err := strconv.Atoi(word[numStart:numEnd])
						if err != nil || number != 1 {
							fmt.Println("Invalid input for bin, expected (bin):", err)
							return
						}

						// Convert the previous word from binary to decimal
						binNumber, err := strconv.ParseInt(datafile[i-2], 2, 64)
						if err != nil {
							fmt.Println("Invalid binary input:", err)
							return
						}
						datafile[i-2] = strconv.FormatInt(binNumber, 10)

						// Clear the instruction
						datafile[i] = ""
						datafile[i+1] = ""
						datafile[i-1] = ""
					}

					// Handle hexadecimal conversion
					if strings.HasPrefix(word, "(hex") {
						numStart := strings.Index(word, ",") + 1
						numEnd := strings.Index(word, ")")
						number, err := strconv.Atoi(word[numStart:numEnd])
						if err != nil || number != 1 {
							fmt.Println("Invalid input for hex, expected (hex,1):", err)
							return
						}

						// Convert the previous word from hexadecimal to decimal
						hexNumber, err := strconv.ParseInt(datafile[i-2], 16, 64)
						if err != nil {
							fmt.Println("Invalid hexadecimal input:", err)
							return
						}
						datafile[i-2] = strconv.FormatInt(hexNumber, 10)

						// Clear the instruction
						datafile[i] = ""
						datafile[i+1] = ""
						datafile[i-1] = ""
					}

					// Clear the transformation instruction (up, low, cap) from the list
					datafile[i] = ""
				}
				if strings.HasPrefix(word, "(bin") {
					numStart := strings.Index(word, ",") + 1
					numEnd := strings.Index(word, ")")
					number, err := strconv.Atoi(word[numStart:numEnd])
					if err != nil || number != 1 {
						fmt.Println("Invalid input for bin, expected (bin,1):", err)
						return
					}

					// Convert the previous word from binary to decimal
					binNumber, err := strconv.ParseInt(datafile[i-1], 2, 64)
					if err != nil {
						fmt.Println("Invalid binary input:", err)
						return
					}
					datafile[i-2] = strconv.FormatInt(binNumber, 10)

					// Clear the instruction
					datafile[i] = ""
					datafile[i+1] = ""
					datafile[i-1] = ""
				}

				// Handle hexadecimal conversion
				if strings.HasPrefix(word, "(hex") {
					numStart := strings.Index(word, ",") + 1
					numEnd := strings.Index(word, ")")
					number, err := strconv.Atoi(word[numStart:numEnd])
					if err != nil || number != 1 {
						fmt.Println("Invalid input for hex, expected (hex,1):", err)
						return
					}

					// Convert the previous word from hexadecimal to decimal
					hexNumber, err := strconv.ParseInt(datafile[i-1], 16, 64)
					if err != nil {
						fmt.Println("Invalid hexadecimal input:", err)
						return
					}
					datafile[i-2] = strconv.FormatInt(hexNumber, 10)

					// Clear the instruction
					datafile[i] = ""
					datafile[i+1] = ""
					datafile[i-1] = ""
				}

			}
			s := ""
			for _, ch := range datafile {
				s += ch + " "
			}
			resultat := strings.Fields(s)
			err = os.WriteFile(output, []byte(strings.Join(resultat, " ")), 0o644)
			if err != nil {
				fmt.Println("Error writing to file:", err)
				return
			}
			fmt.Println("File processed successfully.")
		}

	} else if len(args) > 2 {
		fmt.Println("Too much arguments")
		os.Exit(1)
	} else if len(args) == 1 {
		fmt.Println("less arguments please enter the input and the output files name")
		os.Exit(1)
	}
}
