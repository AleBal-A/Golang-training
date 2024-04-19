package main

import (
	"bufio"
	db "day_01/src/DBReader_lib"
	"fmt"
	"log"
	"os"
)

func main() {
	oldFilePath, oldFileExt, newFilePath, newFileExt := db.GetTwoFilesPaths()

	if *oldFilePath == "" || *newFilePath == "" {
		log.Fatal("One or both file paths are empty")
	}
	if oldFileExt != ".txt" && newFileExt != ".txt" {
		log.Fatal("The file extension must be \".txt\"")
	}

	oldFileStrings := make(map[string]struct{})
	oldFile, err := os.Open(*oldFilePath)
	if err != nil {
		log.Fatal("Failed to open old file: %s", err)
	}
	defer oldFile.Close()

	scanner := bufio.NewScanner(oldFile)
	for scanner.Scan() {
		oldFileStrings[scanner.Text()] = struct{}{}
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("Failed to read \"%s\": %s", *oldFilePath, err)
	}

	newFile, err := os.Open(*newFilePath)
	if err != nil {
		log.Fatal("Failed to open new file: %s", err)
	}
	defer newFile.Close()

	scanner = bufio.NewScanner(newFile)
	for scanner.Scan() {
		line := scanner.Text()
		if _, found := oldFileStrings[line]; found {
			delete(oldFileStrings, line)
		} else {
			fmt.Printf("ADDED %s\n", line)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("Failed to read \"%s\": %s", *newFilePath, err)
	}

	for textLine := range oldFileStrings {
		fmt.Printf("REMOVED %s\n", textLine)
	}

}
