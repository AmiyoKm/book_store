package main

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/AmiyoKm/book_store/internal/store"
	"github.com/go-chi/chi/v5"
)

type bookKey string

const bookCtx bookKey = "book"

type createBookPayload struct {
	Title         string   `json:"title" validate:"required,max=255"`
	Author        string   `json:"author" validate:"required,max=255"`
	ISBN          string   `json:"isbn" validate:"required,numeric,len=13"`
	Price         int      `json:"price" validate:"required,gte=0,lte=100000"`
	Tags          []string `json:"tags" validate:"dive,max=30"`
	Description   string   `json:"description" validate:"max=1000"`
	CoverImageUrl string   `json:"cover_image_url" validate:"url"`
	Pages         int      `json:"pages" validate:"gte=1,lte=100000"`
	Stock         int      `json:"stock" validate:"required,gte=0"`
}

// createBookHandler godoc
//
//	@Summary		Creates a book
//	@Description	Creates a book
//	@Tags			book
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		createBookPayload	true	"Book details"
//	@Success		201		{object}	store.Book			"Book created"
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/books [post]
func (app *Application) createBookHandler(w http.ResponseWriter, r *http.Request) {
	var payload createBookPayload

	if err := readJson(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}
	ctx := r.Context()
	book := &store.Book{
		Title:         payload.Title,
		Author:        payload.Author,
		ISBN:          payload.ISBN,
		Price:         payload.Price,
		Tags:          payload.Tags,
		Description:   payload.Description,
		CoverImageUrl: payload.CoverImageUrl,
		Pages:         payload.Pages,
		Stock:         payload.Stock,
	}
	err := app.store.Books.Create(ctx, book)

	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := jsonResponse(w, http.StatusCreated, book); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// getBookHandler godoc
//
//	@Summary		Get a book
//	@Description	Get a book by its ID
//	@Tags			book
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int			true	"Book ID"
//	@Success		200	{object}	store.Book	"Book created"
//	@Failure		400	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/books/{id} [get]
func (app *Application) getBookHandler(w http.ResponseWriter, r *http.Request) {
	book := getBookFromContext(r)

	if err := jsonResponse(w, http.StatusOK, book); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

type updateBookPayload struct {
	Title         *string   `json:"title" validate:"omitempty,max=50"`
	Author        *string   `json:"author" validate:"omitempty,max=50"`
	ISBN          *string   `json:"isbn" validate:"omitempty,numeric,len=13"`
	Price         *int      `json:"price" validate:"omitempty,gte=0,lte=100000"`
	Tags          *[]string `json:"tags" validate:"omitempty,dive,max=30"`
	Description   *string   `json:"description" validate:"omitempty,max=1000"`
	CoverImageUrl *string   `json:"cover_image_url" validate:"omitempty,url"`
	Pages         *int      `json:"pages" validate:"omitempty,gte=1,lte=100000"`
	Stock         *int      `json:"stock" validate:"omitempty,gte=0"`
}

// updateBookHandler godoc
//
//	@Summary		Update a book
//	@Description	Update a book by its ID
//	@Tags			book
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int					true	"Book ID"
//	@Param			payload	body		updateBookPayload	true	"Update Book Payload"
//	@Success		200		{object}	store.Book
//	@Failure		400		{object}	error
//	@Failure		401		{object}	error
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/books/{id} [patch]
func (app *Application) updateBookHandler(w http.ResponseWriter, r *http.Request) {
	book := getBookFromContext(r)
	ctx := r.Context()

	var payload updateBookPayload

	if err := readJson(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if payload.Title != nil {
		book.Title = *payload.Title
	}
	if payload.Author != nil {
		book.Author = *payload.Author
	}
	if payload.ISBN != nil {
		book.ISBN = *payload.ISBN
	}
	if payload.Price != nil {
		book.Price = *payload.Price
	}
	if payload.Tags != nil {
		book.Tags = *payload.Tags
	}
	if payload.Description != nil {
		book.Description = *payload.Description
	}
	if payload.CoverImageUrl != nil {
		book.CoverImageUrl = *payload.CoverImageUrl
	}
	if payload.Pages != nil {
		book.Pages = *payload.Pages
	}
	if payload.Stock != nil {
		book.Stock = *payload.Stock
	}

	err := app.store.Books.Update(ctx, book)

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

	if err := jsonResponse(w, http.StatusCreated, book); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}

// deleteBookHandler godoc
//
//	@Summary		deletes a book
//	@Description	deletes a book by its ID
//	@Tags			book
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int		true	"Book ID"
//	@Success		204	{object}	string	"Book deleted"
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/books{id} [delete]
func (app *Application) deleteBookHandler(w http.ResponseWriter, r *http.Request) {
	book := getBookFromContext(r)
	ctx := r.Context()
	err := app.store.Books.Delete(ctx, book.ID)

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

func (app *Application) bookContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		paramID := chi.URLParam(r, "bookID")

		bookID, err := strconv.Atoi(paramID)
		if err != nil {
			app.notFoundError(w, r, err)
			return
		}
		ctx := r.Context()
		book, err := app.store.Books.GetByID(ctx, bookID)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrorNotFound):
				app.notFoundError(w, r, err)
			default:
				app.internalServerError(w, r, err)
			}
			return
		}

		ctx = context.WithValue(ctx, bookCtx, book)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getBookFromContext(r *http.Request) *store.Book {
	book, _ := r.Context().Value(bookCtx).(*store.Book)
	return book
}
