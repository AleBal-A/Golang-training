package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// flags
	var showSymLinks, showDirs, showFiles bool
	var fileExtension string
	flag.BoolVar(&showSymLinks, "sl", false, "Show symbolic links")
	flag.BoolVar(&showDirs, "d", false, "Show directories")
	flag.BoolVar(&showFiles, "f", false, "Show only files")
	flag.StringVar(&fileExtension, "ext", "", "Filter files by extension (use with -f)")
	flag.Parse()

	pathArg := flag.Args()
	if len(pathArg) != 1 {
		fmt.Println("Usage: ./myFind [-d] [-sl] [-f] [f:optional -ext extension] <path>")
		os.Exit(1)
	}

	var visitFunc filepath.WalkFunc = func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Failed to access %s: %v\n", path, err)
			return nil
		}

		fileMode := info.Mode()
		switch {
		case fileMode.IsRegular() && showFiles:
			if fileExtension != "" {
				if strings.HasSuffix(info.Name(), "."+fileExtension) {
					fmt.Println(path)
				}
			} else {
				fmt.Println(path)
			}
		case fileMode.IsDir() && showDirs:
			fmt.Println(path)
		case fileMode&os.ModeSymlink != 0 && showSymLinks:
			realPath, err := os.Readlink(path)
			if err != nil {
				fmt.Printf("%s  -> [broken]\n", path)
			} else {
				fmt.Printf("%s  -> %s\n", path, realPath)
			}
		}
		return nil
	}
	// Walk the file tree
	err := filepath.Walk(pathArg[0], visitFunc)
	if err != nil {
		fmt.Printf("Error walking the path %s: %v\n", pathArg[0], err)
	}
}
