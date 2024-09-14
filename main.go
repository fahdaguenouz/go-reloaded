package main

import (
	"fmt"
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
			
			for _, line := range lines {
				res := functions.ApplyParenthesesLogic(line)
				res = functions.ProcessSingleQuotes(res)
				finalLine := functions.ReplaceAWithAn(res)
				finalLine = functions.Ponctuation(finalLine)
				finalResult += finalLine + "\n"
			}

			err = os.WriteFile(output, []byte(finalResult), 0o644)
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
