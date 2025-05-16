package main

import (
	"net/http"
)

func (app *Application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("internal server error :", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	writeJsonError(w, http.StatusInternalServerError, "internal server error")
}
func (app *Application) notFoundError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("not found error :", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	writeJsonError(w, http.StatusNotFound, err.Error())

}
func (app *Application) badRequestError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("bad request error :", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	writeJsonError(w, http.StatusBadRequest, err.Error())

}
func (app *Application) conflictError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("conflict error :", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	writeJsonError(w, http.StatusConflict, err.Error())
}

func (app *Application) unauthorizedError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("unauthorized error :", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	writeJsonError(w, http.StatusUnauthorized, err.Error())
}

func (app *Application) forbiddenError(w http.ResponseWriter, r *http.Request) {
	app.logger.Errorw("forbidden error :", "method", r.Method, "path", r.URL.Path, "error", "lower level role")
	writeJsonError(w, http.StatusForbidden, "lower level role , not allowed")
}
