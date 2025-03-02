package requests

import "io"

type CreateExerciseRequestBody struct {
	TitleEn     string `json:"title_en"`
	TitleRu     string `json:"title_ru"`
	TagIds      []int  `json:"tag_ids"`
	File        io.Reader
	fileSize    int64
	fileName    string
	ContentType string
}

func (req *CreateExerciseRequestBody) GetFileName() string {
	return req.fileName
}

type CreateTagRequestBody struct {
	TitleEn string `json:"title_en"`
	TitleRu string `json:"title_ru"`
}
