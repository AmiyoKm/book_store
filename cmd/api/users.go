package main

import (
	"net/http"

	"github.com/AmiyoKm/book_store/internal/store"
)

func getUserFromContext(r *http.Request) *store.User {
	user , _ := r.Context().Value(userCtx).(*store.User)
	return user
}
