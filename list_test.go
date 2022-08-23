package main

import (
	"log"
	"os"
	"strings"
	"testing"
)


func TestBasic(t *testing.T) {
	t.Errorf("Not ready for testing")
}

func Test0(t *testing.T) {
	lines, err := os.ReadFile("testdata/test0.txt")
	if err != nil {
		t.Fatalf("Cannot read file: %v", err)
	}

	for _, line := range strings.Split(string(lines), "\n") {
		if strings.HasPrefix(line, ";;") {
			log.Printf("%v", line)
		}
	}
}