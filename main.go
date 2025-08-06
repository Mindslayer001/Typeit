package main

import (
	"fmt"
	"os"

	Words "github.com/mindslayer001/typeit/utils"
	"golang.org/x/term"
)

var oldState *term.State

func restoreTerminal() {
	if oldState != nil {
		term.Restore(int(os.Stdin.Fd()), oldState)
	}
}

func main() {
	var err error
	oldState, err = term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	defer restoreTerminal() // Automatic Code end

	fmt.Println("TYPE IT")
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("Terminal size: %dx%d\n", width, height)

	wordList := Words.GetWords()

	for _, word := range wordList {
		fmt.Printf("\r\nType this word: %s\n", word)
		fmt.Print("> ")

		capturedChar := make([]byte, 1)
		i := 0

		for i < len(word) {
			_, err := os.Stdin.Read(capturedChar)
			if err != nil {
				fmt.Println("\nError reading:", err)
				return
			}

			ch := capturedChar[0]

			switch ch {
			case 3: // Ctrl+C
				fmt.Println("\n Caught Ctrl+C")
				restoreTerminal()
				os.Exit(0)
			case 127, 8: // Backspace
				if i > 0 {
					i--
					fmt.Print("\b \b")
				}
				continue
			}

			if ch == word[i] {
				fmt.Printf("%c", ch)
				i++
			} else {
				fmt.Print("\a") // Beep
			}
		}
	}
}
