package main

import (
	"fmt"
	"log"

	db "day_01/src/DBReader_lib"
)

func CompareRecipes(recipes1, recipes2 db.Recipes) {
	cakes1 := make(map[string]db.Cake)
	cakes2 := make(map[string]db.Cake)

	for _, cake := range recipes1.Cakes {
		cakes1[cake.Name] = cake
	}
	for _, cake := range recipes2.Cakes {
		cakes2[cake.Name] = cake
	}

	// Проверяется каждый торт из вторй структуры
	for name := range cakes2 {
		if _, ok := cakes1[name]; !ok {
			fmt.Printf("ADDED cake \"%s\" \n", name)
		}
	}

	// Проверяется каждый торт из первой структуры
	for name := range cakes1 {
		_, ok := cakes2[name]
		if !ok {
			fmt.Printf("REMOVED cake \"%s\"\n", name)
			continue
		}
	}

	// Двойной проход нужен для определенного формата вывода
	for name, cake1 := range cakes1 {
		cake2, ok := cakes2[name]
		if ok {
			if cake1.StoveTime != cake2.StoveTime {
				fmt.Printf("CHANGED cooking time for cake \"%s\" - \"%s\" instead of \"%s\"\n", name, cake2.StoveTime, cake1.StoveTime)
			}
			compareIngredients(cake1, cake2)
		}
	}

}

// Сравнение ингредиентов торта
func compareIngredients(cake1, cake2 db.Cake) {
	ingredientMap1 := make(map[string]db.Ingredient)
	ingredientMap2 := make(map[string]db.Ingredient)

	for _, ingredient := range cake1.Ingredients {
		ingredientMap1[ingredient.ItemName] = ingredient
	}
	for _, ingredient := range cake2.Ingredients {
		ingredientMap2[ingredient.ItemName] = ingredient
	}

	for name, ingredient1 := range ingredientMap1 {
		ingredient2, ok := ingredientMap2[name]
		if !ok {
			fmt.Printf("REMOVED ingredient \"%s\" for cake \"%s\"\n", name, cake1.Name)
			continue
		}

		if ingredient1.ItemCount != ingredient2.ItemCount {
			fmt.Printf("CHANGED unit count for ingredient \"%s\" for cake: \"%s\" - \"%s\" instead of \"%s\"\n", name, cake1.Name, ingredient2.ItemCount, ingredient1.ItemCount)
		}
		if ingredient1.ItemUnit != ingredient2.ItemUnit {
			if ingredient1.ItemUnit != "" && ingredient2.ItemUnit == "" {
				fmt.Printf("REMOVED unit \"%s\" for ingredient \"%s\" for cake \"%s\"\n", ingredient1.ItemUnit, name, cake1.Name)
				continue
			}
			fmt.Printf("CHANGED unit for ingredient \"%s\" for cake \"%s\" - \"%s\" instead of \"%s\"\n", name, cake1.Name, ingredient2.ItemUnit, ingredient1.ItemUnit)
		}

	}

	for name := range ingredientMap2 {
		if _, ok := ingredientMap1[name]; !ok {
			fmt.Printf("ADDED ingredient \"%s\" for cake \"%s\" \n", name, cake1.Name)
		}
	}
}

func main() {
	fileName1, oldFileExt, fileName2, newFileExt := db.GetTwoFilesPaths()

	var reader1 db.DBReader
	if oldFileExt == ".xml" {
		reader1 = db.XMLReader{}
	} else if oldFileExt == ".json" {
		reader1 = db.JSONReader{}
	} else {
		log.Fatalf("Неподдерживаемый формат файла \"%s\"", oldFileExt)
	}

	var reader2 db.DBReader
	if newFileExt == ".xml" {
		reader2 = db.XMLReader{}
	} else if newFileExt == ".json" {
		reader2 = db.JSONReader{}
	} else {
		log.Fatalf("Неподдерживаемый формат файла \"%s\"", newFileExt)
	}

	recipes1, err := reader1.ReadRecipesFromFile(fileName1)
	if err != nil {
		log.Fatalf("Ошибка при чтении файла \"%s\": %v", fileName1, err)
	}

	recipes2, err := reader2.ReadRecipesFromFile(fileName2)
	if err != nil {
		log.Fatalf("Ошибка при чтении файла \"%s\": %v", fileName2, err)
	}

	CompareRecipes(recipes1, recipes2)
}
