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

type ExerciseRouter struct {
	storage   *storage.Storage
	presenter *presenters.Presenter
}

func newExercisesRouter(st *storage.Storage) *ExerciseRouter {
	return &ExerciseRouter{
		storage:   st,
		presenter: presenters.NewPresenter(),
	}
}

func RegisterExercisesRoutes(mux *http.ServeMux, st *storage.Storage) {
	router := newExercisesRouter(st)
	mux.HandleFunc("/exercises/create", router.create)
	mux.HandleFunc("/exercises/list", router.list)
	mux.HandleFunc("/exercises/", router.mutate)
}

func (router *ExerciseRouter) list(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, "No such endpoint", http.StatusNotFound)
		return
	}

	useCase := use_cases.NewExercisesUseCase(router.storage)
	result, err := useCase.List()
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
	w.WriteHeader(http.StatusOK)

	if _, err = w.Write(byteData); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	}

}

func (router *ExerciseRouter) create(w http.ResponseWriter, r *http.Request) {
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

	useCase := use_cases.NewExercisesUseCase(router.storage)
	req := requests.ExerciseRequest{
		Exercise: &models.Exercise{
			TitleEn: r.FormValue("title_en"),
			TitleRu: r.FormValue("title_ru"),
		},
		TagIds:     r.FormValue("tag_ids"),
		File:       &file,
		FileHeader: header,
	}
	result, err := useCase.Create(&req)
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
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	}
}

func (router *ExerciseRouter) mutate(w http.ResponseWriter, r *http.Request) {
	fmt.Println("mutate")
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

func (router *ExerciseRouter) update(id int, w http.ResponseWriter, r *http.Request) {
	fmt.Println("update")

	err := r.ParseMultipartForm(32 << 20) // 32MB limit
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}
	file, header, _ := r.FormFile("file")
	fmt.Println("file", file)
	fmt.Println("header", header)

	if file != nil {
		defer func() {
			err = file.Close()
			if err != nil {
				fmt.Printf("defer file close err: %s", err)
			}
		}()
	}

	useCase := use_cases.NewExercisesUseCase(router.storage)
	req := requests.ExerciseRequest{
		Exercise: &models.Exercise{
			TitleEn: r.FormValue("title_en"),
			TitleRu: r.FormValue("title_ru"),
		},
		TagIds: r.FormValue("tag_ids"),
	}
	if file != nil && header != nil {
		req.File = &file
		req.FileHeader = header
	}

	result, err := useCase.Update(id, &req)
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

func (router *ExerciseRouter) delete(id int, w http.ResponseWriter, r *http.Request) {
	useCase := use_cases.NewExercisesUseCase(router.storage)
	err := useCase.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	w.WriteHeader(http.StatusOK)
	if _, err = w.Write([]byte("successfully deleted")); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	}
}
