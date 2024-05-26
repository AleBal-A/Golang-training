package main

import (
	"day3/src/db"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"strconv"
	"text/template"
	"time"
)

var jwtKey = []byte("your_secret_key")

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type Server struct {
	store db.Store
}

func (s *Server) placesHandler(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	if pageStr == "" {
		pageStr = "1"
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		http.Error(w, fmt.Sprintf("Invalid 'page' value: '%s'", pageStr), http.StatusBadRequest)
		return
	}

	limit := 10
	offset := (page - 1) * limit

	places, total, err := s.store.GetPlaces(limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	funcMap := template.FuncMap{
		"add": func(x, y int) int { return x + y },
		"sub": func(x, y int) int { return x - y },
	}

	tmpl, err := template.New("places").Funcs(funcMap).Parse(`
	<!doctype html>
	<html>
	<head>
		<meta charset="utf-8">
		<title>Places</title>
		<meta name="description" content="">
		<meta name="viewport" content="width=device-width, initial-scale=1">
	</head>

	<body>
	<h5>Total: {{.Total}}</h5>
	<ul>
		{{range .Places}}
		<li>
			<div>{{.Name}}</div>
			<div>{{.Address}}</div>
			<div>{{.Phone}}</div>
		</li>
		{{end}}
	</ul>
	{{if gt .Page 1}}
	<a href="/?page={{sub .Page 1}}">Previous</a>
	{{end}}
	{{if lt .Page .LastPage}}
	<a href="/?page={{add .Page 1}}">Next</a>
	{{end}}
	<a href="/?page={{.LastPage}}">Last</a>
	</body>
	</html>`)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	lastPage := (total + limit - 1) / limit

	data := struct {
		Places   []db.Place
		Total    int
		Page     int
		LastPage int
	}{
		Places:   places,
		Total:    total,
		Page:     page,
		LastPage: lastPage,
	}

	tmpl.Execute(w, data)
}

func (s *Server) apiPlacesHandler(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	if pageStr == "" {
		pageStr = "1"
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": fmt.Sprintf("Invalid 'page' value: '%s'", pageStr),
		})
		return
	}

	limit := 10
	offset := (page - 1) * limit

	places, total, err := s.store.GetPlaces(limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	lastPage := (total + limit - 1) / limit

	data := struct {
		Name     string     `json:"name"`
		Total    int        `json:"total"`
		Places   []db.Place `json:"places"`
		PrevPage int        `json:"prev_page,omitempty"`
		NextPage int        `json:"next_page,omitempty"`
		LastPage int        `json:"last_page"`
	}{
		Name:     "Places",
		Total:    total,
		Places:   places,
		LastPage: lastPage,
	}

	if page > 1 {
		data.PrevPage = page - 1
	}
	if page < lastPage {
		data.NextPage = page + 1
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func (s *Server) apiRecommendHandler(w http.ResponseWriter, r *http.Request) {
	latStr := r.URL.Query().Get("lat")
	lonStr := r.URL.Query().Get("lon")

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": fmt.Sprintf("Invalid 'lat' value: '%s'", latStr),
		})
		return
	}

	lon, err := strconv.ParseFloat(lonStr, 64)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": fmt.Sprintf("Invalid 'lon' value: '%s'", lonStr),
		})
		return
	}

	limit := 3
	places, err := s.store.GetClosestPlaces(lat, lon, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		Name   string     `json:"name"`
		Places []db.Place `json:"places"`
	}{
		Name:   "Recommendation",
		Places: places,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func (s *Server) getTokenHandler(w http.ResponseWriter, r *http.Request) {
	// Создаем временные данные токена
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		Username: "admin",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// Создаем токен
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	// Возвращаем токен клиенту
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": tokenString,
	})
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Извлекаем токен из заголовка
		tokenStr := r.Header.Get("Authorization")
		if tokenStr == "" {
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}

		// Проверяем формат токена
		if len(tokenStr) < 7 || tokenStr[:7] != "Bearer " {
			http.Error(w, "Invalid token format", http.StatusUnauthorized)
			return
		}
		tokenStr = tokenStr[7:]

		// Проверяем токен
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	cfg := elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"},
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	store := db.NewElasticsearchStore(es)

	server := &Server{store: store}

	http.HandleFunc("/", server.placesHandler)
	http.HandleFunc("/api/places", server.apiPlacesHandler)
	http.HandleFunc("/api/get_token", server.getTokenHandler)
	http.Handle("/api/recommend", authMiddleware(http.HandlerFunc(server.apiRecommendHandler)))

	log.Println("Starting server on :8888")
	log.Fatal(http.ListenAndServe(":8888", nil))
}
