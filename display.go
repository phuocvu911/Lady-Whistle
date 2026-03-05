package main

import (
	"fmt"
	"strings"
)

// PrintHeader prints a double-line bordered section header.
func PrintHeader(title string) {
	line := strings.Repeat("═", 60)
	fmt.Printf("\n%s\n  %s\n%s\n", line, title, line)
}

// PrintDivider prints a single-line divider.
func PrintDivider() {
	fmt.Println(strings.Repeat("─", 60))
}
