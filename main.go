package main

import (
	"fmt"
	"log"
	"os"
	"bufio"

	"github.com/Jorghy-Del/gorth/eval"
	"github.com/Jorghy-Del/gorth/lexer"
	"github.com/Jorghy-Del/gorth/word"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal(fmt.Sprintf("usage: %s <filename>", os.Args[0]))
	}
	filename := os.Args[1]

	fh, err := os.Open(filename)
	defer fh.Close()
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(fh)

	var words []word.Word
	dictionary := make(map[string][]word.Word)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Printf("line: %s\n", line)
		lxr := lexer.New(line)

		for {
			w, d := lxr.NextToken()
			if w.Type == word.EOF {
				fmt.Printf("%v\n", w.Literal)
				break
			}
			for udf, def := range d {
				if _, ok := dictionary[udf]; !ok {
					dictionary[udf] = def
				} else {
					words = append(words, def...)
				}
			}
			words = append(words, w)
			if v, ok := dictionary[w.Literal]; ok {
				words = append(words, v...)
			}
		}

		returnStack := eval.Execute(words)
		fmt.Printf("return stack: %v\n\n", returnStack)
	}
}
