package handlers

import (
	"html/template"
	"main/ex01/internal/domain/models"
	"main/ex01/internal/repository"
	"net/http"
	"strconv"
)

type IndexHandler struct {
	repo  repository.ArticleRepository
	tmpl  *template.Template
	limit int
}

type PageData struct {
	Articles     []models.Article
	CurrentPage  int
	TotalPages   int
	PreviousPage int
	NextPage     int
}

func NewIndexHandler(repo repository.ArticleRepository, tmpl *template.Template, limit int) *IndexHandler {
	return &IndexHandler{repo: repo, tmpl: tmpl, limit: limit}
}

func (h *IndexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	offset := (page - 1) * h.limit
	articles, err := h.repo.GetArticles(offset, h.limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	totalArticles, err := h.repo.CountArticles()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	totalPages := (totalArticles + h.limit - 1) / h.limit

	pageData := PageData{
		Articles:     articles,
		CurrentPage:  page,
		TotalPages:   totalPages,
		PreviousPage: page - 1,
		NextPage:     page + 1,
	}

	if pageData.PreviousPage < 1 {
		pageData.PreviousPage = 1
	}
	if pageData.NextPage > totalPages {
		pageData.NextPage = totalPages
	}

	err = h.tmpl.ExecuteTemplate(w, "index.html", pageData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
