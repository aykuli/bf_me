package routes

import (
	"bf_me/internal/requests"
	"bf_me/internal/use_cases"
	"encoding/json"
	"net/http"

	"gorm.io/gorm"
)

type ExerciseRouter struct {
	db *gorm.DB
}

func newRouter(db *gorm.DB) *ExerciseRouter {
	return &ExerciseRouter{db}
}

func RegisterExercisesRoutes(mux *http.ServeMux, db *gorm.DB) {
	router := newRouter(db)
	mux.HandleFunc("/exercises/create", router.create)
	//mux.HandleFunc("/exercises", listExercises)
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

	use_case := use_cases.NewExerciseUseCase(r.db)
	resultExercise, err := use_case.Create(exerciseBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
