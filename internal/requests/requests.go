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

type FilterExercisesRequestBody struct {
	CreatedAt  string `json:"created_at"`
	BlockID    bool   `json:"block_id,omitempty"`
	Suggestion string `json:"suggestion,omitempty"`
}

type CreateTagRequestBody struct {
	TitleEn string `json:"title_en"`
	TitleRu string `json:"title_ru"`
}

type UserRequestBody struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type BlockRequestBody struct {
	TitleEn       string `json:"title_en"`
	TitleRu       string `json:"title_ru"`
	TotalDuration uint8  `json:"total_duration"`
	OnTime        uint8  `json:"on_time"`
	RelaxTime     uint8  `json:"relax_time"`
	Draft         bool   `json:"draft"`
	ExercisesIds  []int  `json:"exercises_ids"`
}

type FilterBlocksRequestBody struct {
	Draft bool `json:"draft,omitempty"`
}
