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
	Price         float32  `json:"price" validate:"required,gte=0,lte=100000"`
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
//	@Success		200	{object}	store.Book	"Get Book "
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
	Price         *float32  `json:"price" validate:"omitempty,gte=0,lte=100000"`
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

// getBooksBySearchHandler godoc
//
//	@Summary		Search books
//	@Description	Search for books using a general query or specific filters such as title, author, tags, price range, and stock status.
//	@Tags			book
//	@Accept			json
//	@Produce		json
//	@Param			query		query		string		false	"Free-text search query"
//	@Param			title		query		string		false	"Filter by book title"
//	@Param			author		query		string		false	"Filter by author name"
//	@Param			tag			query		[]string	false	"Filter by tags"
//	@Param			min_price	query		number		false	"Minimum price filter"
//	@Param			max_price	query		number		false	"Maximum price filter"
//	@Param			in_stock	query		boolean		false	"Filter by stock status (true for in-stock, false for out-of-stock)"
//	@Success		200			{array}		store.Book
//	@Failure		500			{object}	error
//	@Security		ApiKeyAuth
//	@Router			/books/search [get]
func (app *Application) getBooksBySearchHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	filters := store.BooksBySearchPayload{
		Query:  q.Get("query"),
		Title:  q.Get("title"),
		Author: q.Get("author"),
		Tags:   q["tag"],
	}

	if min := q.Get("min_price"); min != "" {
		if v, err := strconv.ParseFloat(min, 32); err != nil {
			filters.MinPrice = float32(v)
		}
	}

	if max := q.Get("max_price"); max != "" {
		if v, err := strconv.ParseFloat(max, 64); err != nil {
			filters.MaxPrice = float32(v)
		}
	}
	if stock := q.Get("in_stock"); stock != "" {
		inStock := stock == "true"
		filters.InStock = &inStock
	}
	ctx := r.Context()
	books, err := app.store.Books.SearchByBooks(ctx, filters)

	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := jsonResponse(w, http.StatusOK, books); err != nil {
		app.internalServerError(w, r, err)
	}

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
