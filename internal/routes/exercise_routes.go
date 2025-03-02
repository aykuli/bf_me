package routes

import (
	"bf_me/internal/requests"
	"bf_me/internal/storage"
	"bf_me/internal/use_cases"
	"encoding/json"
	"net/http"
)

type ExerciseRouter struct {
	storage *storage.Storage
}

func newExercisesRouter(st *storage.Storage) *ExerciseRouter {
	return &ExerciseRouter{storage: st}
}

func RegisterExercisesRoutes(mux *http.ServeMux, st *storage.Storage) {
	router := newExercisesRouter(st)
	mux.HandleFunc("/exercises/create", router.create)
	mux.HandleFunc("/exercises/list", router.list)
}

func (r *ExerciseRouter) list(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, "No such endpoint", http.StatusNotFound)
		return
	}

	useCase := use_cases.NewExercisesUseCase(r.storage)
	exercises, err := useCase.List()
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	byteData, err := json.Marshal(exercises)
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

func (r *ExerciseRouter) create(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "No such endpoint", http.StatusNotFound)
		return
	}

	var exerciseBody requests.CreateExerciseRequestBody
	if err := json.NewDecoder(req.Body).Decode(&exerciseBody); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	useCase := use_cases.NewExercisesUseCase(r.storage)
	resultExercise, err := useCase.Create(exerciseBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	byteData, err := json.Marshal(resultExercise)
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
