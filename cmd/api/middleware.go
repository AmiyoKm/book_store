package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/AmiyoKm/book_store/internal/store"
	"github.com/golang-jwt/jwt/v5"
)

type userCTX string

const userCtx userCTX = "user"

func (app *Application) AuthTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			app.unauthorizedError(w, r, fmt.Errorf("authorization header is missing"))
			return
		}
		parts := strings.Split(authHeader, " ")

		if len(parts) != 2 || parts[0] != "Bearer" {
			app.unauthorizedError(w, r, fmt.Errorf("authorization header is not correct"))
			return
		}

		token := parts[1]

		jwtToken, err := app.auth.ValidateToken(token)
		if err != nil {
			app.unauthorizedError(w, r, err)
			return
		}
		claims := jwtToken.Claims.(jwt.MapClaims)

		userID, err := strconv.ParseInt(fmt.Sprint(claims["sub"]), 10, 32)
		if err != nil {
			app.unauthorizedError(w, r, err)
			return
		}
		ctx := r.Context()
		user, err := app.getUser(ctx, userID)

		if err != nil {
			app.unauthorizedError(w, r, err)
			return
		}
		ctx = context.WithValue(ctx, userCtx, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *Application) getUser(ctx context.Context, userID int64) (*store.User, error) {
	user, err := app.store.Users.GetByID(ctx, int(userID))
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (app *Application) checkBookManipulationAuthority(requiredRole string , next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := getUserFromContext(r)
		allowed , err := app.checkRolePrecedence(r.Context(), user , requiredRole)
		if err != nil {
			app.internalServerError(w,r,err)
			return
		}
		if !allowed {
			app.forbiddenError(w, r)
			return
		}
		next.ServeHTTP(w,r)
	})
}

func (app *Application) checkRolePrecedence(ctx context.Context, user *store.User, roleName string) (bool, error) {
	role, err := app.store.Roles.GetByName(ctx, roleName)
	if err != nil {
		return false, err
	}

	if role.Level <= user.Role.Level {
		return true, nil
	}
	return false, nil
}
