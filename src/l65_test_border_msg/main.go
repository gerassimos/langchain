package main

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
)

// ## ANSI 16 colors (4-bit)
// lipgloss.Color("5")  // magenta
// lipgloss.Color("9")  // red
// lipgloss.Color("12") // light blue

// ## ANSI 256 Colors (8-bit)
// lipgloss.Color("86")  // aqua
// lipgloss.Color("201") // hot pink
// lipgloss.Color("202") // orange
// lipgloss.Color("63") // purple/blue

// ## True Color (16,777,216 colors; 24-bit)
//lipgloss.Color("#0000FF") // good ol' 100% blue
//lipgloss.Color("#04B575") // a green
//lipgloss.Color("#3C3C3C") // a dark gray

func main() {
	fmt.Println("Hello, World!")
	border := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63")).
		Padding(1, 2).
		Margin(1).
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57"))

	title := lipgloss.NewStyle().
		Foreground(lipgloss.Color("63")).
		Bold(true)

	box := border.Render("Hello from Go in a nice box!")
	fmt.Println(title.Render(" INFO ") + "\n" + box)

	printMessageWithBorder("This is a message with a border!")
}

func printMessageWithBorder(msg string) {
	var style = lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("63"))

	style = style.SetString("TEST")

	fmt.Println(style.Render(msg))
}
