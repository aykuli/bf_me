package use_cases

import (
	"bf_me/internal/models"
	"bf_me/internal/storage"
)

type ExercisesUseCase struct {
	storage *storage.Storage
}

func NewExercisesUseCase(st *storage.Storage) *ExercisesUseCase {
	return &ExercisesUseCase{storage: st}
}

// todo docs, pagination, filter by fields, fetch tags
func (euc *ExercisesUseCase) List() ([]*models.Exercise, error) {
	var exercises []*models.Exercise
	result := euc.storage.DB.Order("updated_at DESC").Find(&exercises)
	return exercises, result.Error
}

// @note tagIds is the string, containing tag ids separated by comma
func (euc *ExercisesUseCase) Create(e *models.Exercise, tagIds string) (*models.Exercise, error) {
	//var tags []models.Tag
	//if len(req.TagIds) != 0 {
	//	euc.storage.DB.Find(&tags, req.TagIds)
	//}

	result := euc.storage.DB.Create(e)
	return e, result.Error
}
