package main

import (
	"bufio"
	"log"
	"main/ex01/internal/domain/models"
	"main/ex01/internal/handlers"
	"main/ex01/internal/middleware"
	"main/ex01/internal/repository"
	"net/http"
	"os"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"html/template"
)

const (
	port            = ":8888"
	articlesPerPage = 3
)

func main() {
	err := loadEnv()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dsn := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}

	err = db.AutoMigrate(&models.Article{})
	if err != nil {
		log.Fatal("failed to migrate database:", err)
	}

	articleRepo := repository.NewArticleRepository(db)
	tmpl := template.Must(template.ParseFiles("../../web/templates/index.html", "../../web/templates/admin.html"))

	router := handlers.SetupRoutes(articleRepo, tmpl, articlesPerPage)

	// Все маршруты черз middleware на случай рейтлимита
	loggedRouter := middleware.RateLimiter(router)

	log.Printf("Server starting at port %s", port)
	log.Fatal(http.ListenAndServe(port, loggedRouter))
}

func loadEnv() error {
	file, err := os.Open(".env")
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		err = os.Setenv(parts[0], parts[1])
		if err != nil {
			return err
		}
	}

	return scanner.Err()
}
