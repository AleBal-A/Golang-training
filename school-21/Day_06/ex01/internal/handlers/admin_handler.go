package handlers

import (
	"github.com/dgrijalva/jwt-go"
	"html/template"
	"net/http"
)

type AdminHandler struct {
	tmpl *template.Template
}

func NewAdminHandler(tmpl *template.Template) *AdminHandler {
	return &AdminHandler{tmpl: tmpl}
}

func (h *AdminHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	tokenStr := cookie.Value
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if !token.Valid {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if r.Method == http.MethodGet {
		h.tmpl.ExecuteTemplate(w, "admin.html", nil)
	}
}
