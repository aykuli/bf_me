package routes

import (
	"bf_me/internal/models"
	"bf_me/internal/presenters"
	"bf_me/internal/requests"
	"bf_me/internal/storage"
	"bf_me/internal/use_cases"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strconv"

	"gorm.io/gorm"
)

type TrainingRouter struct {
	presenter   *presenters.Presenter
	useCase     *use_cases.TrainingsUseCase
	authUseCase *use_cases.SessionsUseCase
}

func newTrainingsRouter(st *storage.Storage) *TrainingRouter {
	return &TrainingRouter{
		presenter:   presenters.NewPresenter(),
		useCase:     use_cases.NewTrainingsUseCase(st),
		authUseCase: use_cases.NewSessionsUseCase(st),
	}
}

func RegisterTrainingRouter(mux *http.ServeMux, st *storage.Storage) {
	router := newTrainingsRouter(st)
	mux.HandleFunc("/api/v1/trainings/create", AuthMiddleware(router.authUseCase, router.create))
	mux.HandleFunc("/api/v1/trainings/list", AuthMiddleware(router.authUseCase, router.list))

	// action is enum of ["add", "remove"]
	mux.HandleFunc("/api/v1/trainings/{training_id}/{action}/block/{block_id}", AuthMiddleware(router.authUseCase, router.handleExercise))
	mux.HandleFunc("/api/v1/trainings/{id}/toggle_draft", AuthMiddleware(router.authUseCase, router.toggleDraft))
	mux.HandleFunc("/api/v1/trainings/{id}", AuthMiddleware(router.authUseCase, router.mux))
}

func (router *TrainingRouter) create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "No such endpoint", http.StatusNotFound)
		return
	}

	var req requests.TrainingRequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	result, err := router.useCase.Create(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	byteData, err := json.Marshal(router.presenter.Training(result))
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

func (router *TrainingRouter) list(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "No such endpoint", http.StatusNotFound)
		return
	}

	req := requests.FilterRequestBody{UpdatedAt: "DESC"}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	result, err := router.useCase.List(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	byteData, err := json.Marshal(router.presenter.Blocks(result))
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

func (router *TrainingRouter) handleExercise(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "No such endpoint", http.StatusNotFound)
		return
	}

	action := r.PathValue("action")
	if slices.Contains([]string{"add", "remove"}, action) == false {
		http.Error(w, "No such endpoint", http.StatusNotFound)
		return
	}

	blockIDStr := r.PathValue("block_id")
	exerciseIDStr := r.PathValue("exercise_id")
	blockID, err := strconv.ParseUint(blockIDStr, 10, 8)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	exerciseID, err := strconv.ParseUint(exerciseIDStr, 10, 8)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	var block *models.Block
	if action == "add" {
		block, err = router.useCase.AddBlockExercise(uint(blockID), uint(exerciseID))
	} else {
		block, err = router.useCase.RemoveBlockExercise(uint(blockID), uint(exerciseID))
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	byteData, err := json.Marshal(router.presenter.Block(block))
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

func (router *TrainingRouter) mux(w http.ResponseWriter, r *http.Request) {
	idInt, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, fmt.Errorf("invalid id provided: %s", err).Error(), http.StatusUnprocessableEntity)
		return
	}

	if r.Method == http.MethodGet {
		router.get(idInt, w, r)
		return
	}
	if r.Method == http.MethodPost {
		router.update(idInt, w, r)
		return
	}
	if r.Method == http.MethodDelete {
		router.delete(idInt, w, r)
		return
	}
}

func (router *TrainingRouter) toggleDraft(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "No such endpoint", http.StatusNotFound)
		return
	}

	idInt, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, fmt.Errorf("invalid id provided: %s", err).Error(), http.StatusUnprocessableEntity)
		return
	}

	result, err := router.useCase.ToggleDraft(idInt)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	byteData, err := json.Marshal(router.presenter.Block(result))
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

func (router *TrainingRouter) get(id int, w http.ResponseWriter, _ *http.Request) {
	result, err := router.useCase.Find(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	byteData, err := json.Marshal(router.presenter.Block(result))
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

func (router *TrainingRouter) update(id int, w http.ResponseWriter, r *http.Request) {
	var req requests.BlockRequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	result, err := router.useCase.Update(id, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	byteData, err := json.Marshal(router.presenter.Block(result))
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

func (router *TrainingRouter) delete(id int, w http.ResponseWriter, _ *http.Request) {
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
