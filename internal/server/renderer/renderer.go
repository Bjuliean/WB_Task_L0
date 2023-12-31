package renderer

import (
	"html/template"
	"net/http"
	"wbl0/WB_Task_L0/internal/models"
)

const(
	homePagePath = "templates/home_page.html"
)

func ShowHomePage(w *http.ResponseWriter, Orders []models.Order) error {
	tmpl, err := template.ParseFiles(homePagePath)
	if err != nil {
		return err
	}

	tmpl.Execute(*w, Orders)
	if err != nil {
		return err
	}

	return nil
}