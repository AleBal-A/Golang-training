package main

import (
	db "day_01/src/DBReader_lib"
	"log"
)

func main() {
	fileName, ext := db.GetPath()
	if ext != ".xml" && ext != ".json" {
		log.Fatalf("Unknown file type: %s", ext)
	}

	var reader db.DBReader
	if ext == ".xml" {
		reader = db.XMLReader{}
	} else {
		reader = db.JSONReader{}
	}

	recipes, err := reader.ReadRecipesFromFile(fileName)
	if err != nil {
		log.Fatalf("Error due to file parsing: %v", err)
		return
	}

	err = db.CreateJson(recipes, fileName, ext)
	if err != nil {
		log.Fatalf("Failed to write to file: %v", err)
	}
}
