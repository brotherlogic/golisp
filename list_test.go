package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
)

func Test0(t *testing.T) {
	lines, err := os.ReadFile("testdata/test0.txt")
	if err != nil {
		t.Fatalf("Cannot read file: %v", err)
	}

	buffer := ""
	for _, line := range strings.Split(string(lines), "\n") {
		if strings.HasPrefix(line, ";;") {
			log.Printf("%v", line)
		} else if strings.HasPrefix(line, ";") {
			if buffer != line[1:] {
				t.Errorf("Bad Process: %v vs %v", buffer, line)
			}
		} else if len(line) > 0 {
			buffer = fmt.Sprintf("=>%v", rep(line))
		}
	}
}
