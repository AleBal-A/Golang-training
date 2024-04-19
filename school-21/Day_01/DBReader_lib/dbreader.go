package db

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Recipes ВНИМАНИЕ !!! Структура для корневого элемента XML-документа "recipes"
type Recipes struct {
	Cakes []Cake `xml:"cake" json:"cake"`
}

// Cake ВНИМАНИЕ !!! Структура для каждого элемента "cake"
type Cake struct {
	Name        string       `xml:"name" json:"name"`
	StoveTime   string       `xml:"stovetime" json:"time"`
	Ingredients []Ingredient `xml:"ingredients>item" json:"ingredients"`
}

type Ingredient struct {
	ItemName  string `xml:"itemname" json:"ingredient_name"`
	ItemCount string `xml:"itemcount" json:"ingredient_count"`
	ItemUnit  string `xml:"itemunit,omitempty" json:"ingredient_unit,omitempty"`
}

// DBReader - через этот интрефейс должно читаться два файла. В результате должен выдавать одни и те же типы объектов
type DBReader interface {
	ReadRecipesFromFile(fileName *string) (Recipes, error)
}

// JSONReader реализует интерфейс DBReader для чтения JSON файлов
type JSONReader struct{}

func (jr JSONReader) ReadRecipesFromFile(fileName *string) (Recipes, error) {
	var recipes Recipes
	dataFile, err := os.Open(*fileName)
	if err != nil {
		return recipes, err
	}
	defer dataFile.Close()

	fileData, err := io.ReadAll(dataFile)
	if err != nil {
		return recipes, err
	}

	err = json.Unmarshal(fileData, &recipes)
	return recipes, err
}

// XMLReader реализует интерфейс DBReader для чтения XML файлов
type XMLReader struct{}

func (xr XMLReader) ReadRecipesFromFile(fileName *string) (Recipes, error) {
	var recipes Recipes
	dataFile, err := os.Open(*fileName)
	if err != nil {
		return recipes, err
	}
	defer dataFile.Close()

	fileData, err := io.ReadAll(dataFile)
	if err != nil {
		return recipes, err
	}

	err = xml.Unmarshal(fileData, &recipes)
	return recipes, err
}

func CreateJson(recipes Recipes, fileName *string, ext string) error {
	jsonData, err := json.MarshalIndent(recipes, "", "  ")
	if err != nil {
		log.Fatalf("Error due converting to JSON: %v", err)
	}

	baseName := strings.TrimSuffix(*fileName, ext)
	jsonFileName := baseName + ".json"

	err = os.WriteFile(jsonFileName, jsonData, 0644)
	if err != nil {
		log.Fatalf("Failed to write to file: %v", err)
	}

	log.Printf("JSON data successfully written to: %v", jsonFileName)

	return err
}

func GetPath() (*string, string) {
	fileName := flag.String("f", "", "Path to file")
	flag.Parse()
	if *fileName == "" {
		log.Fatal("-f path is empty")
	}

	ext := filepath.Ext(*fileName)
	return fileName, ext
}

func GetTwoFilesPaths() (*string, string, *string, string) {
	oldFileName := flag.String("old", "", "Path to the original file")
	newFileName := flag.String("new", "", "Path to the new file")
	flag.Parse()

	if *oldFileName == "" || *newFileName == "" {
		log.Fatal("One or both file paths are empty")
	}

	oldFileExt := filepath.Ext(*oldFileName)
	newFileExt := filepath.Ext(*newFileName)

	return oldFileName, oldFileExt, newFileName, newFileExt
}
