package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	Words "github.com/mindslayer001/typeit/utils"
	"golang.org/x/term"
)

var oldState *term.State

func main() {
	// Handle Ctrl+C cleanup
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGINT)
	go func() {
		<-sigChan
		fmt.Println("\nCaught Ctrl+C! Cleaning up...")
		term.Restore(int(os.Stdin.Fd()), oldState)
		os.Exit(0)
	}()

	fmt.Print("TYPE IT\n")
	var err error
	oldState, err = term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("Terminal size: %dx%d\n", width, height)

	wordList := Words.GetWords()

	for _, word := range wordList {
		fmt.Println("\nType this word:", word)

		capturedWord := make([]byte, len(word))
		capturedChar := make([]byte, 1)
		i := 0

		for i < len(word) {
			_, err := os.Stdin.Read(capturedChar)
			if err != nil {
				fmt.Println("Error reading:", err)
				return
			}

			// Handle backspace (both 127 and 8)
			if capturedChar[0] == 127 || capturedChar[0] == 8 {
				if i > 0 {
					i--
					fmt.Print("\b \b") // erase the last character visually
				}
				continue
			}

			// Record character
			capturedWord[i] = capturedChar[0]

			// Compare
			if capturedChar[0] == word[i] {
				fmt.Printf("%c", capturedChar[0]) // print char if correct
			} else {
				fmt.Printf("\nWrong character at position %d: expected '%c', got '%c'\n", i, word[i], capturedChar[0])
			}

			i++
		}
	}
}
