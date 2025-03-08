package requests

import (
	"bf_me/internal/models"
	"mime/multipart"
)

type CreateExerciseRequest struct {
	Exercise   *models.Exercise
	TagIds     string
	File       *multipart.File
	FileHeader *multipart.FileHeader
}

type UpdateExerciseRequestBody struct {
	TitleEn string `json:"title_en"`
	TitleRu string `json:"title_ru"`
}

type CreateTagRequestBody struct {
	TitleEn string `json:"title_en"`
	TitleRu string `json:"title_ru"`
}
