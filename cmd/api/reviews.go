package main

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/AmiyoKm/book_store/internal/store"
	"github.com/go-chi/chi/v5"
)

type reviewCTX string

const reviewCtx reviewCTX = "review"

type CreateReviewPayload struct {
	Content string `json:"content" validate:"required,min=1"`
	Rating  int    `json:"rating" validate:"required,min=1,max=5"`
}

// createReviewHandler godoc
//
//	@Summary		Create an Review
//	@Description	Create an Review
//	@Tags			review
//	@Accept			json
//	@Produce		json
//	@Param			bookID	path		int					true	"Book ID"
//	@Param			payload	body		CreateReviewPayload	true	"Create Review Payload"
//	@Success		201		{object}	store.Review		"Creates an Review"
//	@Failure		400		{object}	error				"Invalid request"
//	@Failure		500		{object}	error				"Server error"
//	@Security		ApiKeyAuth
//	@Router			/books/{bookID}/reviews [post]
func (app *Application) createReviewHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r)
	book := getBookFromContext(r)
	var payload CreateReviewPayload

	if err := readJson(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}
	if err := validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}
	ctx := r.Context()
	review := &store.Review{
		UserID:  user.ID,
		BookID:  book.ID,
		Content: payload.Content,
		Rating:  payload.Rating,
	}

	err := app.store.Reviews.Create(ctx, review)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	if err := jsonResponse(w, http.StatusCreated, review); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// getAllReviewsHandler godoc
//
//	@Summary		Get all Reviews
//	@Description	Get all reviews by Book ID
//	@Tags			review
//	@Accept			json
//	@Produce		json
//	@Param			bookID	path		int				true	"Book ID"
//	@Success		200		{array}		store.Review	"Get all Review"
//	@Failure		400		{object}	error			"Invalid request"
//	@Failure		500		{object}	error			"Server error"
//	@Security		ApiKeyAuth
//	@Router			/books/{bookID}/reviews [get]
func (app *Application) getAllReviewsHandler(w http.ResponseWriter, r *http.Request) {
	book := getBookFromContext(r)

	ctx := r.Context()
	reviews, err := app.store.Reviews.GetByBookID(ctx, book.ID)

	if err != nil {
		switch err {
		case store.ErrorNotFound:
			app.notFoundError(w, r, err)
			return
		default:
			app.notFoundError(w, r, err)
			return
		}
	}
	if err := jsonResponse(w, http.StatusOK, reviews); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// deleteReviewHandler godoc
//
//	@Summary		Delete a Review
//	@Description	Delete a review for a book by ID
//	@Tags			review
//	@Accept			json
//
//	@Param			bookID		path	int	true	"Book ID"
//	@Param			reviewID	path	int	true	"Review ID"
//
//	@Success		204			"Review deleted successfully"
//	@Failure		400			{object}	error	"Invalid request"
//	@Failure		500			{object}	error	"Server error"
//	@Security		ApiKeyAuth
//	@Router			/books/{bookID}/reviews/{reviewID} [delete]
func (app *Application) deleteReviewHandler(w http.ResponseWriter, r *http.Request) {
	review := getReviewFromContext(r)
	user := getUserFromContext(r)
	err := app.store.Reviews.Delete(r.Context(), review.ID, user.ID)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrorNotFound):
			app.notFoundError(w, r, store.ErrorNotFound)
			return
		default:
			app.internalServerError(w, r, err)
			return
		}
	}
	w.WriteHeader(http.StatusNoContent)
}

type updateReviewPayload struct {
	Content *string `json:"content" validate:"omitempty,min=1"`
	Rating  *int    `json:"rating" validate:"omitempty,min=1,max=5"`
}
// updateReviewHandler godoc
//	@Summary		Update a review
//	@Description	Update a review
//	@Tags			review
//	@Accept			json
//	@Produce		json
//	@Param			reviewID	path		int					true	"Review ID"
//	@Param			bookID		path		int					true	"Book ID"
//	@Param			payload		body		updateReviewPayload	true	"Update Review Payload"
//	@Success		200			{object}	store.Review		"Updated Review"
//	@Failure		400			{object}	error				"Invalid request"
//	@Failure		500			{object}	error				"Server error"
//	@Security		ApiKeyAuth
//	@Router			/books/{bookID}/reviews/{reviewID} [patch]
func (app *Application) updateReviewHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r)
	review := getReviewFromContext(r)
	book := getBookFromContext(r)
	var payload updateReviewPayload

	if err := readJson(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}
	if err := validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}
	newReview := &store.Review{
		ID:     review.ID,
		UserID: user.ID,
		BookID: book.ID,
	}
	if payload.Content != nil {
		newReview.Content = *payload.Content
	}
	if payload.Rating != nil {
		newReview.Rating = *payload.Rating
	}

	err := app.store.Reviews.Update(r.Context(), newReview)

	if err != nil {
		switch err {
		case store.ErrorNotFound:
			app.notFoundError(w, r, store.ErrorNotFound)
			return
		default:
			app.internalServerError(w, r, err)
			return
		}
	}
	if err := jsonResponse(w, http.StatusOK, newReview); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}

func (app *Application) reviewContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		paramID := chi.URLParam(r, "reviewID")
		reviewID, err := strconv.ParseInt(paramID, 10, 64)
		if err != nil {
			app.badRequestError(w, r, err)
			return
		}
		review, err := app.store.Reviews.GetByID(r.Context(), int(reviewID))
		if err != nil {
			switch err {
			case store.ErrorNotFound:
				app.notFoundError(w, r, store.ErrorNotFound)
				return
			default:
				app.internalServerError(w, r, err)
				return
			}
		}
		ctx := context.WithValue(r.Context(), reviewCtx, review)
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}

func getReviewFromContext(r *http.Request) *store.Review {
	review, _ := r.Context().Value(reviewCtx).(*store.Review)
	return review
}
