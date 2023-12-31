package saveorder

import (
	"encoding/json"
	"fmt"
	"net/http"
	"wbl0/WB_Task_L0/internal/models"
	resp "wbl0/WB_Task_L0/internal/server/response"

	"github.com/go-chi/render"
)

type OrderSaver interface {
	AddOrder(nworder models.Order) error
}

func New(orderSaver OrderSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const ferr = "internal.server.handlers.save_order.New"
		var nwOrder models.Order

		err := json.NewDecoder(r.Body).Decode(&nwOrder)
		if err != nil {
			msg := fmt.Sprintf("%s: failed to read request body", ferr)
			render.JSON(w, r, resp.Send(resp.InternalError, msg, err.Error()))
			return
		}

		err = orderSaver.AddOrder(nwOrder)
		if err != nil {
			msg := fmt.Sprintf("%s: failed to add new order", ferr)
			render.JSON(w, r, resp.Send(resp.InternalError, msg, err.Error()))
			return
		}

		msg := fmt.Sprintf("order successfully added: %s", nwOrder.OrderUID.String())
		render.JSON(w, r, resp.Send(resp.OK, msg, ""))
	}
}