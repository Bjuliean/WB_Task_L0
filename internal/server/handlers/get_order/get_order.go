package getorder

import (
	"fmt"
	"net/http"
	"wbl0/WB_Task_L0/internal/models"
	"wbl0/WB_Task_L0/internal/server/renderer"
	resp "wbl0/WB_Task_L0/internal/server/response"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type OrderGetter interface {
	GetOrder(uid uuid.UUID) (*models.Order, error)
}

func New(orderGetter OrderGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const ferr = "internal.server.handlers.get_order.New"

		order_uid := chi.URLParam(r, "order_uid")
		if order_uid == "" {
			msg := fmt.Sprintf("%s: order_uid is empty", ferr)
			render.JSON(w, r, resp.Send(resp.BadRequest, msg, ""))
			return
		}

		uid_val, err := uuid.Parse(order_uid)
		if err != nil {
			msg := fmt.Sprintf("%s: order_uid parsing failed", ferr)
			render.JSON(w, r, resp.Send(resp.BadRequest, msg, err.Error()))
			return
		}

		res, err := orderGetter.GetOrder(uid_val)

		if err != nil {
			msg := fmt.Sprintf("%s: order_uid not found", ferr)
			render.JSON(w, r, resp.Send(resp.OK, msg, err.Error()))
			return
		}

		renderer.ShowOrderPage(&w, res)
	}
}
