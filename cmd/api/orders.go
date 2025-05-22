package main

import (
	"context"
	"fmt"
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

// createOrderHandler godoc
//
//	@Summary		Create an order
//	@Description	Create an order
//	@Tags			order
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		createOrderPayload	true	"Create Order Payload"
//	@Success		201		{object}	store.Order			"Creates an order"
//	@Failure		400		{object}	error				"Invalid request"
//	@Failure		500		{object}	error				"Server error"
//	@Security		ApiKeyAuth
//	@Router			/orders [post]
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

// getOrderHandler godoc
//
//	@Summary		Get a order
//	@Description	Get a order by its ID
//	@Tags			order
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int			true	"order ID"
//	@Success		200	{object}	store.Order	"Get Order"
//	@Failure		400	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/orders/{id} [get]
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

// getAllOrderHandler godoc
//
//	@Summary		Get all orders
//	@Description	Get all orders
//	@Tags			order
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}		store.Order	"List of all orders"
//	@Failure		400	{object}	error		"Invalid request"
//	@Failure		500	{object}	error		"Server error"
//	@Security		ApiKeyAuth
//	@Router			/orders [get]

func (app *Application) getAllOrdersHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r)

	if user == nil {
		app.unauthorizedError(w, r, fmt.Errorf("unauthorized user"))
		return
	}
	ctx := r.Context()
	orders, err := app.store.Orders.Get(ctx, user.ID)

	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	if err := jsonResponse(w, http.StatusOK, orders); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

type updateOrderPayload struct {
	ShippingAddress *string `json:"shipping_address" validate:"min=1"`
}

// updateOderHandler godoc
//
//	@Summary		Update an order by User
//	@Description	Update an order by User
//	@Tags			order
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int					true	"Order ID"
//	@Param			payload	body		updateOrderPayload	true	"Update Order Payload by User"
//	@Success		200		{object}	store.Order			"Updated Order"
//	@Failure		400		{object}	error				"Invalid request"
//	@Failure		500		{object}	error				"Server error"
//	@Security		ApiKeyAuth
//	@Router			/orders/{id} [patch]
func (app *Application) updateOderHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r)
	order := getOrderFromContext(r)

	if user == nil {
		app.unauthorizedError(w, r, fmt.Errorf("unauthorized user"))
		return
	}
	if order == nil {
		app.notFoundError(w, r, fmt.Errorf("order not found"))
		return
	}
	var payload updateOrderPayload

	if err := readJson(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}
	if err := validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if payload.ShippingAddress != nil {
		order.ShippingAddress = *payload.ShippingAddress
	}

	err := app.store.Orders.Update(r.Context(), order)
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
	if err := jsonResponse(w, http.StatusOK, order); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}

type updateOrderAdminPayload struct {
	ShippingAddress *string `json:"shipping_address" validate:"omitempty,min=1"`
	PaymentMethod   *string `json:"payment_method" validate:"omitempty,oneof=cash_on_delivery Bkash credit_card"`
	Status          *string `json:"status" validate:"omitempty,oneof=pending processing shipped delivered cancelled returned failed refunded"`
}

// updateAdminOrderHandler godoc
//
//	@Summary		Update an order by Admin
//	@Description	Update an order by Admin
//	@Tags			order
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int						true	"Order ID"
//	@Param			payload	body		updateOrderAdminPayload	true	"Update Order Payload by Admin"
//	@Success		200		{object}	store.Order				"Updated Order"
//	@Failure		400		{object}	error					"Invalid request"
//	@Failure		500		{object}	error					"Server error"
//	@Security		ApiKeyAuth
//	@Router			/admin/orders/{id} [patch]
func (app *Application) updateAdminOrderHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r)
	order := getOrderFromContext(r)

	if user == nil {
		app.unauthorizedError(w, r, fmt.Errorf("unauthorized user"))
		return
	}
	if order == nil {
		app.notFoundError(w, r, fmt.Errorf("order not found"))
		return
	}
	var payload updateOrderAdminPayload

	if err := readJson(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}
	if err := validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if payload.ShippingAddress != nil {
		order.ShippingAddress = *payload.ShippingAddress
	}
	if payload.Status != nil {
		order.Status = *payload.Status
	}
	if payload.PaymentMethod != nil {
		order.PaymentMethod = *payload.PaymentMethod
	}

	err := app.store.Orders.Update(r.Context(), order)
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
	if err := jsonResponse(w, http.StatusOK, order); err != nil {
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
