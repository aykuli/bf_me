package use_cases

import (
	"bf_me/internal/models"
	"bf_me/internal/requests"
	"gorm.io/gorm"
)

type ExerciseUseCase struct {
	db *gorm.DB
}

func NewExerciseUseCase(db *gorm.DB) *ExerciseUseCase {
	return &ExerciseUseCase{db}
}

func (euc *ExerciseUseCase) Create(req requests.CreateExerciseRequestBody) (*models.Exercise, error) {
	var tags []models.Tag
	euc.db.Find(&tags, req.Tag_ids)

	var e = &models.Exercise{
		TitleEn:  req.TitleEn,
		TitleRu:  req.TitleRu,
		FileUUID: "",
		Tags:     tags,
	}
	euc.db.Create(e)

	return e, nil
}
