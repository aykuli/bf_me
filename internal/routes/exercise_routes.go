package routes

import (
	"bf_me/internal/models"
	"bf_me/internal/storage"
	"bf_me/internal/use_cases"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"
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

func (router *ExerciseRouter) list(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, "No such endpoint", http.StatusNotFound)
		return
	}

	useCase := use_cases.NewExercisesUseCase(router.storage)
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

func (router *ExerciseRouter) create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "No such endpoint", http.StatusNotFound)
		return
	}

	err := r.ParseMultipartForm(36 << 20)
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}
	titleEn := r.FormValue("title_en")
	titleRu := r.FormValue("title_ru")
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

	path, err := router.storage.S3.Upload(router.makeFilename(titleEn, header.Filename), file, header.Header.Get("Content-Type"))
	if err != nil {
		http.Error(w, fmt.Sprintf("minio upload file err: %s", err.Error()), http.StatusBadRequest)
		return
	}

	useCase := use_cases.NewExercisesUseCase(router.storage)
	result, err := useCase.Create(&models.Exercise{TitleEn: titleEn, TitleRu: titleRu, Filename: path}, r.FormValue("tag_ids"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	byteData, err := json.Marshal(result)
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

func (router *ExerciseRouter) sanitizeFilename(filename string) string {
	// Replace unsupported characters with underscores
	reg := regexp.MustCompile(`[^a-zA-Z0-9_.-]`)
	sanitized := reg.ReplaceAllString(filename, "_")
	sanitized = strings.ToLower(sanitized)

	// Trim leading and trailing spaces
	sanitized = strings.TrimSpace(sanitized)

	// Ensure the filename doesn't exceed the maximum length (255 characters)
	if len(sanitized) > 255 {
		sanitized = sanitized[:255]
	}

	// Ensure the filename is not empty
	if sanitized == "" {
		sanitized = "unnamed_file"
	}

	return sanitized
}

func (router *ExerciseRouter) makeFilename(title, filename string) string {
	sanitized := router.sanitizeFilename(title)
	return fmt.Sprintf("%s%s", sanitized, filepath.Ext(filename))
}
