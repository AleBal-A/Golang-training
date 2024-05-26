package main

import (
	"bufio"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
)

func main() {

	cfg := elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"},
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating client: %s", err)
	}

	indexName := "places"

	deleteIndex(es, indexName)

	res, err := es.Indices.Exists([]string{indexName})
	if err != nil {
		log.Fatalf("Error checking if index exists: %s", err)
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		log.Printf("Index %s already exists", indexName)
	} else {
		createIndex(es, "places")
	}

	if err := processCSVFile("./data.csv", es); err != nil {
		log.Fatalf("Error processing CSV file: %s", err)
	}
}

func processCSVFile(filePath string, es *elasticsearch.Client) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}

	defer file.Close()

	reader := csv.NewReader(bufio.NewReader(file))
	reader.Comma = '\t' // Указываем, что разделитель поля - табуляция

	// Пропуск заголовка
	if _, err := reader.Read(); err != nil {
		return err
	}

	var buffer strings.Builder

	// Парсинг каждой строки
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		doc := map[string]interface{}{
			"id":      record[0],
			"name":    record[1],
			"address": record[2],
			"phone":   record[3],
			"location": map[string]interface{}{
				"lat": parseFloat(record[5]),
				"lon": parseFloat(record[4]),
			},
		}

		jsonDoc, err := json.Marshal(doc)
		if err != nil {
			return err
		}

		// Метаданные индексации с указанием идентификатора документа
		meta := fmt.Sprintf(`{"index":{"_index":"places", "_id": %q}}`, record[0])
		buffer.WriteString(meta + "\n" + string(jsonDoc) + "\n")
	}

	res, err := es.Bulk(strings.NewReader(buffer.String()), es.Bulk.WithIndex("places"))
	if err != nil {
		return err
	}
	defer res.Body.Close()
	log.Println("Datas successfully added!")
	return nil
}

func parseFloat(str string) float64 {
	f, err := strconv.ParseFloat(str, 64)
	if err != nil {
		log.Fatalf("Error parsing float: %s", err)
	}
	return f
}

func createIndex(es *elasticsearch.Client, indexName string) {
	// Определение маппинга для индекса с дополнительным полем id
	mapping := `{
        "settings": {
            "number_of_shards": 1,
            "number_of_replicas": 1
        },
        "mappings": {
            "properties": {
                "id": {"type": "keyword"},  // Добавляем id как ключевое поле
                "name": {"type": "text"},
                "address": {"type": "text"},
                "phone": {"type": "text"},
                "location": {"type": "geo_point"}
            }
        }
    }`

	// Создаем индекс
	res, err := es.Indices.Create(
		indexName,
		es.Indices.Create.WithBody(strings.NewReader(mapping)),
		es.Indices.Create.WithContext(context.Background()),
	)
	if err != nil {
		log.Fatalf("Cannot create index: %s", err)
	}
	defer res.Body.Close()
	log.Println("Index created successfully")
}

func deleteIndex(es *elasticsearch.Client, indexName string) {
	res, err := es.Indices.Delete([]string{indexName}, es.Indices.Delete.WithContext(context.Background()))
	if err != nil {
		log.Fatalf("Error deleting index: %s", err)
	}
	defer res.Body.Close()

	// Проверка статуса HTTP для определения успешности удаления
	if res.StatusCode == 200 {
		log.Printf("Index %s successfully deleted.", indexName)
	} else {
		log.Printf("Failed to delete index %s. Status: %d", indexName, res.StatusCode)
	}
}
