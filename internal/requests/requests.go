package requests

import (
	"bf_me/internal/models"
	"mime/multipart"
)

// @note Tips should be sent in form `str1,str2,str3`
type CreateExerciseRequest struct {
	Exercise   *models.Exercise
	TagIds     string
	File       *multipart.File
	FileHeader *multipart.FileHeader
	Tips       []string
}

// @note Tips should be sent in form `str1,str2,str3`
type UpdateExerciseRequestBody struct {
	TitleEn string   `json:"titleEn"`
	TitleRu string   `json:"titleRu"`
	Tips    []string `json:"tips"`
}

type FilterExercisesRequestBody struct {
	UpdatedAt  string `json:"updatedAt"`
	CreatedAt  string `json:"createdAt"`
	Ids        bool   `json:"ids,omitempty"`
	BlockIDs   []uint `json:"blockIds"`
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
}

type FilterBlocksRequestBody struct {
	BlockType  string `json:"blockType"` // draft, ready
	UpdatedAt  string `json:"updatedAt"`
	Suggestion string `json:"suggestion,omitempty"`
}
