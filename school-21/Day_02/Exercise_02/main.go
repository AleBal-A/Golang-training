package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: myXargs <command>")
		os.Exit(1)
	}

	programName := os.Args[1]
	args := os.Args[2:]

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		args = append(args, input.Text())
	}

	if err := input.Err(); err != nil {
		log.Fatal("Error while reading stdin:", err)
	}

	cmd := exec.Command(programName, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatalf("Command execution failed: %v", err)
	}
}
