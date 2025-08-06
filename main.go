package main

import (
	"fmt"
	"os"

	"github.com/mindslayer001/typeit/utils"
	Words "github.com/mindslayer001/typeit/words"
	"golang.org/x/term"
)

var oldState *term.State

func restoreTerminal() {
	if oldState != nil {
		term.Restore(int(os.Stdin.Fd()), oldState)
	}
}

func GetResult(capturedWord []byte, target string) string {
	result := ""
	i := 0
	for i < len(capturedWord) {
		switch capturedWord[i] {
		case target[i]:
			result += utils.BasicColors.Green + string(target[i])
		case '\x00':
			result += utils.BasicColors.Reset + string(target[i:])
			return result
		default:
			result += utils.BasicColors.Red + string(target[i])
		}
		i++
	}
	result += utils.BasicColors.Reset
	return result
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
	fmt.Print("\r\nGame Started:\n")
	for _, word := range wordList {
		fmt.Print("\r\033[K")
		fmt.Print(word)

		capturedChar := make([]byte, 1)
		capturedWord := make([]byte, len(word))
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
				if i > 0 && len(capturedWord) > 0 {
					capturedWord[i-1] = '\x00'
					i--
				} else {
					capturedWord[i] = '\x00'
				}
			default:
				capturedWord[i] = capturedChar[0]
				i++
			}
			fmt.Print("\r" + GetResult(capturedWord, word))
		}
	}
	//fmt.Print("\033[H\033[2J") // Clear screen
}
