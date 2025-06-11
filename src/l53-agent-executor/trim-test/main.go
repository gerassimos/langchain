package main

import (
	"fmt"
	"strings"
)

func main() {
	s := strings.TrimSpace(`
    Ice Cream	 


`)
	fmt.Printf(">>%s<<", s)
}
