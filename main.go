package main

import (
	"fmt"
	"os"

	"github.com/zhcHoward/MyDictionary-Go/api"
)

func main() {
	var word, dictName string
	switch len(os.Args) {
	case 1:
		fmt.Fprintln(os.Stderr, "Please type in the word you want to look up")
		os.Exit(0)
	case 2:
		word = os.Args[1]
		dictName = "iciba"
	case 3:
		word = os.Args[1]
		dictName = os.Args[2]
	}
	dict, err := api.GetService(dictName)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	dict.Search(word)
}
