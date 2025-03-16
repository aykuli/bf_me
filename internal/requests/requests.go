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
	TitleEn string `json:"titleEn"`
	TitleRu string `json:"titleRu"`
}

type FilterExercisesRequestBody struct {
	UpdatedAt  string `json:"updatedAt"`
	CreatedAt  string `json:"createdAt"`
	Ids        bool   `json:"ids,omitempty"`
	Suggestion string `json:"suggestion,omitempty"`
}

type CreateTagRequestBody struct {
	TitleEn string `json:"titleEn"`
	TitleRu string `json:"titleRu"`
}

type UserRequestBody struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type BlockRequestBody struct {
	TitleEn       string `json:"titleEn"`
	TitleRu       string `json:"titleRu"`
	TotalDuration uint8  `json:"totalDuration"`
	OnTime        uint8  `json:"onTime"`
	RelaxTime     uint8  `json:"relaxTime"`
	Draft         bool   `json:"draft"`
}

type FilterBlocksRequestBody struct {
	Draft bool `json:"draft,omitempty"`
}
