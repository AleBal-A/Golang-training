package handlers

import (
	"github.com/gorilla/mux"
	"html/template"
	"main/ex01/internal/repository"
	"net/http"
)

func SetupRoutes(articleRepo repository.ArticleRepository, tmpl *template.Template, articlesPerPage int) *mux.Router {
	r := mux.NewRouter()
	r.Handle("/", NewIndexHandler(articleRepo, tmpl, articlesPerPage)).Methods("GET")
	r.Handle("/admin", NewAdminHandler(tmpl)).Methods("GET")
	r.Handle("/login", NewAuthHandler()).Methods("GET", "POST")
	r.Handle("/submit", NewSubmitHandler(articleRepo)).Methods("POST")
	r.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir("../../web/static/css/"))))
	r.PathPrefix("/images/").Handler(http.StripPrefix("/images/", http.FileServer(http.Dir("../../web/static/images/"))))
	//r.PathPrefix("/js/").Handler(http.StripPrefix("/js/", http.FileServer(http.Dir("../js/"))))
	return r
}
