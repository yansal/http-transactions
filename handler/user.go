package handler

import (
	"fmt"
	"net/http"

	"github.com/yansal/http-transactions/httpx"
	"github.com/yansal/http-transactions/manager"
	"github.com/yansal/http-transactions/payload"
)

type User struct{ m *manager.User }

func NewUser(manager *manager.User) *User {
	return &User{m: manager}
}

func (h *User) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h.serveHTTP(w, r)
	if err == nil {
		return
	}

	herr, ok := err.(httpx.Error)
	if !ok {
		herr = httpx.Error{Err: err, Code: http.StatusInternalServerError}
	}
	http.Error(w, fmt.Sprintf("%+v", herr.Error()), herr.Code)
}

func (h *User) serveHTTP(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPut {
		return httpx.Error{Code: http.StatusMethodNotAllowed}
	}
	p, err := payload.NewUser(r)
	if err != nil {
		return err
	}
	ctx := r.Context()
	return h.m.CreateUser(ctx, p)
}
