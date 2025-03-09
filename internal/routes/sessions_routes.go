package routes

import (
	"bf_me/internal/presenters"
	"bf_me/internal/requests"
	"bf_me/internal/storage"
	"bf_me/internal/use_cases"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type SessionRouter struct {
	presenter *presenters.Presenter
	useCase   *use_cases.SessionsUseCase
}

func NewSessionsRouter(st *storage.Storage) *SessionRouter {
	return &SessionRouter{
		useCase:   use_cases.NewSessionsUseCase(st),
		presenter: presenters.NewPresenter(),
	}
}

func RegisterSessionsRoutes(mux *http.ServeMux, st *storage.Storage) {
	router := NewSessionsRouter(st)
	mux.HandleFunc("/vragneproidet", router.register)
	mux.HandleFunc("/login", router.login)
	mux.HandleFunc("/logout", AuthMiddleware(router.useCase, router.logout))
}

func (router *SessionRouter) register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "No such endpoint", http.StatusNotFound)
		return
	}

	var req requests.UserRequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	session, err := router.useCase.CreateUser(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	byteData, err := json.Marshal(router.presenter.Session(session))
	if err != nil {
		http.Error(w, fmt.Sprintf("json encoding err: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(byteData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (router *SessionRouter) login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "No such endpoint", http.StatusNotFound)
		return
	}

	var req requests.UserRequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	session, err := router.useCase.Create(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	byteData, err := json.Marshal(router.presenter.Session(session))
	if err != nil {
		http.Error(w, fmt.Sprintf("json encoding err: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(byteData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (router *SessionRouter) logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "No such endpoint", http.StatusNotFound)
		return
	}

	bearerToken := r.Header.Get("Authorization")
	sessionId := strings.TrimPrefix(bearerToken, "Bearer token=")

	err := router.useCase.Delete(sessionId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write([]byte("ok")); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
