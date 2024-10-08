package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"go-reloaded/functions"
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
			data, err := os.ReadFile(input)
			if err != nil {
				fmt.Println("Error reading the file:", err)
				return
			}
			finalResult := ""
			
			rawData := string(data)
			// Split data into lines
			lines := strings.Split(rawData, "\n")
			
			
			for i, line := range lines {
				res := functions.ApplyParenthesesLogic(line)
				res= functions.ReplaceAWithAn(res)
				finalLine := functions.Ponctuation(res)
				 finalLine = functions.ProcessSingleQuotes(finalLine)
				finalResult += finalLine
				if i != len(lines) -1 {
					finalResult+="\n"
				}
			}

			file, err := os.Create(output)
			if err != nil {
				fmt.Println("Error creating the file:", err)
				return
			}
			defer file.Close()
			

			// Use io.WriteString to write the processed result to the file
			_, err = io.WriteString(file, finalResult) // Pass the file as io.Writer and finalResult as the string
			if err != nil {
				fmt.Println("Error writing to file:", err)
				return
			}
			
			fmt.Println("File processed successfully 'NADII' !!!.")
		}

	} else if len(args) > 2 {
		fmt.Println("Too much arguments")
		os.Exit(1)
	} else if len(args) == 1 {
		fmt.Println("Less arguments please enter the input and the output files name")
		os.Exit(1)
	}
}
