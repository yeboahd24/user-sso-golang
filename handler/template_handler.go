package handler

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type TemplateHandler struct {
	templates *template.Template
}

func NewTemplateHandler() (*TemplateHandler, error) {
	// Parse all templates from the template directory
	tmpl, err := template.ParseGlob(filepath.Join("template", "*.html"))
	if err != nil {
		log.Fatalf("Failed to parse templates: %v", err)
		return nil, err
	}

	return &TemplateHandler{
		templates: tmpl,
	}, nil
}

func (h *TemplateHandler) LoginPage(c *gin.Context) {
	data := gin.H{
		"Error": c.Query("error"),
	}

	c.HTML(http.StatusOK, "login.html", data)
}
