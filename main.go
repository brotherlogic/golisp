package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		text, err := reader.ReadString('\n')
		text = strings.TrimRight(text, "\n")
		if err != nil {
			return
		}
		fmt.Println(rep(text))
	}
}
