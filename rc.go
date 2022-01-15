package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"unicode/utf8"
	"github.com/pkg/errors"
    "golang.org/x/crypto/ssh/terminal"
)

const (
	Args  = iota
	Stdin = iota
)

func main() {

	var error error
	var status int
	defer exit(&status, &error)

	args := os.Args[1:]
	stdin := os.Stdin

	var useType int
	useType, status, error = validateInputs(args, stdin)
	if error != nil {
		return
	}

	var runeCount int
	switch useType {
	case Args:
		runeCount = countArgs(args)
	case Stdin:
		runeCount, status, error = countFile(stdin)
		if error != nil {
			return
		}
	}

	output := strconv.Itoa(runeCount)
	fmt.Println(output)
}

func validateInputs(args []string, stdin *os.File) (useType int, status int, error error) {
	hasArgs := len(args) > 0
	hasStream := !terminal.IsTerminal(int(os.Stdin.Fd()))

	if hasArgs {
		if hasStream {
			error = errors.New("Cannot mix arguments and standard input.")
			status = 2
		} else {
			useType = Args
		}
	} else {
		if hasStream {
			useType = Stdin
		} else {
			error = errors.New("Either arguments or standard input must be provided.")
			status = 2
		}
	}

	return
}

func countArgs(args []string) int {

	runeCount := 0

	for _, arg := range args {
		runeCount += utf8.RuneCountInString(arg)
	}

	return runeCount
}

func countFile(file *os.File) (count int, status int, error error) {

	reader := bufio.NewReader(file)

	var last rune
	for {
		var rune rune
		rune, _, error = reader.ReadRune()
		if error != nil {
			if error == io.EOF {
				if last == '\n' {
					count -= 1
				}
				error = nil
				break
			} else {
				status = 1
				return
			}
		}
		count += 1
		last = rune
	}

	return
}

func exit(status *int, error *error) {

	switch *status {
	case 1:
		fmt.Printf("%+v\n", *error)
	case 2:
		fmt.Printf("%s\n", *error)
	}
	os.Exit(*status)
}
