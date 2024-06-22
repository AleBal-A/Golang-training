package handlers

import (
	"main/ex01/internal/domain/models"
	"main/ex01/internal/repository"
	"net/http"
)

type SubmitHandler struct {
	repo repository.ArticleRepository
}

func NewSubmitHandler(repo repository.ArticleRepository) *SubmitHandler {
	return &SubmitHandler{repo: repo}
}

func (h *SubmitHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		title := r.FormValue("title")
		content := r.FormValue("content")

		article := models.Article{Title: title, Content: content}
		err := h.repo.CreateArticle(&article)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
