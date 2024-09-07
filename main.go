package main

import (
	"fmt"
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
			parenth:=false
			temp := "" 
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
						if ch == '.' || ch == ',' || ch == '!' || ch == '?' || ch == ':' || ch == ';' {
							// Trim any trailing space in res before punctuation
							if len(res) > 0 && res[len(res)-1] == ' ' {
								res = res[:len(res)-1]
							}
							res += string(ch) // Add punctuation
							if i < len(data)-1 && data[i+1] != ' ' { 
								// Ensure there's a space after punctuation if it's not the last character
								res += " "
							}
						} else if ch == ' ' {
							// Only add a space if the previous character is not a space or punctuation
							if len(res) > 0 && res[len(res)-1] != ' ' && !strings.Contains(".,!?;:", string(res[len(res)-1])) {
								res += " "
							}
						} else {
							res += string(ch)
						}
					}
					// if parenth {
						
					// 	if ch != ' ' {
					// 		temp += string(ch)
					// 	}
					// } else {
				
					// 	if ch > 32 && ch < 48 {
					// 		res += " " + string(ch) + " "
					// 	} else {
					// 		res += string(ch)
					// 	}
					// }
				}
				// if ch == '(' {
				// 	parenth = true
				// 	res += " " + string(ch) 
				// } else if ch == ')' {
				// 	parenth = false
				// 	res += string(ch) + " " 
				// } else {
				// 	if parenth {
				// 		if ch == ',' {
				// 			res += string(ch)
				// 		} else if ch != ' ' {
				// 			res += string(ch)
				// 		}
				// 	} else {
				// 		if ch > 32 && ch < 48 {
				// 			res += " " + string(ch) + " "
				// 		} else {
				// 			res += string(ch)
				// 		}
				// 	}
				// }
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
							break // Ensure index doesn't go negative
						}
			
						if strings.HasPrefix(word, "(up") {
							datafile[i-j-1] = ToUpper(datafile[i-j-1])
						} else if strings.HasPrefix(word, "(low") {
							datafile[i-j-1] = ToLower(datafile[i-j-1])
						} else if strings.HasPrefix(word, "(cap") {
							datafile[i-j-1] = ToUpper(string(datafile[i-j-1][0])) + ToLower(datafile[i-j-1][1:])
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
		
				// if word == "(up)" || word == "(low)" || word == "(cap)" {

				// 	if i+2 < len(datafile) && datafile[i+1] == "," {
				// 		number, err := strconv.Atoi(datafile[i+2])
				// 		if err != nil {
				// 			fmt.Println("Invalid input for number:", err)
				// 			return
				// 		}
				// 		fmt.Print(number)
				// 		for j := i - 2; number != 0; j-- {
							
				// 			if word == "up" {
				// 				datafile[i-number] = ToUpper(datafile[i-number])
				// 				datafile[i] = ""
				// 				datafile[i+1] = ""
				// 				datafile[i+2] = ""
				// 				datafile[i+3] = ""

				// 				datafile[i-1] = ""
				// 			}
				// 			if word == "low" {
				// 				datafile[i-2] = ToLower(datafile[i-2])
				// 				datafile[i] = ""
				// 				datafile[i+1] = ""
				// 				datafile[i-1] = ""

				// 			}
				// 			if word == "cap" {
				// 				c := datafile[i-2]
				// 				datafile[i-2] = ToUpper(string(c[0])) + ToLower(string(c[1:]))
				// 				datafile[i] = ""
				// 				datafile[i+1] = ""
				// 				datafile[i-1] = ""
				// 			}
						
				// 		}
				// 	}
				// }
				// if word == "(bin,1)" {
				// 	number, err := strconv.ParseInt(datafile[i-2], 2, 64)
				// 	if err != nil {
				// 		fmt.Println("Invalid input bin:", err)
				// 		return
				// 	}
				// 	datafile[i-2] = strconv.FormatInt(number, 10)
				// 	datafile[i] = ""
				// 	datafile[i+1] = ""
				// 	datafile[i-1] = ""
				// }
				// if word == "(hex,1)" {
				// 	number, err := strconv.ParseInt(datafile[i-2], 16, 64)
				// 	if err != nil {
				// 		fmt.Println("Invalid input hex:", err)
				// 		return
				// 	}
				// 	datafile[i-2] = strconv.FormatInt(number, 10)
				// 	datafile[i] = ""
				// 	datafile[i+1] = ""
				// 	datafile[i-1] = ""
				// }

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
