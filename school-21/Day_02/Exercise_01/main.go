package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
)

func main() {
	var countLines, countSymbols, countWords bool
	flag.BoolVar(&countLines, "l", false, "Show the number of lines")
	flag.BoolVar(&countSymbols, "m", false, "Show the number of symbols")
	flag.BoolVar(&countWords, "w", false, "Show the number of word")
	flag.Parse()

	if countLines && countSymbols || countLines && countWords || countSymbols && countWords {
		log.Fatal("Error: Flags -l, -m, and -w are mutually exclusive.")
	}
	if !countLines && !countSymbols && !countWords {
		countWords = true // default behavior
	}

	// Files to process
	files := flag.Args()
	if len(files) == 0 {
		log.Fatal("Usage: ./myWc [-l] [-m] [-w] <file1> [file2] ...")
	}

	var wg sync.WaitGroup
	for _, path := range files {
		wg.Add(1)
		go func(path string) {
			defer wg.Done()
			if err := ProcessFile(path, countLines, countSymbols, countWords); err != nil {
				log.Printf("Error processing file %s: %v\n", path, err)
			}
		}(path)
	}
	wg.Wait()
}

func ProcessFile(path string, countLines, countSymbols, countWords bool) error {

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	var lines, symbols, words int
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines++
		line := scanner.Text()
		symbols += len(line)
		words += len(strings.Fields(line))
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	switch {
	case countLines:
		fmt.Printf("%d\t%s\n", lines, path)
	case countSymbols:
		fmt.Printf("%d\t%s\n\n", symbols, path)
	case countWords:
		fmt.Printf("%d\t%s\n", words, path)
	}

	return nil
}
