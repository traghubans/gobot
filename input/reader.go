package input

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Reader handles reading user input from the console
type Reader struct {
	reader *bufio.Reader
}

// NewReader creates a new Reader instance
func NewReader() *Reader {
	return &Reader{
		reader: bufio.NewReader(os.Stdin),
	}
}

// ReadLine reads a single line of input from the user
func (r *Reader) ReadLine() string {
	input, err := r.reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Error reading input: %v\n", err)
		return ""
	}
	return strings.TrimSpace(input)
}

// ReadPrompt reads input with a prompt
func (r *Reader) ReadPrompt(prompt string) string {
	fmt.Print(prompt)
	return r.ReadLine()
}
