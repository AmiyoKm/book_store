package main

import (
	"context"
	"net/http"
	"strconv"

	"github.com/AmiyoKm/book_store/internal/store"
	"github.com/go-chi/chi/v5"
)

type itemCTX string

const itemCtx itemCTX = "item"

type addToCartPayload struct {
	BookID   int `json:"book_id" validate:"required"`
	Quantity int `json:"quantity" validate:"required,min=1,max=10"`
}
type CartResponse struct {
	CartID int                      `json:"cart_id"`
	Items  []store.CartItemWithBook `json:"items"`
}

// addToCartHandler godoc
//
//	@Summary		Add book to cart
//	@Description	Add book to cart
//	@Tags			cart
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		addToCartPayload	true	"Add to Cart Payload"
//	@Success		201		{object}	map[string]string
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/carts [post]
func (app *Application) addToCartHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r)

	var payload addToCartPayload
	if err := readJson(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}
	if err := validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}
	ctx := r.Context()
	cart, err := app.store.Carts.GetOrCreateCart(ctx, user.ID)
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
	err = app.store.Carts.InsertOrUpdateCartItem(ctx, cart.ID, payload.BookID, payload.Quantity)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := jsonResponse(w, http.StatusCreated, map[string]string{
		"message": "Item added to cart",
	}); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// getCartHandler godoc
//
//	@Summary		Get Cart Items
//	@Description	Get Cart Items
//	@Tags			cart
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	CartResponse
//	@Failure		400	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/carts [get]
func (app *Application) getCartHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r)

	cart, err := app.store.Carts.GetOrCreateCart(r.Context(), user.ID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	items, err := app.store.Carts.GetCartItemsWithBooks(r.Context(), cart.ID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	returnCart := CartResponse{
		CartID: cart.ID,
		Items:  items,
	}

	if err := jsonResponse(w, http.StatusOK, returnCart); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// deleteItemHandler godoc
//
//	@Summary		Delete item from cart
//	@Description	Delete a specific item from the user's cart by item ID
//	@Tags			cart
//	@Param			itemID	path	int	true	"Cart Item ID"
//	@Produce		json
//	@Success		204	{string}	string	"Item deleted successfully"
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/carts/items/{itemID} [delete]
func (app *Application) deleteItemHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r)
	item := getItemFromContext(r)
	err := app.store.Carts.DeleteCartItem(r.Context(), user.ID, item.ID)

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
	w.WriteHeader(http.StatusNoContent)
}

// deleteCartHandler godoc
//
//	@Summary		Delete entire cart
//	@Description	Delete all items from the authenticated user's cart
//	@Tags			cart
//	@Produce		json
//	@Success		204	{string}	string	"Cart deleted successfully"
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/carts [delete]
func (app *Application) deleteCartHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r)
	err := app.store.Carts.DeleteCart(r.Context(), user.ID)

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
	w.WriteHeader(http.StatusNoContent)
}

type updateItemPayload struct {
	Quantity int `json:"quantity" validate:"required,min=1,max=10"`
}

// updateItemHandler godoc
//
//	@Summary		Update Cart Item Quantity
//	@Description	Update Cart Item Quantity
//	@Tags			cart
//	@Produce		json
//	@Param			itemID	path		int					true	"Item ID"
//	@Param			payload	body		updateItemPayload	true	"Quantity Payload"
//	@Success		200		{object}	store.CartItem		"CartItem Updated successfully"
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/carts/items/{itemID} [patch]
func (app *Application) updateItemHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r)
	item := getItemFromContext(r)

	var payload updateItemPayload
	if err := readJson(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}
	if err := validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	err := app.store.Carts.UpdateQuantity(r.Context(), payload.Quantity, item.ID, user.ID)
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
	item.Quantity = payload.Quantity
	if err := jsonResponse(w, http.StatusOK, item); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *Application) itemContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "itemID")
		itemID, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			app.notFoundError(w, r, err)
			return
		}
		ctx := r.Context()
		item, err := app.store.Carts.GetCartItem(ctx, int(itemID))

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
		ctx = context.WithValue(ctx, itemCtx, item)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getItemFromContext(r *http.Request) *store.CartItem {
	item, _ := r.Context().Value(itemCtx).(*store.CartItem)
	return item
}
