package routes

import (
	"bf_me/internal/models"
	"bf_me/internal/presenters"
	"bf_me/internal/requests"
	"bf_me/internal/storage"
	"bf_me/internal/use_cases"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type ExercisesRouter struct {
	presenter   *presenters.Presenter
	useCase     *use_cases.ExercisesUseCase
	authUseCase *use_cases.SessionsUseCase
}

func newExercisesRouter(st *storage.Storage) *ExercisesRouter {
	return &ExercisesRouter{
		presenter:   presenters.NewPresenter(),
		useCase:     use_cases.NewExercisesUseCase(st),
		authUseCase: use_cases.NewSessionsUseCase(st),
	}
}

func RegisterExercisesRoutes(mux *http.ServeMux, st *storage.Storage) {
	router := newExercisesRouter(st)
	mux.HandleFunc("/exercises/create", AuthMiddleware(router.authUseCase, router.create))
	mux.HandleFunc("/exercises/list", AuthMiddleware(router.authUseCase, router.list))
	mux.HandleFunc("/exercises/", AuthMiddleware(router.authUseCase, router.mutate))
}

func (router *ExercisesRouter) list(w http.ResponseWriter, r *http.Request) {
	result, err := router.useCase.List()
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	byteData, err := json.Marshal(router.presenter.Exercises(result))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)

	if _, err = w.Write(byteData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (router *ExercisesRouter) create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "No such endpoint", http.StatusNotFound)
		return
	}

	err := r.ParseMultipartForm(32 << 20) // 32MB limit
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, fmt.Sprintf("retrieving file err: %s", err.Error()), http.StatusBadRequest)
		return
	}
	defer func() {
		if err = file.Close(); err != nil {
			fmt.Printf("defer file close err: %s", err)
		}
	}()

	req := requests.CreateExerciseRequest{
		Exercise: &models.Exercise{
			TitleEn: r.FormValue("title_en"),
			TitleRu: r.FormValue("title_ru"),
		},
		TagIds:     r.FormValue("tag_ids"),
		File:       &file,
		FileHeader: header,
	}
	result, err := router.useCase.Create(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	byteData, err := json.Marshal(router.presenter.Exercise(result))
	if err != nil {
		http.Error(w, fmt.Sprintf("json encoding err: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if _, err = w.Write(byteData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (router *ExercisesRouter) mutate(w http.ResponseWriter, r *http.Request) {
	// Extract the ID from the URL path
	path := strings.TrimPrefix(r.URL.Path, "/exercises/")
	id := strings.TrimSuffix(path, "/")

	if id == "" {
		http.Error(w, "invalid id provided", http.StatusBadRequest)
		return
	}
	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, fmt.Errorf("invalid id provided: %s", err).Error(), http.StatusUnprocessableEntity)
		return
	}

	if r.Method == http.MethodDelete {
		router.delete(idInt, w, r)
		return
	}
	if r.Method == http.MethodPost {
		router.update(idInt, w, r)
		return
	}
}

func (router *ExercisesRouter) update(id int, w http.ResponseWriter, r *http.Request) {
	var req requests.UpdateExerciseRequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	result, err := router.useCase.Update(id, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	byteData, err := json.Marshal(router.presenter.Exercise(result))
	if err != nil {
		http.Error(w, fmt.Sprintf("json encoding err: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if _, err = w.Write(byteData); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	}
}

func (router *ExercisesRouter) delete(id int, w http.ResponseWriter, r *http.Request) {
	err := router.useCase.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	w.WriteHeader(http.StatusOK)
	if _, err = w.Write([]byte("successfully deleted")); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
