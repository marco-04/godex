package main

import (
	"fmt"
	"os"
	"bufio"
)

func main() {
	stdin := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		for !stdin.Scan() {
			fmt.Print("Pokedex > ")
		}
		cleanedText := cleanInput(stdin.Text())
		dispatchCommand(cleanedText)
	}
}
