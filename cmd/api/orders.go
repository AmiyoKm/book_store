package main

import (
	"context"
	"net/http"
	"strconv"

	"github.com/AmiyoKm/book_store/internal/store"
	"github.com/go-chi/chi/v5"
)

type orderCTX string

const orderCtx orderCTX = "order"

type createOrderPayload struct {
	TotalAmount     float64            `json:"total_amount" validate:"required,gt=0"`
	PaymentMethod   string             `json:"payment_method" validate:"required,oneof=cash_on_delivery Bkash credit_card"`
	ShippingAddress string             `json:"shipping_address" validate:"required,min=5"`
	Items           []OrderItemPayload `json:"items" validate:"required,dive"`
}
type OrderItemPayload struct {
	BookID   int     `json:"book_id" validate:"required,min=1"`
	Quantity int     `json:"quantity" validate:"required,min=1"`
	Price    float64 `json:"price" validate:"required,gt=0"`
}

func (app *Application) createOrderHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r)

	var payload createOrderPayload

	if err := readJson(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}
	if err := validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}
	order := &store.Order{
		UserID:          user.ID,
		TotalAmount:     payload.TotalAmount,
		PaymentMethod:   payload.PaymentMethod,
		ShippingAddress: payload.ShippingAddress,
		Items:           make([]store.OrderItem, len(payload.Items)),
	}

	for i, item := range payload.Items {
		order.Items[i] = store.OrderItem{
			BookID:   item.BookID,
			Quantity: item.Quantity,
			Price:    item.Price,
		}
	}

	ctx := r.Context()
	if err := app.store.Orders.Create(ctx, order); err != nil {
		app.internalServerError(w, r, err)
		return
	}
	if err := jsonResponse(w, http.StatusCreated, order); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
func (app *Application) getOrderHandler(w http.ResponseWriter, r *http.Request) {
	order := getOrderFromContext(r)
	
	if order == nil {
		app.notFoundError(w, r, nil)
		return
	}
	if err := jsonResponse(w, http.StatusAccepted, order); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *Application) orderContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "orderID")
		orderID, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			app.notFoundError(w, r, err)
			return
		}
		ctx := r.Context()
		order, err := app.store.Orders.GetByID(ctx, int(orderID))

		if err != nil {
			switch err {
			case store.ErrorNotFound:
				app.notFoundError(w, r, err)
				return
			default:
				app.internalServerError(w, r, err)
				return
			}
		}
		ctx = context.WithValue(ctx, orderCtx, order)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getOrderFromContext(r *http.Request) *store.Order {
	order, _ := r.Context().Value(orderCtx).(*store.Order)
	return order
}
