package getorders

import (
	"net/http"
	"wbl0/WB_Task_L0/internal/models"
	"wbl0/WB_Task_L0/internal/server/renderer"
)

type OrdersGetter interface {
	GetOrders() ([]models.Order, error)
}

func New(ordersGetter OrdersGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, _ := ordersGetter.GetOrders()

		renderer.ShowHomePage(&w, res)
	}
}
