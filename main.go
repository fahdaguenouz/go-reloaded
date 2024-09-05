package main

import (
	"fmt"
	"io"
	"os"
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
			datafile := string(data)
			instance := ""
			inpar := false
			for i := 0; i < len(datafile); i++ {
				if datafile[i] == '(' {
					inpar = true
					continue
				}
				if datafile[i] == ')' {
					inpar = false
					continue
				}
				if inpar {
					instance += string(datafile[i])
				}
			}
			fmt.Println(instance)
		}

	} else if len(args) > 2 {
		fmt.Println("Too much arguments")
		os.Exit(1)
	} else if len(args) == 1 {
		fmt.Println("less arguments please enter the input and the output files name")
		os.Exit(1)
	}
}
