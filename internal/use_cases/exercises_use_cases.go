package use_cases

import (
	"bf_me/internal/models"
	"bf_me/internal/requests"
	"bf_me/internal/storage"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

type ExercisesUseCase struct {
	storage *storage.Storage
}

func NewExercisesUseCase(st *storage.Storage) *ExercisesUseCase {
	return &ExercisesUseCase{storage: st}
}

func (euc *ExercisesUseCase) List(req *requests.FilterExercisesRequestBody) ([]*models.Exercise, error) {
	var exercises []*models.Exercise

	if len(req.BlockIDs) != 0 {
		result := euc.storage.DB.Joins("INNER JOIN exercise_blocks ON exercise_blocks.exercise_id = exercises.id").
			Where("exercise_blocks.block_id IN ?", req.BlockIDs).Find(&exercises)
		return exercises, result.Error
	}

	result := euc.storage.DB.Order(fmt.Sprintf("updated_at %s", req.UpdatedAt)).Find(&exercises)
	return exercises, result.Error
}

func (euc *ExercisesUseCase) Create(req *requests.CreateExerciseRequest) (*models.Exercise, error) {
	e := req.Exercise
	// todo if filename is already used, try another one. Maxtries = 5
	path, err := euc.storage.S3.Upload(euc.makeFilename(e.TitleEn, req.FileHeader.Filename), *req.File, req.FileHeader.Header.Get("Content-Type"))
	if err != nil {
		return nil, fmt.Errorf("minio upload file err: %s", err)
	}
	e.Filename = path
	result := euc.storage.DB.Create(e)
	return e, result.Error
}

func (euc *ExercisesUseCase) Find(id int) (*models.Exercise, error) {
	var e models.Exercise
	result := euc.storage.DB.First(&e, id)
	return &e, result.Error
}

func (euc *ExercisesUseCase) Update(id int, req *requests.UpdateExerciseRequestBody) (*models.Exercise, error) {
	var e *models.Exercise
	result := euc.storage.DB.First(&e, id)
	if req.TitleRu != "" {
		e.TitleRu = req.TitleRu
	}
	if req.TitleEn != "" {
		e.TitleEn = req.TitleEn
	}
	if len(req.Tips) != 0 {
		e.Tips = req.Tips
	}

	result = euc.storage.DB.Save(e)
	return e, result.Error
}

func (euc *ExercisesUseCase) makeFilename(title, filename string) string {
	sanitized := euc.sanitizeFilename(title)
	return fmt.Sprintf("%s%s", sanitized, filepath.Ext(filename))
}

func (euc *ExercisesUseCase) sanitizeFilename(filename string) string {
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

func (euc *ExercisesUseCase) Delete(id int) error {
	var e *models.Exercise
	result := euc.storage.DB.First(&e, id)

	spl := strings.Split(e.Filename, "/")
	err := euc.storage.S3.Delete(spl[0])
	if err != nil {
		fmt.Printf("minio file delete err: %s", err)
	}

	result = euc.storage.DB.Delete(&models.Exercise{}, id)
	return result.Error
}
