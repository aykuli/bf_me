package routes

import (
	"bf_me/internal/requests"
	"bf_me/internal/use_cases"
	"encoding/json"
	"net/http"

	"gorm.io/gorm"
)

type TagsRouter struct {
	db *gorm.DB
}

func newTagsRouter(db *gorm.DB) *TagsRouter {
	return &TagsRouter{db}
}

func RegisterTagsRoutes(mux *http.ServeMux, db *gorm.DB) {
	router := newTagsRouter(db)
	mux.HandleFunc("/api/v1", router.create)
	mux.HandleFunc("/api/v1/tags/list", router.list)
}

func (r *TagsRouter) list(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, "No such endpoint", http.StatusNotFound)
		return
	}

	useCase := use_cases.NewTagsUseCase(r.db)
	result, err := useCase.List()
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	byteData, err := json.Marshal(result)
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

func (r *TagsRouter) create(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "No such endpoint", http.StatusNotFound)
		return
	}

	var tBody requests.CreateTagRequestBody
	if err := json.NewDecoder(req.Body).Decode(&tBody); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	useCase := use_cases.NewTagsUseCase(r.db)
	result, err := useCase.Create(tBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	byteData, err := json.Marshal(result)
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
