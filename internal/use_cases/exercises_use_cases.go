package use_cases

import (
	"bf_me/internal/models"
	"bf_me/internal/requests"
	"gorm.io/gorm"
)

type ExercisesUseCase struct {
	db *gorm.DB
}

func NewExercisesUseCase(db *gorm.DB) *ExercisesUseCase {
	return &ExercisesUseCase{db}
}

// todo docs, pagination, filter by fields, fetch tags
func (euc *ExercisesUseCase) List() ([]*models.Exercise, error) {
	var exercises []*models.Exercise
	result := euc.db.Order("updated_at DESC").Find(&exercises)
	return exercises, result.Error
}

func (euc *ExercisesUseCase) Create(req requests.CreateExerciseRequestBody) (*models.Exercise, error) {
	var tags []models.Tag
	if len(req.Tag_ids) != 0 {
		euc.db.Find(&tags, req.Tag_ids)
	}

	var e = &models.Exercise{
		TitleEn:  req.TitleEn,
		TitleRu:  req.TitleRu,
		FileUUID: "",
		Tags:     tags,
	}
	result := euc.db.Create(e)
	return e, result.Error
}
