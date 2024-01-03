package renderer

import (
	"html/template"
	"net/http"
	"wbl0/WB_Task_L0/internal/models"
)

const (
	homePagePath  = "templates/home_page.html"
	orderPagePath = "templates/order_page.html"
)

func ShowOrderPage(w *http.ResponseWriter, order *models.Order) error {
	tmpl, err := template.ParseFiles(orderPagePath)
	if err != nil {
		return err
	}

	tmpl.Execute(*w, order)
	if err != nil {
		return err
	}

	return nil
}

func ShowHomePage(w *http.ResponseWriter, orders []models.Order) error {
	tmpl, err := template.ParseFiles(homePagePath)
	if err != nil {
		return err
	}

	tmpl.Execute(*w, orders)
	if err != nil {
		return err
	}

	return nil
}
