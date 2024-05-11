package main

import (
	"archive/tar"
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
)

func main() {
	var archiveDir string
	flag.StringVar(&archiveDir, "a", "", "Archives all files in a directory")
	flag.Parse()

	files := flag.Args()
	if len(files) == 0 {
		log.Fatal("Usage: \"./myRotate [-a] <directory>\"\n or\n \"/myRotate [-a] <directory>...\"")
	}

	if archiveDir != "" && exists(archiveDir) {
		log.Fatalf("Archive directory does not exist: %s\n", archiveDir)
	}

	var wg sync.WaitGroup
	for _, path := range files {
		if archiveDir == "" {
			archiveDir = filepath.Dir(path)
		}

		wg.Add(1)
		go func(path, archiveDir string) {
			defer wg.Done()
			err := rotateLog(path, archiveDir)
			if err != nil {
				fmt.Printf("Failed to rotate log %s: %v\n", path, err)
			}
		}(path, archiveDir)
	}
	wg.Wait()
}

func rotateLog(logFile, archiveDir string) error {
	fileInfo, err := os.Stat(logFile)
	if err != nil {
		log.Fatalf("Error while retrieving file information: %s", err)
	}

	modTimeUnix := fileInfo.ModTime().Unix()
	baseName := filepath.Base(logFile)
	newFileName := fmt.Sprintf("%s_%d.tar.gz", baseName, modTimeUnix)

	targetPath := filepath.Join(archiveDir, newFileName)

	file, err := os.Create(targetPath)
	if err != nil {
		return err
	}
	defer file.Close()

	gzipWriter := gzip.NewWriter(file)
	defer gzipWriter.Close()

	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	return addToArchive(tarWriter, logFile, baseName)
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func addToArchive(tarWriter *tar.Writer, logFile, baseName string) error {
	file, err := os.Open(logFile)
	if err != nil {
		return err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	header, err := tar.FileInfoHeader(fileInfo, fileInfo.Name())
	if err != nil {
		return err
	}
	header.Name = baseName

	if err := tarWriter.WriteHeader(header); err != nil {
		return err
	}

	if _, err := io.Copy(tarWriter, file); err != nil {
		return err
	}

	return nil
}
